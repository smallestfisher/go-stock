<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MTabs from '../../components/base/MTabs.vue'
import StockCard from '../../components/cards/StockCard.vue'

// 指标选股（对齐桌面端 SelectStock.vue）：热门策略 + 自建策略切换
// 数据后续接 GetHotStrategy / GetAllCustomStrategies 填充，结果接 SearchStock
const strategyTabs = [
  { label: '热门策略', value: 'hot' },
  { label: '我的策略', value: 'custom' },
]
const activeStrategyTab = ref('hot')

// 策略列表
const hotStrategies = ref([])
const customStrategies = ref([])

// 选股结果
const resultStocks = ref([])

async function handleRefresh() {
  // 后续接真实 API
}

function runStrategy(strategy) {
  // 后续接 SearchStock 执行选股
}
</script>

<template>
  <div class="page">
    <div class="indicator-header">
      <MTabs v-model="activeStrategyTab" :tabs="strategyTabs" />
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 热门策略 -->
        <template v-if="activeStrategyTab === 'hot'">
          <MCard v-if="hotStrategies.length" padding="none">
            <div class="strategy-list">
              <div
                v-for="(s, i) in hotStrategies"
                :key="i"
                class="strategy-item"
                @click="runStrategy(s)"
              >
                <div class="strategy-name">{{ s.name }}</div>
                <p v-if="s.desc" class="strategy-desc">{{ s.desc }}</p>
              </div>
            </div>
          </MCard>
          <MEmpty v-else description="暂无热门策略" />
        </template>

        <!-- 我的策略 -->
        <template v-else>
          <MCard v-if="customStrategies.length" padding="none">
            <div class="strategy-list">
              <div
                v-for="(s, i) in customStrategies"
                :key="i"
                class="strategy-item"
                @click="runStrategy(s)"
              >
                <div class="strategy-name">{{ s.name }}</div>
                <p v-if="s.desc" class="strategy-desc">{{ s.desc }}</p>
              </div>
            </div>
          </MCard>
          <MEmpty v-else description="还没有自定义策略" />
        </template>

        <!-- 选股结果 -->
        <div v-if="resultStocks.length" class="result-section">
          <h4 class="result-title">选股结果 ({{ resultStocks.length }})</h4>
          <div class="stock-list">
            <div v-for="stock in resultStocks" :key="stock.code" class="stock-item">
              <StockCard :stock="stock" :show-sparkline="false" />
            </div>
          </div>
        </div>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; display: flex; flex-direction: column; overflow: hidden; }
.indicator-header { background: var(--m-bg-card); border-bottom: 1px solid var(--m-divider-color); }
.container { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.strategy-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.strategy-item { padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.strategy-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); margin-bottom: var(--m-space-xs); }
.strategy-desc { font-size: var(--m-font-sm); color: var(--m-text-secondary); }
.result-section { margin-top: var(--m-space-md); }
.result-title { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-md); }
.stock-list { display: flex; flex-direction: column; gap: var(--m-space-sm); }
</style>
