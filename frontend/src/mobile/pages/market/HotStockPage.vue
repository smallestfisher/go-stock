<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 对齐桌面端 market.vue「当前热门」Tab 的 6 个子分类
// 数据后续分别接：全球/沪深/港股/美股 → HotStock(marketType)；热门话题 → HotTopic；重大事件 → InvestCalendarTimeLine
const hotTabs = [
  { label: '全球', value: '10' },
  { label: '沪深', value: '12' },
  { label: '港股', value: '13' },
  { label: '美股', value: '11' },
  { label: '热门话题', value: 'topic' },
  { label: '重大事件', value: 'event' },
]

const activeTab = ref('10')

// 各 Tab 数据（后续接 API 填充）
const stocksData = ref([])
const topicsData = ref([])
const eventsData = ref([])

async function handleRefresh() {
  // 后续根据 activeTab 调对应 API
}
</script>

<template>
  <div class="page">
    <div class="hot-tabs">
      <MTabs v-model="activeTab" :tabs="hotTabs" />
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 热门股票排行（全球/沪深/港股/美股） -->
        <template v-if="['10', '11', '12', '13'].includes(activeTab)">
          <MEmpty v-if="!stocksData.length" description="暂无热门股票数据" />
        </template>

        <!-- 热门话题 -->
        <template v-else-if="activeTab === 'topic'">
          <MEmpty v-if="!topicsData.length" description="暂无热门话题" />
        </template>

        <!-- 重大事件时间轴 -->
        <template v-else-if="activeTab === 'event'">
          <MEmpty v-if="!eventsData.length" description="暂无重大事件" />
        </template>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.hot-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.container {
  padding: var(--m-space-md);
}
</style>
