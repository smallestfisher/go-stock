<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount, watch } from 'vue'
import { HotStock, HotTopic, InvestCalendarTimeLine, ClsCalendar } from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import HotStockCard from '../../components/cards/HotStockCard.vue'

// 对齐桌面端 market.vue「当前热门」的 7 个子分类。
// 全球/沪深/港股/美股 → HotStock(marketType)（雪球热榜）
// 热门话题 → HotTopic；重大事件 → InvestCalendarTimeLine（韭研公社）；财经日历 → ClsCalendar（财联社）
const STOCK_TYPES = { '10': '全球', '12': '沪深', '13': '港股', '11': '美股' }
const tabs = [
  { label: '全球', value: '10' },
  { label: '沪深', value: '12' },
  { label: '港股', value: '13' },
  { label: '美股', value: '11' },
  { label: '热门话题', value: 'topic' },
  { label: '重大事件', value: 'event' },
  { label: '财经日历', value: 'calendar' },
]

const activeTab = ref('10')
const isStockTab = computed(() => activeTab.value in STOCK_TYPES)

// 各类数据独立缓存
const stockData = ref({}) // { marketType: HotItem[] }
const topicData = ref([])
const eventData = ref([]) // [{ date, list:[{title, like_count}] }]
const calendarData = ref([]) // [{ calendar_day, week, items:[{title, event, economic}] }]
const loading = ref(false)

const today = new Date()
const formattedYM = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}`
const formattedDate = `${formattedYM}-${String(today.getDate()).padStart(2, '0')}`

const currentStocks = computed(() => stockData.value[activeTab.value] || [])
const sortedEventData = computed(() => {
  return [...eventData.value].sort((a, b) => {
    const at = new Date(a.date).getTime()
    const bt = new Date(b.date).getTime()
    return (Number.isFinite(bt) ? bt : 0) - (Number.isFinite(at) ? at : 0)
  })
})

async function loadStock(marketType) {
  try {
    const res = await HotStock(marketType)
    stockData.value = { ...stockData.value, [marketType]: Array.isArray(res) ? res : [] }
  } catch (e) {
    console.error('加载热门股票失败:', e)
  }
}

async function loadTopic() {
  try {
    const res = await HotTopic(20)
    topicData.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error('加载热门话题失败:', e)
  }
}

async function loadEvent() {
  try {
    const res = await InvestCalendarTimeLine(formattedYM)
    eventData.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error('加载重大事件失败:', e)
  }
}

async function loadCalendar() {
  try {
    const res = await ClsCalendar()
    calendarData.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error('加载财经日历失败:', e)
  }
}

// 按当前 Tab 拉取（带 loading），切 Tab 时懒加载，已有缓存则跳过。
async function loadCurrent(force = false) {
  const tab = activeTab.value
  if (tab in STOCK_TYPES) {
    if (!force && (stockData.value[tab]?.length)) return
    loading.value = true
    await loadStock(tab)
  } else if (tab === 'topic') {
    if (!force && topicData.value.length) return
    loading.value = true
    await loadTopic()
  } else if (tab === 'event') {
    if (!force && eventData.value.length) return
    loading.value = true
    await loadEvent()
  } else if (tab === 'calendar') {
    if (!force && calendarData.value.length) return
    loading.value = true
    await loadCalendar()
  }
  loading.value = false
}

async function handleRefresh() {
  await loadCurrent(true)
}

watch(activeTab, () => loadCurrent())

// 热门股票随行情变化快，开市时段 5s 轮询当前股票 Tab（对齐桌面端 HotStockList feed）。
function refreshStockFeed() {
  if (isStockTab.value) loadStock(activeTab.value)
}

onBeforeMount(() => {
  loadCurrent()
  registerFeed('mobile-hot-stock', {
    fetch: refreshStockFeed,
    intervalMs: 5000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-hot-stock')
})

// 话题原文链接（对齐桌面端 gubatopic 链接）
function topicUrl(item) {
  return item.htid ? `https://gubatopic.eastmoney.com/topic_v3.html?htid=${item.htid}` : ''
}
function bigNumber(n) {
  const v = Number(n)
  if (!Number.isFinite(v)) return n ?? '-'
  if (v >= 10000) return (v / 10000).toFixed(1) + '万'
  return v.toLocaleString('en-US')
}
</script>

<template>
  <div class="page">
    <div class="hot-tabs">
      <MTabs v-model="activeTab" :tabs="tabs" />
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading" class="hot-loading">
          <MLoading text="加载中..." vertical />
        </div>

        <!-- 热门股票排行（全球/沪深/港股/美股） -->
        <template v-else-if="isStockTab">
          <div v-if="currentStocks.length" class="hot-list">
            <HotStockCard
              v-for="(item, index) in currentStocks"
              :key="item.code || index"
              :item="item"
              :rank="index + 1"
            />
          </div>
          <MEmpty v-else description="暂无数据，可能为接口失败或休市" />
        </template>

        <!-- 热门话题（股吧） -->
        <template v-else-if="activeTab === 'topic'">
          <div v-if="topicData.length" class="topic-list">
            <a
              v-for="(item, index) in topicData"
              :key="item.htid || index"
              class="topic-item"
              :href="topicUrl(item)"
              target="_blank"
              rel="noopener"
            >
              <img v-if="item.squareImg" class="topic-img" :src="item.squareImg" alt="" />
              <div class="topic-body">
                <div class="topic-title">{{ item.nickname }}</div>
                <div v-if="item.desc" class="topic-desc">{{ item.desc }}</div>
                <div v-if="item.stock_list && item.stock_list.length" class="topic-stocks">
                  <span v-for="(s, i) in item.stock_list" :key="i" class="topic-stock-tag">{{ s.name }}</span>
                </div>
                <div class="topic-meta">
                  <span v-if="item.postNumber != null">💬 {{ bigNumber(item.postNumber) }}</span>
                  <span v-if="item.clickNumber != null">👁 {{ bigNumber(item.clickNumber) }}</span>
                </div>
              </div>
            </a>
          </div>
          <MEmpty v-else description="暂无热门话题" />
        </template>

        <!-- 重大事件时间轴（韭研公社） -->
        <template v-else-if="activeTab === 'event'">
          <div v-if="sortedEventData.length" class="cal-list">
            <div
              v-for="day in sortedEventData"
              :key="day.date"
              class="cal-day"
              :class="{ 'cal-day--today': day.date === formattedDate }"
            >
              <div class="cal-date">{{ day.date }}</div>
              <div class="cal-items">
                <div v-for="(l, i) in day.list" :key="l.article_id || i" class="cal-item">
                  <span class="cal-idx">{{ i + 1 }}</span>
                  <span class="cal-text">{{ l.title }}</span>
                </div>
              </div>
            </div>
          </div>
          <MEmpty v-else description="暂无重大事件" />
        </template>

        <!-- 财经日历（财联社） -->
        <template v-else-if="activeTab === 'calendar'">
          <div v-if="calendarData.length" class="cal-list">
            <div
              v-for="day in calendarData"
              :key="day.calendar_day"
              class="cal-day"
              :class="{ 'cal-day--today': day.calendar_day === formattedDate }"
            >
              <div class="cal-date">{{ day.calendar_day }} {{ day.week }}</div>
              <div class="cal-items">
                <div v-for="(l, i) in day.items" :key="l.id || i" class="cal-item cal-item--block">
                  <div class="cal-item-head">
                    <span class="cal-idx">{{ i + 1 }}</span>
                    <span class="cal-text">{{ l.title }}</span>
                    <span v-if="l.event" class="cal-tag cal-tag--event">事件</span>
                    <span v-if="l.economic" class="cal-tag cal-tag--data">数据</span>
                  </div>
                  <div v-if="l.economic" class="cal-econ">
                    <span>公布 {{ l.economic.actual ?? '-' }}</span>
                    <span>预测 {{ l.economic.consensus ?? '-' }}</span>
                    <span>前值 {{ l.economic.front ?? '-' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <MEmpty v-else description="暂无财经日历" />
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

.hot-loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

/* 热门股票列表 */
.hot-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.hot-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

/* 话题列表 */
.topic-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.topic-item {
  display: flex;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  text-decoration: none;
  color: inherit;
}

.topic-item:active {
  background: var(--m-bg-primary);
}

.topic-img {
  width: 44px;
  height: 44px;
  border-radius: var(--m-radius-sm);
  object-fit: cover;
  flex-shrink: 0;
}

.topic-body {
  flex: 1;
  min-width: 0;
}

.topic-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.topic-desc {
  margin-top: 2px;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.topic-stocks {
  margin-top: var(--m-space-xs);
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
}

.topic-stock-tag {
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
  background: var(--m-bg-primary);
  color: var(--m-text-secondary);
  border: 1px solid var(--m-divider-color);
}

.topic-meta {
  margin-top: var(--m-space-xs);
  display: flex;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

/* 日历 / 时间轴 */
.cal-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.cal-day {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.cal-date {
  padding: var(--m-space-sm) var(--m-space-md);
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-secondary);
  background: var(--m-bg-primary);
}

.cal-day--today .cal-date {
  color: var(--m-color-rise);
}

.cal-items {
  padding: var(--m-space-xs) var(--m-space-md);
}

.cal-item {
  display: flex;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) 0;
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-normal);
}

.cal-item:not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.cal-item--block {
  flex-direction: column;
  gap: var(--m-space-xs);
}

.cal-item-head {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  flex-wrap: wrap;
}

.cal-idx {
  flex-shrink: 0;
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.cal-text {
  flex: 1;
  min-width: 0;
  color: var(--m-text-primary);
}

.cal-tag {
  flex-shrink: 0;
  padding: 0 var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
  line-height: 1.6;
}

.cal-tag--event {
  background: var(--m-color-fall-light);
  color: var(--m-color-fall);
}

.cal-tag--data {
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.cal-econ {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-md);
  padding-left: calc(var(--m-font-sm) + var(--m-space-sm));
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}
</style>
