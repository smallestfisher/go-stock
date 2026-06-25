<script setup>
import { ref, computed, defineAsyncComponent } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MTabs from '../../components/base/MTabs.vue'

// 研究中心子页面Tab配置
const researchTabs = [
  { label: 'AI分析', value: 'ai' },
  { label: '行业分析', value: 'industry' },
  { label: '技术分析', value: 'technical' },
  { label: '基本面', value: 'fundamental' },
  { label: '情绪分析', value: 'sentiment' },
  { label: '研报精选', value: 'reports' },
]

const activeTab = ref('ai')

// 动态加载子页面组件
const componentMap = {
  'ai': defineAsyncComponent(() => import('./AiAnalysisPage.vue')),
  'industry': defineAsyncComponent(() => import('./IndustryAnalysisPage.vue')),
  'technical': defineAsyncComponent(() => import('./TechnicalAnalysisPage.vue')),
  'fundamental': defineAsyncComponent(() => import('./FundamentalAnalysisPage.vue')),
  'sentiment': defineAsyncComponent(() => import('./SentimentAnalysisPage.vue')),
  'reports': defineAsyncComponent(() => import('./ReportsPage.vue')),
}

const currentComponent = computed(() => {
  return componentMap[activeTab.value]
})
</script>

<template>
  <div class="research-page">
    <!-- 顶部导航 -->
    <PageHeader title="研究中心" show-back />

    <!-- Tab切换栏 -->
    <div class="research-tabs">
      <MTabs v-model="activeTab" :tabs="researchTabs" />
    </div>

    <!-- 内容区 -->
    <div class="research-content">
      <KeepAlive>
        <component :is="currentComponent" />
      </KeepAlive>
    </div>
  </div>
</template>

<style scoped>
.research-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.research-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  position: sticky;
  top: 56px;
  z-index: var(--m-z-sticky);
}

.research-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}
</style>
