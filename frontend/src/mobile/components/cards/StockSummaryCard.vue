<script setup>
import { computed } from 'vue'
import MCard from '../base/MCard.vue'

const props = defineProps({
  // 股票数据
  stocks: {
    type: Array,
    default: () => []
  },
  // 最多显示数量
  maxShow: {
    type: Number,
    default: 3
  }
})

const emit = defineEmits(['viewAll', 'stockClick'])

// 显示的股票列表
const displayStocks = computed(() => {
  return props.stocks.slice(0, props.maxShow)
})

// 统计数据
const stats = computed(() => {
  if (!props.stocks.length) {
    return { total: 0, rise: 0, fall: 0, flat: 0 }
  }

  const rise = props.stocks.filter(s => s.changePercent > 0).length
  const fall = props.stocks.filter(s => s.changePercent < 0).length
  const flat = props.stocks.length - rise - fall

  return {
    total: props.stocks.length,
    rise,
    fall,
    flat
  }
})

// 涨跌方向
function dir(stock) {
  const p = Number(stock.changePercent) || 0
  if (p > 0) return 'rise'
  if (p < 0) return 'fall'
  return 'flat'
}

// 价格格式化
function fmt(v) {
  const n = Number(v) || 0
  return n ? n.toFixed(2) : '--'
}

// 量额格式化（万/亿）
function big(n) {
  n = Number(n) || 0
  if (!n) return '--'
  if (n >= 100000000) return `${(n / 100000000).toFixed(2)}亿`
  if (n >= 10000) return `${(n / 10000).toFixed(2)}万`
  return String(n)
}

// 时间格式化：兼容「16:14:27」纯时间和「2026-06-18T22:56:02+08:00」ISO 串
function fmtTime(t) {
  if (!t) return ''
  const s = String(t).trim()
  // 纯时间 HH:mm:ss（含 HH:mm）—— 去掉秒更紧凑
  if (/^\d{1,2}:\d{2}(:\d{2})?$/.test(s)) {
    return s.replace(/:(\d{2})$/, m => '')  // 16:14:27 → 16:14
  }
  // ISO / Date 串 —— 取 月-日 时:分
  const d = new Date(s)
  if (!isNaN(d.getTime())) {
    const mm = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    const hh = String(d.getHours()).padStart(2, '0')
    const mi = String(d.getMinutes()).padStart(2, '0')
    return `${mm}-${dd} ${hh}:${mi}`
  }
  return s
}

function handleViewAll() {
  emit('viewAll')
}

function handleStockClick(stock) {
  emit('stockClick', stock)
}
</script>

<template>
  <MCard>
    <!-- 头部 -->
    <div class="stock-summary-header">
      <div class="header-left">
        <h3 class="header-title">📊 我的自选</h3>
        <div class="header-stats">
          <span class="stat-item m-rise">{{ stats.rise }}涨</span>
          <span class="stat-item m-fall">{{ stats.fall }}跌</span>
          <span class="stat-item m-flat">{{ stats.flat }}平</span>
        </div>
      </div>
      <button class="header-action" @click="handleViewAll">
        查看全部 {{ stats.total }} 只 →
      </button>
    </div>

    <!-- 股票列表 -->
    <div v-if="displayStocks.length" class="stock-list">
      <div
        v-for="stock in displayStocks"
        :key="stock.code"
        class="stock-item"
        :class="`stock-item--${dir(stock)}`"
        @click="handleStockClick(stock)"
      >
        <div class="stock-info">
          <div class="stock-name-row">
            <span class="stock-name">{{ stock.name }}</span>
            <span class="stock-code">{{ stock.code }}</span>
            <span v-if="stock.time" class="stock-time">{{ fmtTime(stock.time) }}</span>
          </div>
          <div class="stock-sub">
            <span>开 <b>{{ fmt(stock.open) }}</b></span>
            <span>高 <b class="m-rise">{{ fmt(stock.high) }}</b></span>
            <span>低 <b class="m-fall">{{ fmt(stock.low) }}</b></span>
            <span>昨收 <b>{{ fmt(stock.preClose) }}</b></span>
          </div>
          <div class="stock-sub">
            <span>量 {{ big(stock.volume) }}</span>
            <span>额 {{ big(stock.turnover) }}</span>
          </div>
        </div>
        <div class="stock-price">
          <div class="price-now" :class="`m-${dir(stock)}`">{{ fmt(stock.price) }}</div>
          <div class="price-change" :class="`bg-${dir(stock)}`">
            {{ (Number(stock.changePercent) || 0) > 0 ? '+' : '' }}{{ (Number(stock.changePercent) || 0).toFixed(2) }}%
          </div>
          <div class="price-amt" :class="`m-${dir(stock)}`">
            {{ (Number(stock.changeAmount) || 0) > 0 ? '+' : '' }}{{ fmt(stock.changeAmount) }}
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="stock-empty">
      <p>还没有添加自选股票</p>
      <button class="empty-action" @click="handleViewAll">
        去添加 →
      </button>
    </div>
  </MCard>
</template>

<style scoped>
.stock-summary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--m-space-lg);
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.header-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.header-stats {
  display: flex;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
}

.stat-item {
  font-weight: var(--m-font-weight-medium);
}

.header-action {
  flex-shrink: 0;
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

.stock-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.stock-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.stock-item:active {
  transform: scale(0.98);
  background: var(--m-divider-color);
}

.stock-info {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
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

.stock-name-row {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  flex-wrap: wrap;
}

.stock-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-quaternary, var(--m-text-tertiary));
  margin-left: auto;
  font-variant-numeric: tabular-nums;
}

.stock-sub {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs) var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.stock-sub b {
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-secondary);
  font-variant-numeric: tabular-nums;
}

.price-now {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
}

.price-amt {
  font-size: var(--m-font-xs);
  font-variant-numeric: tabular-nums;
}

/* 涨跌色 */
.m-rise { color: var(--m-color-rise); }
.m-fall { color: var(--m-color-fall); }
.m-flat { color: var(--m-text-secondary); }
.bg-rise { background: var(--m-color-rise); }
.bg-fall { background: var(--m-color-fall); }
.bg-flat { background: var(--m-text-tertiary); }

.price-change {
  color: #fff;
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-bold);
  padding: 2px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  min-width: 60px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.stock-price {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--m-space-xs);
}

.stock-empty {
  text-align: center;
  padding: var(--m-space-xl) 0;
  color: var(--m-text-secondary);
}

.stock-empty p {
  margin-bottom: var(--m-space-md);
}

.empty-action {
  padding: var(--m-space-sm) var(--m-space-lg);
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
  border: none;
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  cursor: pointer;
}

.empty-action:active {
  opacity: 0.7;
}
</style>
