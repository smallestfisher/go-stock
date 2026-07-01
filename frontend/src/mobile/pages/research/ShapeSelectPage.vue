<script setup>
import { ref, reactive, computed, onBeforeMount } from 'vue'
import { GetAllStocks } from '../../../api/app'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MSheet from '../../components/base/MSheet.vue'
import MButton from '../../components/base/MButton.vue'
import StockDetailSheet from '../../components/sheets/StockDetailSheet.vue'
import MIcon from '../../components/base/MIcon.vue'
import { toast } from '../../composables/useToast'

// 对齐桌面端 allStockList.vue（研究中心→形态选股）：
// 全部股票分页列表 + 关键词搜索 + 形态/技术指标筛选。
// GetAllStocks(page, pageSize, keyword, technicalIndicator) -> { result: { data, count } }
// technicalIndicator 字段对齐桌面 technicalIndicatorReactive 与后端 models.TechnicalIndicators：
//   形态(布尔) + 人气/关注/连涨/连跌(数值)，false/0 表示不筛选。
// 字段：SECUCODE / SECURITY_NAME_ABBR / NEW_PRICE / CHANGE_RATE / HIGH_PRICE / LOW_PRICE / VOLUME

const PAGE_SIZE = 20

const keyword = ref('')
const list = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / PAGE_SIZE) || 1)
const hasMore = computed(() => page.value < totalPages.value)

// ===== 形态/技术指标筛选（对齐桌面 allStockList.vue 的 technicalIndicatorReactive）=====
// 形态：布尔，多选；数值：单选(0=不筛选)。key 必须与后端 models.TechnicalIndicators 的 json tag 完全一致。
const FILTER_GROUPS = [
  {
    title: '金叉 & 资金',
    items: [
      { key: 'MACD_GOLDEN_FORK', label: 'MACD金叉' },
      { key: 'KDJ_GOLDEN_FORK', label: 'KDJ金叉' },
      { key: 'LOW_FUNDS_INFLOW', label: '低位资金净流入' },
      { key: 'HIGH_FUNDS_OUTFLOW', label: '高位资金净流出' },
    ],
  },
  {
    title: '均线',
    items: [
      { key: 'BREAKUP_MA_5DAYS', label: '向上突破5日均线' },
      { key: 'LONG_AVG_ARRAY', label: '均线多头排列' },
      { key: 'SHORT_AVG_ARRAY', label: '均线空头排列' },
    ],
  },
  {
    title: '量能',
    items: [
      { key: 'BREAK_THROUGH', label: '放量突破' },
      { key: 'UPSIDE_VOLUME', label: '放量上攻' },
      { key: 'UPPER_LARGE_VOLUME', label: '连涨放量' },
      { key: 'DOWN_NARROW_VOLUME', label: '下跌无量' },
      { key: 'HEAVEN_RULE', label: '天量法则' },
    ],
  },
  {
    title: '大阳线 / 连阳',
    items: [
      { key: 'ONE_DAYANG_LINE', label: '一根大阳线' },
      { key: 'TWO_DAYANG_LINES', label: '两根大阳线' },
      { key: 'RISE_SUN', label: '旭日东升' },
      { key: 'POWER_FULGUN', label: '强势多方炮' },
      { key: 'RESTORE_JUSTICE', label: '拨云见日' },
      { key: 'UPPER_4DAYS', label: '四串阳' },
      { key: 'UPPER_8DAYS', label: '八仙过海(八连阳)' },
      { key: 'UPPER_9DAYS', label: '九阳神功(九连阳)' },
    ],
  },
  {
    title: '阴线 / 反转 / 整理',
    items: [
      { key: 'DOWN_7DAYS', label: '七仙女下凡(七连阴)' },
      { key: 'BEARISH_ENGULFING', label: '穿头破脚' },
      { key: 'REVERSING_HAMMER', label: '倒转锤头' },
      { key: 'SHOOTING_STAR', label: '射击之星' },
      { key: 'EVENING_STAR', label: '黄昏之星' },
      { key: 'FIRST_DAWN', label: '曙光初现' },
      { key: 'PREGNANT', label: '身怀六甲' },
      { key: 'BLACK_CLOUD_TOPS', label: '乌云盖顶' },
      { key: 'MORNING_STAR', label: '早晨之星' },
      { key: 'NARROW_FINISH', label: '窄幅整理' },
    ],
  },
]

// 数值型单选：key 对齐后端，value=0 表示不筛选
const RADIO_GROUPS = [
  {
    key: 'UPP_DAYS', title: '人气排名连涨', options: [
      { value: 3, label: '3天及以上' },
      { value: 5, label: '5天及以上' },
      { value: 7, label: '7天及以上' },
    ],
  },
  {
    key: 'CONCERN_RANK_7DAYS', title: '7日关注排名', options: [
      { value: 10, label: '前10名' },
      { value: 50, label: '前50名' },
      { value: 100, label: '前100名' },
    ],
  },
  {
    key: 'UPNDAY', title: '连涨天数', options: [
      { value: 3, label: '3天及以上' },
      { value: 5, label: '5天及以上' },
      { value: 8, label: '8天及以上' },
    ],
  },
  {
    key: 'DOWNNDAY', title: '连跌天数', options: [
      { value: 3, label: '3天及以上' },
      { value: 5, label: '5天及以上' },
      { value: 8, label: '8天及以上' },
      { value: 10, label: '10天及以上' },
      { value: 14, label: '14天及以上' },
    ],
  },
]

function makeDefaults() {
  const o = {}
  for (const g of FILTER_GROUPS) for (const it of g.items) o[it.key] = false
  for (const g of RADIO_GROUPS) o[g.key] = 0
  return o
}

// 实际生效的筛选条件（fetchPage 读取它）
const technicalIndicator = reactive(makeDefaults())
// 抽屉内的编辑草稿：确定时才写回 technicalIndicator，取消则丢弃
const draft = reactive(makeDefaults())

const activeFilterCount = computed(() => {
  let n = 0
  for (const g of FILTER_GROUPS) for (const it of g.items) if (technicalIndicator[it.key]) n++
  for (const g of RADIO_GROUPS) if (technicalIndicator[g.key]) n++
  return n
})

const filterVisible = ref(false)

function openFilter() {
  Object.assign(draft, technicalIndicator)
  filterVisible.value = true
}

function resetDraft() {
  Object.assign(draft, makeDefaults())
}

function toggleRadio(key, value) {
  draft[key] = draft[key] === value ? 0 : value
}

function applyFilter() {
  Object.assign(technicalIndicator, draft)
  filterVisible.value = false
  fetchPage(1, false)
}

function clearAll() {
  Object.assign(technicalIndicator, makeDefaults())
  fetchPage(1, false)
}

async function fetchPage(targetPage, append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const res = await GetAllStocks(targetPage, PAGE_SIZE, keyword.value, technicalIndicator)
    if (res && res.result && Array.isArray(res.result.data)) {
      const rows = res.result.data
      list.value = append ? list.value.concat(rows) : rows
      total.value = res.result.count || 0
      page.value = targetPage
    } else if (!append) {
      list.value = []
      total.value = 0
    }
  } catch (e) {
    console.error('加载股票列表失败:', e)
    toast.error('加载股票列表失败')
  } finally {
    loading.value = false
  }
}

function handleRefresh() {
  return fetchPage(1, false)
}

function loadMore() {
  if (hasMore.value && !loading.value) fetchPage(page.value + 1, true)
}

function applySearch() {
  fetchPage(1, false)
}

// 东财代码：SECUCODE 形如 600519.SH -> sh600519，供 StockDetailSheet 使用
function toEastMoneyCode(secucode) {
  if (!secucode) return ''
  const c = String(secucode).trim().toUpperCase()
  if (c.endsWith('.SH')) return 'sh' + c.slice(0, -3).toLowerCase()
  if (c.endsWith('.SZ')) return 'sz' + c.slice(0, -3).toLowerCase()
  if (c.endsWith('.BJ')) return 'bj' + c.slice(0, -3).toLowerCase()
  if (c.endsWith('.HK')) return 'hk' + c.slice(0, -3).toLowerCase()
  if (c.startsWith('6')) return 'sh' + c.toLowerCase()
  if (c.startsWith('0') || c.startsWith('3')) return 'sz' + c.toLowerCase()
  if (c.startsWith('8') || c.startsWith('9')) return 'bj' + c.toLowerCase()
  return c.toLowerCase()
}

function isNum(v) {
  return v !== null && v !== undefined && v !== '' && !isNaN(Number(v))
}

function fmtPrice(v) {
  return isNum(v) ? Number(v).toFixed(2) : '-'
}

function fmtRate(v) {
  if (!isNum(v)) return '-'
  const n = Number(v)
  return `${n >= 0 ? '+' : ''}${n.toFixed(2)}%`
}

function rateClass(v) {
  if (!isNum(v)) return 'flat'
  const n = Number(v)
  if (n > 0) return 'rise'
  if (n < 0) return 'fall'
  return 'flat'
}

function fmtVolume(v) {
  const n = Number(v) || 0
  if (!n) return '-'
  if (n >= 100000000) return `${(n / 100000000).toFixed(2)}亿`
  if (n >= 10000) return `${(n / 10000).toFixed(2)}万`
  return String(n)
}

// ===== 点击进入个股详情 =====
const detailVisible = ref(false)
const selectedStock = ref(null)

function openDetail(row) {
  selectedStock.value = {
    code: toEastMoneyCode(row.SECUCODE),
    name: row.SECURITY_NAME_ABBR || '',
  }
  detailVisible.value = true
}

onBeforeMount(() => {
  fetchPage(1, false)
})
</script>

<template>
  <div class="page">
    <div class="op-bar">
      <div class="search-row">
        <input
          v-model="keyword"
          class="search-input"
          type="text"
          placeholder="输入股票代码/名称筛选..."
          @keyup.enter="applySearch"
        />
        <button type="button" class="search-btn" @click="applySearch">搜索</button>
      </div>
      <div class="filter-row">
        <button
          type="button"
          class="filter-btn"
          :class="{ 'filter-btn--active': activeFilterCount > 0 }"
          @click="openFilter"
        >
          <span><MIcon name="target" :size="16" /> 形态筛选</span>
          <span v-if="activeFilterCount > 0" class="badge">{{ activeFilterCount }}</span>
        </button>
        <button v-if="activeFilterCount > 0" type="button" class="clear-link" @click="clearAll">清除</button>
        <span class="total">共 {{ total }} 只</span>
      </div>
    </div>

    <MPullRefresh class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="list.length" class="stock-list">
          <div
            v-for="row in list"
            :key="row.SECUCODE"
            class="stock-row"
            @click="openDetail(row)"
          >
            <div class="stock-main">
              <span class="stock-name">{{ row.SECURITY_NAME_ABBR }}</span>
              <span class="stock-code">{{ row.SECUCODE }}</span>
            </div>
            <div class="stock-quote">
              <span class="stock-price" :class="rateClass(row.CHANGE_RATE)">{{ fmtPrice(row.NEW_PRICE) }}</span>
              <span class="stock-extra">高 {{ fmtPrice(row.HIGH_PRICE) }} · 低 {{ fmtPrice(row.LOW_PRICE) }} · 量 {{ fmtVolume(row.VOLUME) }}</span>
            </div>
            <span class="stock-rate" :class="rateClass(row.CHANGE_RATE)">{{ fmtRate(row.CHANGE_RATE) }}</span>
          </div>

          <div class="load-more">
            <button v-if="hasMore" type="button" class="load-more-btn" :disabled="loading" @click="loadMore">
              {{ loading ? '加载中...' : '加载更多' }}
            </button>
            <span v-else class="load-more-end">共 {{ total }} 只</span>
          </div>
        </div>

        <MLoading v-else-if="loading" text="加载中..." vertical class="page-loading" />
        <MEmpty v-else description="暂无股票数据" />
      </div>
    </MPullRefresh>

    <!-- 形态筛选抽屉 -->
    <MSheet v-model:show="filterVisible" title="形态筛选" height="85vh">
      <div class="filter-sheet">
        <section v-for="g in FILTER_GROUPS" :key="g.title" class="filter-group">
          <h5 class="group-title">{{ g.title }}</h5>
          <div class="chip-wrap">
            <button
              v-for="it in g.items"
              :key="it.key"
              type="button"
              class="chip"
              :class="{ 'chip--active': draft[it.key] }"
              @click="draft[it.key] = !draft[it.key]"
            >{{ it.label }}</button>
          </div>
        </section>

        <section v-for="g in RADIO_GROUPS" :key="g.key" class="filter-group">
          <h5 class="group-title">{{ g.title }}</h5>
          <div class="chip-wrap">
            <button
              v-for="opt in g.options"
              :key="opt.value"
              type="button"
              class="chip"
              :class="{ 'chip--active': draft[g.key] === opt.value }"
              @click="toggleRadio(g.key, opt.value)"
            >{{ opt.label }}</button>
          </div>
        </section>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="resetDraft">重置</MButton>
          <MButton type="primary" block @click="applyFilter">确定</MButton>
        </div>
      </template>
    </MSheet>

    <StockDetailSheet v-model:show="detailVisible" :stock="selectedStock" />
  </div>
</template>

<style scoped>
.page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.op-bar {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.search-row {
  display: flex;
  gap: var(--m-space-sm);
}

.search-input {
  flex: 1;
  min-width: 0;
  padding: var(--m-space-sm) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font: inherit;
  font-size: var(--m-font-sm);
  outline: none;
}

.search-input:focus {
  border-color: var(--m-color-rise);
}

.search-btn {
  flex-shrink: 0;
  padding: var(--m-space-sm) var(--m-space-lg);
  border: none;
  border-radius: var(--m-radius-full);
  background: var(--m-color-rise);
  color: #fff;
  font: inherit;
  font-size: var(--m-font-sm);
}

.filter-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.filter-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.filter-btn--active {
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
  background: var(--m-color-rise-light);
}

.badge {
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: var(--m-color-rise);
  color: #fff;
  border-radius: var(--m-radius-full);
  font-size: 10px;
  line-height: 1;
}

.clear-link {
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-color-rise);
}

.total {
  margin-left: auto;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
}

.page-loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

.stock-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.stock-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  box-shadow: var(--m-shadow-sm);
}

.stock-row:active {
  background: var(--m-bg-primary);
}

.stock-main {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
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
  font-variant-numeric: tabular-nums;
}

.stock-quote {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  flex-shrink: 0;
}

.stock-price {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
}

.stock-extra {
  font-size: 10px;
  color: var(--m-text-tertiary);
  white-space: nowrap;
}

.stock-rate {
  flex-shrink: 0;
  min-width: 64px;
  text-align: right;
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
}

.rise { color: var(--m-color-rise); }
.fall { color: var(--m-color-fall); }
.flat { color: var(--m-text-secondary); }

.load-more {
  display: flex;
  justify-content: center;
  padding: var(--m-space-md) 0 var(--m-space-xs);
}

.load-more-btn {
  padding: var(--m-space-xs) var(--m-space-xl);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-card);
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.load-more-end {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

/* ===== 筛选抽屉 ===== */
.filter-sheet {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.group-title {
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-secondary);
}

.chip-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-sm);
}

.chip {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
}

.chip:active {
  opacity: 0.7;
}

.chip--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.sheet-footer {
  display: flex;
  gap: var(--m-space-md);
}
</style>
