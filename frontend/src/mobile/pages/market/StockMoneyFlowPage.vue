<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import StockCard from '../../components/cards/StockCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

const stocks = ref([
  { code: '600519', name: '贵州茅台', price: 1820.50, changePercent: 2.34, changeAmount: 41.60, volume: 125000, turnover: 228250000, sparklineData: [1780, 1790, 1785, 1800, 1810, 1805, 1815, 1820], netInflow: 12500000 },
  { code: '000858', name: '五粮液', price: 156.80, changePercent: -1.23, changeAmount: -1.95, volume: 3450000, turnover: 540860000, sparklineData: [159, 158, 157.5, 157, 156.5, 157, 156.8, 156.8], netInflow: -8900000 },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}

function formatMoney(value) {
  if (value >= 100000000) return `${(value / 100000000).toFixed(2)}亿`
  if (value >= 10000) return `${(value / 10000).toFixed(2)}万`
  return value
}
</script>

<template>
  <div class="stock-money-flow-page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="stocks.length">
          <div class="stock-list">
            <div v-for="stock in stocks" :key="stock.code" class="flow-item">
              <StockCard :stock="stock" :show-sparkline="false" />
              <div class="flow-info">
                <span :class="stock.netInflow > 0 ? 'm-rise' : 'm-fall'">
                  净流入: {{ stock.netInflow > 0 ? '+' : '' }}{{ formatMoney(stock.netInflow) }}
                </span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.stock-money-flow-page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.stock-list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.flow-item { display: flex; flex-direction: column; gap: var(--m-space-xs); }
.flow-info { padding: 0 var(--m-space-md); font-size: var(--m-font-sm); font-weight: var(--m-font-weight-medium); }
</style>
