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

// 是否持仓（有成本+持仓才显示盈亏，对齐桌面端 costPrice>0）
function hasPosition(stock) {
  return (Number(stock.costPrice) || 0) > 0 && (Number(stock.costVolume) || 0) > 0
}

// 盈亏方向：盈 rise(红) / 亏 fall(绿)
function pnlDir(stock) {
  const a = Number(stock.profitAmount) || 0
  if (a > 0) return 'rise'
  if (a < 0) return 'fall'
  return 'flat'
}

// 总盈亏率（带符号）
function pnlRate(stock) {
  const r = Number(stock.profit) || 0
  return `${r > 0 ? '+' : ''}${r.toFixed(2)}%`
}

// 总盈亏额（带符号，取整 ¥）
function pnlAmount(stock) {
  const a = Number(stock.profitAmount) || 0
  return `${a > 0 ? '+' : ''}${a.toFixed(0)}¥`
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

function fmtDatePart(date) {
  if (!date) return ''
  const s = String(date).trim()
  if (/^\d{4}-\d{2}-\d{2}/.test(s)) return s.slice(5, 10)
  if (/^\d{8}$/.test(s)) return `${s.slice(4, 6)}-${s.slice(6, 8)}`
  const d = new Date(s)
  if (isNaN(d.getTime())) return ''
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${mm}-${dd}`
}

function fmtTimePart(t) {
  if (!t) return ''
  const s = String(t).trim()
  if (/^\d{1,2}:\d{2}(:\d{2})?$/.test(s)) {
    return s
  }
  const d = new Date(s)
  if (!isNaN(d.getTime())) {
    const hh = String(d.getHours()).padStart(2, '0')
    const mi = String(d.getMinutes()).padStart(2, '0')
    const ss = String(d.getSeconds()).padStart(2, '0')
    return `${hh}:${mi}:${ss}`
  }
  return s
}

// 行情时间：后端行情源返回的日期/时间，不是页面刷新时间。
function fmtQuoteTime(stock) {
  const time = fmtTimePart(stock?.time)
  if (!time) return ''

  const date = fmtDatePart(stock?.quoteDate || stock?.date)
  if (date) return `行情 ${date} ${time}`

  const dateFromTime = fmtDatePart(stock?.time)
  if (dateFromTime) return `行情 ${dateFromTime} ${time}`

  return `行情 ${time}`
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
            <span v-if="stock.time" class="stock-time">{{ fmtQuoteTime(stock) }}</span>
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
          <!-- 持仓盈亏（有成本+持仓才显示，对齐桌面端） -->
          <div v-if="hasPosition(stock)" class="stock-pnl">
            <span class="pnl-tag" :class="`m-${pnlDir(stock)}`">{{ pnlRate(stock) }}</span>
            <span class="pnl-tag" :class="`m-${pnlDir(stock)}`">{{ pnlAmount(stock) }}</span>
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

/* 持仓盈亏（精简：盈亏率% + 盈亏额¥） */
.stock-pnl {
  display: flex;
  gap: var(--m-space-xs);
}

.pnl-tag {
  font-size: var(--m-font-xs);
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
  padding: 1px var(--m-space-xs);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-card);
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
