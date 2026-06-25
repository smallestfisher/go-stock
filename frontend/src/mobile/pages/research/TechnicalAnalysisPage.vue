<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import StockCard from '../../components/cards/StockCard.vue'

const technicalSignals = ref([
  { type: '金叉', desc: 'MA5上穿MA10', stocks: 3, color: 'rise' },
  { type: '突破', desc: '突破前期高点', stocks: 5, color: 'rise' },
  { type: '死叉', desc: 'MA5下穿MA10', stocks: 2, color: 'fall' },
])

const signalStocks = ref([
  { code: '600519', name: '贵州茅台', price: 1820.50, changePercent: 2.34, changeAmount: 41.60, volume: 125000, turnover: 228250000, sparklineData: [1780, 1790, 1785, 1800, 1810, 1805, 1815, 1820], signal: 'MA5金叉MA10' },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 技术信号统计 -->
        <MCard>
          <h3 class="section-title">今日技术信号</h3>
          <div class="signal-list">
            <div v-for="(signal, i) in technicalSignals" :key="i" class="signal-item">
              <div class="signal-info">
                <div class="signal-type" :class="`m-${signal.color}`">{{ signal.type }}</div>
                <div class="signal-desc">{{ signal.desc }}</div>
              </div>
              <div class="signal-count">{{ signal.stocks }} 只</div>
            </div>
          </div>
        </MCard>

        <!-- 信号个股 -->
        <MCard>
          <h3 class="section-title">金叉信号个股</h3>
          <div class="stock-list">
            <div v-for="stock in signalStocks" :key="stock.code" class="stock-item">
              <StockCard :stock="stock" />
              <div class="stock-signal">{{ stock.signal }}</div>
            </div>
          </div>
        </MCard>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.section-title { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-lg); }
.signal-list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.signal-item { display: flex; justify-content: space-between; align-items: center; padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.signal-info { flex: 1; }
.signal-type { font-weight: var(--m-font-weight-bold); margin-bottom: var(--m-space-xs); }
.signal-desc { font-size: var(--m-font-sm); color: var(--m-text-secondary); }
.signal-count { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-bold); color: var(--m-text-primary); }
.stock-list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.stock-item { display: flex; flex-direction: column; gap: var(--m-space-xs); }
.stock-signal { padding: 0 var(--m-space-md); font-size: var(--m-font-sm); color: var(--m-color-rise); font-weight: var(--m-font-weight-medium); }
</style>
