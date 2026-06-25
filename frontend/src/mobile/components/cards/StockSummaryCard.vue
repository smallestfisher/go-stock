<script setup>
import { ref, computed } from 'vue'
import MCard from '../base/MCard.vue'
import PriceTag from '../widgets/PriceTag.vue'
import PercentTag from '../widgets/PercentTag.vue'

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
        @click="handleStockClick(stock)"
      >
        <div class="stock-info">
          <div class="stock-name">{{ stock.name }}</div>
          <div class="stock-code">{{ stock.code }}</div>
        </div>
        <div class="stock-price">
          <PriceTag
            :price="stock.price"
            :change="stock.changePercent"
            size="large"
            bold
            prefix="¥"
          />
          <PercentTag
            :value="stock.changePercent"
            size="small"
          />
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
