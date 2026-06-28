<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount } from 'vue'
import { GetBKFundFlowTopListByDate, GetBKFundFlowListByDate } from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MiniSparkline from '../../components/charts/MiniSparkline.vue'
import { formatMoney } from '../../composables/useFormat'

// 对齐桌面端 market.vue「板块资金流向」(bkFundFlowChart.vue)。
// 桌面端是 echarts 多板块折线对比；移动端改为「流入 TOP + 流出 TOP」两段列表，
// 每行带迷你折线(当日 netInflow 走势)，信息可操作、窄屏可读。

const TOP_LIMIT = 20 // 两段各取前 20，避免 sparkline 的 N 次请求过多

const todayStr = (() => {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
})()

const all = ref([]) // 全量 TopList（已按 netInflow DESC）
const sparkMap = ref({}) // { code: netInflow[] }
const loading = ref(false)

// 流入 TOP（netInflow>0，已是降序）、流出 TOP（<0，反转使流出最大者居顶）
const inflowList = computed(() => all.value.filter(it => Number(it.netInflow) > 0).slice(0, TOP_LIMIT))
const outflowList = computed(() =>
  all.value.filter(it => Number(it.netInflow) < 0).slice(-TOP_LIMIT).reverse()
)

function pickArray(res) {
  return (Array.isArray(res) ? res : []).filter(it => it && typeof it === 'object')
}

// 拉取每个 TOP 行的当日折线（流入+流出各 TOP_LIMIT 条 → 最多 2*TOP_LIMIT 次请求）
async function loadSparklines(list) {
  await Promise.all(list.map(async it => {
    if (!it.code || sparkMap.value[it.code]) return
    try {
      const points = await GetBKFundFlowListByDate(it.code, todayStr)
      sparkMap.value = {
        ...sparkMap.value,
        [it.code]: pickArray(points).map(p => Number(p.netInflow)),
      }
    } catch { /* 单条失败不影响整体 */ }
  }))
}

async function loadData(silent = false) {
  if (!silent) loading.value = true
  try {
    const res = await GetBKFundFlowTopListByDate(todayStr, 500)
    all.value = pickArray(res)
    // 先派生 TOP 再拉折线
    await loadSparklines([...inflowList.value, ...outflowList.value])
  } catch (e) {
    console.error('加载板块资金失败:', e)
    all.value = []
  } finally {
    if (!silent) loading.value = false
  }
}

async function handleRefresh() {
  // 刷新时清缓存重拉折线
  sparkMap.value = {}
  await loadData(true)
}

onBeforeMount(() => {
  loadData()
  registerFeed('mobile-sector-flow', {
    fetch: () => loadData(true),
    intervalMs: 60 * 1000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-sector-flow')
})
</script>

<template>
  <div class="page">
    <MPullRefresh class="flow-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !all.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else-if="all.length">
          <!-- 流入 TOP -->
          <section v-if="inflowList.length" class="block">
            <div class="block-title block-title--rise">🔴 资金流入 TOP{{ inflowList.length }}</div>
            <div class="rank-list">
              <div v-for="(it, i) in inflowList" :key="it.code || i" class="rank-item">
                <span class="rank-idx">{{ i + 1 }}</span>
                <span class="rank-name">{{ it.name }}</span>
                <MiniSparkline
                  v-if="sparkMap[it.code]?.length"
                  class="rank-spark"
                  :data="sparkMap[it.code]"
                  :width="56"
                  :height="22"
                />
                <span class="rank-amt m-rise">{{ formatMoney(it.netInflow) }}</span>
              </div>
            </div>
          </section>

          <!-- 流出 TOP -->
          <section v-if="outflowList.length" class="block">
            <div class="block-title block-title--fall">🟢 资金流出 TOP{{ outflowList.length }}</div>
            <div class="rank-list">
              <div v-for="(it, i) in outflowList" :key="it.code || i" class="rank-item">
                <span class="rank-idx">{{ i + 1 }}</span>
                <span class="rank-name">{{ it.name }}</span>
                <MiniSparkline
                  v-if="sparkMap[it.code]?.length"
                  class="rank-spark"
                  :data="sparkMap[it.code]"
                  :width="56"
                  :height="22"
                />
                <span class="rank-amt m-fall">{{ formatMoney(it.netInflow) }}</span>
              </div>
            </div>
          </section>

          <MEmpty v-if="!inflowList.length && !outflowList.length" description="暂无数据" />
        </template>
        <MEmpty v-else description="暂无数据" />
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

.flow-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

.block {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.block-title {
  padding: var(--m-space-sm) var(--m-space-md);
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  background: var(--m-bg-primary);
}

.block-title--rise {
  color: var(--m-color-rise);
}

.block-title--fall {
  color: var(--m-color-fall);
}

.rank-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.rank-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
}

.rank-idx {
  width: 20px;
  flex-shrink: 0;
  text-align: center;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.rank-name {
  flex: 1;
  min-width: 0;
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-spark {
  flex-shrink: 0;
}

.rank-amt {
  flex-shrink: 0;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
}

.rank-amt.m-rise {
  color: var(--m-color-rise);
}

.rank-amt.m-fall {
  color: var(--m-color-fall);
}
</style>
