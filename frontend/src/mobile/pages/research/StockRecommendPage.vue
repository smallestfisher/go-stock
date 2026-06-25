<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'

// AI 推荐股票记录（后续接 GetAiRecommendStocksList 填充）
const recommends = ref([])

async function handleRefresh() {
  // 后续接真实 API
}

// 字段映射参考桌面端 aiRecommendStocksList.vue：
// { stockCode, stockName, price, changePercent, recommendDate, reason }
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="recommends.length" padding="none">
          <div class="recommend-list">
            <div v-for="(item, i) in recommends" :key="i" class="recommend-item">
              <div class="item-header">
                <div class="item-name">{{ item.stockName }}</div>
                <PercentTag :value="item.changePercent" size="small" />
              </div>
              <div class="item-code">{{ item.stockCode }}</div>
              <p v-if="item.reason" class="item-reason">{{ item.reason }}</p>
              <div class="item-footer">
                <span>{{ item.recommendDate }}</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无 AI 推荐记录" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.recommend-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-lg); }
.recommend-item { padding-bottom: var(--m-space-lg); border-bottom: 1px solid var(--m-divider-color); }
.recommend-item:last-child { border-bottom: none; }
.item-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--m-space-xs); }
.item-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.item-code { font-size: var(--m-font-xs); color: var(--m-text-tertiary); margin-bottom: var(--m-space-sm); }
.item-reason { font-size: var(--m-font-sm); color: var(--m-text-secondary); line-height: var(--m-line-height-normal); margin-bottom: var(--m-space-sm); }
.item-footer { display: flex; justify-content: flex-end; font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
</style>
