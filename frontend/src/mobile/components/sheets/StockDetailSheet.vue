<script setup>
import { ref, computed, watch } from 'vue'
import MSheet from '../base/MSheet.vue'
import MTabs from '../base/MTabs.vue'
import MButton from '../base/MButton.vue'
import PriceTag from '../widgets/PriceTag.vue'
import PercentTag from '../widgets/PercentTag.vue'
import FenshiChart from '../charts/FenshiChart.vue'
import FullKlineChart from '../charts/FullKlineChart.vue'
import MoneyTrendChart from '../charts/MoneyTrendChart.vue'
import { toast } from '../../composables/useToast'
import {
  Greet,
  GetStockMinutePriceLineData,
  GetStockKLine,
  GetStockMoneyTrendByDay,
  SetCostPriceAndVolume,
  SetAlarmChangePercent,
  SetTradingPrice,
  SetStockSort
} from '../../../api/app'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  stock: {
    type: Object,
    default: () => null
  }
})

const emit = defineEmits(['update:show', 'saved'])

// 详情标签页
const detailTabs = [
  { label: '分时', value: 'fenshi' },
  { label: '日K', value: 'kline' },
  { label: '盘口', value: 'info' },
  { label: '资金', value: 'fund' },
]
const activeTab = ref('fenshi')

// ===== 分时数据 =====
const fenshiData = ref([])        // priceData
const fenshiDate = ref('')        // 交易日期
const fenshiPreClose = ref(0)     // 昨收（基准线）
const fenshiLoading = ref(false)
let fenshiLoadedFor = ''         // 已加载该 code 的分时，避免重复拉取

async function loadFenshi() {
  const s = props.stock
  if (!s || !s.code) return
  if (fenshiLoadedFor === s.code && fenshiData.value.length) return
  fenshiLoading.value = true
  try {
    // 拿昨收基准线（Greet 的「昨日收盘价」）
    const [quote, minute] = await Promise.all([
      Greet(s.code).catch(() => null),
      GetStockMinutePriceLineData(s.code, s.name || '').catch(() => null)
    ])
    if (quote && quote['昨日收盘价']) {
      fenshiPreClose.value = Number(quote['昨日收盘价']) || 0
    }
    if (minute && Array.isArray(minute.priceData)) {
      fenshiData.value = minute.priceData
      fenshiDate.value = minute.date || ''
      fenshiLoadedFor = s.code
    } else {
      fenshiData.value = []
    }
  } catch (e) {
    console.error('加载分时数据失败:', e)
    fenshiData.value = []
  } finally {
    fenshiLoading.value = false
  }
}

// ===== 日K数据 =====
const klineData = ref([])         // [{ day, open, close, low, high, volume }]
const klineLoading = ref(false)
let klineLoadedFor = ''           // 已加载该 code 的日K，避免重复拉取

async function loadKline() {
  const s = props.stock
  if (!s || !s.code) return
  if (klineLoadedFor === s.code && klineData.value.length) return
  klineLoading.value = true
  try {
    // 对齐桌面端 handleKLine：GetStockKLine(code, name, 365)
    const result = await GetStockKLine(s.code, s.name || '', 365).catch(() => null)
    if (Array.isArray(result) && result.length) {
      klineData.value = result
      klineLoadedFor = s.code
    } else {
      klineData.value = []
    }
  } catch (e) {
    console.error('加载日K数据失败:', e)
    klineData.value = []
  } finally {
    klineLoading.value = false
  }
}

// ===== 资金趋势数据 =====
const moneyData = ref([])         // [{ opendate, netamount, r0_net, trade }]
const moneyLoading = ref(false)
let moneyLoadedFor = ''           // 已加载该 code 的资金趋势，避免重复拉取

async function loadMoney() {
  const s = props.stock
  if (!s || !s.code) return
  if (moneyLoadedFor === s.code && moneyData.value.length) return
  moneyLoading.value = true
  try {
    // 对齐桌面端 moneyTrend：GetStockMoneyTrendByDay(code, 360)
    const result = await GetStockMoneyTrendByDay(s.code, 360).catch(() => null)
    if (Array.isArray(result) && result.length) {
      moneyData.value = result
      moneyLoadedFor = s.code
    } else {
      moneyData.value = []
    }
  } catch (e) {
    console.error('加载资金趋势失败:', e)
    moneyData.value = []
  } finally {
    moneyLoading.value = false
  }
}

// ===== 盘口（资料 tab）：买卖五档 =====
const quote = ref(null)          // Greet 完整行情
let quoteLoadedFor = ''
async function loadQuote() {
  const s = props.stock
  if (!s || !s.code) return
  if (quoteLoadedFor === s.code && quote.value) return
  try {
    const r = await Greet(s.code)
    quote.value = r || null
    quoteLoadedFor = s.code
  } catch (e) {
    console.error('加载盘口失败:', e)
  }
}

// 买卖五档（卖盘倒序：卖五在下、卖一在上紧贴成交价）
const bidAsk = computed(() => {
  const q = quote.value
  if (!q) return { asks: [], bids: [] }
  const asks = [5, 4, 3, 2, 1].map(i => ({
    label: `卖${['一','二','三','四','五'][i-1]}`,
    price: q[`卖${'一二三四五'[i-1]}报价`],
    volume: q[`卖${'一二三四五'[i-1]}申报`]
  }))
  const bids = [1, 2, 3, 4, 5].map(i => ({
    label: `买${'一二三四五'[i-1]}`,
    price: q[`买${'一二三四五'[i-1]}报价`],
    volume: q[`买${'一二三四五'[i-1]}申报`]
  }))
  return { asks, bids }
})

// 资料：基础行情字段
const info = computed(() => {
  const q = quote.value
  const s = props.stock
  const g = (k) => (q && q[k] != null && q[k] !== '') ? q[k] : ''
  return {
    price: g('当前价格') || s?.price || 0,
    open: g('今日开盘价'),
    high: g('今日最高价'),
    low: g('今日最低价'),
    preClose: g('昨日收盘价'),
    volume: g('成交的股票数'),
    turnover: g('成交金额')
  }
})

// ===== 持仓盈亏（来自 Greet，对齐桌面端 profit/profitAmount/profitAmountToday） =====
const hasPosition = computed(() => {
  const s = props.stock || {}
  return (Number(s.costPrice) || 0) > 0 && (Number(s.costVolume) || 0) > 0
})

// 盈亏方向：盈 rise(红) / 亏 fall(绿)
function pnlDir(amt) {
  const a = Number(amt) || 0
  if (a > 0) return 'rise'
  if (a < 0) return 'fall'
  return 'flat'
}

function signFixed(v, digits = 2) {
  const n = Number(v) || 0
  return `${n > 0 ? '+' : ''}${n.toFixed(digits)}`
}

const pnl = computed(() => {
  const s = props.stock || {}
  return {
    costPrice: Number(s.costPrice) || 0,
    costVolume: Number(s.costVolume) || 0,
    profit: Number(s.profit) || 0,                 // 总盈亏率 %
    profitAmount: Number(s.profitAmount) || 0,     // 总盈亏额 ¥
    profitToday: Number(s.profitToday) || 0        // 今日盈亏额 ¥
  }
})

// 价格相对昨收的方向（着色）
function priceDir(p) {
  const pc = Number(info.value.preClose) || 0
  if (!pc || !p) return 'flat'
  return Number(p) > pc ? 'rise' : Number(p) < pc ? 'fall' : 'flat'
}

function fmtVol(n) {
  const v = Number(n) || 0
  if (!v) return '--'
  if (v >= 100000000) return (v / 100000000).toFixed(2) + '亿'
  if (v >= 10000) return (v / 10000).toFixed(2) + '万'
  return String(v)
}

// 切换到分时 tab 才加载分时；日K / 盘口同理（按需懒加载）
watch(activeTab, (t) => {
  if (t === 'fenshi') loadFenshi()
  if (t === 'kline') loadKline()
  if (t === 'fund') loadMoney()
  if (t === 'info') loadQuote()
})

// 打开抽屉时：默认分时 tab，预加载分时 + 盘口
watch(() => props.show, (v) => {
  if (v && props.stock) {
    activeTab.value = 'fenshi'
    fenshiLoadedFor = ''
    quoteLoadedFor = ''
    klineLoadedFor = ''
    moneyLoadedFor = ''
    fenshiData.value = []
    quote.value = null
    klineData.value = []
    moneyData.value = []
    loadFenshi()
    loadQuote()
  }
})

function handleClose() {
  emit('update:show', false)
}

// ===== 设置弹窗（对齐桌面端：成本/数量/提醒/止盈止损/排序） =====
const settingVisible = ref(false)
const saving = ref(false)
const form = ref({
  costPrice: 0,      // 股票成本
  holdVolume: 0,     // 持仓数量
  alarm: 0,          // 涨跌提醒(%)
  alarmPrice: 0,     // 股价提醒
  entryPrice: 0,     // 开仓价
  takeProfitPrice: 0, // 止盈价
  stopLossPrice: 0,  // 止损价
  sort: 0            // 排序
})

function openSetting() {
  const s = props.stock || {}
  form.value = {
    costPrice: Number(s.costPrice) || 0,
    holdVolume: Number(s.holdVolume) || 0,
    alarm: Number(s.alarm) || 0,
    alarmPrice: Number(s.alarmPrice) || 0,
    entryPrice: Number(s.entryPrice) || 0,
    takeProfitPrice: Number(s.takeProfitPrice) || 0,
    stopLossPrice: Number(s.stopLossPrice) || 0,
    sort: Number(s.sort) || 0
  }
  settingVisible.value = true
}

function closeSetting() {
  settingVisible.value = false
}

async function saveSetting() {
  const s = props.stock
  if (!s || !s.code || saving.value) return
  saving.value = true
  const code = s.code
  const f = form.value
  try {
    // 排序
    if (f.sort) await SetStockSort(f.sort, code).catch(() => {})
    // 涨跌/股价提醒
    if (f.alarm || f.alarmPrice) {
      await SetAlarmChangePercent(f.alarm, f.alarmPrice, code).catch(() => {})
    }
    // 交易价格（开仓/止盈/止损/成本）
    if (f.entryPrice || f.takeProfitPrice || f.stopLossPrice || f.costPrice) {
      await SetTradingPrice(code, f.entryPrice || 0, f.takeProfitPrice || 0, f.stopLossPrice || 0, f.costPrice || 0).catch(() => {})
    }
    // 成本价 + 持仓数量（最后调，桌面端以它的返回为准）
    await SetCostPriceAndVolume(code, f.costPrice || 0, f.holdVolume || 0)
    // 回写到 stock 对象，下次打开预填最新值
    Object.assign(s, {
      costPrice: f.costPrice,
      holdVolume: f.holdVolume,
      alarm: f.alarm,
      alarmPrice: f.alarmPrice,
      entryPrice: f.entryPrice,
      takeProfitPrice: f.takeProfitPrice,
      stopLossPrice: f.stopLossPrice,
      sort: f.sort
    })
    settingVisible.value = false
    emit('saved', s)
    toast.success('保存成功')
  } catch (e) {
    console.error('保存设置失败:', e)
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <MSheet
    :show="show"
    :title="stock?.name"
    height="85vh"
    @update:show="handleClose"
  >
    <div v-if="stock" class="stock-detail">
      <!-- 股票基本信息 -->
      <div class="detail-header">
        <div class="header-left">
          <div class="stock-code">{{ stock.code }}</div>
        </div>
        <div class="header-right">
          <PriceTag
            :price="stock.price"
            :change="stock.changePercent"
            size="large"
            bold
          />
          <div class="price-change">
            <PercentTag :value="stock.changePercent" />
            <span
              class="change-amount"
              :class="{
                'm-rise': stock.changeAmount > 0,
                'm-fall': stock.changeAmount < 0
              }"
            >
              {{ stock.changeAmount > 0 ? '+' : '' }}{{ stock.changeAmount?.toFixed(2) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 持仓盈亏（有成本+持仓才显示，对齐桌面端今日盈亏/总盈亏） -->
      <div v-if="hasPosition" class="pnl-panel">
        <div class="pnl-cell">
          <span class="pnl-label">今日盈亏</span>
          <span class="pnl-value" :class="`m-${pnlDir(pnl.profitToday)}`">
            {{ signFixed(pnl.profitToday) }}
          </span>
        </div>
        <div class="pnl-divider" />
        <div class="pnl-cell">
          <span class="pnl-label">总盈亏</span>
          <span class="pnl-value" :class="`m-${pnlDir(pnl.profitAmount)}`">
            {{ signFixed(pnl.profitAmount) }}
          </span>
          <span class="pnl-sub" :class="`m-${pnlDir(pnl.profitAmount)}`">
            {{ signFixed(pnl.profit) }}%
          </span>
        </div>
        <div class="pnl-divider" />
        <div class="pnl-cell">
          <span class="pnl-label">成本/持仓</span>
          <span class="pnl-value pnl-cost">{{ pnl.costPrice.toFixed(2) }}</span>
          <span class="pnl-sub">{{ pnl.costVolume }}股</span>
        </div>
      </div>

      <!-- 标签页 -->
      <MTabs v-model="activeTab" :tabs="detailTabs" />

      <!-- 内容区 -->
      <div class="detail-content">
        <!-- 分时图 -->
        <div v-if="activeTab === 'fenshi'" class="chart-area">
          <div v-if="fenshiLoading" class="chart-loading">加载中...</div>
          <FenshiChart
            :data="fenshiData"
            :pre-close="fenshiPreClose"
            :date="fenshiDate"
          />
        </div>

        <!-- 日K：K线 + MA + 成交量 + 缩放（对齐桌面端 handleKLine） -->
        <div v-else-if="activeTab === 'kline'" class="chart-area">
          <div v-if="klineLoading" class="chart-loading">加载中...</div>
          <FullKlineChart :data="klineData" />
        </div>

        <!-- 资料：盘口 + 基础行情 -->
        <div v-else-if="activeTab === 'info'" class="info-panel">
          <!-- 买卖五档盘口 -->
          <div class="panel-title">买卖五档</div>
          <div class="bid-ask-panel">
            <div class="bid-ask-header">
              <span>档位</span>
              <span>价格</span>
              <span>挂单量</span>
            </div>
            <div class="ask-list">
              <div
                v-for="(item, index) in bidAsk.asks"
                :key="`ask-${index}`"
                class="bid-ask-item"
              >
                <span class="ba-label">{{ item.label }}</span>
                <span class="m-fall">{{ item.price || '--' }}</span>
                <span class="ba-vol">{{ fmtVol(item.volume) }}</span>
              </div>
            </div>
            <div class="bid-list">
              <div
                v-for="(item, index) in bidAsk.bids"
                :key="`bid-${index}`"
                class="bid-ask-item"
              >
                <span class="ba-label">{{ item.label }}</span>
                <span class="m-rise">{{ item.price || '--' }}</span>
                <span class="ba-vol">{{ fmtVol(item.volume) }}</span>
              </div>
            </div>
          </div>

          <!-- 基础行情 -->
          <div class="panel-title">基础行情</div>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">今开</span>
              <span class="info-value" :class="`m-${priceDir(info.open)}`">{{ info.open || '--' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">昨收</span>
              <span class="info-value">{{ info.preClose || '--' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">最高</span>
              <span class="info-value m-rise">{{ info.high || '--' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">最低</span>
              <span class="info-value m-fall">{{ info.low || '--' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">成交量</span>
              <span class="info-value">{{ fmtVol(info.volume) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">成交额</span>
              <span class="info-value">{{ fmtVol(info.turnover) }}</span>
            </div>
          </div>
        </div>

        <!-- 资金：当日净流入 + 主力净流入 + 股价 + 累计净流入副图（对齐桌面端 moneyTrend） -->
        <div v-else-if="activeTab === 'fund'" class="chart-area">
          <div v-if="moneyLoading" class="chart-loading">加载中...</div>
          <MoneyTrendChart :data="moneyData" :name="stock?.name" />
        </div>
      </div>
    </div>

    <!-- 底部操作：设置（成本/持仓/提醒/止盈止损） -->
    <template #footer>
      <div class="detail-actions">
        <MButton type="primary" block @click="openSetting">
          设置
        </MButton>
      </div>
    </template>
  </MSheet>

  <!-- 设置弹窗：对齐桌面端字段 -->
  <MSheet
    :show="settingVisible"
    title="设置"
    height="auto"
    @update:show="closeSetting"
  >
    <div class="setting-form">
      <div class="setting-row">
        <label class="setting-item">
          <span class="setting-label">股票成本</span>
          <div class="setting-input">
            <input v-model.number="form.costPrice" type="number" step="0.01" inputmode="decimal">
            <span class="setting-unit">¥</span>
          </div>
        </label>
        <label class="setting-item">
          <span class="setting-label">持仓数量</span>
          <div class="setting-input">
            <input v-model.number="form.holdVolume" type="number" step="100" inputmode="numeric">
            <span class="setting-unit">股</span>
          </div>
        </label>
      </div>
      <div class="setting-row">
        <label class="setting-item">
          <span class="setting-label">涨跌提醒</span>
          <div class="setting-input">
            <input v-model.number="form.alarm" type="number" step="0.1" inputmode="decimal">
            <span class="setting-unit">%</span>
          </div>
        </label>
        <label class="setting-item">
          <span class="setting-label">股价提醒</span>
          <div class="setting-input">
            <input v-model.number="form.alarmPrice" type="number" step="0.01" inputmode="decimal">
            <span class="setting-unit">¥</span>
          </div>
        </label>
      </div>
      <div class="setting-row">
        <label class="setting-item">
          <span class="setting-label">开仓价</span>
          <div class="setting-input">
            <input v-model.number="form.entryPrice" type="number" step="0.01" inputmode="decimal">
            <span class="setting-unit">¥</span>
          </div>
        </label>
        <label class="setting-item">
          <span class="setting-label">股票排序</span>
          <div class="setting-input">
            <input v-model.number="form.sort" type="number" step="1" inputmode="numeric">
          </div>
        </label>
      </div>
      <div class="setting-row">
        <label class="setting-item">
          <span class="setting-label">止盈价</span>
          <div class="setting-input">
            <input v-model.number="form.takeProfitPrice" type="number" step="0.01" inputmode="decimal">
            <span class="setting-unit">¥</span>
          </div>
        </label>
        <label class="setting-item">
          <span class="setting-label">止损价</span>
          <div class="setting-input">
            <input v-model.number="form.stopLossPrice" type="number" step="0.01" inputmode="decimal">
            <span class="setting-unit">¥</span>
          </div>
        </label>
      </div>
    </div>
    <template #footer>
      <div class="detail-actions">
        <MButton type="primary" block :loading="saving" @click="saveSetting">
          {{ saving ? '保存中...' : '保存' }}
        </MButton>
      </div>
    </template>
  </MSheet>
</template>

<style scoped>
.stock-detail {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding-bottom: var(--m-space-lg);
  border-bottom: 1px solid var(--m-divider-color);
}

.stock-code {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

.header-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--m-space-xs);
}

.price-change {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  font-size: var(--m-font-sm);
}

.change-amount {
  font-variant-numeric: tabular-nums;
}

/* ===== 持仓盈亏面板 ===== */
.pnl-panel {
  display: flex;
  align-items: stretch;
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md) var(--m-space-sm);
}

.pnl-cell {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  min-width: 0;
}

.pnl-divider {
  width: 1px;
  background: var(--m-divider-color);
  margin: 0 var(--m-space-xs);
}

.pnl-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.pnl-value {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
}

.pnl-cost {
  color: var(--m-text-primary);
}

.pnl-sub {
  font-size: var(--m-font-xs);
  font-variant-numeric: tabular-nums;
  color: var(--m-text-secondary);
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
  min-height: 320px;
}

.chart-area {
  position: relative;
  min-height: 300px;
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-sm);
}

.chart-loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--m-text-tertiary);
  font-size: var(--m-font-sm);
  z-index: 2;
}

.chart-placeholder {
  height: 250px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
  color: var(--m-text-secondary);
}

.chart-placeholder p {
  margin: var(--m-space-xs) 0;
}

.placeholder-tip {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

/* ===== 盘口 + 基础行情容器 ===== */
.info-panel {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
}

.panel-title {
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-secondary);
  margin-bottom: var(--m-space-sm);
}

.bid-ask-panel {
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
}

.bid-ask-header,
.bid-ask-item {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  align-items: center;
  padding: var(--m-space-xs) 0;
  font-size: var(--m-font-sm);
  font-variant-numeric: tabular-nums;
}

.bid-ask-header {
  color: var(--m-text-tertiary);
  font-weight: var(--m-font-weight-medium);
  border-bottom: 1px solid var(--m-divider-color);
}

.bid-ask-item .ba-label {
  color: var(--m-text-secondary);
}

.bid-ask-item .ba-vol {
  color: var(--m-text-secondary);
  text-align: right;
}

.bid-ask-item span:nth-child(2) {
  font-weight: var(--m-font-weight-medium);
}

.ask-list,
.bid-list {
  display: flex;
  flex-direction: column;
}

/* ===== 基础行情网格 ===== */
.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--m-space-sm);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-sm);
}

.info-label {
  color: var(--m-text-secondary);
}

.info-value {
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
}

.detail-actions {
  display: flex;
  gap: var(--m-space-md);
}

.detail-actions .m-button {
  flex: 1;
}

/* ===== 设置表单 ===== */
.setting-form {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  padding: var(--m-space-xs) 0;
}

.setting-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--m-space-sm);
}

.setting-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.setting-label {
  flex-shrink: 0;
  width: 4em;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.setting-input {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  height: 36px;
  padding: 0 var(--m-space-sm);
  background: var(--m-bg-primary);
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-sm);
}

.setting-input input {
  flex: 1;
  min-width: 0;
  width: 100%;
  border: none;
  background: transparent;
  outline: none;
  text-align: right;
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
}

.setting-unit {
  flex-shrink: 0;
  margin-left: 4px;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}
</style>
