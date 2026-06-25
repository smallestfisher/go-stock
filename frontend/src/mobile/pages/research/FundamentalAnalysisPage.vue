<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'

const valuationData = ref([
  { name: '贵州茅台', code: '600519', pe: 35.6, pb: 12.3, roe: 28.5, rating: '低估' },
  { name: '五粮液', code: '000858', pe: 28.9, pb: 9.8, roe: 24.3, rating: '合理' },
])

const performanceData = ref([
  { quarter: '2024Q1', revenue: 125.6, profit: 45.3, yoy: 15.6 },
  { quarter: '2023Q4', revenue: 138.2, profit: 52.1, yoy: 12.3 },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 估值分析 -->
        <MCard>
          <h3 class="section-title">估值分析</h3>
          <div class="table">
            <div class="table-header">
              <span>股票</span>
              <span>PE</span>
              <span>PB</span>
              <span>ROE(%)</span>
              <span>评级</span>
            </div>
            <div v-for="item in valuationData" :key="item.code" class="table-row">
              <span class="stock-name">{{ item.name }}</span>
              <span>{{ item.pe }}</span>
              <span>{{ item.pb }}</span>
              <span>{{ item.roe }}</span>
              <span class="rating m-rise">{{ item.rating }}</span>
            </div>
          </div>
        </MCard>

        <!-- 业绩表现 -->
        <MCard>
          <h3 class="section-title">业绩表现</h3>
          <div class="performance-list">
            <div v-for="item in performanceData" :key="item.quarter" class="performance-item">
              <div class="performance-quarter">{{ item.quarter }}</div>
              <div class="performance-data">
                <span>营收: {{ item.revenue }}亿</span>
                <span>净利: {{ item.profit }}亿</span>
                <PercentTag :value="item.yoy" size="small" />
              </div>
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
.table { display: flex; flex-direction: column; }
.table-header, .table-row { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr 1fr; gap: var(--m-space-sm); padding: var(--m-space-md); font-size: var(--m-font-sm); }
.table-header { background: var(--m-bg-primary); font-weight: var(--m-font-weight-medium); border-radius: var(--m-radius-sm); }
.table-row { border-bottom: 1px solid var(--m-divider-color); }
.table-row:last-child { border-bottom: none; }
.stock-name { font-weight: var(--m-font-weight-medium); }
.rating { font-weight: var(--m-font-weight-medium); }
.performance-list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.performance-item { padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.performance-quarter { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-sm); }
.performance-data { display: flex; gap: var(--m-space-lg); font-size: var(--m-font-sm); color: var(--m-text-secondary); }
</style>
