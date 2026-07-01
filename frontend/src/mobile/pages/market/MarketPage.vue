<script setup>
import { ref, computed, defineAsyncComponent } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MIcon from '../../components/base/MIcon.vue'

// 市场行情分组导航：严格对齐桌面端 market.vue 的 marketMobileGroups
// 两级（分类格子 → 功能 chip），与 ResearchPage 同范式，统一移动端导航风格。
// 已按需求移除"名站优选"。
const groups = [
  {
    category: '行情',
    icon: 'news',
    tabs: [
      { label: '市场快讯', value: 'news' },
      { label: '当前热门', value: 'hot-stock' },
    ],
  },
  {
    category: '指数',
    icon: 'globe',
    tabs: [
      { label: '全球股指', value: 'global' },
      { label: '重大指数', value: 'major' },
    ],
  },
  {
    category: '资金',
    icon: 'money',
    tabs: [
      { label: '行业排名', value: 'industry' },
      { label: '个股资金', value: 'stock-flow' },
      { label: '板块资金', value: 'sector-flow' },
      { label: '概念资金', value: 'concept-flow' },
    ],
  },
  {
    category: '研报',
    icon: 'chart',
    tabs: [
      { label: '龙虎榜', value: 'dragon-tiger' },
      { label: '个股研报', value: 'stock-report' },
      { label: '公司公告', value: 'company-notice' },
      { label: '行业研究', value: 'industry-research' },
    ],
  },
]

const activeGroup = ref('行情')
const activeTab = ref('news')

// 动态加载子页面组件
const componentMap = {
  'news': defineAsyncComponent(() => import('./NewsListPage.vue')),
  'hot-stock': defineAsyncComponent(() => import('./HotStockPage.vue')),
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
}

const currentComponent = computed(() => componentMap[activeTab.value])

// 当前分类下的功能项
const currentTabs = computed(() => {
  const g = groups.find(g => g.category === activeGroup.value)
  return g ? g.tabs : []
})

function selectGroup(category) {
  activeGroup.value = category
  // 切换分类时跳到该分类第一个功能，避免内容停留在旧 Tab（对齐桌面端 selectMarketMobileGroup）
  const g = groups.find(g => g.category === category)
  if (g && g.tabs.length && !g.tabs.some(t => t.value === activeTab.value)) {
    activeTab.value = g.tabs[0].value
  }
}

function selectTab(value) {
  activeTab.value = value
}
</script>

<template>
  <div class="market-page">
    <!-- 顶部导航 -->
    <PageHeader title="市场行情" />

    <!-- 分组导航：分类格子 + 功能 chip（与 ResearchPage 同范式） -->
    <div class="market-nav">
      <div class="market-groups">
        <button
          v-for="g in groups"
          :key="g.category"
          type="button"
          class="group-tile"
          :class="{ 'group-tile--active': activeGroup === g.category }"
          @click="selectGroup(g.category)"
        >
          <span class="group-tile__icon"><MIcon :name="g.icon" :size="18" /></span>
          <span class="group-tile__name">{{ g.category }}</span>
        </button>
      </div>

      <!-- 当前分类下的功能 chip -->
      <div class="market-tabs">
        <button
          v-for="tab in currentTabs"
          :key="tab.value"
          type="button"
          class="tab-chip"
          :class="{ 'tab-chip--active': activeTab === tab.value }"
          @click="selectTab(tab.value)"
        >
          {{ tab.label }}
        </button>
      </div>
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

.market-nav {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.market-groups {
  display: grid;
  gap: var(--m-space-sm);
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.group-tile {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: var(--m-space-sm) var(--m-space-xs);
  background: var(--m-bg-card);
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-md);
  color: var(--m-text-primary);
  font: inherit;
}

.group-tile--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.group-tile__icon {
  font-size: 18px;
}

.group-tile__name {
  font-size: var(--m-font-xs);
  font-weight: var(--m-font-weight-medium);
}

.market-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-sm);
}

.tab-chip {
  padding: var(--m-space-xs) var(--m-space-md);
  background: var(--m-bg-card);
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-full);
  color: var(--m-text-primary);
  font: inherit;
  font-size: var(--m-font-sm);
}

.tab-chip--active {
  background: var(--m-color-rise);
  border-color: var(--m-color-rise);
  color: #fff;
  font-weight: var(--m-font-weight-medium);
}

.market-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}
</style>
