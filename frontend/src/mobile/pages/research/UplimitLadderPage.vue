<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount } from 'vue'
import { GetUplimitHot, IsTradingDay, GetLatestTradingDay } from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'

// 对齐桌面端 uplimitLadder.vue 的「涨停高度」ladder 视图：GetUplimitHot(date,20) 按连板次数分组。
const MAX_BACKTRACK = 7
const TOP_N = 20

const todayStr = (() => {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
})()

const rawData = ref(null)
const loading = ref(false)
const selectedDate = ref(todayStr)

// 偏移日期
function shiftDate(dateStr, days) {
  const d = new Date(dateStr)
  d.setDate(d.getDate() + days)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 取数（带空数据回退）
async function fetchData(date, back = 0) {
  loading.value = true
  try {
    const res = await GetUplimitHot(date, TOP_N)
    if (res && res.code === 20000) {
      const data = res.data
      const hasData = data && data.plate?.length > 0 && data.stocks && data.stocks.trim() !== ''
      if (hasData) {
        rawData.value = data
        selectedDate.value = date
        return
      }
      if (back < MAX_BACKTRACK) {
        return fetchData(shiftDate(date, -1), back + 1)
      }
      rawData.value = null
      selectedDate.value = date
      return
    }
    rawData.value = null
    selectedDate.value = date
  } catch (e) {
    console.error('加载涨停梯队失败:', e)
    rawData.value = null
  } finally {
    loading.value = false
  }
}

// 初始日期：交易日取今天，否则取最近交易日
async function initDate() {
  try {
    const isToday = await IsTradingDay(todayStr)
    selectedDate.value = isToday ? todayStr : (await GetLatestTradingDay() || todayStr)
  } catch {
    selectedDate.value = todayStr
  }
  fetchData(selectedDate.value)
}

function changeDate(delta) {
  selectedDate.value = shiftDate(selectedDate.value, delta)
  fetchData(selectedDate.value)
}

async function handleRefresh() {
  await fetchData(selectedDate.value)
}

// ===== 派生数据（对齐桌面 uplimitLadder.vue:150-215） =====
const maxCount = computed(() => rawData.value?.max_count || 0)

const totalZtCount = computed(() => {
  if (!rawData.value?.stocks) return 0
  return rawData.value.stocks.split(',').filter(s => s.trim()).length
})

// 按连板次数分组
const ladderData = computed(() => {
  if (!rawData.value?.plate_stocks) return {}
  const allCodes = (rawData.value.stocks || '').split(',').filter(s => s.trim())
  const stockInfo = rawData.value.stock_info || {}
  const ladder = {}
  for (let i = maxCount.value; i >= 1; i--) ladder[i] = []

  const seen = new Set()
  for (const code of allCodes) {
    if (seen.has(code)) continue
    seen.add(code)
    let stockData = null
    for (const plateCode of Object.keys(rawData.value.plate_stocks)) {
      const found = rawData.value.plate_stocks[plateCode].find(s => s.stock_code === code)
      if (found) { stockData = found; break }
    }
    if (!stockData) continue
    const keepTimes = stockData.up_limit_keep_times || 0
    if (keepTimes >= 1 && ladder[keepTimes]) {
      ladder[keepTimes].push({ ...stockData, plates: stockInfo[code]?.plates || [] })
    }
  }
  for (const key of Object.keys(ladder)) {
    ladder[key].sort((a, b) =>
      (b.up_limit_keep_times - a.up_limit_keep_times) ||
      (a.up_limit_time || '').localeCompare(b.up_limit_time || '')
    )
  }
  return ladder
})

// 分组列表（降序）：[{ level, count, list }]
const ladderGroups = computed(() => {
  const groups = []
  const banInfo = rawData.value?.ban_info || {}
  for (let i = maxCount.value; i >= 1; i--) {
    const list = ladderData.value[i] || []
    if (!list.length) continue
    groups.push({
      level: i,
      count: banInfo[String(i)]?.count || list.length,
      list,
    })
  }
  return groups
})

// 涨停类型标签/颜色（对齐桌面 getTypeLabel/getTypeColor）
function typeLabel(type) {
  if (!type) return ''
  if (type === '一') return '一字'
  if (type === 'T') return 'T字'
  if (type === '自') return '自然'
  if (type.startsWith('烂')) return '烂' + type.slice(1) + '板'
  if (type === '炸') return '炸板'
  return type
}
function typeColor(type) {
  if (!type) return ''
  if (type === '一') return '#e03030'
  if (type === 'T' || type.startsWith('烂')) return '#f0a020'
  if (type === '自') return '#2080f0'
  if (type === '炸') return '#999'
  return ''
}

function levelLabel(level) {
  return level === 1 ? '首板' : `${level}连板`
}

onBeforeMount(() => {
  initDate()
  // 仅当天+开市 60s 轮询
  registerFeed('mobile-uplimit-ladder', {
    fetch: () => { if (selectedDate.value === todayStr) fetchData(selectedDate.value) },
    intervalMs: 60 * 1000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-uplimit-ladder')
})
</script>

<template>
  <div class="page">
    <!-- 日期切换条 -->
    <div class="date-bar">
      <button class="date-btn" type="button" @click="changeDate(-1)">‹</button>
      <span class="date-text">{{ selectedDate }}</span>
      <button class="date-btn" type="button" @click="changeDate(1)">›</button>
    </div>

    <MPullRefresh class="ladder-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !rawData" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else-if="rawData">
          <!-- 头部统计 -->
          <div class="stat-bar">
            <span class="stat-item">涨停 <b class="stat-num">{{ totalZtCount }}</b> 只</span>
            <span class="stat-item">最高 <b class="stat-num stat-num--rise">{{ maxCount }}</b> 连板</span>
          </div>

          <!-- 连板梯队分组 -->
          <div v-if="ladderGroups.length" class="group-list">
            <div v-for="g in ladderGroups" :key="g.level" class="group">
              <div class="group-title">
                <span class="group-level">{{ levelLabel(g.level) }}</span>
                <span class="group-count">{{ g.count }} 只</span>
              </div>
              <div class="group-stocks">
                <div v-for="(s, i) in g.list" :key="s.stock_code || i" class="stock-item">
                  <div class="stock-left">
                    <span class="stock-name">{{ s.stock_name }}</span>
                    <span class="stock-code">{{ s.stock_code }}</span>
                  </div>
                  <div class="stock-mid">
                    <span
                      v-if="s.up_limit_type"
                      class="type-tag"
                      :style="{ color: typeColor(s.up_limit_type) }"
                    >
                      {{ typeLabel(s.up_limit_type) }}
                    </span>
                    <span v-if="s.up_limit_time" class="stock-time">{{ s.up_limit_time }}</span>
                  </div>
                  <div class="stock-right">
                    <span v-if="s.amount" class="stock-amt">{{ s.amount }}亿</span>
                    <span v-if="s.fd_close != null" class="stock-fd">封{{ s.fd_close }}%</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <MEmpty v-else description="暂无连板数据" />
        </template>
        <MEmpty v-else description="暂无涨停数据" />
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

.date-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--m-space-xl);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.date-btn {
  width: 36px;
  height: 36px;
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-lg);
  line-height: 1;
}

.date-btn:active {
  transform: scale(0.92);
}

.date-text {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
  min-width: 110px;
  text-align: center;
}

.ladder-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

.stat-bar {
  display: flex;
  gap: var(--m-space-xl);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.stat-num {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
}

.stat-num--rise {
  color: var(--m-color-rise);
}

/* 分组 */
.group-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.group {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.group-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-primary);
}

.group-level {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-color-rise);
}

.group-count {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.group-stocks > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.stock-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
}

.stock-left {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-xs);
  flex: 1;
  min-width: 0;
}

.stock-name {
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stock-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  flex-shrink: 0;
}

.stock-mid {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  flex-shrink: 0;
}

.type-tag {
  font-size: var(--m-font-xs);
  font-weight: var(--m-font-weight-medium);
}

.stock-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.stock-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1px;
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.stock-amt {
  color: var(--m-text-secondary);
}
</style>
