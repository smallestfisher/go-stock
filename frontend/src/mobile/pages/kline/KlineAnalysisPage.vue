<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import FenshiChart from '../../components/charts/FenshiChart.vue'
import FullKlineChart from '../../components/charts/FullKlineChart.vue'
import ChartControlSheet from '../../components/sheets/ChartControlSheet.vue'
import PriceTag from '../../components/widgets/PriceTag.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import { useSwipe } from '../../composables/useSwipe'

const router = useRouter()

// 模拟股票数据
const stockInfo = ref({
  code: '600519',
  name: '贵州茅台',
  price: 1820.50,
  preClose: 1778.90,
  changePercent: 2.34,
  changeAmount: 41.60,
})

// 当前周期
const currentPeriod = ref('day')

// 当前指标
const currentIndicators = ref(['MA', 'VOL'])

// 控制抽屉
const controlVisible = ref(false)

// 容器引用
const chartContainer = ref(null)

// 是否分时图
const isFenshi = computed(() => currentPeriod.value === 'fenshi')

// 模拟分时数据
const fenshiData = ref([
  { time: '09:30', price: 1780, avgPrice: 1780, volume: 1000 },
  { time: '10:00', price: 1790, avgPrice: 1785, volume: 2000 },
  { time: '10:30', price: 1800, avgPrice: 1790, volume: 1500 },
  { time: '11:00', price: 1810, avgPrice: 1795, volume: 1800 },
  { time: '11:30', price: 1805, avgPrice: 1797, volume: 1200 },
  { time: '13:00', price: 1815, avgPrice: 1800, volume: 2200 },
  { time: '14:00', price: 1820, avgPrice: 1803, volume: 1900 },
  { time: '15:00', price: 1820.50, avgPrice: 1805, volume: 2500 },
])

// 模拟K线数据
const klineData = ref([
  { time: '2024-06-10', open: 1750, close: 1760, high: 1770, low: 1740, volume: 100000 },
  { time: '2024-06-11', open: 1760, close: 1755, high: 1765, low: 1750, volume: 95000 },
  { time: '2024-06-12', open: 1755, close: 1770, high: 1780, low: 1755, volume: 110000 },
  { time: '2024-06-13', open: 1770, close: 1785, high: 1790, low: 1765, volume: 120000 },
  { time: '2024-06-14', open: 1785, close: 1780, high: 1795, low: 1775, volume: 105000 },
  { time: '2024-06-17', open: 1780, close: 1795, high: 1800, low: 1778, volume: 115000 },
  { time: '2024-06-18', open: 1795, close: 1810, high: 1815, low: 1790, volume: 125000 },
  { time: '2024-06-19', open: 1810, close: 1805, high: 1820, low: 1800, volume: 108000 },
  { time: '2024-06-20', open: 1805, close: 1820.50, high: 1825, low: 1800, volume: 130000 },
])

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
    // 下一个周期
    if (currentPeriodIndex.value < periods.length - 1) {
      currentPeriod.value = periods[currentPeriodIndex.value + 1]
    }
  },
  onSwipeRight: () => {
    // 上一个周期
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

// 返回
function handleBack() {
  router.back()
}
</script>

<template>
  <div class="kline-analysis-page">
    <!-- 顶部信息栏 -->
    <div class="page-header">
      <button class="back-btn" @click="handleBack">←</button>
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
      <button class="control-btn" @click="handleOpenControl">⚙️</button>
    </div>

    <!-- 图表区 -->
    <div ref="chartContainer" class="chart-area">
      <FenshiChart
        v-if="isFenshi"
        :data="fenshiData"
        :pre-close="stockInfo.preClose"
        :width="375"
        :height="300"
      />
      <FullKlineChart
        v-else
        :data="klineData"
        :width="375"
        :height="400"
        :ma-lines="[5, 10, 20, 30]"
      />

      <div class="swipe-hint">← 左右滑动切换周期 →</div>
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

.page-header {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.back-btn {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 24px;
  color: var(--m-text-primary);
  cursor: pointer;
}

.back-btn:active {
  opacity: 0.6;
}

.stock-info {
  flex: 1;
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

.control-btn {
  padding: var(--m-space-xs) var(--m-space-md);
  background: transparent;
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  font-size: 16px;
  cursor: pointer;
}

.control-btn:active {
  opacity: 0.6;
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
