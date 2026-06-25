<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 全球股指数据（后续接 GlobalStockIndexes，按 common/america/europe/asia/other 区域分组填充）
const indexData = ref([])

const refreshCount = ref(0)

// 下拉刷新（后续接真实 API）
async function handleRefresh() {
  refreshCount.value++
}

// 点击指数
function handleIndexClick(index) {
  console.log('点击指数:', index)
  // TODO: 打开指数详情
}

// 地区标签
const regionMap = {
  'CN': '🇨🇳 A股',
  'HK': '🇭🇰 港股',
  'US': '🇺🇸 美股',
  'JP': '🇯🇵 日经',
  'UK': '🇬🇧 英国',
  'DE': '🇩🇪 德国',
}
</script>

<template>
  <div class="global-index-page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="index-container">
        <MCard v-if="indexData.length">
          <div class="index-list">
            <div
              v-for="index in indexData"
              :key="index.code"
              class="index-item"
              @click="handleIndexClick(index)"
            >
              <div class="index-info">
                <div class="index-name">{{ index.name }}</div>
                <div class="index-meta">
                  <span class="index-region">{{ regionMap[index.region] }}</span>
                  <span class="index-code">{{ index.code }}</span>
                </div>
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
.global-index-page {
  height: 100%;
  overflow: hidden;
}

.index-container {
  padding: var(--m-space-md);
}

.index-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.index-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.index-item:active {
  transform: scale(0.98);
  background: var(--m-divider-color);
}

.index-info {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.index-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.index-meta {
  display: flex;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
}

.index-region {
  color: var(--m-text-secondary);
}

.index-code {
  color: var(--m-text-tertiary);
}

.index-price {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--m-space-xs);
}

.price-value {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
}
</style>
