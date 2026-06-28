<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount, watch } from 'vue'
import { GetMoneyRankSina } from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import StockMoneyFlowCard from '../../components/cards/StockMoneyFlowCard.vue'

// 对齐桌面端 market.vue「个股资金流向」Tab 的 9 个子分类（rankTable.vue / GetMoneyRankSina）。
// 每个 sort 维度不同：净流入/流出/率、主力(r0)、散户(r3)。卡片按当前 sort 高亮对应指标。
const flowTabs = [
  { label: '净流入额', value: 'netamount' },
  { label: '流出额', value: 'outamount' },
  { label: '净流入率', value: 'ratioamount' },
  { label: '主力净流入额', value: 'r0_net' },
  { label: '主力流出', value: 'r0_out' },
  { label: '主力净流入率', value: 'r0_ratio' },
  { label: '散户净流入额', value: 'r3_net' },
  { label: '散户流出', value: 'r3_out' },
  { label: '散户净流入率', value: 'r3_ratio' },
]

// sort → 高亮指标定义（label/kind/get）。kind:'money' 用 formatMoney，'rate' 用 PercentTag。
// 比率字段 0~1，get 里 *100 转百分数。
const METRIC_MAP = {
  netamount: { label: '净流入', kind: 'money', get: it => Number(it.netamount) },
  outamount: { label: '流出额', kind: 'money', get: it => Number(it.outamount) },
  ratioamount: { label: '净流入率', kind: 'rate', get: it => Number(it.ratioamount) * 100 },
  r0_net: { label: '主力净流入', kind: 'money', get: it => Number(it.r0_net) },
  r0_out: { label: '主力流出', kind: 'money', get: it => Number(it.r0_out) },
  r0_ratio: { label: '主力净流入率', kind: 'rate', get: it => Number(it.r0_ratio) * 100 },
  r3_net: { label: '散户净流入', kind: 'money', get: it => Number(it.r3_net) },
  r3_out: { label: '散户流出', kind: 'money', get: it => Number(it.r3_out) },
  r3_ratio: { label: '散户净流入率', kind: 'rate', get: it => Number(it.r3_ratio) * 100 },
}

const activeTab = ref('netamount')
const currentMetric = computed(() => METRIC_MAP[activeTab.value])

const flowData = ref({}) // { sort: item[] }
const loading = ref(false)

const currentList = computed(() => flowData.value[activeTab.value] || [])

// 数据安全提取：result 本身是数组
function pickArray(res) {
  return (Array.isArray(res) ? res : []).filter(it => it && typeof it === 'object')
}

async function loadFlow(sort) {
  try {
    const res = await GetMoneyRankSina(sort)
    flowData.value = { ...flowData.value, [sort]: pickArray(res) }
  } catch (e) {
    console.error('加载个股资金失败:', e)
  }
}

// 按当前 tab 拉取（带 loading），切 tab 懒加载，已有缓存则跳过。
async function loadCurrent(force = false) {
  const sort = activeTab.value
  if (!force && flowData.value[sort]?.length) return
  loading.value = true
  await loadFlow(sort)
  loading.value = false
}

async function handleRefresh() {
  await loadCurrent(true)
}

watch(activeTab, () => loadCurrent())

// 轮询只刷当前 tab（对齐桌面 rankTable 60s）
function refreshFeed() {
  loadFlow(activeTab.value)
}

onBeforeMount(() => {
  loadCurrent()
  registerFeed('mobile-stock-money-flow', {
    fetch: refreshFeed,
    intervalMs: 60 * 1000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-stock-money-flow')
})
</script>

<template>
  <div class="page">
    <div class="flow-tabs">
      <MTabs v-model="activeTab" :tabs="flowTabs" />
    </div>

    <MPullRefresh class="flow-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !currentList.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else>
          <div v-if="currentList.length" class="flow-list">
            <StockMoneyFlowCard
              v-for="(item, i) in currentList"
              :key="item.symbol || i"
              :item="item"
              :metric="currentMetric"
            />
          </div>
          <MEmpty v-else description="暂无数据" />
        </template>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.flow-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

/* MPullRefresh 占满 tab 之外的剩余高度，避免底部被 overflow:hidden 裁切。 */
.flow-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

/* 卡片列表：卡片间分隔线（对齐 HotStockPage .hot-list） */
.flow-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.flow-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}
</style>
