<script setup>
import {computed, h, onBeforeMount, onBeforeUnmount, onMounted,onUnmounted, ref,reactive} from 'vue'
import {GetAIResponseResultList} from "../api/app";
import {NButton, NEllipsis, NText} from "naive-ui";
import ResearchReport from "./researchReport.vue";
import AiRecommendStocksList from "./aiRecommendStocksList.vue";
import PromptTemplateList from "./promptTemplateList.vue";
import AllStockList from "./allStockList.vue";
import AllStockInfoList from "./allStockInfoList.vue";
import CronTaskManager from "./cron-task-manager.vue";
import TradingRecordManager from "./TradingRecordManager.vue";
import StockChangesMonitor from "./stockChangesMonitor.vue";
import MCPServiceManager from "./mcp-server-manager.vue";
import SkillManager from "./skill-manager.vue";
import UplimitLadder from "./uplimitLadder.vue";
import PromptPlaza from "./promptPlaza.vue";
import PromptQa from "./promptQa.vue";
import SelectStock from "./SelectStock.vue";
import {EventsOff, EventsOn} from "../api/runtime";
import {useRoute} from 'vue-router'
import {useDevice} from "../composables/useDevice";

const {isMobile} = useDevice()


const nowTab = ref("AI分析报告")
const route = useRoute()
// 移动端：分组导航（两级：分类 → 具体功能），避免一行十几个横向滚动标签
const mobileGroups = [
  {
    category: 'AI分析',
    icon: '📈',
    tabs: ['AI分析报告', '股票推荐记录', '异动监控', '涨停梯队'],
  },
  {
    category: '提示词',
    icon: '💬',
    tabs: ['提示词模板', '提示词广场', '问答广场'],
  },
  {
    category: '选股',
    icon: '🔍',
    tabs: ['形态选股', '指标选股'],
  },
  {
    category: '系统',
    icon: '⚙️',
    tabs: ['定时任务', '交易日志', 'MCP服务'],
  },
]
const mobileActiveGroup = ref('AI分析')
function mobileGroupOf(tabName) {
  const g = mobileGroups.find(g => g.tabs.includes(tabName))
  return g ? g.category : mobileGroups[0].category
}
function selectMobileGroup(category) {
  mobileActiveGroup.value = category
}

// 移动端当前分类下的功能项
const mobileCurrentTabs = computed(() => {
  const g = mobileGroups.find(g => g.category === mobileActiveGroup.value)
  return g ? g.tabs : []
})
onBeforeMount(() => {
  nowTab.value = route.query.name
  mobileActiveGroup.value = mobileGroupOf(nowTab.value)
})

onBeforeUnmount(() => {
  EventsOff("changeResearchTab")
})

onUnmounted(() => {

});

EventsOn("changeResearchTab", async (msg) => {
  console.log("changeResearchTab", msg)
  updateTab(msg.name)
})
function updateTab(name) {
  nowTab.value = name
  mobileActiveGroup.value = mobileGroupOf(name)
}
</script>

<template>
  <n-card class="research-page-shell">
    <!-- 移动端：分类 + 功能两级菜单（替换原生横向滚动标签栏） -->
    <div v-if="isMobile" class="research-mobile-nav">
      <div class="research-mobile-nav__groups">
        <button
            v-for="g in mobileGroups"
            :key="g.category"
            type="button"
            class="research-group-tile"
            :class="{ 'research-group-tile--active': mobileActiveGroup === g.category }"
            @click="selectMobileGroup(g.category)"
        >
          <span class="research-group-tile__icon">{{ g.icon }}</span>
          <span class="research-group-tile__name">{{ g.category }}</span>
        </button>
      </div>
      <div class="research-mobile-nav__tabs">
        <button
            v-for="tab in mobileCurrentTabs"
            :key="tab"
            type="button"
            class="research-tab-chip"
            :class="{ 'research-tab-chip--active': nowTab === tab }"
            @click="updateTab(tab)"
        >
          {{ tab }}
        </button>
      </div>
    </div>
    <n-tabs class="research-mobile-tabs" :class="{ 'research-mobile-tabs--native-hidden': isMobile }" type="line" animated @update-value="updateTab" :value="nowTab" style="">
      <n-tab-pane name="AI分析报告">
        <ResearchReport/>
      </n-tab-pane>
      <n-tab-pane name="股票推荐记录">
        <AiRecommendStocksList/>
      </n-tab-pane>
      <n-tab-pane name="异动监控">
        <StockChangesMonitor/>
      </n-tab-pane>
      <n-tab-pane name="涨停梯队">
        <UplimitLadder/>
      </n-tab-pane>
      <n-tab-pane name="提示词模板">
        <PromptTemplateList/>
      </n-tab-pane>
      <n-tab-pane name="提示词广场">
        <PromptPlaza/>
      </n-tab-pane>
      <n-tab-pane name="问答广场">
        <PromptQa/>
      </n-tab-pane>
      <n-tab-pane name="形态选股">
        <AllStockList/>
      </n-tab-pane>
      <n-tab-pane name="指标选股">
        <SelectStock/>
      </n-tab-pane>
      <n-tab-pane name="定时任务">
        <CronTaskManager />
      </n-tab-pane>
      <n-tab-pane name="交易日志">
        <TradingRecordManager />
      </n-tab-pane>
<!--      <n-tab-pane name="全部股票信息">-->
<!--        <AllStockInfoList/>-->
<!--      </n-tab-pane>-->
      <n-tab-pane name="MCP服务">
        <MCPServiceManager/>
      </n-tab-pane>
<!--      <n-tab-pane name="技能管理">-->
<!--        <SkillManager/>-->
<!--      </n-tab-pane>-->
    </n-tabs>
  </n-card>
</template>

<style scoped>
@media (max-width: 768px) {
  .research-page-shell {
    margin: 0 8px calc(var(--mobile-bottom-nav-height) + var(--safe-bottom) + 8px);
  }

  .research-page-shell :deep(.n-card__content) {
    padding: 10px;
  }

  /* 移动端用自定义两级菜单，隐藏原生横向滚动标签栏 */
  .research-mobile-tabs--native-hidden :deep(.n-tabs-nav) {
    display: none !important;
  }

  .research-mobile-tabs :deep(.n-tab-pane) {
    min-width: 0;
    padding-top: 8px;
  }
}

/* ============ 移动端两级导航（仅在 isMobile 渲染） ============ */
.research-mobile-nav {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 6px;
}

.research-mobile-nav__groups {
  display: grid;
  gap: 8px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.research-group-tile {
  align-items: center;
  appearance: none;
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 10px;
  color: inherit;
  display: flex;
  flex-direction: column;
  font: inherit;
  gap: 2px;
  padding: 8px 4px;
}

.research-group-tile--active {
  background: var(--n-color-target, rgba(24, 160, 88, 0.1));
  border-color: #18a058;
  color: #18a058;
}

.research-group-tile__icon {
  font-size: 18px;
}

.research-group-tile__name {
  font-size: 12px;
  font-weight: 600;
}

.research-mobile-nav__tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.research-tab-chip {
  appearance: none;
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 16px;
  color: var(--n-text-color, #333);
  font: inherit;
  font-size: 13px;
  padding: 6px 14px;
}

.research-tab-chip--active {
  background: #18a058;
  border-color: #18a058;
  color: #fff;
  font-weight: 700;
}
</style>
