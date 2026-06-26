<script setup>
import { computed } from 'vue'

const props = defineProps({
  // 股票数据（字段由 StockListPage 的 fetchRealtime 填充，对齐桌面端 StockMobileList）
  stock: {
    type: Object,
    required: true
    // { code, name, price, changePercent, changeAmount, open, preClose, high, low, volume, turnover, time, profitToday }
  }
})

const emit = defineEmits(['click'])

function handleClick() {
  emit('click')
}

// 涨跌方向：rise 涨 / fall 跌 / flat 平（决定颜色，对齐桌面端 result.type 的红绿语义）
const direction = computed(() => {
  const p = Number(props.stock.changePercent) || 0
  if (p > 0) return 'rise'
  if (p < 0) return 'fall'
  return 'flat'
})

// 格式化价格（保留 2-3 位小数）
function formatPrice(v) {
  const n = Number(v) || 0
  return n ? n.toFixed(2) : '--'
}

// 格式化成交量/额（万、亿）
function formatBig(n) {
  n = Number(n) || 0
  if (!n) return '--'
  if (n >= 100000000) return `${(n / 100000000).toFixed(2)}亿`
  if (n >= 10000) return `${(n / 10000).toFixed(2)}万`
  return String(n)
}
</script>

<template>
  <div class="stock-card" :class="`stock-card--${direction}`" @click="handleClick">
    <!-- 左侧：股票信息 -->
    <div class="stock-info">
      <div class="stock-name-row">
        <h4 class="stock-name">{{ stock.name }}</h4>
        <span class="stock-code">{{ stock.code }}</span>
      </div>

      <!-- OHLC 行情次要信息 -->
      <div class="stock-ohlc">
        <span class="ohlc-item">开 <b>{{ formatPrice(stock.open) }}</b></span>
        <span class="ohlc-item">高 <b class="m-rise">{{ formatPrice(stock.high) }}</b></span>
        <span class="ohlc-item">低 <b class="m-fall">{{ formatPrice(stock.low) }}</b></span>
        <span class="ohlc-item">昨收 <b>{{ formatPrice(stock.preClose) }}</b></span>
      </div>

      <!-- 量额 -->
      <div class="stock-meta">
        <span class="meta-item">量 {{ formatBig(stock.volume) }}</span>
        <span class="meta-item">额 {{ formatBig(stock.turnover) }}</span>
        <span v-if="stock.time" class="meta-item meta-time">{{ stock.time }}</span>
      </div>
    </div>

    <!-- 右侧：价格 + 涨跌 -->
    <div class="stock-price">
      <div class="price-now" :class="`m-${direction}`">
        {{ formatPrice(stock.price) }}
      </div>
      <!-- 涨跌幅红绿背景块（对齐桌面端 bg-success/error） -->
      <div class="price-change" :class="`bg-${direction}`">
        <span v-if="direction === 'rise'">+</span>{{ (Number(stock.changePercent) || 0).toFixed(2) }}%
      </div>
      <div class="change-amount" :class="`m-${direction}`">
        {{ Number(stock.changeAmount) > 0 ? '+' : '' }}{{ formatPrice(stock.changeAmount) }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.stock-card {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  border-left: 4px solid var(--m-divider-color);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.stock-card--rise {
  border-left-color: var(--m-color-rise);
}

.stock-card--fall {
  border-left-color: var(--m-color-fall);
}

.stock-card--flat {
  border-left-color: var(--m-text-tertiary);
}

.stock-card:active {
  transform: scale(0.98);
  background: var(--m-bg-primary);
}

.stock-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.stock-name-row {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
}

.stock-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
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

.stock-ohlc {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-sm) var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.ohlc-item b {
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-secondary);
  font-variant-numeric: tabular-nums;
}

.stock-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-sm) var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.meta-item {
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

.meta-time {
  color: var(--m-text-quaternary, var(--m-text-tertiary));
}

.stock-price {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 3px;
  min-width: 84px;
}

.price-now {
  font-size: var(--m-font-xl);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
}

.price-change {
  color: #fff;
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-bold);
  padding: 2px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  min-width: 64px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.change-amount {
  font-size: var(--m-font-xs);
  font-variant-numeric: tabular-nums;
}

/* 涨跌色 */
.m-rise {
  color: var(--m-color-rise);
}

.m-fall {
  color: var(--m-color-fall);
}

.m-flat {
  color: var(--m-text-secondary);
}

.bg-rise {
  background: var(--m-color-rise);
}

.bg-fall {
  background: var(--m-color-fall);
}

.bg-flat {
  background: var(--m-text-tertiary);
}
</style>
