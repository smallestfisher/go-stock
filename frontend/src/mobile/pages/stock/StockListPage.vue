<script setup>
import { ref, computed } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'
import StockCard from '../../components/cards/StockCard.vue'
import VirtualList from '../../components/widgets/VirtualList.vue'
import StockDetailSheet from '../../components/sheets/StockDetailSheet.vue'

// 分组数据（模拟）
const groups = ref([
  { label: '全部', value: 'all' },
  { label: '自选', value: 'favorite' },
  { label: '持仓', value: 'holding' },
  { label: '关注', value: 'watch' },
])

const activeGroup = ref('all')

// 模拟股票数据（大量数据测试虚拟滚动）
const stocksData = ref([
  {
    code: '600519',
    name: '贵州茅台',
    price: 1820.50,
    changePercent: 2.34,
    changeAmount: 41.60,
    volume: 125000,
    turnover: 228250000,
    sparklineData: [1780, 1790, 1785, 1800, 1810, 1805, 1815, 1820]
  },
  {
    code: '000858',
    name: '五粮液',
    price: 156.80,
    changePercent: -1.23,
    changeAmount: -1.95,
    volume: 3450000,
    turnover: 540860000,
    sparklineData: [159, 158, 157.5, 157, 156.5, 157, 156.8, 156.8]
  },
  {
    code: '00700',
    name: '腾讯控股',
    price: 358.20,
    changePercent: 0.85,
    changeAmount: 3.02,
    volume: 8920000,
    turnover: 3194644000,
    sparklineData: [355, 356, 357, 358, 359, 358.5, 358, 358.2]
  },
  {
    code: '688256',
    name: '寒武纪',
    price: 89.50,
    changePercent: 10.00,
    changeAmount: 8.14,
    volume: 15600000,
    turnover: 1396200000,
    sparklineData: [81.5, 83, 85, 87, 88, 89, 89.5, 89.5]
  },
  {
    code: '002594',
    name: '比亚迪',
    price: 256.30,
    changePercent: 3.45,
    changeAmount: 8.55,
    volume: 6780000,
    turnover: 1737714000,
    sparklineData: [248, 250, 252, 254, 255, 256, 256.5, 256.3]
  },
  {
    code: '300750',
    name: '宁德时代',
    price: 198.60,
    changePercent: 5.23,
    changeAmount: 9.87,
    volume: 12340000,
    turnover: 2450964000,
    sparklineData: [189, 192, 194, 196, 197, 198, 198.5, 198.6]
  },
])

// 当前分组的股票列表
const currentStocks = computed(() => {
  // 这里后续可以根据 activeGroup 过滤
  return stocksData.value
})

const refreshCount = ref(0)
const selectedStock = ref(null)
const detailVisible = ref(false)

// 下拉刷新
async function handleRefresh() {
  return new Promise(resolve => {
    setTimeout(() => {
      refreshCount.value++
      // TODO: 调用真实API刷新数据
      resolve()
    }, 1500)
  })
}

// 点击股票
function handleStockClick(stock) {
  selectedStock.value = stock
  detailVisible.value = true
}

// 交易操作
function handleTrade({ stock, type }) {
  console.log('交易操作:', type, stock)
  // TODO: 跳转到交易页面
}

// 搜索股票
function handleSearch() {
  console.log('搜索股票')
  // TODO: 跳转到搜索页
}

// 添加股票
function handleAddStock() {
  console.log('添加股票')
  // TODO: 跳转到添加页面或打开搜索
}
</script>

<template>
  <div class="stock-list-page">
    <!-- 顶部分组标签 -->
    <div class="stock-list-header">
      <MTabs v-model="activeGroup" :tabs="groups" />
      <div class="header-actions">
        <button class="action-btn" @click="handleSearch">🔍</button>
        <button class="action-btn" @click="handleAddStock">➕</button>
      </div>
    </div>

    <!-- 股票列表 -->
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="stock-list-container">
        <!-- 有数据：虚拟滚动列表 -->
        <VirtualList
          v-if="currentStocks.length"
          :items="currentStocks"
          :item-height="96"
          class="stock-list"
        >
          <template #default="{ item }">
            <div class="stock-item-wrapper">
              <StockCard :stock="item" @click="handleStockClick(item)" />
            </div>
          </template>
        </VirtualList>

        <!-- 空状态 -->
        <MEmpty v-else description="还没有添加自选股票">
          <MButton type="primary" @click="handleAddStock">
            添加自选
          </MButton>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 底部提示 -->
    <div v-if="currentStocks.length" class="stock-list-footer">
      <p>共 {{ currentStocks.length }} 只股票</p>
      <p class="footer-tip">下拉刷新 · 已刷新 {{ refreshCount }} 次</p>
    </div>

    <!-- 股票详情抽屉 -->
    <StockDetailSheet
      v-model:show="detailVisible"
      :stock="selectedStock"
      @trade="handleTrade"
    />
  </div>
</template>

<style scoped>
.stock-list-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.stock-list-header {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  position: sticky;
  top: 0;
  z-index: var(--m-z-sticky);
  display: flex;
  align-items: center;
}

.header-actions {
  display: flex;
  padding: 0 var(--m-space-md);
  gap: var(--m-space-sm);
}

.action-btn {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 18px;
  cursor: pointer;
}

.action-btn:active {
  opacity: 0.6;
}

.stock-list-container {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.stock-list {
  height: 100%;
}

.stock-item-wrapper {
  padding: var(--m-space-xs) var(--m-space-md);
}

.stock-list-footer {
  text-align: center;
  padding: var(--m-space-lg);
  background: var(--m-bg-card);
  border-top: 1px solid var(--m-divider-color);
  color: var(--m-text-tertiary);
  font-size: var(--m-font-sm);
}

.stock-list-footer p {
  margin: var(--m-space-xs) 0;
}

.footer-tip {
  font-size: var(--m-font-xs);
}
</style>
