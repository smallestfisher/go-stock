<script setup>
import MCard from '../base/MCard.vue'
import MIcon from '../base/MIcon.vue'
import PercentTag from '../widgets/PercentTag.vue'

defineProps({
  // 异动数据
  alerts: {
    type: Array,
    default: () => []
    // { stockName, stockCode, type, changePercent, time }
  }
})

const emit = defineEmits(['alertClick', 'viewAll'])

// 异动类型映射
const alertTypeMap = {
  rapid_rise: { label: '急速拉升', icon: 'trend-up', color: 'var(--m-color-rise)' },
  rapid_fall: { label: '急速下跌', icon: 'trend-down', color: 'var(--m-color-fall)' },
  volume_surge: { label: '放量异动', icon: 'burst', color: 'var(--m-color-rise)' },
  limit_up: { label: '涨停', icon: 'rocket', color: 'var(--m-color-rise)' },
  limit_down: { label: '跌停', icon: 'warning', color: 'var(--m-color-fall)' },
  break_high: { label: '突破新高', icon: 'arrow-up', color: 'var(--m-color-rise)' },
  break_low: { label: '跌破新低', icon: 'arrow-down', color: 'var(--m-color-fall)' },
}

function getAlertType(type) {
  return alertTypeMap[type] || { label: '异动', icon: 'bolt', color: 'var(--m-text-secondary)' }
}

function handleAlertClick(alert) {
  emit('alertClick', alert)
}

function handleViewAll() {
  emit('viewAll')
}

// 格式化时间
function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <MCard>
    <!-- 头部 -->
    <div class="alert-header">
      <h3 class="header-title">异动监控</h3>
      <button class="header-action" @click="handleViewAll">
        查看全部 <MIcon name="arrow-right" :size="14" />
      </button>
    </div>

    <!-- 异动列表 -->
    <div v-if="alerts.length" class="alert-list">
      <div
        v-for="(alert, index) in alerts"
        :key="index"
        class="alert-item"
        @click="handleAlertClick(alert)"
      >
        <div class="alert-icon" :style="{ color: getAlertType(alert.type).color }">
          <MIcon :name="getAlertType(alert.type).icon" :size="22" />
        </div>
        <div class="alert-info">
          <div class="alert-stock">
            <span class="stock-name">{{ alert.stockName }}</span>
            <span class="stock-code">{{ alert.stockCode }}</span>
          </div>
          <div class="alert-meta">
            <span class="alert-type" :style="{ color: getAlertType(alert.type).color }">
              {{ getAlertType(alert.type).label }}
            </span>
            <span class="alert-time">{{ formatTime(alert.time) }}</span>
          </div>
        </div>
        <div class="alert-change">
          <PercentTag :value="alert.changePercent" size="large" bold />
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="alert-empty">
      <p>暂无异动提醒</p>
      <p class="empty-tip">开盘后会自动监控自选股票</p>
    </div>
  </MCard>
</template>

<style scoped>
.alert-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--m-space-lg);
}

.header-title {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  letter-spacing: 0.2px;
}

.header-title::before {
  content: '';
  width: 3px;
  height: 15px;
  border-radius: var(--m-radius-full);
  background: var(--m-color-rise);
}

.header-action {
  padding: var(--m-space-xs) var(--m-space-sm);
  background: transparent;
  border: none;
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
  cursor: pointer;
}

.header-action:active {
  opacity: 0.6;
}

.alert-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.alert-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.alert-item:active {
  transform: scale(0.98);
  background: var(--m-divider-color);
}

.alert-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.alert-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.alert-stock {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
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

.alert-meta {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
}

.alert-type {
  font-weight: var(--m-font-weight-medium);
}

.alert-time {
  color: var(--m-text-tertiary);
}

.alert-change {
  flex-shrink: 0;
}

.alert-empty {
  text-align: center;
  padding: var(--m-space-xl) 0;
  color: var(--m-text-secondary);
}

.alert-empty p {
  margin-bottom: var(--m-space-xs);
}

.empty-tip {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}
</style>
