<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import AiSuggestCard from '../../components/cards/AiSuggestCard.vue'
import MCard from '../../components/base/MCard.vue'
import MButton from '../../components/base/MButton.vue'

// AI 分析结果（后续接 SummaryStockNews / ChatWithAgent 流式生成，先留空）
const aiAnalysis = ref(null)

const aiLoading = ref(false)

// 历史分析记录（后续接 GetAIResponseResultList 填充）
const analysisHistory = ref([])

async function handleRefresh() {
  // 后续接真实 API
}

function handleAiRefresh() {
  aiLoading.value = true
  setTimeout(() => {
    aiLoading.value = false
  }, 2000)
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- AI分析卡片 -->
        <AiSuggestCard
          :suggestion="aiAnalysis"
          :loading="aiLoading"
          @refresh="handleAiRefresh"
        />

        <!-- 历史分析 -->
        <MCard>
          <h3 class="section-title">历史分析</h3>
          <div v-if="analysisHistory.length" class="history-list">
            <div v-for="(item, i) in analysisHistory" :key="i" class="history-item">
              <div class="history-date">{{ item.date }}</div>
              <div class="history-title">{{ item.title }}</div>
              <div class="history-summary">{{ item.summary }}</div>
            </div>
          </div>
          <p v-else class="empty-tip">暂无历史分析记录</p>
        </MCard>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.section-title { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-lg); }
.history-list { display: flex; flex-direction: column; gap: var(--m-space-lg); }
.history-item { padding-bottom: var(--m-space-lg); border-bottom: 1px solid var(--m-divider-color); }
.history-item:last-child { border-bottom: none; }
.history-date { font-size: var(--m-font-xs); color: var(--m-text-tertiary); margin-bottom: var(--m-space-xs); }
.history-title { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-sm); }
.history-summary { font-size: var(--m-font-sm); color: var(--m-text-secondary); }
.empty-tip { text-align: center; padding: var(--m-space-lg) 0; color: var(--m-text-tertiary); font-size: var(--m-font-sm); }
</style>
