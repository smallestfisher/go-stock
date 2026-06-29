<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import {
  GetTradingRecordList,
  GetTradingRecordStatistics,
  AddTradingRecord,
  UpdateTradingRecord,
  DeleteTradingRecord,
  GetAllStockInfoList,
  GetStockRealTimePrice,
} from '../../../api/app'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MButton from '../../components/base/MButton.vue'
import MSheet from '../../components/base/MSheet.vue'
import { formatMoney } from '../../composables/useFormat'
import { toast } from '../../composables/useToast'

// 对齐桌面端 TradingRecordManager.vue：统计概览 + 记录列表(搜索/方向筛选/分页) + 新增/编辑/删除
// 统计字段：totalBuyAmount/totalSellAmount/totalProfit/profitRate/holdingsAmount/currentValue/stockCount
// 列表字段：StockCode/StockName/Direction(买入/卖出)/Price/Volume/Amount/TradingTime/Reason，归一化 closePrice/profitAmount/profitPercent

const PAGE_SIZE = 15

const DIRECTIONS = [
  { label: '全部', value: '' },
  { label: '买入', value: '买入' },
  { label: '卖出', value: '卖出' },
]

const statistics = ref(null)
const records = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const totalPages = ref(1)

const keyword = ref('')
const direction = ref('')

const hasMore = computed(() => page.value < totalPages.value)

async function loadStatistics() {
  try {
    const res = await GetTradingRecordStatistics()
    statistics.value = res || null
  } catch (e) {
    console.error('加载交易统计失败:', e)
  }
}

// 列表行字段可能是 PascalCase，统一归一化
function normalizeRow(row) {
  if (!row || typeof row !== 'object') return row
  return {
    ...row,
    profitAmount: Number(row.profitAmount ?? row.ProfitAmount ?? 0),
    profitPercent: Number(row.profitPercent ?? row.ProfitPercent ?? 0),
    closePrice: Number(row.closePrice ?? row.ClosePrice ?? 0),
  }
}

async function loadRecords(targetPage = 1, append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const res = await GetTradingRecordList({
      page: targetPage,
      pageSize: PAGE_SIZE,
      keyword: keyword.value,
      direction: direction.value,
      startDate: '',
      endDate: '',
    })
    const rows = (res && Array.isArray(res.list) ? res.list : []).map(normalizeRow)
    records.value = append ? records.value.concat(rows) : rows
    total.value = (res && res.total) || 0
    totalPages.value = (res && res.totalPages) || 1
    page.value = targetPage
  } catch (e) {
    console.error('加载交易记录失败:', e)
    toast.error('加载失败')
  } finally {
    loading.value = false
  }
}

function handleRefresh() {
  return Promise.all([loadStatistics(), loadRecords(1, false)])
}

function loadMore() {
  if (hasMore.value && !loading.value) loadRecords(page.value + 1, true)
}

function applySearch() {
  loadRecords(1, false)
}

function selectDirection(value) {
  direction.value = value
  loadRecords(1, false)
}

// ===== 新增 / 编辑 =====
const editVisible = ref(false)
const editSaving = ref(false)
const blankForm = () => ({
  ID: 0,
  StockCode: '',
  StockName: '',
  Direction: '买入',
  Price: 0,
  Volume: 0,
  Reason: '',
  StopLossPrice: 0,
  TakeProfitPrice: 0,
  Mindset: '',
  TradingTime: Date.now(),
})
const editForm = ref(blankForm())
const editTitle = computed(() => (editForm.value.ID > 0 ? '编辑交易记录' : '新增交易记录'))

// 把时间戳/字符串转 datetime-local 输入值（本地时区）
function toLocalInput(ts) {
  const d = new Date(ts)
  if (Number.isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}
const tradingTimeInput = computed({
  get: () => toLocalInput(editForm.value.TradingTime),
  set: (v) => { editForm.value.TradingTime = v ? new Date(v).getTime() : Date.now() },
})

function openCreate() {
  editForm.value = blankForm()
  stockResults.value = []
  stockKeyword.value = ''
  editVisible.value = true
}

function openEdit(item) {
  editForm.value = {
    ID: item.ID,
    StockCode: item.StockCode,
    StockName: item.StockName,
    Direction: item.Direction || '买入',
    Price: item.Price || 0,
    Volume: item.Volume || 0,
    Reason: item.Reason || '',
    StopLossPrice: item.StopLossPrice || 0,
    TakeProfitPrice: item.TakeProfitPrice || 0,
    Mindset: item.Mindset || '',
    TradingTime: item.TradingTime ? new Date(item.TradingTime).getTime() : Date.now(),
  }
  stockResults.value = []
  stockKeyword.value = ''
  editVisible.value = true
}

// 股票搜索（对齐桌面 GetAllStockInfoList → SECUCODE/SECURITY_NAME_ABBR/MARKET）
const stockKeyword = ref('')
const stockResults = ref([])
const stockSearching = ref(false)

async function searchStock() {
  const kw = stockKeyword.value.trim()
  if (!kw) {
    stockResults.value = []
    return
  }
  stockSearching.value = true
  try {
    const res = await GetAllStockInfoList({ searchKeyWord: kw })
    stockResults.value = (res && Array.isArray(res.list) ? res.list : []).slice(0, 20)
  } catch (e) {
    console.error('搜索股票失败:', e)
  } finally {
    stockSearching.value = false
  }
}

// 东财代码格式（用于查实时价）：6→sh，0/3→sz，8/4→bj
function toEastMoneyCode(secucode, market) {
  const c = String(secucode || '').toUpperCase()
  if (c.endsWith('.SH')) return 'sh' + c.slice(0, -3).toLowerCase()
  if (c.endsWith('.SZ')) return 'sz' + c.slice(0, -3).toLowerCase()
  if (c.endsWith('.BJ')) return 'bj' + c.slice(0, -3).toLowerCase()
  const code = c.replace(/\..*$/, '')
  if (code.startsWith('6')) return 'sh' + code.toLowerCase()
  if (code.startsWith('0') || code.startsWith('3')) return 'sz' + code.toLowerCase()
  if (code.startsWith('8') || code.startsWith('4')) return 'bj' + code.toLowerCase()
  return code.toLowerCase()
}

async function pickStock(item) {
  editForm.value.StockCode = item.SECUCODE
  editForm.value.StockName = item.SECURITY_NAME_ABBR
  stockResults.value = []
  stockKeyword.value = `${item.SECURITY_NAME_ABBR} (${item.SECUCODE})`
  // 拉实时价填入
  try {
    const res = await GetStockRealTimePrice(toEastMoneyCode(item.SECUCODE, item.MARKET))
    if (res && res.code === 0 && res.price > 0) editForm.value.Price = res.price
  } catch (_) {}
}

async function saveRecord() {
  const f = editForm.value
  if (!f.StockCode || !f.StockName) {
    toast.warning('请选择股票')
    return
  }
  if (!f.Price || !f.Volume) {
    toast.warning('请填写价格和数量')
    return
  }
  editSaving.value = true
  try {
    const payload = {
      ...f,
      Amount: Number(f.Price) * Number(f.Volume),
      TradingTime: new Date(f.TradingTime),
    }
    const apiCall = f.ID > 0 ? UpdateTradingRecord : AddTradingRecord
    await apiCall(payload)
    toast.success('保存成功')
    editVisible.value = false
    handleRefresh()
  } catch (e) {
    console.error('保存交易记录失败:', e)
    toast.error('保存失败')
  } finally {
    editSaving.value = false
  }
}

async function removeRecord(item) {
  if (!confirm(`确定删除「${item.StockName}」这条交易记录吗？`)) return
  try {
    await DeleteTradingRecord(item.ID)
    toast.success('删除成功')
    handleRefresh()
  } catch (e) {
    console.error('删除交易记录失败:', e)
    toast.error('删除失败')
  }
}

function fmtTime(t) {
  if (!t) return '-'
  return String(t).substring(0, 19).replace('T', ' ')
}

function profitClass(v) {
  const n = Number(v) || 0
  if (n > 0) return 'm-rise'
  if (n < 0) return 'm-fall'
  return ''
}

const profitRateText = computed(() => {
  const r = Number(statistics.value?.profitRate) || 0
  return `${r > 0 ? '+' : ''}${r.toFixed(2)}%`
})

onBeforeMount(() => {
  loadStatistics()
  loadRecords(1, false)
})
</script>

<template>
  <div class="page">
    <!-- 操作条 -->
    <div class="op-bar">
      <div class="search-row">
        <input
          v-model="keyword"
          class="search-input"
          type="text"
          placeholder="搜索股票名称/代码..."
          @keyup.enter="applySearch"
        />
        <MButton type="primary" size="small" @click="openCreate">＋ 记录</MButton>
      </div>
      <div class="chip-row">
        <button
          v-for="d in DIRECTIONS"
          :key="d.value"
          type="button"
          class="chip"
          :class="{ 'chip--active': direction === d.value }"
          @click="selectDirection(d.value)"
        >{{ d.label }}</button>
      </div>
    </div>

    <MPullRefresh class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <!-- 统计概览 -->
        <div v-if="statistics" class="stat-card">
          <div class="stat-grid">
            <div class="stat-item">
              <div class="stat-value" :class="profitClass(statistics.totalProfit)">
                {{ formatMoney(statistics.totalProfit) }}
              </div>
              <div class="stat-label">总盈亏</div>
            </div>
            <div class="stat-item">
              <div class="stat-value" :class="profitClass(statistics.profitRate)">{{ profitRateText }}</div>
              <div class="stat-label">收益率</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ statistics.stockCount || 0 }}</div>
              <div class="stat-label">持仓股数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ formatMoney(statistics.totalBuyAmount) }}</div>
              <div class="stat-label">买入总额</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ formatMoney(statistics.totalSellAmount) }}</div>
              <div class="stat-label">卖出总额</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ formatMoney(statistics.currentValue) }}</div>
              <div class="stat-label">当前市值</div>
            </div>
          </div>
        </div>

        <!-- 记录列表 -->
        <div v-if="records.length" class="record-list">
          <div v-for="item in records" :key="item.ID" class="record-card">
            <div class="record-head">
              <div class="record-name-wrap">
                <span class="dir-tag" :class="item.Direction === '买入' ? 'dir-tag--buy' : 'dir-tag--sell'">
                  {{ item.Direction }}
                </span>
                <span class="record-name">{{ item.StockName }}</span>
                <span class="record-code">{{ item.StockCode }}</span>
              </div>
              <button type="button" class="del-btn" @click="removeRecord(item)">✕</button>
            </div>

            <div class="record-body" @click="openEdit(item)">
              <div class="record-fields">
                <div class="rf">
                  <span class="rf-label">价格</span>
                  <span class="rf-value">{{ item.Price }}</span>
                </div>
                <div class="rf">
                  <span class="rf-label">数量</span>
                  <span class="rf-value">{{ item.Volume }}</span>
                </div>
                <div class="rf">
                  <span class="rf-label">金额</span>
                  <span class="rf-value">{{ formatMoney(item.Amount) }}</span>
                </div>
                <div v-if="item.profitAmount" class="rf">
                  <span class="rf-label">盈亏</span>
                  <span class="rf-value" :class="profitClass(item.profitAmount)">
                    {{ formatMoney(item.profitAmount) }}
                    <template v-if="item.profitPercent">({{ item.profitPercent.toFixed(2) }}%)</template>
                  </span>
                </div>
              </div>
              <div class="record-foot">
                <span class="record-time">{{ fmtTime(item.TradingTime) }}</span>
                <span v-if="item.Reason" class="record-reason">{{ item.Reason }}</span>
              </div>
            </div>
          </div>

          <div class="load-more">
            <button v-if="hasMore" type="button" class="load-more-btn" :disabled="loading" @click="loadMore">
              {{ loading ? '加载中...' : '加载更多' }}
            </button>
            <span v-else class="load-more-end">共 {{ total }} 条</span>
          </div>
        </div>

        <MLoading v-else-if="loading" text="加载中..." vertical class="page-loading" />
        <MEmpty v-else description="还没有交易记录">
          <MButton type="primary" @click="openCreate">添加记录</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 新增/编辑抽屉 -->
    <MSheet v-model:show="editVisible" :title="editTitle">
      <div class="form">
        <!-- 股票选择 -->
        <div class="field">
          <span class="field-label">股票 <i class="req">*</i></span>
          <input
            v-model="stockKeyword"
            class="field-input"
            type="text"
            placeholder="输入代码/名称搜索"
            @input="searchStock"
          />
          <div v-if="stockSearching" class="stock-hint">搜索中...</div>
          <div v-else-if="stockResults.length" class="stock-results">
            <button
              v-for="s in stockResults"
              :key="s.SECUCODE"
              type="button"
              class="stock-result"
              @click="pickStock(s)"
            >
              <span class="sr-name">{{ s.SECURITY_NAME_ABBR }}</span>
              <span class="sr-code">{{ s.SECUCODE }}</span>
            </button>
          </div>
          <div v-if="editForm.StockCode" class="picked-stock">
            已选：{{ editForm.StockName }}（{{ editForm.StockCode }}）
          </div>
        </div>

        <!-- 方向 -->
        <div class="field">
          <span class="field-label">方向 <i class="req">*</i></span>
          <div class="chip-row">
            <button
              type="button"
              class="chip"
              :class="{ 'chip--active': editForm.Direction === '买入' }"
              @click="editForm.Direction = '买入'"
            >买入</button>
            <button
              type="button"
              class="chip"
              :class="{ 'chip--active': editForm.Direction === '卖出' }"
              @click="editForm.Direction = '卖出'"
            >卖出</button>
          </div>
        </div>

        <div class="field-grid">
          <label class="field">
            <span class="field-label">价格 <i class="req">*</i></span>
            <input v-model.number="editForm.Price" class="field-input" type="number" inputmode="decimal" placeholder="0.00" />
          </label>
          <label class="field">
            <span class="field-label">数量(股) <i class="req">*</i></span>
            <input v-model.number="editForm.Volume" class="field-input" type="number" inputmode="numeric" placeholder="0" />
          </label>
        </div>

        <div class="field-grid">
          <label class="field">
            <span class="field-label">止损价</span>
            <input v-model.number="editForm.StopLossPrice" class="field-input" type="number" inputmode="decimal" placeholder="0.00" />
          </label>
          <label class="field">
            <span class="field-label">止盈价</span>
            <input v-model.number="editForm.TakeProfitPrice" class="field-input" type="number" inputmode="decimal" placeholder="0.00" />
          </label>
        </div>

        <label class="field">
          <span class="field-label">交易时间</span>
          <input v-model="tradingTimeInput" class="field-input" type="datetime-local" />
        </label>

        <label class="field">
          <span class="field-label">交易理由</span>
          <textarea v-model="editForm.Reason" class="field-textarea" rows="2" placeholder="为什么买/卖" />
        </label>

        <label class="field">
          <span class="field-label">心态记录</span>
          <textarea v-model="editForm.Mindset" class="field-textarea" rows="2" placeholder="当时的心态" />
        </label>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="editVisible = false">取消</MButton>
          <MButton type="primary" block :loading="editSaving" @click="saveRecord">保存</MButton>
        </div>
      </template>
    </MSheet>
  </div>
</template>

<style scoped>
.page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 操作条 */
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

.chip-row {
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
  color: var(--m-text-secondary);
}

.chip--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.page-loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

/* 统计 */
.stat-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--m-space-md) var(--m-space-sm);
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
}

.stat-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  margin-top: 2px;
}

/* 记录列表 */
.record-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.record-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.record-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
}

.record-name-wrap {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  min-width: 0;
}

.dir-tag {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  align-self: center;
}

.dir-tag--buy {
  color: var(--m-color-rise);
  background: var(--m-color-rise-light);
}

.dir-tag--sell {
  color: var(--m-color-fall);
  background: var(--m-color-fall-light);
}

.record-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.record-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.del-btn {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  border: none;
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  color: var(--m-text-tertiary);
  font-size: var(--m-font-xs);
}

.record-body {
  margin-top: var(--m-space-sm);
}

.record-fields {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--m-space-xs) var(--m-space-md);
}

.rf {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-xs);
}

.rf-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.rf-value {
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
}

.record-foot {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  margin-top: var(--m-space-sm);
  padding-top: var(--m-space-sm);
  border-top: 1px solid var(--m-divider-color);
}

.record-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
  flex-shrink: 0;
}

.record-reason {
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 加载更多 */
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

/* 表单 */
.form {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
}

.field {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.field-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--m-space-md);
}

.field-label {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.req {
  color: var(--m-color-rise);
  font-style: normal;
}

.field-input,
.field-textarea {
  width: 100%;
  padding: var(--m-space-sm) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font: inherit;
  font-size: var(--m-font-md);
  outline: none;
}

.field-textarea {
  resize: vertical;
  line-height: var(--m-line-height-normal);
}

.field-input:focus,
.field-textarea:focus {
  border-color: var(--m-color-rise);
}

/* 股票搜索结果 */
.stock-hint {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.stock-results {
  display: flex;
  flex-direction: column;
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
}

.stock-result {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border: none;
  border-bottom: 1px solid var(--m-divider-color);
  font: inherit;
  text-align: left;
}

.stock-result:last-child {
  border-bottom: none;
}

.sr-name {
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
}

.sr-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.picked-stock {
  font-size: var(--m-font-xs);
  color: var(--m-color-rise);
}

.sheet-footer {
  display: flex;
  gap: var(--m-space-md);
}
</style>
