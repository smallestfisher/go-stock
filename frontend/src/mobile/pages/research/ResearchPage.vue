<script setup>
import { ref, computed, defineAsyncComponent } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MIcon from '../../components/base/MIcon.vue'

// 研究中心功能分组（严格对齐桌面端 researchIndex.vue 的 mobileGroups）
// 每个分组下若干功能 Tab，value 对应下方 componentMap
const groups = [
  {
    category: 'AI分析',
    icon: 'trend-up',
    tabs: [
      { label: 'AI分析报告', value: 'ai' },
      { label: '股票推荐记录', value: 'recommend' },
      { label: '异动监控', value: 'changes' },
      { label: '涨停梯队', value: 'uplimit' },
    ],
  },
  {
    category: '提示词',
    icon: 'comment',
    tabs: [
      { label: '提示词模板', value: 'template' },
      { label: '提示词广场', value: 'plaza' },
      { label: '问答广场', value: 'qa' },
    ],
  },
  {
    category: '选股',
    icon: 'search',
    tabs: [
      { label: '形态选股', value: 'shape' },
      { label: '指标选股', value: 'indicator' },
    ],
  },
  {
    category: '系统',
    icon: 'settings',
    tabs: [
      { label: '定时任务', value: 'cron' },
      { label: '交易日志', value: 'trade' },
      { label: 'MCP服务', value: 'mcp' },
    ],
  },
]

const activeGroup = ref('AI分析')
const activeTab = ref('ai')

// 动态加载子页面组件
const componentMap = {
  // AI分析组
  'ai': defineAsyncComponent(() => import('./AiAnalysisPage.vue')),
  'recommend': defineAsyncComponent(() => import('./StockRecommendPage.vue')),
  'changes': defineAsyncComponent(() => import('./StockChangesPage.vue')),
  'uplimit': defineAsyncComponent(() => import('./UplimitLadderPage.vue')),
  // 提示词组
  'template': defineAsyncComponent(() => import('./PromptTemplatePage.vue')),
  'plaza': defineAsyncComponent(() => import('./PromptPlazaPage.vue')),
  'qa': defineAsyncComponent(() => import('./PromptQaPage.vue')),
  // 选股组
  'shape': defineAsyncComponent(() => import('./ShapeSelectPage.vue')),
  'indicator': defineAsyncComponent(() => import('./IndicatorSelectPage.vue')),
  // 系统组
  'cron': defineAsyncComponent(() => import('./CronTaskPage.vue')),
  'trade': defineAsyncComponent(() => import('./TradingRecordPage.vue')),
  'mcp': defineAsyncComponent(() => import('./McpServerPage.vue')),
}

const currentComponent = computed(() => componentMap[activeTab.value])

// 当前分类下的功能项
const currentTabs = computed(() => {
  const g = groups.find(g => g.category === activeGroup.value)
  return g ? g.tabs : []
})

function selectGroup(category) {
  activeGroup.value = category
  const g = groups.find(g => g.category === category)
  // 切换分类时跳到该分类第一个功能，避免内容停留在旧 Tab
  if (g && g.tabs.length && !g.tabs.some(t => t.value === activeTab.value)) {
    activeTab.value = g.tabs[0].value
  }
}

function selectTab(value) {
  activeTab.value = value
}
</script>

<template>
  <div class="research-page">
    <!-- 顶部导航 -->
    <PageHeader title="研究中心" />

    <!-- 分组导航：紧凑单行分段栏（横向可滚动），把首屏高度让给内容 -->
    <div class="research-nav">
      <div class="research-groups">
        <button
          v-for="g in groups"
          :key="g.category"
          type="button"
          class="group-seg"
          :class="{ 'group-seg--active': activeGroup === g.category }"
          @click="selectGroup(g.category)"
        >
          <MIcon class="group-seg__icon" :name="g.icon" :size="16" />
          <span class="group-seg__name">{{ g.category }}</span>
        </button>
      </div>

      <!-- 当前分类下的功能 chip：仅当有 2 个及以上功能才显示，避免冗余单 chip 行 -->
      <div v-if="currentTabs.length > 1" class="research-tabs">
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

.research-nav {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

/* 分类分段栏：单行、横向可滚动，替代原先占一整行的 4 大方格 */
.research-groups {
  display: flex;
  gap: var(--m-space-xs);
  overflow-x: auto;
  scrollbar-width: none;
}

.research-groups::-webkit-scrollbar {
  display: none;
}

.group-seg {
  flex: 1 0 auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--m-space-xs);
  padding: var(--m-space-xs) var(--m-space-md);
  background: var(--m-bg-primary);
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-full);
  color: var(--m-text-secondary);
  font: inherit;
  font-size: var(--m-font-sm);
  white-space: nowrap;
}

.group-seg--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.group-seg__icon {
  flex-shrink: 0;
}

.research-tabs {
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

.research-content {
  flex: 1;
  overflow: hidden;
  position: relative;
}
</style>
