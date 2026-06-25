<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'
import StockCard from '../../components/cards/StockCard.vue'
import VirtualList from '../../components/widgets/VirtualList.vue'
import StockDetailSheet from '../../components/sheets/StockDetailSheet.vue'

const router = useRouter()

// 分组数据（后续接 GetGroupList 填充，先留「全部」）
const groups = ref([
  { label: '全部', value: 0 },
])

const activeGroup = ref(0)

// 股票列表（后续接 GetFollowList 填充）
const stocksData = ref([])

// 当前分组的股票列表
const currentStocks = computed(() => {
  return stocksData.value
})

const refreshCount = ref(0)
const selectedStock = ref(null)
const detailVisible = ref(false)
const loading = ref(false)

// 下拉刷新（后续接真实 API）
async function handleRefresh() {
  refreshCount.value++
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
    <!-- 顶部导航（只保留标题） -->
    <PageHeader title="自选股票" />

    <!-- 分组标签 + 添加 -->
    <div class="stock-list-tabs">
      <MTabs v-model="activeGroup" :tabs="groups" />
      <button class="add-stock-btn" @click="handleAddStock">➕</button>
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
        <MEmpty v-else-if="!loading" description="还没有添加自选股票">
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

.stock-list-tabs {
  display: flex;
  align-items: center;
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  position: sticky;
  top: 56px;
  z-index: var(--m-z-sticky);
}

/* 让 MTabs 占满宽度 */
.stock-list-tabs :deep(.m-tabs) {
  flex: 1;
  min-width: 0;
}

.add-stock-btn {
  flex-shrink: 0;
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 20px;
  color: var(--m-text-primary);
  cursor: pointer;
}

.add-stock-btn:active {
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
