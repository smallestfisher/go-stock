<script setup>
import { ref } from 'vue'
import MSheet from '../base/MSheet.vue'
import MTabs from '../base/MTabs.vue'
import MButton from '../base/MButton.vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  // 当前周期
  period: {
    type: String,
    default: 'day'
  },
  // 当前指标
  indicators: {
    type: Array,
    default: () => ['MA']
  }
})

const emit = defineEmits(['update:show', 'changePeriod', 'toggleIndicator'])

// 周期选项
const periodOptions = [
  { label: '分时', value: 'fenshi' },
  { label: '1分', value: '1min' },
  { label: '5分', value: '5min' },
  { label: '15分', value: '15min' },
  { label: '30分', value: '30min' },
  { label: '60分', value: '60min' },
  { label: '日K', value: 'day' },
  { label: '周K', value: 'week' },
  { label: '月K', value: 'month' },
]

// 指标选项
const indicatorOptions = [
  { label: 'MA', value: 'MA', desc: '移动平均线' },
  { label: 'BOLL', value: 'BOLL', desc: '布林带' },
  { label: 'MACD', value: 'MACD', desc: '指数平滑异同移动平均线' },
  { label: 'KDJ', value: 'KDJ', desc: '随机指标' },
  { label: 'RSI', value: 'RSI', desc: '相对强弱指标' },
  { label: 'VOL', value: 'VOL', desc: '成交量' },
]

function handlePeriodChange(period) {
  emit('changePeriod', period)
  emit('update:show', false)
}

function handleIndicatorToggle(indicator) {
  emit('toggleIndicator', indicator)
}

function handleClose() {
  emit('update:show', false)
}
</script>

<template>
  <MSheet
    :show="show"
    title="图表设置"
    height="60vh"
    @update:show="handleClose"
  >
    <div class="chart-control">
      <!-- 周期选择 -->
      <div class="control-section">
        <h4 class="section-title">周期</h4>
        <div class="period-grid">
          <button
            v-for="option in periodOptions"
            :key="option.value"
            class="period-btn"
            :class="{ 'period-btn--active': period === option.value }"
            @click="handlePeriodChange(option.value)"
          >
            {{ option.label }}
          </button>
        </div>
      </div>

      <!-- 指标选择 -->
      <div class="control-section">
        <h4 class="section-title">指标</h4>
        <div class="indicator-list">
          <div
            v-for="option in indicatorOptions"
            :key="option.value"
            class="indicator-item"
            @click="handleIndicatorToggle(option.value)"
          >
            <div class="indicator-info">
              <span class="indicator-name">{{ option.label }}</span>
              <span class="indicator-desc">{{ option.desc }}</span>
            </div>
            <div
              class="indicator-checkbox"
              :class="{ 'indicator-checkbox--active': indicators.includes(option.value) }"
            >
              <span v-if="indicators.includes(option.value)">✓</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </MSheet>
</template>

<style scoped>
.chart-control {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xl);
}

.control-section {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.section-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.period-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--m-space-md);
}

.period-btn {
  min-height: var(--m-touch-min);
  padding: var(--m-space-sm);
  background: var(--m-bg-primary);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
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

.indicator-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.indicator-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  cursor: pointer;
  transition: background var(--m-duration-fast);
}

.indicator-item:active {
  background: var(--m-divider-color);
}

.indicator-info {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.indicator-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.indicator-desc {
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.indicator-checkbox {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  color: white;
  transition: all var(--m-duration-fast);
}

.indicator-checkbox--active {
  background: var(--m-color-rise);
  border-color: var(--m-color-rise);
}
</style>
