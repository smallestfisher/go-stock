<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MTabs from '../../components/base/MTabs.vue'
import StockCard from '../../components/cards/StockCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import { formatMoney } from '../../composables/useFormat'

// 对齐桌面端 market.vue「个股资金流向」Tab 的 9 个子分类
// 数据后续接 GetMoneyRankSina(rankType) 按 rankType 取不同维度
const flowTabs = [
  { label: '净流入额', value: 'netamount' },
  { label: '流出额', value: 'outamount' },
  { label: '净流入率', value: 'ratioamount' },
  { label: '主力净流入', value: 'r0_net' },
  { label: '主力流出', value: 'r0_out' },
  { label: '主力净流入率', value: 'r0_ratio' },
  { label: '散户净流入', value: 'r3_net' },
  { label: '散户流出', value: 'r3_out' },
  { label: '散户净流入率', value: 'r3_ratio' },
]

const activeTab = ref('netamount')

// 股票资金流数据（后续接 API 填充）
const stocks = ref([])

async function handleRefresh() {
  // 后续根据 activeTab 调对应 API
}
</script>

<template>
  <div class="stock-money-flow-page">
    <div class="flow-tabs">
      <MTabs v-model="activeTab" :tabs="flowTabs" />
    </div>

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
.stock-money-flow-page { height: 100%; display: flex; flex-direction: column; overflow: hidden; }
.flow-tabs { background: var(--m-bg-card); border-bottom: 1px solid var(--m-divider-color); }
.container { padding: var(--m-space-md); }
.stock-list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.flow-item { display: flex; flex-direction: column; gap: var(--m-space-xs); }
.flow-info { padding: 0 var(--m-space-md); font-size: var(--m-font-sm); font-weight: var(--m-font-weight-medium); }
</style>
