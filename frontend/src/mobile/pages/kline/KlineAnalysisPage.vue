<script setup>
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MultiPeriodKlineChart from '../../components/charts/MultiPeriodKlineChart.vue'
import FenshiChart from '../../components/charts/FenshiChart.vue'
import KlineSignalSummary from '../../components/charts/KlineSignalSummary.vue'
import ChipDistribution from '../../components/charts/ChipDistribution.vue'
import LongPositionSheet from '../../components/sheets/LongPositionSheet.vue'
import IndicatorPanelSheet from '../../components/sheets/IndicatorPanelSheet.vue'
import { GetStockList, GetStockKLineWithFallback, GetStockMinutePriceLineData } from '../../../api/app'
import { evaluateIndicatorSignals, summarizeSignals } from '../../composables/indicatorSignals'

// 当前选中的股票 { code, name }
const stockInfo = ref(null)

// 搜索
const searchKeyword = ref('')
const searchResults = ref([])
const stockList = ref([])

// 最近浏览股票（对齐桌面端 kline-analysis.vue，localStorage 持久化）
const recentStocks = ref([])

// 当前周期
const currentPeriod = ref('day')

// 图表数据
const fenshiData = ref([])
const fenshiDate = ref('')
const fenshiPreClose = ref(0)
const klineData = ref([])
const loading = ref(false)

// 弹层显隐
const longPosVisible = ref(false)
const indicatorPanelVisible = ref(false)

// 已叠加的指标集合（传给图表控制副图/叠加）
const activeIndicators = ref(['MA', 'VOL'])

// 最新收盘价（供多单计算器一键填充）
const latestClose = computed(() => {
  const arr = klineData.value
  if (!arr.length) return ''
  return Number(arr[arr.length - 1].close) || ''
})

// 是否分时图
const isFenshi = computed(() => currentPeriod.value === 'fenshi')

// 周期标签映射
const periodLabels = {
  'fenshi': '分时',
  '1min': '1分',
  '5min': '5分',
  '15min': '15分',
  '30min': '30分',
  '60min': '60分',
  'day': '日K',
  'week': '周K',
  'month': '月K',
}

// 周期 → klt 参数映射（对齐桌面端 kLineWithFallback，东方财富格式）
const periodToKlt = {
  '1min': '1',
  '5min': '5',
  '15min': '15',
  '30min': '30',
  '60min': '60',
  'day': '101',
  'week': '102',
  'month': '103',
}

// 指标信号汇总（基于当前 K 线数据；分时图不评估）
const signalSummary = computed(() => {
  if (!klineData.value || !klineData.value.length || isFenshi.value) return null
  const signals = evaluateIndicatorSignals(klineData.value)
  return summarizeSignals(signals)
})

// 初始化：加载股票列表供搜索过滤
async function initStockList() {
  try {
    const list = await GetStockList('')
    stockList.value = Array.isArray(list) ? list : []
  } catch (error) {
    console.error('加载股票列表失败:', error)
  }
}

// 搜索股票（本地过滤）
function handleSearch() {
  const keyword = searchKeyword.value.trim().toLowerCase()
  if (!keyword) {
    searchResults.value = []
    return
  }
  searchResults.value = stockList.value.filter(s =>
    s.name?.toLowerCase().includes(keyword) ||
    s.ts_code?.toLowerCase().includes(keyword) ||
    s.symbol?.toLowerCase().includes(keyword) ||
    s.cnspell?.toLowerCase().includes(keyword)
  ).slice(0, 30)
}

// 选择股票
function handleSelectStock(stock) {
  stockInfo.value = {
    code: stock.ts_code,
    name: stock.name,
  }
  searchKeyword.value = ''
  searchResults.value = []
  addToRecent(stock.ts_code, stock.name)
  loadData()
}

// 最近浏览：写入 localStorage（去重、置顶、最多 10 条）
function addToRecent(code, name) {
  if (!code) return
  const list = recentStocks.value.filter(s => s.code !== code)
  list.unshift({ code, name: name || '' })
  if (list.length > 10) list.length = 10
  recentStocks.value = list
  try {
    localStorage.setItem('kline-recent-stocks', JSON.stringify(list))
  } catch { /* 忽略隐私模式/配额 */ }
}

function loadRecentStocks() {
  try {
    const raw = localStorage.getItem('kline-recent-stocks')
    if (raw) recentStocks.value = JSON.parse(raw) || []
  } catch { /* 忽略解析失败 */ }
}

// 点击最近浏览快速切换
function selectRecent(stock) {
  stockInfo.value = {
    code: stock.code,
    name: stock.name,
  }
  addToRecent(stock.code, stock.name)
  loadData()
}

// 加载数据（根据周期）
async function loadData() {
  if (!stockInfo.value) return
  refreshToken++            // 切换股票/周期：使在途的实时刷新响应失效
  loading.value = true
  try {
    if (isFenshi.value) {
      await loadFenshi()
    } else {
      await loadKline()
    }
  } finally {
    loading.value = false
    setupRefreshTimer()     // 按当前模式启停实时轮询
  }
}

// 加载分时数据
async function loadFenshi() {
  try {
    const result = await GetStockMinutePriceLineData(stockInfo.value.code, stockInfo.value.name)
    fenshiData.value = result?.priceData || []
    fenshiDate.value = result?.date || ''
    // 昨收价：分时接口不含昨收，取最近 2 根日 K 的前一根 close 作为基准
    fenshiPreClose.value = await fetchPreClose()
  } catch (error) {
    console.error('加载分时数据失败:', error)
    fenshiData.value = []
    fenshiPreClose.value = 0
  }
}

// 取昨收价（最近 2 根日 K 的前一根收盘）
async function fetchPreClose() {
  try {
    const result = await GetStockKLineWithFallback(stockInfo.value.code, stockInfo.value.name, '101', 2)
    const list = result?.data || []
    // 2 根时前一根是昨日；仅 1 根（新股）时无昨收，回退 0
    if (list.length >= 2) return Number(list[list.length - 2].close) || 0
    return 0
  } catch {
    return 0
  }
}

// 加载 K 线数据
async function loadKline() {
  const klt = periodToKlt[currentPeriod.value] || '101'
  try {
    // GetStockKLineWithFallback(code, name, klt, limit)
    const result = await GetStockKLineWithFallback(stockInfo.value.code, stockInfo.value.name, klt, 500)
    klineData.value = (result?.data) || []
  } catch (error) {
    console.error('加载K线数据失败:', error)
    klineData.value = []
  }
}

// 实时刷新：对齐桌面端 StockLightweightKlineChart 的 realtimeIntervalMs 轮询。
// 仅 K 线模式下定时拉取最新 K 线，静默更新（失败不打断看盘）。
const REALTIME_INTERVAL_MS = 60 * 1000
let refreshTimer = null
// 请求令牌：切换股票/周期时自增，丢弃过期响应，避免覆盖到新股票
let refreshToken = 0

function clearRefreshTimer() {
  if (refreshTimer) { clearInterval(refreshTimer); refreshTimer = null }
}

// 按当前周期静默刷新最新 K 线（不切换 loading 态）
async function refreshLatestKline() {
  if (!stockInfo.value || isFenshi.value) return
  const token = refreshToken
  const klt = periodToKlt[currentPeriod.value] || '101'
  try {
    const result = await GetStockKLineWithFallback(stockInfo.value.code, stockInfo.value.name, klt, 500)
    if (token !== refreshToken) return          // 已切换股票/周期，丢弃
    const list = (result?.data) || []
    if (list.length) klineData.value = list
  } catch {
    /* 静默 */
  }
}

// 进入 K 线模式且已选股时启动轮询；否则清理
function setupRefreshTimer() {
  clearRefreshTimer()
  if (stockInfo.value && !isFenshi.value) {
    refreshTimer = setInterval(refreshLatestKline, REALTIME_INTERVAL_MS)
  }
}

// 周期切换
watch(currentPeriod, () => {
  if (stockInfo.value) loadData()
})

// 初始化
initStockList()
loadRecentStocks()

// 卸载时清理实时轮询
onBeforeUnmount(clearRefreshTimer)
</script>

<template>
  <div class="kline-analysis-page">
    <PageHeader title="K线分析" />

    <!-- 股票搜索区 -->
    <div class="search-area">
      <input
        v-model="searchKeyword"
        class="search-input"
        type="text"
        placeholder="输入股票代码/名称"
        @input="handleSearch"
      >
      <!-- 搜索结果下拉 -->
      <div v-if="searchKeyword && searchResults.length" class="search-results">
        <div
          v-for="stock in searchResults"
          :key="stock.ts_code"
          class="search-item"
          @click="handleSelectStock(stock)"
        >
          <span class="search-name">{{ stock.name }}</span>
          <span class="search-code">{{ stock.ts_code }}</span>
        </div>
      </div>

      <!-- 最近浏览（无搜索词且未选中股票时展示） -->
      <div v-if="!searchKeyword && !stockInfo && recentStocks.length" class="recent-area">
        <span class="recent-label">最近</span>
        <button
          v-for="s in recentStocks.slice(0, 8)"
          :key="s.code"
          class="recent-chip"
          @click="selectRecent(s)"
        >{{ s.name || s.code }}</button>
      </div>
    </div>

    <!-- 选中股票信息 -->
    <div v-if="stockInfo" class="stock-info-bar">
      <div class="stock-name">{{ stockInfo.name }}</div>
      <div class="stock-code">{{ stockInfo.code }}</div>
    </div>

    <!-- 周期切换栏 -->
    <div v-if="stockInfo" class="period-bar">
      <button
        v-for="(label, period) in periodLabels"
        :key="period"
        class="period-btn"
        :class="{ 'period-btn--active': currentPeriod === period }"
        @click="currentPeriod = period"
      >
        {{ label }}
      </button>
      <button class="period-btn tool-btn" @click="longPosVisible = true">多单</button>
      <button
        class="period-btn tool-btn"
        :class="{ 'period-btn--active': indicatorPanelVisible }"
        @click="indicatorPanelVisible = true"
      >指标</button>
    </div>

    <!-- 图表区 -->
    <div class="chart-area">
      <MEmpty v-if="!stockInfo" description="请先搜索选择一只股票" />
      <div v-else-if="loading" class="loading-hint">加载中...</div>
      <template v-else>
        <!-- 分时图 -->
        <FenshiChart
          v-if="isFenshi"
          :data="fenshiData"
          :pre-close="fenshiPreClose"
          :date="fenshiDate"
          :height="320"
        />
        <!-- K线图 + OHLC 信息条 -->
        <MultiPeriodKlineChart
          v-else
          :data="klineData"
          :height="320"
          :period="currentPeriod"
          :indicators="activeIndicators"
        />

        <!-- 指标信号汇总（仅非分时且有数据时显示） -->
        <KlineSignalSummary
          v-if="!isFenshi && signalSummary"
          :summary="signalSummary"
          class="signal-section"
        />

        <!-- 筹码分布（仅非分时且有数据时显示） -->
        <ChipDistribution
          v-if="!isFenshi && klineData.length"
          :data="klineData"
          :height="160"
          class="signal-section"
        />
      </template>
    </div>

    <!-- 多单仓位计算器 -->
    <LongPositionSheet
      v-model:show="longPosVisible"
      :latest-close="latestClose"
    />

    <!-- 指标开关面板 -->
    <IndicatorPanelSheet
      v-model:show="indicatorPanelVisible"
      v-model:indicators="activeIndicators"
    />
  </div>
</template>

<style scoped>
.kline-analysis-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.search-area {
  position: relative;
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.search-input {
  width: 100%;
  height: 36px;
  padding: 0 var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-md);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-sm);
  outline: none;
}

.search-input:focus {
  border-color: var(--m-color-rise);
}

.search-results {
  position: absolute;
  top: calc(100% - 1px);
  left: var(--m-space-md);
  right: var(--m-space-md);
  background: var(--m-bg-card);
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-sm);
  box-shadow: var(--m-shadow-md);
  z-index: var(--m-z-dropdown);
  max-height: 240px;
  overflow-y: auto;
}

.search-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--m-space-sm) var(--m-space-md);
  border-bottom: 1px solid var(--m-divider-color);
}

.search-item:last-child {
  border-bottom: none;
}

.search-item:active {
  background: var(--m-bg-primary);
}

.search-name {
  color: var(--m-text-primary);
  font-size: var(--m-font-sm);
}

.search-code {
  color: var(--m-text-tertiary);
  font-size: var(--m-font-xs);
}

/* 最近浏览 */
.recent-area {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
  margin-top: var(--m-space-sm);
}

.recent-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  white-space: nowrap;
}

.recent-chip {
  height: 26px;
  padding: 0 var(--m-space-sm);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-xs);
  white-space: nowrap;
}

.recent-chip:active {
  transform: scale(0.95);
}

.stock-info-bar {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.stock-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.stock-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

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
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  white-space: nowrap;
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.period-btn--active {
  background: var(--m-color-rise);
  color: white;
  border-color: var(--m-color-rise);
}

.period-btn:active {
  transform: scale(0.95);
}

.chart-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.loading-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

.signal-section {
  margin: var(--m-space-md);
}
</style>
