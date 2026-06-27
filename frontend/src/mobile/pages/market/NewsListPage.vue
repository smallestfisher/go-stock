<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount } from 'vue'
import { GetTelegraphList, ReFleshTelegraphList } from '../../../api/app'
import { EventsOn, EventsOff } from '../../../api/runtime'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import TelegraphCard from '../../components/cards/TelegraphCard.vue'

// 三个资讯源对齐桌面端 market.vue：财联社电报 / 新浪财经 / 外媒。
// 外媒（TradingView）常常无数据，故仅在拉到内容后才显示该 Tab（与桌面端 v-if 一致）。
const SOURCES = [
  { label: '财联社电报', value: '财联社电报', event: 'newTelegraph' },
  { label: '新浪财经', value: '新浪财经', event: 'newSinaNews' },
  { label: '外媒', value: '外媒', event: 'tradingViewNews' },
]

const activeSource = ref('财联社电报')
const loading = ref(true)
// 各源数据独立缓存：{ 来源: 列表 }
const newsBySource = ref({})

const tabs = computed(() =>
  SOURCES
    // 外媒无数据则不展示该 Tab
    .filter(s => s.value !== '外媒' || (newsBySource.value['外媒']?.length))
    .map(s => ({ label: s.label, value: s.value }))
)

const currentList = computed(() => newsBySource.value[activeSource.value] || [])

// 按日期分组：返回 [{ date, items }]，每组内按时间倒序（后端已 data_time desc）。
const groupedList = computed(() => {
  const groups = []
  let last = null
  for (const item of currentList.value) {
    const date = dateOf(item)
    if (!last || last.date !== date) {
      last = { date, items: [] }
      groups.push(last)
    }
    last.items.push(item)
  }
  return groups
})

function dateOf(item) {
  // dataTime 形如 "2026-06-27T10:20:00+08:00"，取日期部分；缺失则用今天。
  const raw = item.dataTime || ''
  if (typeof raw === 'string' && raw.length >= 10) return raw.slice(0, 10)
  return new Date().toISOString().slice(0, 10)
}

function itemKey(item) {
  return item.ID ?? item.id ?? `${item.time}-${item.title || item.content || ''}`.slice(0, 64)
}

async function loadSource(source) {
  try {
    const res = await GetTelegraphList(source)
    newsBySource.value = {
      ...newsBySource.value,
      [source]: Array.isArray(res) ? res : [],
    }
  } catch (e) {
    console.error(`加载${source}快讯失败:`, e)
  }
}

// 首屏：三源一起拉，外媒 Tab 是否出现取决于其是否有数据。
async function loadAll() {
  await Promise.all(SOURCES.map(s => loadSource(s.value)))
  loading.value = false
}

// 下拉刷新：触发一次后端三源抓取，再读取当前源缓存。
async function handleRefresh() {
  try {
    await ReFleshTelegraphList('')
  } catch (e) {
    // 抓取失败不打断，仍读现有缓存
  }
  await Promise.all(SOURCES.map(s => loadSource(s.value)))
}

// 实时推送：新条目插到对应源头部，并裁掉等量尾部，保持长度稳定（对齐桌面端）。
function makeHandler(source) {
  return (data) => {
    if (!Array.isArray(data) || !data.length) return
    const prev = newsBySource.value[source] || []
    const merged = [...data, ...prev]
    // 去重（按 key），保留先到的（即新推送的）
    const seen = new Set()
    const deduped = []
    for (const it of merged) {
      const k = itemKey(it)
      if (seen.has(k)) continue
      seen.add(k)
      deduped.push(it)
    }
    newsBySource.value = {
      ...newsBySource.value,
      [source]: deduped.slice(0, 50),
    }
  }
}

const handlers = {}

onBeforeMount(() => {
  loadAll()

  // 注册实时事件
  for (const s of SOURCES) {
    handlers[s.value] = makeHandler(s.value)
    EventsOn(s.event, handlers[s.value])
  }

  // 开市时段轮询后端抓取（与桌面端 market.industry feed 同 10s 节奏）
  registerFeed('mobile-market-news', {
    fetch: () => ReFleshTelegraphList(''),
    intervalMs: 10000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-market-news')
  for (const s of SOURCES) {
    EventsOff(s.event)
  }
})
</script>

<template>
  <div class="news-list-page">
    <!-- 资讯源切换 -->
    <MTabs v-model="activeSource" :tabs="tabs" />

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="news-container">
        <div v-if="loading && !currentList.length" class="news-loading">
          <MLoading text="加载快讯..." vertical />
        </div>

        <template v-else-if="currentList.length">
          <div
            v-for="group in groupedList"
            :key="group.date"
            class="news-group"
          >
            <div class="news-date">{{ group.date }}</div>
            <div class="news-items">
              <TelegraphCard
                v-for="item in group.items"
                :key="itemKey(item)"
                :item="item"
              />
            </div>
          </div>
        </template>

        <MEmpty v-else description="暂无快讯" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.news-list-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.news-container {
  padding: var(--m-space-md);
}

.news-loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

.news-group {
  margin-bottom: var(--m-space-md);
}

.news-date {
  position: sticky;
  top: 0;
  z-index: 1;
  padding: var(--m-space-xs) var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  background: var(--m-bg-primary);
}

.news-items {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

/* 卡片之间的细分隔线（最后一条不画） */
.news-items > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}
</style>
