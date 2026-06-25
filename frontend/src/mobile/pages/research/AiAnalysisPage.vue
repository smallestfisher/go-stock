<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import AiSuggestCard from '../../components/cards/AiSuggestCard.vue'
import MCard from '../../components/base/MCard.vue'
import MButton from '../../components/base/MButton.vue'

const aiAnalysis = ref({
  title: '基于你的持仓，AI建议关注半导体板块机会',
  content: '近期AI芯片需求激增，半导体板块持续强势。结合你的自选股腾讯控股与多家半导体公司的业务合作，建议关注相关概念股机会。同时注意控制仓位，当前市场波动较大。',
  stocks: [
    { code: '688256', name: '寒武纪' },
    { code: '688981', name: '中芯国际' },
    { code: '603986', name: '兆易创新' },
  ],
  action: '查看详细分析报告'
})

const aiLoading = ref(false)

const analysisHistory = ref([
  { date: '2024-06-20', title: '本周投资策略：关注AI芯片板块', summary: '综合技术面和基本面分析...' },
  { date: '2024-06-17', title: '风险提示：警惕短期回调风险', summary: '市场情绪过热，建议适当减仓...' },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
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
          <div class="history-list">
            <div v-for="(item, i) in analysisHistory" :key="i" class="history-item">
              <div class="history-date">{{ item.date }}</div>
              <div class="history-title">{{ item.title }}</div>
              <div class="history-summary">{{ item.summary }}</div>
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
.history-list { display: flex; flex-direction: column; gap: var(--m-space-lg); }
.history-item { padding-bottom: var(--m-space-lg); border-bottom: 1px solid var(--m-divider-color); }
.history-item:last-child { border-bottom: none; }
.history-date { font-size: var(--m-font-xs); color: var(--m-text-tertiary); margin-bottom: var(--m-space-xs); }
.history-title { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-sm); }
.history-summary { font-size: var(--m-font-sm); color: var(--m-text-secondary); }
</style>
