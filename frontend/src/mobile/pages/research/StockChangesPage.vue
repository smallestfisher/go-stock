<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount, watch } from 'vue'
import { GetStockChanges, GetStockChangeHistory } from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import { formatMoney } from '../../composables/useFormat'

// 对齐桌面端 stockChangesMonitor.vue：实时 GetStockChanges(types,0,n) / 历史 GetStockChangeHistory。
// 异动类型 22 种分利好/利空（注意：空壳 typeTabs 的占位全错，已重做）。
const modeTabs = [
  { label: '实时', value: 'realtime' },
  { label: '历史', value: 'history' },
]
const activeMode = ref('realtime')

// 22 种类型（对齐桌面 bullishTypes/bearishTypes）
const BULLISH = [
  { label: '火箭发射', value: '8201' },
  { label: '快速反弹', value: '8202' },
  { label: '大笔买入', value: '8193' },
  { label: '封涨停板', value: '4' },
  { label: '打开跌停板', value: '32' },
  { label: '有大买盘', value: '64' },
  { label: '竞价上涨', value: '8207' },
  { label: '高开5日线', value: '8209' },
  { label: '向上缺口', value: '8211' },
  { label: '60日新高', value: '8213' },
  { label: '60日大幅上涨', value: '8215' },
  { label: '打开涨停板', value: '16' },
]
const BEARISH = [
  { label: '加速下跌', value: '8204' },
  { label: '高台跳水', value: '8203' },
  { label: '大笔卖出', value: '8194' },
  { label: '封跌停板', value: '8' },
  { label: '有大卖盘', value: '128' },
  { label: '竞价下跌', value: '8208' },
  { label: '低开5日线', value: '8210' },
  { label: '向下缺口', value: '8212' },
  { label: '60日新低', value: '8214' },
  { label: '60日大幅下跌', value: '8216' },
]
const ALL_TYPES = [...BULLISH, ...BEARISH]
const BULLISH_NAMES = new Set(BULLISH.map(t => t.label))
const BEARISH_NAMES = new Set(BEARISH.map(t => t.label))

const selectedTypes = ref(ALL_TYPES.map(t => t.value)) // 默认全选

const list = ref([])
const loading = ref(false)

function pickList(res) {
  return (res && Array.isArray(res.list) ? res.list : []).filter(it => it && typeof it === 'object')
}
function pickData(res) {
  return (res && Array.isArray(res.data) ? res.data : []).filter(it => it && typeof it === 'object')
}

// 实时异动
async function fetchRealtime() {
  loading.value = true
  try {
    const types = selectedTypes.value.map(t => Number(t))
    const res = await GetStockChanges(types, 0, 80)
    list.value = pickData(res)
  } catch (e) {
    console.error('加载实时异动失败:', e)
    list.value = []
  } finally {
    loading.value = false
  }
}

// 历史异动
async function fetchHistory() {
  loading.value = true
  try {
    const res = await GetStockChangeHistory({
      page: 1,
      pageSize: 80,
      changeTypes: selectedTypes.value.map(t => Number(t)),
    })
    list.value = pickList(res)
  } catch (e) {
    console.error('加载历史异动失败:', e)
    list.value = []
  } finally {
    loading.value = false
  }
}

async function loadCurrent(force = false) {
  if (!force && list.value.length) return
  if (activeMode.value === 'realtime') await fetchRealtime()
  else await fetchHistory()
}

// 快捷选择
function selectAll() { selectedTypes.value = ALL_TYPES.map(t => t.value) }
function selectBullish() { selectedTypes.value = BULLISH.map(t => t.value) }
function selectBearish() { selectedTypes.value = BEARISH.map(t => t.value) }
function toggleType(value) {
  const set = new Set(selectedTypes.value)
  if (set.has(value)) set.delete(value)
  else set.add(value)
  selectedTypes.value = ALL_TYPES.map(t => t.value).filter(v => set.has(v))
}

// 切换模式 / 类型变化 → 重新拉取
watch(activeMode, () => { list.value = []; loadCurrent(true) })
watch(selectedTypes, () => loadCurrent(true), { deep: true })

async function handleRefresh() {
  await loadCurrent(true)
}

// 字段双 key 兜底 getter（对齐桌面 getChange*）
function getCode(it) { return it.stockCode || it.StockCode || it.code || '' }
function getName(it) { return it.stockName || it.StockName || it.name || '-' }
function getTypeName(it) { return it.typeName || it.TypeName || '-' }
function getRate(it) {
  const r = Number(it.changeRate || it.ChangeRate)
  return Number.isFinite(r) ? r : NaN
}
function getDateTime(it) {
  const d = it.changeDate || it.ChangeDate
  const t = it.changeTime || it.ChangeTime || it.time || ''
  return d ? `${d} ${t}` : t
}
function getPrice(it) {
  const p = Number(it.price || it.Price)
  return Number.isFinite(p) && p > 0 ? p.toFixed(2) : '--'
}
// 类型着色：利好红 / 利空绿（对齐桌面 getChangeTypeTag）
function typeLevel(name) {
  if (BULLISH_NAMES.has(name)) return 'bull'
  if (BEARISH_NAMES.has(name)) return 'bear'
  return 'flat'
}

onBeforeMount(() => {
  loadCurrent(true)
  registerFeed('mobile-stock-changes', {
    fetch: () => { if (activeMode.value === 'realtime') fetchRealtime() },
    intervalMs: 10 * 1000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-stock-changes')
})
</script>

<template>
  <div class="page">
    <!-- 模式切换 -->
    <div class="mode-bar">
      <MTabs v-model="activeMode" :tabs="modeTabs" />
    </div>

    <!-- 类型筛选：快捷 + 利好/利空 chip -->
    <div class="type-filter">
      <div class="quick-row">
        <button class="quick-btn" type="button" @click="selectAll">全部</button>
        <button class="quick-btn quick-btn--bull" type="button" @click="selectBullish">利好</button>
        <button class="quick-btn quick-btn--bear" type="button" @click="selectBearish">利空</button>
      </div>
      <div class="chip-row">
        <button
          v-for="t in BULLISH"
          :key="'b' + t.value"
          type="button"
          class="type-chip type-chip--bull"
          :class="{ 'type-chip--on': selectedTypes.includes(t.value) }"
          @click="toggleType(t.value)"
        >
          {{ t.label }}
        </button>
      </div>
      <div class="chip-row">
        <button
          v-for="t in BEARISH"
          :key="'s' + t.value"
          type="button"
          class="type-chip type-chip--bear"
          :class="{ 'type-chip--on': selectedTypes.includes(t.value) }"
          @click="toggleType(t.value)"
        >
          {{ t.label }}
        </button>
      </div>
    </div>

    <MPullRefresh class="change-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !list.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else>
          <div v-if="list.length" class="change-list">
            <div v-for="(it, i) in list" :key="(getCode(it)) + i" class="change-item">
              <!-- 第一行：名称+代码 / 涨跌幅 -->
              <div class="row-top">
                <div class="stock-left">
                  <span class="stock-name">{{ getName(it) }}</span>
                  <span class="stock-code">{{ getCode(it) }}</span>
                </div>
                <PercentTag :value="getRate(it)" size="small" />
              </div>
              <!-- 第二行：异动类型标签 + 时间 + 价格 -->
              <div class="row-sub">
                <span class="type-tag" :class="`type-tag--${typeLevel(getTypeName(it))}`">
                  {{ getTypeName(it) }}
                </span>
                <span class="time">{{ getDateTime(it) }}</span>
                <span class="price">{{ getPrice(it) }}</span>
              </div>
            </div>
          </div>
          <MEmpty v-else description="暂无异动数据" />
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

.mode-bar {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.type-filter {
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.quick-row {
  display: flex;
  gap: var(--m-space-xs);
}

.quick-btn {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.quick-btn--bull {
  color: var(--m-color-rise);
  border-color: var(--m-color-rise);
}

.quick-btn--bear {
  color: var(--m-color-fall);
  border-color: var(--m-color-fall);
}

.chip-row {
  display: flex;
  gap: var(--m-space-xs);
  overflow-x: auto;
  scrollbar-width: none;
}

.chip-row::-webkit-scrollbar {
  display: none;
}

.type-chip {
  flex-shrink: 0;
  padding: 2px var(--m-space-sm);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  white-space: nowrap;
}

.type-chip--bull.type-chip--on {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.type-chip--bear.type-chip--on {
  background: var(--m-color-fall-light);
  border-color: var(--m-color-fall);
  color: var(--m-color-fall);
}

.change-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

.change-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.change-item {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
}

.row-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-md);
}

.stock-left {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-xs);
  min-width: 0;
}

.stock-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
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

.row-sub {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.type-tag {
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
}

.type-tag--bull {
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.type-tag--bear {
  background: var(--m-color-fall-light);
  color: var(--m-color-fall);
}

.type-tag--flat {
  background: var(--m-bg-primary);
  color: var(--m-text-secondary);
}

.time {
  font-variant-numeric: tabular-nums;
}

.price {
  margin-left: auto;
  color: var(--m-text-secondary);
  font-variant-numeric: tabular-nums;
}
</style>
