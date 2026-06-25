<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import StockCard from '../../components/cards/StockCard.vue'

// 形态选股：全部股票信息列表（对齐桌面端 allStockList.vue）
// 数据后续接 GetAllStockInfoList / GetAllStocks 填充
const keyword = ref('')
const stocks = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="page">
    <div class="shape-header">
      <input
        v-model="keyword"
        class="search-input"
        type="text"
        placeholder="输入股票代码/名称筛选..."
      >
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="stocks.length" padding="none">
          <div class="stock-list">
            <div v-for="stock in stocks" :key="stock.code" class="stock-item">
              <StockCard :stock="stock" :show-sparkline="false" />
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无股票数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; display: flex; flex-direction: column; overflow: hidden; }
.shape-header { padding: var(--m-space-md); background: var(--m-bg-card); border-bottom: 1px solid var(--m-divider-color); }
.search-input { width: 100%; padding: var(--m-space-sm) var(--m-space-md); border: 1px solid var(--m-divider-color); border-radius: var(--m-radius-md); font-size: var(--m-font-sm); background: var(--m-bg-primary); }
.search-input:focus { outline: none; border-color: var(--m-color-rise); }
.container { padding: var(--m-space-md); }
.stock-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-sm); }
.stock-item { padding: 0; }
</style>
