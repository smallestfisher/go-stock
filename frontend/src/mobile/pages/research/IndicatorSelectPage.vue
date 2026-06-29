<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import { SearchStock, GetHotStrategy, GetAllCustomStrategies, SaveCustomStrategy, DeleteCustomStrategy, Follow } from '../../../api/app'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MButton from '../../components/base/MButton.vue'
import MSheet from '../../components/base/MSheet.vue'
import StockDetailSheet from '../../components/sheets/StockDetailSheet.vue'
import { toast } from '../../composables/useToast'

// 对齐桌面端 SelectStock.vue：热门策略 / 我的策略 + 自然语言选股(SearchStock) + 结果列表 + 关注 + 保存策略
// GetHotStrategy → {code:1, data:[{rank, question}]}
// GetAllCustomStrategies → [{id, name, query, description, sortOrder}]
// SearchStock(query) → {code:100, data:{traceInfo:{showText}, result:{columns, dataList}}}，行字段 SECURITY_CODE/SECURITY_SHORT_NAME/MARKET_SHORT_NAME

const TABS = [
  { label: '热门策略', value: 'hot' },
  { label: '我的策略', value: 'custom' },
]
const activeTab = ref('hot')

const hotStrategies = ref([])       // [{rank, question}]
const customStrategies = ref([])    // [{id, name, query, description}]
const strategyLoading = ref(false)

// 选股
const query = ref('')
const searching = ref(false)
const traceInfo = ref('')
const columns = ref([])             // 原始列定义（含 children）
const dataList = ref([])            // 结果行

const detailVisible = ref(false)
const selectedStock = ref(null)

// 展平列（去掉操作列与 children 层级），用于卡片展示关键字段
function flattenColumns(list) {
  const result = []
  for (const col of list || []) {
    if (col.key === 'actions') continue
    if (col.children && col.children.length) {
      for (const child of col.children) {
        result.push({ title: `${col.title}/${child.title}`, key: child.key })
      }
    } else {
      result.push({ title: col.title, key: col.key })
    }
  }
  return result
}

// 取每行用于展示的前 6 个关键字段（排除代码/名称/市场等标识列）
function rowFields(row) {
  return flattenColumns(columns.value)
    .filter(col => !['SECURITY_CODE', 'SECURITY_SHORT_NAME', 'MARKET_SHORT_NAME', 'SERIAL'].includes(col.key))
    .map(col => ({ ...col, value: row[col.key] }))
    .filter(item => item.value !== undefined && item.value !== null && item.value !== '')
    .slice(0, 6)
}

async function loadHotStrategies() {
  try {
    const res = await GetHotStrategy()
    if (res && res.code === 1 && Array.isArray(res.data)) {
      hotStrategies.value = res.data
    }
  } catch (e) {
    console.error('加载热门策略失败:', e)
  }
}

async function loadCustomStrategies() {
  try {
    const res = await GetAllCustomStrategies()
    customStrategies.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error('加载自定义策略失败:', e)
  }
}

async function handleRefresh() {
  strategyLoading.value = true
  try {
    await Promise.all([loadHotStrategies(), loadCustomStrategies()])
  } finally {
    strategyLoading.value = false
  }
}

// 执行选股
async function runSearch() {
  const q = query.value.trim()
  if (!q) {
    toast.warning('请输入选股指标或要求')
    return
  }
  searching.value = true
  traceInfo.value = ''
  try {
    const res = await SearchStock(q)
    if (res && res.code === 100 && res.data && res.data.result) {
      traceInfo.value = res.data.traceInfo?.showText || ''
      columns.value = (res.data.result.columns || []).filter(
        c => !c.hiddenNeed && c.title !== '市场码' && c.title !== '市场简称'
      )
      dataList.value = res.data.result.dataList || []
      if (!dataList.value.length) toast.info('没有符合条件的股票')
    } else {
      const msg = (res && (res.msg || res.message)) || '选股失败'
      toast.error(msg)
      columns.value = []
      dataList.value = []
    }
  } catch (e) {
    console.error('选股失败:', e)
    toast.error('选股失败: ' + (e.message || ''))
  } finally {
    searching.value = false
  }
}

// 点策略 → 填入并执行
function applyStrategy(text) {
  query.value = text
  runSearch()
}

// 关注（对齐桌面：Follow(market.toLowerCase() + code)）
async function followStock(row) {
  const code = (row.MARKET_SHORT_NAME || '').toLowerCase() + row.SECURITY_CODE
  try {
    const res = await Follow(code)
    if (res === '关注成功') toast.success(res)
    else toast.warning(typeof res === 'string' ? res : '关注完成')
  } catch (e) {
    toast.error('关注失败: ' + (e.message || ''))
  }
}

// 看详情：转东财格式（sh600519），交给 StockDetailSheet
function openDetail(row) {
  const market = (row.MARKET_SHORT_NAME || '').toLowerCase()
  const code = row.SECURITY_CODE
  let full = ''
  if (market === 'sh' || market === 'sz' || market === 'bj') full = market + code
  else if (market === 'hk') full = 'hk' + code
  else if (market === 'us') full = 'gb_' + String(code).toLowerCase()
  else full = code
  selectedStock.value = { code: full, name: row.SECURITY_SHORT_NAME || '' }
  detailVisible.value = true
}

// ===== 保存 / 删除 自定义策略 =====
const saveVisible = ref(false)
const saveForm = ref({ id: 0, name: '', query: '', description: '' })
const saving = ref(false)

function openSave() {
  if (!query.value.trim()) {
    toast.warning('请先输入选股条件')
    return
  }
  saveForm.value = { id: 0, name: '', query: query.value, description: '' }
  saveVisible.value = true
}

function openEditStrategy(s) {
  saveForm.value = { id: s.id, name: s.name, query: s.query, description: s.description || '' }
  saveVisible.value = true
}

async function submitSave() {
  const f = saveForm.value
  if (!f.name.trim()) {
    toast.warning('请输入策略名称')
    return
  }
  if (!f.query.trim()) {
    toast.warning('请输入选股条件')
    return
  }
  saving.value = true
  try {
    const res = await SaveCustomStrategy({ id: f.id || 0, name: f.name, query: f.query, description: f.description, sortOrder: 0 })
    toast.success(typeof res === 'string' ? res : '保存成功')
    saveVisible.value = false
    loadCustomStrategies()
  } catch (e) {
    toast.error('保存失败: ' + (e.message || ''))
  } finally {
    saving.value = false
  }
}

async function removeStrategy(s) {
  if (!confirm(`确定删除策略「${s.name}」吗？`)) return
  try {
    const res = await DeleteCustomStrategy(s.id)
    toast.success(typeof res === 'string' ? res : '删除成功')
    loadCustomStrategies()
  } catch (e) {
    toast.error('删除失败: ' + (e.message || ''))
  }
}

const hasResult = computed(() => dataList.value.length > 0)

onBeforeMount(() => {
  loadHotStrategies()
  loadCustomStrategies()
})
</script>

<template>
  <div class="page">
    <!-- 选股输入 -->
    <div class="op-bar">
      <div class="search-row">
        <input
          v-model="query"
          class="search-input"
          type="text"
          placeholder="输入选股条件，如：市盈率小于20且换手率大于5%"
          @keyup.enter="runSearch"
        />
        <button type="button" class="search-btn" :disabled="searching" @click="runSearch">
          {{ searching ? '...' : '选股' }}
        </button>
      </div>
      <div class="tab-row">
        <button
          v-for="t in TABS"
          :key="t.value"
          type="button"
          class="tab"
          :class="{ 'tab--active': activeTab === t.value }"
          @click="activeTab = t.value"
        >{{ t.label }}</button>
        <button type="button" class="save-link" @click="openSave">＋ 存为策略</button>
      </div>
    </div>

    <MPullRefresh class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <!-- 策略列表 -->
        <div class="strategy-block">
          <template v-if="activeTab === 'hot'">
            <div v-if="hotStrategies.length" class="strategy-list">
              <button
                v-for="s in hotStrategies"
                :key="s.rank"
                type="button"
                class="strategy-chip"
                @click="applyStrategy(s.question)"
              >
                <span class="strategy-rank">#{{ s.rank }}</span>
                <span class="strategy-text">{{ s.question }}</span>
              </button>
            </div>
            <MLoading v-else-if="strategyLoading" text="加载中..." vertical class="block-loading" />
            <MEmpty v-else description="暂无热门策略" />
          </template>

          <template v-else>
            <div v-if="customStrategies.length" class="custom-list">
              <div v-for="s in customStrategies" :key="s.id" class="custom-item">
                <div class="custom-main" @click="applyStrategy(s.query)">
                  <div class="custom-name">{{ s.name }}</div>
                  <p class="custom-query">{{ s.query }}</p>
                  <p v-if="s.description" class="custom-desc">{{ s.description }}</p>
                </div>
                <div class="custom-acts">
                  <button type="button" class="mini-btn" @click="openEditStrategy(s)">编辑</button>
                  <button type="button" class="mini-btn mini-btn--danger" @click="removeStrategy(s)">删除</button>
                </div>
              </div>
            </div>
            <MEmpty v-else description="还没有自定义策略，选股后可保存为策略" />
          </template>
        </div>

        <!-- 选股结果 -->
        <div v-if="searching" class="result-loading">
          <MLoading text="正在选股..." vertical />
        </div>
        <div v-else-if="hasResult" class="result-section">
          <div class="result-head">
            <h4 class="result-title">选股结果（{{ dataList.length }}）</h4>
          </div>
          <p v-if="traceInfo" class="trace-info">{{ traceInfo }}</p>
          <div class="result-list">
            <div v-for="(row, i) in dataList" :key="row.SECURITY_CODE || i" class="result-card">
              <div class="result-card-head" @click="openDetail(row)">
                <div class="result-name-wrap">
                  <span class="result-name">{{ row.SECURITY_SHORT_NAME }}</span>
                  <span class="result-code">{{ row.SECURITY_CODE }}</span>
                </div>
                <button type="button" class="follow-btn" @click.stop="followStock(row)">＋关注</button>
              </div>
              <div v-if="rowFields(row).length" class="result-fields">
                <div v-for="f in rowFields(row)" :key="f.key" class="field-cell">
                  <span class="field-name">{{ f.title }}</span>
                  <span class="field-value">{{ f.value }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </MPullRefresh>

    <!-- 保存策略抽屉 -->
    <MSheet v-model:show="saveVisible" :title="saveForm.id > 0 ? '编辑策略' : '保存为策略'" height="auto">
      <div class="form">
        <label class="field">
          <span class="field-label">策略名称 <i class="req">*</i></span>
          <input v-model="saveForm.name" class="field-input" type="text" placeholder="给策略起个名字" />
        </label>
        <label class="field">
          <span class="field-label">选股条件 <i class="req">*</i></span>
          <textarea v-model="saveForm.query" class="field-textarea" rows="3" placeholder="选股条件" />
        </label>
        <label class="field">
          <span class="field-label">描述</span>
          <textarea v-model="saveForm.description" class="field-textarea" rows="2" placeholder="可选" />
        </label>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="saveVisible = false">取消</MButton>
          <MButton type="primary" block :loading="saving" @click="submitSave">保存</MButton>
        </div>
      </template>
    </MSheet>

    <!-- 个股详情 -->
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

.search-btn {
  flex-shrink: 0;
  padding: 0 var(--m-space-lg);
  border: none;
  border-radius: var(--m-radius-full);
  background: var(--m-color-rise);
  color: #fff;
  font: inherit;
  font-size: var(--m-font-sm);
}

.search-btn:disabled {
  opacity: 0.6;
}

.tab-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.tab {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.tab--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.save-link {
  margin-left: auto;
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-sm);
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
  gap: var(--m-space-lg);
}

.block-loading,
.result-loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-xl) 0;
}

/* 策略 */
.strategy-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.strategy-chip {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  width: 100%;
  text-align: left;
  padding: var(--m-space-sm) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-md);
  background: var(--m-bg-card);
  font: inherit;
}

.strategy-chip:active {
  background: var(--m-bg-primary);
}

.strategy-rank {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  color: var(--m-color-rise);
  background: var(--m-color-rise-light);
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
}

.strategy-text {
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  line-height: var(--m-line-height-normal);
}

.custom-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.custom-item {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.custom-main {
  margin-bottom: var(--m-space-sm);
}

.custom-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.custom-query {
  margin-top: 2px;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
}

.custom-desc {
  margin-top: 2px;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.custom-acts {
  display: flex;
  gap: var(--m-space-md);
  border-top: 1px solid var(--m-divider-color);
  padding-top: var(--m-space-sm);
}

.mini-btn {
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.mini-btn--danger {
  color: var(--m-color-rise);
}

/* 结果 */
.result-section {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.result-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.trace-info {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  line-height: var(--m-line-height-normal);
  padding: var(--m-space-sm);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-sm);
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.result-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.result-card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
}

.result-name-wrap {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  min-width: 0;
}

.result-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.result-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.follow-btn {
  flex-shrink: 0;
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-color-rise);
  border-radius: var(--m-radius-full);
  background: transparent;
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-color-rise);
}

.result-fields {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--m-space-xs) var(--m-space-md);
  margin-top: var(--m-space-sm);
  padding-top: var(--m-space-sm);
  border-top: 1px solid var(--m-divider-color);
}

.field-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-xs);
  min-width: 0;
}

.field-name {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.field-value {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
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

.sheet-footer {
  display: flex;
  gap: var(--m-space-md);
}
</style>
