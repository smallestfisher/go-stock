<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import IndustryCard from '../../components/cards/IndustryCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 对齐桌面端 market.vue「行业排名」Tab 的 4 个子分类
// 数据后续分别接：行业涨幅 → GetIndustryRank；行业资金 → GetIndustryMoneyRankSina；
// 证监会行业 → GetBKFundFlowList；概念板块 → GetConceptFundFlowTopList
const rankTabs = [
  { label: '行业涨幅', value: 'change' },
  { label: '行业资金', value: 'money' },
  { label: '证监会行业', value: 'csrc' },
  { label: '概念板块', value: 'concept' },
]

const activeTab = ref('change')

// 行业数据（后续接 API 填充）
const industries = ref([])

async function handleRefresh() {
  // 后续根据 activeTab 调对应 API
}
</script>

<template>
  <div class="industry-rank-page">
    <div class="rank-tabs">
      <MTabs v-model="activeTab" :tabs="rankTabs" />
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <IndustryCard
          v-if="industries.length"
          :industries="industries"
          :max-show="20"
          @view-all="() => {}"
        />
        <MEmpty v-else description="暂无行业排名数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.industry-rank-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.rank-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.container {
  padding: var(--m-space-md);
}
</style>
