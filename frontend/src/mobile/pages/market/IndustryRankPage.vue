<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount, watch } from 'vue'
import { GetIndustryRank, GetIndustryMoneyRankSina } from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import IndustryMoneyCard from '../../components/cards/IndustryMoneyCard.vue'

// 对齐桌面端 market.vue「行业排名」Tab 的 4 个子分类：
// 行业涨幅 → GetIndustryRank(sort,cnt)；行业资金/证监会/概念 → GetIndustryMoneyRankSina(fenlei,'netamount')
const rankTabs = [
  { label: '行业涨幅', value: 'change' },
  { label: '行业资金', value: 'money' },
  { label: '证监会行业', value: 'csrc' },
  { label: '概念板块', value: 'concept' },
]

// money/csrc/concept → 新浪 fenlei 参数
const FENLEI_MAP = { money: '0', concept: '1', csrc: '2' }

const activeTab = ref('change')
const rankData = ref({}) // { tab: item[] }
const loading = ref(false)

const currentList = computed(() => rankData.value[activeTab.value] || [])
const isChangeTab = computed(() => activeTab.value === 'change')

// 数据安全提取
// GetIndustryRank 返回 map，行数组在 result.data（后端测试 market_news_api_test.go:54 确认）；
// GetIndustryMoneyRankSina 返回数组本身。
function pickArray(res) {
  const arr = Array.isArray(res) ? res : (res?.data ?? [])
  return (Array.isArray(arr) ? arr : []).filter(it => it && typeof it === 'object')
}

async function loadChange() {
  try {
    const res = await GetIndustryRank('0', 150)
    rankData.value = { ...rankData.value, change: pickArray(res) }
  } catch (e) {
    console.error('加载行业涨幅失败:', e)
  }
}

async function loadMoney(fenlei) {
  try {
    const res = await GetIndustryMoneyRankSina(fenlei, 'netamount')
    rankData.value = { ...rankData.value, [activeTab.value]: pickArray(res) }
  } catch (e) {
    console.error('加载行业资金失败:', e)
  }
}

// 按当前 tab 拉取（带 loading），切 tab 懒加载，已有缓存则跳过。
async function loadCurrent(force = false) {
  const tab = activeTab.value
  if (!force && rankData.value[tab]?.length) return
  loading.value = true
  if (tab === 'change') await loadChange()
  else await loadMoney(FENLEI_MAP[tab])
  loading.value = false
}

async function handleRefresh() {
  await loadCurrent(true)
}

watch(activeTab, () => loadCurrent())

// 轮询只刷当前 tab（行业涨幅为主场景，对齐桌面 market.industry 10s）
function refreshFeed() {
  if (activeTab.value === 'change') loadChange()
  else loadMoney(FENLEI_MAP[activeTab.value])
}

onBeforeMount(() => {
  loadCurrent()
  registerFeed('mobile-industry-rank', {
    fetch: refreshFeed,
    intervalMs: 10 * 1000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-industry-rank')
})
</script>

<template>
  <div class="page">
    <div class="rank-tabs">
      <MTabs v-model="activeTab" :tabs="rankTabs" />
    </div>

    <MPullRefresh class="rank-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !currentList.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else>
          <div v-if="currentList.length" class="rank-list">
            <!-- 行业涨幅：页内简洁行列表（腾讯数据 bd_*） -->
            <template v-if="isChangeTab">
              <div v-for="(it, i) in currentList" :key="it.bd_code || i" class="change-item">
                <span class="change-idx">{{ i + 1 }}</span>
                <div class="change-main">
                  <div class="change-name">{{ it.bd_name }}</div>
                  <div class="change-leader" v-if="it.nzg_name">
                    领涨 {{ it.nzg_name }}
                    <PercentTag :value="Number(it.nzg_zdf)" size="small" />
                  </div>
                </div>
                <div class="change-pct">
                  <PercentTag :value="Number(it.bd_zdf)" size="medium" bold />
                </div>
              </div>
            </template>

            <!-- 行业资金/证监会/概念：IndustryMoneyCard（新浪数据） -->
            <template v-else>
              <IndustryMoneyCard
                v-for="(it, i) in currentList"
                :key="it.name || i"
                :item="it"
              />
            </template>
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

.rank-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.rank-refresh {
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

/* 列表容器：卡片间分隔线（对齐 HotStockPage .hot-list） */
.rank-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.rank-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

/* 行业涨幅行 */
.change-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  padding: var(--m-space-md);
}

.change-idx {
  width: 22px;
  flex-shrink: 0;
  text-align: center;
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.change-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.change-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.change-leader {
  display: flex;
  align-items: center;
  gap: var(--m-space-xs);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.change-pct {
  flex-shrink: 0;
}
</style>
