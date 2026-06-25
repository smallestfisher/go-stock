<script setup>
import PriceTag from '../widgets/PriceTag.vue'
import PercentTag from '../widgets/PercentTag.vue'
import MiniSparkline from '../charts/MiniSparkline.vue'

defineProps({
  // 股票数据
  stock: {
    type: Object,
    required: true
    // { code, name, price, changePercent, changeAmount, volume, turnover, sparklineData }
  },
  // 是否显示迷你K线
  showSparkline: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['click'])

function handleClick() {
  emit('click')
}

// 格式化成交量
function formatVolume(volume) {
  if (!volume) return '-'
  if (volume >= 100000000) return `${(volume / 100000000).toFixed(2)}亿`
  if (volume >= 10000) return `${(volume / 10000).toFixed(2)}万`
  return volume
}

// 格式化成交额
function formatTurnover(turnover) {
  if (!turnover) return '-'
  if (turnover >= 100000000) return `${(turnover / 100000000).toFixed(2)}亿`
  if (turnover >= 10000) return `${(turnover / 10000).toFixed(2)}万`
  return turnover
}
</script>

<template>
  <div class="stock-card" @click="handleClick">
    <!-- 左侧：股票信息 -->
    <div class="stock-info">
      <div class="stock-name-row">
        <h4 class="stock-name">{{ stock.name }}</h4>
        <span class="stock-code">{{ stock.code }}</span>
      </div>
      <div class="stock-meta">
        <span class="meta-item">量: {{ formatVolume(stock.volume) }}</span>
        <span class="meta-item">额: {{ formatTurnover(stock.turnover) }}</span>
      </div>
    </div>

    <!-- 中间：迷你走势图 -->
    <div v-if="showSparkline && stock.sparklineData" class="stock-chart">
      <MiniSparkline
        :data="stock.sparklineData"
        :width="80"
        :height="40"
      />
    </div>

    <!-- 右侧：价格信息 -->
    <div class="stock-price">
      <PriceTag
        :price="stock.price"
        :change="stock.changePercent"
        size="large"
        bold
      />
      <div class="price-change">
        <PercentTag
          :value="stock.changePercent"
          size="small"
        />
        <span
          class="change-amount"
          :class="{
            'm-rise': stock.changeAmount > 0,
            'm-fall': stock.changeAmount < 0,
            'm-flat': stock.changeAmount === 0
          }"
        >
          {{ stock.changeAmount > 0 ? '+' : '' }}{{ stock.changeAmount?.toFixed(2) }}
        </span>
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
  cursor: pointer;
  transition: all var(--m-duration-fast);
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
  font-weight: var(--m-font-weight-medium);
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

.stock-meta {
  display: flex;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.meta-item {
  white-space: nowrap;
}

.stock-chart {
  flex-shrink: 0;
}

.stock-price {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--m-space-xs);
}

.price-change {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.change-amount {
  font-size: var(--m-font-xs);
  font-variant-numeric: tabular-nums;
}
</style>
