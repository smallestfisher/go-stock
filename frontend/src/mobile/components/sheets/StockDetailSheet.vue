<script setup>
import { ref, computed } from 'vue'
import MSheet from '../base/MSheet.vue'
import MTabs from '../base/MTabs.vue'
import MButton from '../base/MButton.vue'
import PriceTag from '../widgets/PriceTag.vue'
import PercentTag from '../widgets/PercentTag.vue'

const props = defineProps({
  // 是否显示
  show: {
    type: Boolean,
    default: false
  },
  // 股票数据
  stock: {
    type: Object,
    default: () => null
  }
})

const emit = defineEmits(['update:show', 'trade'])

// 详情标签页
const detailTabs = [
  { label: '分时', value: 'fenshi' },
  { label: '日K', value: 'kline' },
  { label: '资料', value: 'info' },
  { label: '资金', value: 'fund' },
]

const activeTab = ref('fenshi')

// 模拟盘口数据
const bidAskData = computed(() => {
  if (!props.stock) return { bid: [], ask: [] }

  const basePrice = props.stock.price
  return {
    ask: [
      { price: (basePrice + 0.05).toFixed(2), volume: 1250 },
      { price: (basePrice + 0.04).toFixed(2), volume: 2340 },
      { price: (basePrice + 0.03).toFixed(2), volume: 1890 },
      { price: (basePrice + 0.02).toFixed(2), volume: 3210 },
      { price: (basePrice + 0.01).toFixed(2), volume: 2560 },
    ],
    bid: [
      { price: basePrice.toFixed(2), volume: 2890 },
      { price: (basePrice - 0.01).toFixed(2), volume: 3450 },
      { price: (basePrice - 0.02).toFixed(2), volume: 1670 },
      { price: (basePrice - 0.03).toFixed(2), volume: 2340 },
      { price: (basePrice - 0.04).toFixed(2), volume: 1890 },
    ]
  }
})

function handleClose() {
  emit('update:show', false)
}

function handleTrade(type) {
  emit('trade', { stock: props.stock, type })
  console.log(`交易操作: ${type}`, props.stock)
}
</script>

<template>
  <MSheet
    :show="show"
    :title="stock?.name"
    height="85vh"
    @update:show="handleClose"
  >
    <div v-if="stock" class="stock-detail">
      <!-- 股票基本信息 -->
      <div class="detail-header">
        <div class="header-left">
          <div class="stock-code">{{ stock.code }}</div>
        </div>
        <div class="header-right">
          <PriceTag
            :price="stock.price"
            :change="stock.changePercent"
            size="large"
            bold
          />
          <div class="price-change">
            <PercentTag :value="stock.changePercent" />
            <span
              class="change-amount"
              :class="{
                'm-rise': stock.changeAmount > 0,
                'm-fall': stock.changeAmount < 0
              }"
            >
              {{ stock.changeAmount > 0 ? '+' : '' }}{{ stock.changeAmount?.toFixed(2) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 标签页 -->
      <MTabs v-model="activeTab" :tabs="detailTabs" />

      <!-- 内容区 -->
      <div class="detail-content">
        <!-- 分时/K线图占位 -->
        <div v-if="activeTab === 'fenshi' || activeTab === 'kline'" class="chart-placeholder">
          <p>📈 {{ activeTab === 'fenshi' ? '分时图' : 'K线图' }}</p>
          <p class="placeholder-tip">图表组件开发中...</p>
        </div>

        <!-- 盘口数据 -->
        <div v-if="activeTab === 'fenshi'" class="bid-ask-panel">
          <div class="bid-ask-header">
            <span>价格</span>
            <span>挂单量</span>
          </div>
          <!-- 卖盘 -->
          <div class="ask-list">
            <div
              v-for="(item, index) in bidAskData.ask.slice().reverse()"
              :key="`ask-${index}`"
              class="bid-ask-item ask-item"
            >
              <span class="m-fall">{{ item.price }}</span>
              <span>{{ item.volume }}</span>
            </div>
          </div>
          <!-- 买盘 -->
          <div class="bid-list">
            <div
              v-for="(item, index) in bidAskData.bid"
              :key="`bid-${index}`"
              class="bid-ask-item bid-item"
            >
              <span class="m-rise">{{ item.price }}</span>
              <span>{{ item.volume }}</span>
            </div>
          </div>
        </div>

        <!-- 资料 -->
        <div v-if="activeTab === 'info'" class="info-panel">
          <div class="info-item">
            <span class="info-label">成交量:</span>
            <span class="info-value">{{ stock.volume }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">成交额:</span>
            <span class="info-value">{{ stock.turnover }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">总市值:</span>
            <span class="info-value">-</span>
          </div>
          <div class="info-item">
            <span class="info-label">流通市值:</span>
            <span class="info-value">-</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作按钮 -->
    <template #footer>
      <div class="detail-actions">
        <MButton type="danger" @click="handleTrade('sell')">
          卖出
        </MButton>
        <MButton type="success" @click="handleTrade('buy')">
          买入
        </MButton>
      </div>
    </template>
  </MSheet>
</template>

<style scoped>
.stock-detail {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding-bottom: var(--m-space-lg);
  border-bottom: 1px solid var(--m-divider-color);
}

.stock-code {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

.header-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--m-space-xs);
}

.price-change {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  font-size: var(--m-font-sm);
}

.change-amount {
  font-variant-numeric: tabular-nums;
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
  min-height: 300px;
}

.chart-placeholder {
  height: 250px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
  color: var(--m-text-secondary);
}

.chart-placeholder p {
  margin: var(--m-space-xs) 0;
}

.placeholder-tip {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

.bid-ask-panel {
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
}

.bid-ask-header {
  display: flex;
  justify-content: space-between;
  padding: var(--m-space-sm) 0;
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
  font-weight: var(--m-font-weight-medium);
  border-bottom: 1px solid var(--m-divider-color);
}

.ask-list,
.bid-list {
  display: flex;
  flex-direction: column;
}

.bid-ask-item {
  display: flex;
  justify-content: space-between;
  padding: var(--m-space-xs) 0;
  font-size: var(--m-font-sm);
  font-variant-numeric: tabular-nums;
}

.info-panel {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
}

.info-label {
  color: var(--m-text-secondary);
}

.info-value {
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
}

.detail-actions {
  display: flex;
  gap: var(--m-space-md);
}

.detail-actions .m-button {
  flex: 1;
}
</style>
