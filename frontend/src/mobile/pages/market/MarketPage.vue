<script setup>
import { ref, computed, watch, defineAsyncComponent } from 'vue'
import MTabs from '../../components/base/MTabs.vue'

// 市场子页面Tab配置
const marketTabs = [
  { label: '快讯', value: 'news' },
  { label: '全球股指', value: 'global' },
  { label: '重大指数', value: 'major' },
  { label: '行业排名', value: 'industry' },
  { label: '个股资金', value: 'stock-flow' },
  { label: '板块资金', value: 'sector-flow' },
  { label: '概念资金', value: 'concept-flow' },
  { label: '龙虎榜', value: 'dragon-tiger' },
  { label: '个股研报', value: 'stock-report' },
  { label: '公司公告', value: 'company-notice' },
  { label: '行业研究', value: 'industry-research' },
  { label: '当前热门', value: 'hot-stock' },
  { label: '名站优选', value: 'featured' },
]

const activeTab = ref('news')

// 动态加载子页面组件
const componentMap = {
  'news': defineAsyncComponent(() => import('./NewsListPage.vue')),
  'global': defineAsyncComponent(() => import('./GlobalIndexPage.vue')),
  'major': defineAsyncComponent(() => import('./MajorIndexPage.vue')),
  'industry': defineAsyncComponent(() => import('./IndustryRankPage.vue')),
  'stock-flow': defineAsyncComponent(() => import('./StockMoneyFlowPage.vue')),
  'sector-flow': defineAsyncComponent(() => import('./SectorMoneyFlowPage.vue')),
  'concept-flow': defineAsyncComponent(() => import('./ConceptMoneyFlowPage.vue')),
  'dragon-tiger': defineAsyncComponent(() => import('./DragonTigerPage.vue')),
  'stock-report': defineAsyncComponent(() => import('./StockReportPage.vue')),
  'company-notice': defineAsyncComponent(() => import('./CompanyNoticePage.vue')),
  'industry-research': defineAsyncComponent(() => import('./IndustryResearchPage.vue')),
  'hot-stock': defineAsyncComponent(() => import('./HotStockPage.vue')),
  'featured': defineAsyncComponent(() => import('./FeaturedSitesPage.vue')),
}

// 当前激活的组件
const currentComponent = computed(() => {
  return componentMap[activeTab.value]
})

// 切换Tab时记录位置（后续可用于缓存滚动位置）
watch(activeTab, (newVal) => {
  console.log('切换到:', newVal)
})
</script>

<template>
  <div class="market-page">
    <!-- Tab切换栏 -->
    <div class="market-tabs">
      <MTabs v-model="activeTab" :tabs="marketTabs" />
    </div>

    <!-- 内容区 -->
    <div class="market-content">
      <KeepAlive>
        <component :is="currentComponent" />
      </KeepAlive>
    </div>
  </div>
</template>

<style scoped>
.market-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.market-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  position: sticky;
  top: 0;
  z-index: var(--m-z-sticky);
}

.market-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}
</style>
