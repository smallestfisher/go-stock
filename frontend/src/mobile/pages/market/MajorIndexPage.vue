<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount, onMounted, nextTick, watch } from 'vue'
import { GetStockKLineWithFallback } from '../../../api/app'
import MTabs from '../../components/base/MTabs.vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MultiPeriodKlineChart from '../../components/charts/MultiPeriodKlineChart.vue'
import IndicatorPanelSheet from '../../components/sheets/IndicatorPanelSheet.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'

// 对齐桌面端 market.vue「重大指数」子页：桌面端铺了 20 个指数 + 全功能 K 线（StockLightweightKlineChart）。
// 移动端重构为「指数看盘页」：精选 8 个核心指数，指数 Tab + 实时报价头 + 多周期 + 指标切换 + K 线填满，
// 去掉桌面端冗余指数与个股特有功能。代码与桌面端 StockLightweightKlineChart 完全一致，复用 GetStockKLineWithFallback。
const INDEXES = [
  { label: '上证指数', code: '000001.SH', name: '上证指数' },
  { label: '深证成指', code: '399001.SZ', name: '深证成指' },
  { label: '创业板指', code: '399006.SZ', name: '创业板指' },
  { label: '沪深300', code: '000300.SH', name: '沪深300' },
  { label: '恒生指数', code: '100.HSI', name: '恒生指数' },
  { label: '道琼斯', code: '100.DJIA', name: '道琼斯' },
  { label: '标普500', code: '100.SPX', name: '标普500' },
  { label: '纳斯达克', code: '100.NDX', name: '纳斯达克' },
]

const tabs = INDEXES.map(i => ({ label: i.label, value: i.code }))
const activeCode = ref(INDEXES[0].code)
const activeIndex = computed(() => INDEXES.find(i => i.code === activeCode.value) || INDEXES[0])

const klineData = ref([])
const loading = ref(false)

// 周期（指数不提供分时，仅日/周/月）；映射对齐桌面端东方财富 klt
const periodLabels = { day: '日K', week: '周K', month: '月K' }
const periodToKlt = { day: '101', week: '102', month: '103' }
const currentPeriod = ref('day')

// 指标（对齐桌面端指标切换能力；图表据此叠加均线/副图）
const activeIndicators = ref(['MA', 'VOL'])
const indicatorPanelVisible = ref(false)

// ===== 实时报价（取 K 线最后一根：close=最新点位，涨跌额=close-前收，涨跌幅=changePercent） =====
const lastBar = computed(() => (klineData.value.length ? klineData.value[klineData.value.length - 1] : null))
const prevBar = computed(() => (klineData.value.length > 1 ? klineData.value[klineData.value.length - 2] : null))
const latestPrice = computed(() => Number(lastBar.value?.close))

const changeAbs = computed(() => {
  const c = Number(lastBar.value?.close)
  const p = Number(prevBar.value?.close)
  if (!Number.isFinite(c) || !Number.isFinite(p)) return NaN
  return c - p
})

const changePct = computed(() => {
  if (!lastBar.value) return NaN
  const n = Number(lastBar.value.changePercent)
  if (Number.isFinite(n)) return n
  // 兜底：相邻两根收盘推算
  if (Number.isFinite(changeAbs.value)) {
    const p = Number(prevBar.value?.close)
    if (Number.isFinite(p) && p !== 0) return (changeAbs.value / p) * 100
  }
  return NaN
})

// 涨跌方向：red 涨 / green 跌 / gray 平
const dirClass = computed(() => {
  const n = changePct.value
  if (!Number.isFinite(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
})

// 千分位 2 位小数；非数回退 '--'
function fmtPrice(v) {
  const n = Number(v)
  if (!Number.isFinite(n)) return '--'
  return n.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 涨跌额：带正负号
function fmtChangeAbs() {
  if (!Number.isFinite(changeAbs.value)) return '--'
  const v = changeAbs.value
  return (v > 0 ? '+' : '') + v.toFixed(2)
}

// ===== 图表高度：实测容器，避免魔法数字导致底部被裁（移动端嵌套在 MarketPage 导航之下） =====
const chartHostRef = ref(null)
const panelHeight = ref(360) // 占位初值，onMounted 后由 ResizeObserver 纠正
let resizeObserver = null

function applyHeight() {
  const el = chartHostRef.value
  if (!el) return
  const h = el.clientHeight
  if (h > 0 && h !== panelHeight.value) {
    panelHeight.value = h
    // canvas 高度随 props.height 更新后，通知 echarts 按新尺寸重布局
    // （MultiPeriodKlineChart 内部监听 window resize 调用 chart.resize()）
    nextTick(() => window.dispatchEvent(new Event('resize')))
  }
}

// ===== 取数 =====
// 请求令牌：切换指数/周期时自增，丢弃在途响应，避免覆盖到新指数（对齐 KlineAnalysisPage）
let refreshToken = 0

async function loadKline(silent = false) {
  const idx = activeIndex.value
  const token = refreshToken
  const klt = periodToKlt[currentPeriod.value] || '101'
  if (!silent) loading.value = true
  try {
    const result = await GetStockKLineWithFallback(idx.code, idx.name, klt, 500)
    if (token !== refreshToken) return
    klineData.value = (result?.data) || []
  } catch (e) {
    console.error('加载K线失败:', e)
    if (token === refreshToken) klineData.value = []
  } finally {
    if (!silent && token === refreshToken) loading.value = false
  }
}

// 实时刷新：对齐桌面 StockLightweightKlineChart 的 realtimeIntervalMs 轮询，
// 60s 静默拉取当前周期最新 K 线（失败不打断看盘）。
const REALTIME_INTERVAL_MS = 60 * 1000
let refreshTimer = null
function clearTimer() { if (refreshTimer) { clearInterval(refreshTimer); refreshTimer = null } }
function setupTimer() {
  clearTimer()
  refreshTimer = setInterval(() => loadKline(true), REALTIME_INTERVAL_MS)
}

// 切指数 / 切周期：作废在途请求并重新加载
watch([activeCode, currentPeriod], () => { refreshToken++; loadKline() })

// 下拉刷新：强制重新拉取当前指数/周期的 K 线（带 loading 态）
async function handleRefresh() {
  refreshToken++
  await loadKline(true)
}

onBeforeMount(() => {
  loadKline()
  setupTimer()
})

onMounted(() => {
  resizeObserver = new ResizeObserver(applyHeight)
  if (chartHostRef.value) resizeObserver.observe(chartHostRef.value)
  applyHeight()
})

onBeforeUnmount(() => {
  clearTimer()
  resizeObserver?.disconnect()
})
</script>

<template>
  <div class="major-index-page">
    <!-- 指数切换 Tab -->
    <div class="index-tabs">
      <MTabs v-model="activeCode" :tabs="tabs" />
    </div>

    <!-- 实时报价头：指数名/代码 + 最新点位/涨跌额/涨跌幅 -->
    <div class="quote-bar">
      <div class="quote-left">
        <div class="quote-name">{{ activeIndex.name }}</div>
        <div class="quote-code">{{ activeIndex.code }}</div>
      </div>
      <div class="quote-right" :class="dirClass">
        <div class="quote-price">{{ fmtPrice(latestPrice) }}</div>
        <div class="quote-sub">
          <span class="quote-change">{{ fmtChangeAbs() }}</span>
          <PercentTag :value="changePct" size="small" />
        </div>
      </div>
    </div>

    <!-- 周期 + 指标 -->
    <div class="period-bar">
      <button
        v-for="(label, key) in periodLabels"
        :key="key"
        type="button"
        class="period-btn"
        :class="{ 'period-btn--active': currentPeriod === key }"
        @click="currentPeriod = key"
      >
        {{ label }}
      </button>
      <button
        type="button"
        class="period-btn tool-btn"
        :class="{ 'period-btn--active': indicatorPanelVisible }"
        @click="indicatorPanelVisible = true"
      >
        指标
      </button>
    </div>

    <!-- 图表区：flex:1 占满剩余高度；MPullRefresh 提供下拉刷新 + 内容滚动
         （多副图时 canvas 可能高于一屏，靠 track 滚动，避免被裁切）。 -->
    <div ref="chartHostRef" class="chart-area">
      <MPullRefresh :on-refresh="handleRefresh">
        <div class="chart-inner">
          <div v-if="loading && !klineData.length" class="loading">
            <MLoading text="加载中..." vertical />
          </div>
          <MultiPeriodKlineChart
            v-else-if="klineData.length"
            :data="klineData"
            :height="panelHeight"
            :period="currentPeriod"
            :indicators="activeIndicators"
            :fields="['open', 'close', 'high', 'low', 'changePercent', 'changeValue', 'amplitude', 'turnoverRate']"
          />
          <MEmpty v-else description="暂无数据" />
        </div>
      </MPullRefresh>
    </div>

    <!-- 指标切换面板 -->
    <IndicatorPanelSheet v-model:show="indicatorPanelVisible" v-model:indicators="activeIndicators" />
  </div>
</template>

<style scoped>
.major-index-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.index-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

/* 实时报价头 */
.quote-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.quote-left {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.quote-name {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quote-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.quote-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  flex-shrink: 0;
}

.quote-price {
  font-size: var(--m-font-2xl);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
  color: var(--m-text-primary);
}

.quote-sub {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  font-variant-numeric: tabular-nums;
}

.quote-change {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

/* 涨跌色（作用于价格；涨跌幅由 PercentTag 自带颜色） */
.quote-right.m-rise .quote-price,
.quote-right.m-rise .quote-change {
  color: var(--m-color-rise);
}

.quote-right.m-fall .quote-price,
.quote-right.m-fall .quote-change {
  color: var(--m-color-fall);
}

/* 周期 + 指标条（对齐 KlineAnalysisPage 的 period-bar） */
.period-bar {
  display: flex;
  gap: var(--m-space-xs);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  overflow-x: auto;
  scrollbar-width: none;
}

.period-bar::-webkit-scrollbar {
  display: none;
}

.period-btn {
  padding: var(--m-space-xs) var(--m-space-md);
  background: transparent;
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  white-space: nowrap;
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.period-btn--active {
  background: var(--m-color-rise);
  color: #fff;
  border-color: var(--m-color-rise);
}

.period-btn:active {
  transform: scale(0.95);
}

.tool-btn {
  margin-left: auto;
}

/* 图表区：占满剩余高度，作为图表一屏高度的测量基准（chartHostRef 绑在此）。 */
.chart-area {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

/* MPullRefresh 高度填满，内部 track 提供滚动 + 下拉刷新 */
.chart-area :deep(.m-pull-refresh) {
  height: 100%;
}

.chart-inner {
  /* 自然高度：canvas 可能更高（多副图），由 MPullRefresh 的 track 负责滚动 */
  min-height: 100%;
  padding: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}
</style>
