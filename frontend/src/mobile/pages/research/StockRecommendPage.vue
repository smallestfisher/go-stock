<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount } from 'vue'
import {
  GetAiRecommendStocksList,
  DeleteAiRecommendStocks,
  UpdateAiRecommendStocksAlert,
} from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'

// 对齐桌面端 aiRecommendStocksList.vue：GetAiRecommendStocksList 分页列表。
// 移动端精简：一次取 30 条滚动浏览，预警开关 + 删除保留，K 线详情不做。
const list = ref([])
const loading = ref(false)

function pickList(res) {
  return (res && Array.isArray(res.list) ? res.list : []).filter(it => it && typeof it === 'object')
}

async function loadList() {
  loading.value = true
  try {
    const res = await GetAiRecommendStocksList({
      page: 1,
      pageSize: 30,
      startDate: '',
      endDate: '',
      enableAlert: null,
    })
    list.value = pickList(res)
  } catch (e) {
    console.error('加载推荐记录失败:', e)
    list.value = []
  } finally {
    loading.value = false
  }
}

// 涨跌幅 = (current - pre)/pre*100（无 changePercent 字段）
function changePct(it) {
  const c = Number(it.stockCurrentPrice)
  const p = Number(it.stockPrePrice)
  if (!Number.isFinite(c) || !Number.isFinite(p) || p === 0) return NaN
  return (c - p) / p * 100
}

// 盈亏：current vs stockPrice（推荐时价）
function profit(it) {
  const c = Number(it.stockCurrentPrice)
  const p = Number(it.stockPrice)
  if (!Number.isFinite(c) || !Number.isFinite(p) || p === 0) return null
  const diff = (c - p) / p * 100
  if (diff > 0) return { text: `暂赢 ${diff.toFixed(2)}%`, rise: true }
  if (diff < 0) return { text: `暂亏 ${diff.toFixed(2)}%`, rise: false }
  return { text: '暂平', rise: null }
}

function fmtTime(s) {
  return s ? String(s).substring(0, 19).replace('T', ' ') : ''
}

// 预警开关
async function toggleAlert(it) {
  try {
    await UpdateAiRecommendStocksAlert(it.ID, it.enableAlert)
  } catch (e) {
    console.error('更新预警失败:', e)
    it.enableAlert = !it.enableAlert // 回滚
  }
}

// 删除
async function removeItem(it) {
  try {
    await DeleteAiRecommendStocks(it.ID)
    list.value = list.value.filter(x => x.ID !== it.ID)
  } catch (e) {
    console.error('删除失败:', e)
  }
}

async function handleRefresh() {
  await loadList()
}

onBeforeMount(() => {
  loadList()
  // 开市 60s 静默刷最新价（不打断 loading）
  registerFeed('mobile-ai-recommend', {
    fetch: loadList,
    intervalMs: 60 * 1000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-ai-recommend')
})

const hasData = computed(() => list.value.length > 0)
</script>

<template>
  <div class="page">
    <MPullRefresh class="rec-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !hasData" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else>
          <div v-if="hasData" class="rec-list">
            <div v-for="it in list" :key="it.ID" class="rec-item">
              <!-- 头行：股票名+代码 / 盈亏 / 删除 -->
              <div class="row-top">
                <div class="stock-left">
                  <span class="stock-name">{{ it.stockName }}</span>
                  <span class="stock-code">{{ it.stockCode }}</span>
                </div>
                <span
                  v-if="profit(it)"
                  class="profit"
                  :class="{
                    'm-rise': profit(it).rise === true,
                    'm-fall': profit(it).rise === false,
                  }"
                >
                  {{ profit(it).text }}
                </span>
                <button class="del-btn" type="button" @click="removeItem(it)">✕</button>
              </div>

              <!-- 模型 / 板块 / 评级 / 时间 -->
              <div class="row-tags">
                <span v-if="it.modelName" class="tag">{{ it.modelName }}</span>
                <span v-if="it.bkName" class="tag">{{ it.bkName }}</span>
                <span v-if="it.rating" class="tag tag--rate">{{ it.rating }}</span>
                <span class="time">{{ fmtTime(it.CreatedAt) }}</span>
              </div>

              <!-- 价格行：最新价 + 涨跌幅 -->
              <div class="row-price">
                <span class="price-label">最新</span>
                <span class="price-val">{{ it.stockCurrentPrice || '--' }}</span>
                <PercentTag :value="changePct(it)" size="small" />
                <span class="rec-price">推荐时 {{ it.stockPrice || '--' }}</span>
              </div>

              <!-- 建仓价格区间 -->
              <div v-if="it.recommendBuyPrice || it.recommendStopProfitPrice || it.recommendStopLossPrice" class="row-levels">
                <span v-if="it.recommendBuyPrice" class="level"><i>开仓</i><b>{{ it.recommendBuyPrice }}</b></span>
                <span v-if="it.recommendStopProfitPrice" class="level level--profit"><i>止盈</i><b>{{ it.recommendStopProfitPrice }}</b></span>
                <span v-if="it.recommendStopLossPrice" class="level level--loss"><i>止损</i><b>{{ it.recommendStopLossPrice }}</b></span>
              </div>

              <!-- 推荐理由 -->
              <p v-if="it.recommendReason" class="reason">{{ it.recommendReason }}</p>

              <!-- 风险提示 -->
              <p v-if="it.riskRemarks" class="risk">⚠ {{ it.riskRemarks }}</p>

              <!-- 预警开关 -->
              <div class="row-alert">
                <span class="alert-label">监控预警</span>
                <button
                  type="button"
                  class="alert-switch"
                  :class="{ 'alert-switch--on': it.enableAlert }"
                  @click="it.enableAlert = !it.enableAlert; toggleAlert(it)"
                >
                  <span class="alert-knob" />
                </button>
              </div>
            </div>
          </div>
          <MEmpty v-else description="暂无 AI 推荐记录" />
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

.rec-refresh {
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

.rec-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.rec-item {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
}

.row-top {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.stock-left {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-xs);
  flex: 1;
  min-width: 0;
}

.stock-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stock-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  flex-shrink: 0;
}

.profit {
  flex-shrink: 0;
  padding: 1px var(--m-space-sm);
  font-size: var(--m-font-xs);
  border-radius: var(--m-radius-sm);
  font-variant-numeric: tabular-nums;
}

.profit.m-rise {
  color: var(--m-color-rise);
  background: var(--m-color-rise-light);
}

.profit.m-fall {
  color: var(--m-color-fall);
  background: var(--m-color-fall-light);
}

.del-btn {
  width: 24px;
  height: 24px;
  border: none;
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  color: var(--m-text-tertiary);
  font-size: var(--m-font-xs);
  line-height: 1;
  flex-shrink: 0;
}

/* 标签行 */
.row-tags {
  display: flex;
  align-items: center;
  gap: var(--m-space-xs);
  flex-wrap: wrap;
}

.tag {
  padding: 1px var(--m-space-sm);
  font-size: var(--m-font-xs);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-secondary);
}

.tag--rate {
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.time {
  margin-left: auto;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

/* 价格行 */
.row-price {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
}

.price-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.price-val {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
  color: var(--m-text-primary);
}

.rec-price {
  margin-left: auto;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

/* 建仓价区间 */
.row-levels {
  display: flex;
  gap: var(--m-space-md);
  flex-wrap: wrap;
}

.level {
  display: inline-flex;
  align-items: baseline;
  gap: var(--m-space-xs);
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.level i {
  font-style: normal;
  color: var(--m-text-tertiary);
}

.level b {
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
}

.level--profit b {
  color: var(--m-color-rise);
}

.level--loss b {
  color: var(--m-color-fall);
}

.reason {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
}

.risk {
  font-size: var(--m-font-xs);
  color: var(--m-color-fall);
  line-height: var(--m-line-height-normal);
}

/* 预警开关 */
.row-alert {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: var(--m-space-xs);
  border-top: 1px solid var(--m-divider-color);
}

.alert-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.alert-switch {
  width: 36px;
  height: 20px;
  border: none;
  border-radius: var(--m-radius-full);
  background: var(--m-text-disabled);
  position: relative;
  transition: background var(--m-duration-fast);
}

.alert-switch--on {
  background: var(--m-color-rise);
}

.alert-knob {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #fff;
  transition: transform var(--m-duration-fast);
}

.alert-switch--on .alert-knob {
  transform: translateX(16px);
}
</style>
