<script setup>
import { ref, computed, defineAsyncComponent } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 研究中心功能分组（严格对齐桌面端 researchIndex.vue 的 mobileGroups）
// 每个分组下若干功能 Tab，value 对应下方 componentMap
const groups = [
  {
    category: 'AI分析',
    icon: '📈',
    tabs: [
      { label: 'AI分析报告', value: 'ai' },
      { label: '股票推荐记录', value: 'recommend' },
      { label: '异动监控', value: 'changes' },
      { label: '涨停梯队', value: 'uplimit' },
    ],
  },
  {
    category: '提示词',
    icon: '💬',
    tabs: [
      { label: '提示词模板', value: 'template' },
      { label: '提示词广场', value: 'plaza' },
      { label: '问答广场', value: 'qa' },
    ],
  },
  {
    category: '选股',
    icon: '🔍',
    tabs: [
      { label: '形态选股', value: 'shape' },
      { label: '指标选股', value: 'indicator' },
    ],
  },
  {
    category: '系统',
    icon: '⚙️',
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

    <!-- 分组导航：分类格子 -->
    <div class="research-nav">
      <div class="research-groups">
        <button
          v-for="g in groups"
          :key="g.category"
          type="button"
          class="group-tile"
          :class="{ 'group-tile--active': activeGroup === g.category }"
          @click="selectGroup(g.category)"
        >
          <span class="group-tile__icon">{{ g.icon }}</span>
          <span class="group-tile__name">{{ g.category }}</span>
        </button>
      </div>

      <!-- 当前分类下的功能 chip -->
      <div class="research-tabs">
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
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.research-groups {
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
