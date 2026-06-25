<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 重大指数数据（后续接 GetStockRealTimePrice 批量获取，覆盖沪深300/中证500/中证1000/科创50/上证50 等）
const indexData = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="major-index-page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="indexData.length">
          <div class="index-list">
            <div v-for="index in indexData" :key="index.code" class="index-item">
              <div class="index-info">
                <div class="index-name">{{ index.name }}</div>
                <div class="index-code">{{ index.code }}</div>
              </div>
              <div class="index-price">
                <div class="price-value" :class="{
                  'm-rise': index.changePercent > 0,
                  'm-fall': index.changePercent < 0
                }">
                  {{ index.price.toFixed(2) }}
                </div>
                <PercentTag :value="index.changePercent" size="small" />
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
.major-index-page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.index-list { display: flex; flex-direction: column; gap: var(--m-space-sm); }
.index-item { display: flex; justify-content: space-between; align-items: center; padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.index-info { display: flex; flex-direction: column; gap: var(--m-space-xs); }
.index-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.index-code { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.index-price { display: flex; flex-direction: column; align-items: flex-end; gap: var(--m-space-xs); }
.price-value { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-bold); font-variant-numeric: tabular-nums; }
</style>
