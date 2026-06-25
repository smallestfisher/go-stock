<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import ChartControlSheet from '../../components/sheets/ChartControlSheet.vue'
import PriceTag from '../../components/widgets/PriceTag.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import { useSwipe } from '../../composables/useSwipe'
import { SearchStock } from '../../../api/app'

const router = useRouter()

// 当前选中的股票（后续接 SearchStock 选择后填充，先留空）
// 字段：{ code, name, price, preClose, changePercent, changeAmount }
const stockInfo = ref(null)

// 搜索
const searchKeyword = ref('')
const searchResults = ref([])
const searching = ref(false)

// 当前周期
const currentPeriod = ref('day')

// 当前指标
const currentIndicators = ref(['MA', 'VOL'])

// 控制抽屉
const controlVisible = ref(false)

// 容器引用
const chartContainer = ref(null)

// 图表数据（后续接 GetStockKLine / GetStockMinutePriceLineData 填充，先留空）
const fenshiData = ref([])
const klineData = ref([])

// 是否分时图
const isFenshi = computed(() => currentPeriod.value === 'fenshi')

// 搜索股票
async function handleSearch() {
  const keyword = searchKeyword.value.trim()
  if (!keyword) {
    searchResults.value = []
    return
  }
  searching.value = true
  try {
    // 后续接真实 API
    searchResults.value = []
  } catch (error) {
    console.error('搜索股票失败:', error)
  } finally {
    searching.value = false
  }
}

// 选择股票
function handleSelectStock(stock) {
  stockInfo.value = stock
  searchKeyword.value = ''
  searchResults.value = []
  // TODO: 加载该股票的 K线/分时数据
}

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

// 左右滑动切换周期
const periods = Object.keys(periodLabels)
const currentPeriodIndex = computed(() => periods.indexOf(currentPeriod.value))

useSwipe({
  containerRef: chartContainer,
  onSwipeLeft: () => {
    if (currentPeriodIndex.value < periods.length - 1) {
      currentPeriod.value = periods[currentPeriodIndex.value + 1]
    }
  },
  onSwipeRight: () => {
    if (currentPeriodIndex.value > 0) {
      currentPeriod.value = periods[currentPeriodIndex.value - 1]
    }
  }
})

// 打开设置
function handleOpenControl() {
  controlVisible.value = true
}

// 切换周期
function handlePeriodChange(period) {
  currentPeriod.value = period
}

// 切换指标
function handleIndicatorToggle(indicator) {
  const index = currentIndicators.value.indexOf(indicator)
  if (index > -1) {
    currentIndicators.value.splice(index, 1)
  } else {
    currentIndicators.value.push(indicator)
  }
}
</script>

<template>
  <div class="kline-analysis-page">
    <!-- 顶部：只保留标题 -->
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
          :key="stock.code"
          class="search-item"
          @click="handleSelectStock(stock)"
        >
          <span class="search-name">{{ stock.name }}</span>
          <span class="search-code">{{ stock.code }}</span>
        </div>
      </div>
    </div>

    <!-- 选中股票信息 -->
    <div v-if="stockInfo" class="stock-info-bar">
      <div class="stock-info">
        <div class="stock-name">{{ stockInfo.name }}</div>
        <div class="stock-code">{{ stockInfo.code }}</div>
      </div>
      <div class="stock-price">
        <PriceTag :price="stockInfo.price" :change="stockInfo.changePercent" size="large" bold />
        <PercentTag :value="stockInfo.changePercent" size="small" />
      </div>
    </div>

    <!-- 周期切换栏 -->
    <div class="period-bar">
      <button
        v-for="(label, period) in periodLabels"
        :key="period"
        class="period-btn"
        :class="{ 'period-btn--active': currentPeriod === period }"
        @click="currentPeriod = period"
      >
        {{ label }}
      </button>
      <!-- 图表控制入口（指标设置） -->
      <button class="period-btn control-entry" @click="handleOpenControl">⚙️</button>
    </div>

    <!-- 图表区 -->
    <div ref="chartContainer" class="chart-area">
      <MEmpty v-if="!stockInfo" description="请先搜索选择一只股票" />
      <MEmpty v-else-if="isFenshi ? !fenshiData.length : !klineData.length" description="暂无K线数据" />
      <div v-else class="swipe-hint">← 左右滑动切换周期 →</div>
    </div>

    <!-- 图表控制抽屉 -->
    <ChartControlSheet
      v-model:show="controlVisible"
      :period="currentPeriod"
      :indicators="currentIndicators"
      @change-period="handlePeriodChange"
      @toggle-indicator="handleIndicatorToggle"
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

.action-btn {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 18px;
  color: var(--m-text-primary);
  cursor: pointer;
}

.action-btn:active {
  opacity: 0.6;
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
  justify-content: space-between;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.stock-info {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.stock-name {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.stock-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.stock-price {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--m-space-xs);
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

.chart-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--m-space-md);
  overflow: hidden;
  position: relative;
}

.swipe-hint {
  position: absolute;
  bottom: var(--m-space-md);
  left: 50%;
  transform: translateX(-50%);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  padding: var(--m-space-xs) var(--m-space-md);
  background: rgba(0, 0, 0, 0.5);
  border-radius: var(--m-radius-full);
  white-space: nowrap;
}
</style>
