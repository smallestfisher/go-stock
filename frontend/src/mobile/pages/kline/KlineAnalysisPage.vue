<script setup>
import { ref, computed, watch } from 'vue'
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

// 当前周期
const currentPeriod = ref('day')

// 图表数据
const fenshiData = ref([])
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
  loadData()
}

// 加载数据（根据周期）
async function loadData() {
  if (!stockInfo.value) return
  loading.value = true
  try {
    if (isFenshi.value) {
      await loadFenshi()
    } else {
      await loadKline()
    }
  } finally {
    loading.value = false
  }
}

// 加载分时数据
async function loadFenshi() {
  try {
    const result = await GetStockMinutePriceLineData(stockInfo.value.code, stockInfo.value.name)
    fenshiData.value = result?.priceData || []
  } catch (error) {
    console.error('加载分时数据失败:', error)
    fenshiData.value = []
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

// 周期切换
watch(currentPeriod, () => {
  if (stockInfo.value) loadData()
})

// 初始化
initStockList()
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
