<script setup>
import { ref, computed, watch, onBeforeMount } from 'vue'
import {
  GetCronTaskList,
  GetCronTaskByID,
  GetCronTaskTypes,
  CreateCronTask,
  UpdateCronTask,
  DeleteCronTask,
  EnableCronTask,
  ExecuteCronTaskNow,
  ValidateCronExpr,
  CalculateNextRunTimes,
  GetAiConfigs,
  GetPromptTemplates,
} from '../../../api/app'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MButton from '../../components/base/MButton.vue'
import MSheet from '../../components/base/MSheet.vue'
import MIcon from '../../components/base/MIcon.vue'
import { toast } from '../../composables/useToast'

// 对齐桌面端 cron-task-manager.vue：任务列表 + 搜索/筛选 + 启用切换 + 立即执行 + 新建/编辑/删除
// 移动端把桌面的可视化 Cron 构建器简化为：表达式直接输入 + 校验 + 下次执行预览。
// 股票/市场分析任务保留参数表单（提示词/AI配置/系统提示词/思考/Agent模式/股票代码）。

const PAGE_SIZE = 12

const STATUS_FILTERS = [
  { label: '全部', value: '' },
  { label: '活跃', value: 'active' },
  { label: '暂停', value: 'paused' },
  { label: '错误', value: 'error' },
]

const AGENT_MODES = [
  { label: '🤖 自动选择', value: '' },
  { label: '⚡ 快速模式', value: 'react' },
  { label: '🧠 规划模式', value: 'plan_execute' },
]

const list = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const totalPages = ref(1)

const search = ref({ name: '', taskType: '', status: '' })
const filterVisible = ref(false)
const hasFilter = computed(() => !!(search.value.name || search.value.taskType || search.value.status))
const hasMore = computed(() => page.value < totalPages.value)

// 选项
const taskTypeOptions = ref([])  // [{label, value}]
const aiConfigOptions = ref([])  // [{label, value}]
const userPromptOptions = ref([])
const sysPromptOptions = ref([])

function taskTypeLabel(value) {
  const o = taskTypeOptions.value.find(t => t.value === value)
  return o ? o.label : value
}

async function fetchPage(targetPage, append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const res = await GetCronTaskList({
      page: targetPage,
      pageSize: PAGE_SIZE,
      name: search.value.name,
      taskType: search.value.taskType,
      status: search.value.status,
    })
    const rows = res && Array.isArray(res.data) ? res.data : []
    list.value = append ? list.value.concat(rows) : rows
    total.value = (res && res.total) || 0
    totalPages.value = Math.ceil(((res && res.total) || 0) / PAGE_SIZE) || 1
    page.value = targetPage
  } catch (e) {
    console.error('加载定时任务失败:', e)
    toast.error('加载失败')
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
  filterVisible.value = false
  fetchPage(1, false)
}

function resetSearch() {
  search.value = { name: '', taskType: '', status: '' }
  filterVisible.value = false
  fetchPage(1, false)
}

async function loadOptions() {
  try {
    const types = await GetCronTaskTypes()
    // lo.Tuple2 序列化为 {A:value, B:label}
    taskTypeOptions.value = (Array.isArray(types) ? types : []).map(t => ({ label: t.B, value: t.A }))
  } catch (e) {
    console.error('加载任务类型失败:', e)
  }
  try {
    const configs = await GetAiConfigs()
    aiConfigOptions.value = (Array.isArray(configs) ? configs : []).map(c => ({
      label: c.name + (c.modelName ? `[${c.modelName}]` : ''),
      value: c.ID,
    }))
  } catch (e) {
    console.error('加载AI配置失败:', e)
  }
  try {
    const [users, sys] = await Promise.all([
      GetPromptTemplates('', '模型用户Prompt'),
      GetPromptTemplates('', '模型系统Prompt'),
    ])
    userPromptOptions.value = (Array.isArray(users) ? users : []).map(t => ({ label: t.name, value: t.ID }))
    sysPromptOptions.value = (Array.isArray(sys) ? sys : []).map(t => ({ label: t.name, value: t.ID }))
  } catch (e) {
    console.error('加载提示词失败:', e)
  }
}

// ===== 启用切换 / 立即执行 / 删除 =====
async function toggleEnable(task) {
  try {
    const next = !task.enable
    const res = await EnableCronTask(task.id, next)
    if (res === '操作成功') {
      task.enable = next
      toast.success(next ? '任务已启用' : '任务已禁用')
    } else {
      toast.warning(typeof res === 'string' ? res : '操作完成')
      fetchPage(page.value, false)
    }
  } catch (e) {
    toast.error('操作失败: ' + (e.message || ''))
  }
}

async function executeNow(task) {
  try {
    const res = await ExecuteCronTaskNow(task.id)
    toast.success(typeof res === 'string' ? res : '任务执行中')
  } catch (e) {
    toast.error('执行失败: ' + (e.message || ''))
  }
}

async function removeTask(task) {
  if (!confirm(`确定删除任务「${task.name}」吗？`)) return
  try {
    const res = await DeleteCronTask(task.id)
    toast.success(typeof res === 'string' ? res : '删除成功')
    fetchPage(1, false)
  } catch (e) {
    toast.error('删除失败: ' + (e.message || ''))
  }
}

// ===== 新建 / 编辑 =====
const editVisible = ref(false)
const editSaving = ref(false)
const editForm = ref(emptyForm())
const editTitle = computed(() => (editForm.value.id ? '编辑任务' : '新建任务'))
const nextRuns = ref([])         // 下次执行时间预览
const validating = ref(false)

// 股票/市场分析参数
const stockParams = ref({ promptId: null, aiConfigId: null, sysPromptId: null, thinking: true, stockCode: '', stockName: '', agentMode: '' })
const marketParams = ref({ promptId: null, aiConfigId: null, sysPromptId: null, thinking: true, agentMode: '' })

function emptyForm() {
  return { id: 0, name: '', cronExpr: '', taskType: 'market_analysis', target: '', params: '', enable: true, status: 'active', description: '' }
}

function resetParams() {
  stockParams.value = { promptId: null, aiConfigId: null, sysPromptId: null, thinking: true, stockCode: '', stockName: '', agentMode: '' }
  marketParams.value = { promptId: null, aiConfigId: null, sysPromptId: null, thinking: true, agentMode: '' }
}

function openCreate() {
  editForm.value = emptyForm()
  resetParams()
  nextRuns.value = []
  editVisible.value = true
}

async function openEdit(task) {
  try {
    const full = await GetCronTaskByID(task.id)
    const t = full || task
    editForm.value = {
      id: t.id,
      name: t.name,
      cronExpr: (t.cronExpr ?? t.CronExpr ?? '').trim(),
      taskType: t.taskType,
      target: t.target || '',
      params: t.params || '',
      enable: t.enable,
      status: t.status || 'active',
      description: t.description || '',
    }
    resetParams()
    // 回填参数
    if (t.params) {
      try {
        const p = JSON.parse(t.params)
        if (t.taskType === 'stock_analysis') {
          stockParams.value = {
            promptId: p.promptId ?? null,
            aiConfigId: p.aiConfigId ?? null,
            sysPromptId: p.sysPromptId ?? null,
            thinking: p.thinking ?? true,
            stockCode: p.stockCode || '',
            stockName: p.stockName || '',
            agentMode: p.agentMode || '',
          }
        } else if (t.taskType === 'market_analysis') {
          marketParams.value = {
            promptId: p.promptId ?? null,
            aiConfigId: p.aiConfigId ?? null,
            sysPromptId: p.sysPromptId ?? null,
            thinking: p.thinking ?? true,
            agentMode: p.agentMode || '',
          }
        }
      } catch (_) {}
    }
    previewNextRuns()
    editVisible.value = true
  } catch (e) {
    toast.error('获取任务详情失败: ' + (e.message || ''))
  }
}

// 根据任务类型组装 params JSON
function buildParams() {
  if (editForm.value.taskType === 'stock_analysis') {
    const p = stockParams.value
    return JSON.stringify({
      promptId: p.promptId, aiConfigId: p.aiConfigId, sysPromptId: p.sysPromptId,
      thinking: p.thinking, stockCode: p.stockCode, stockName: p.stockName, agentMode: p.agentMode,
    })
  }
  if (editForm.value.taskType === 'market_analysis') {
    const p = marketParams.value
    return JSON.stringify({
      promptId: p.promptId, aiConfigId: p.aiConfigId, sysPromptId: p.sysPromptId,
      thinking: p.thinking, agentMode: p.agentMode,
    })
  }
  return editForm.value.params || ''
}

// 预览下次执行时间（同时起到校验作用）
async function previewNextRuns() {
  const expr = editForm.value.cronExpr.trim()
  if (!expr) { nextRuns.value = []; return }
  try {
    const res = await CalculateNextRunTimes(expr, 5)
    nextRuns.value = Array.isArray(res) ? res : []
  } catch (_) {
    nextRuns.value = []
  }
}

// 输入 cron 后防抖预览
let previewTimer = null
watch(() => editForm.value.cronExpr, () => {
  clearTimeout(previewTimer)
  previewTimer = setTimeout(previewNextRuns, 400)
})

// 校验间隔 ≥ 60 秒（对齐桌面端）
function checkInterval(times) {
  if (!Array.isArray(times) || times.length < 2) return true
  let min = Infinity
  for (let i = 1; i < times.length; i++) {
    const prev = new Date(times[i - 1]).getTime()
    const curr = new Date(times[i]).getTime()
    if (!Number.isNaN(prev) && !Number.isNaN(curr)) min = Math.min(min, (curr - prev) / 1000)
  }
  return min === Infinity || min >= 60
}

async function saveTask() {
  const f = editForm.value
  if (!f.name.trim()) { toast.warning('请输入任务名称'); return }
  if (!f.cronExpr.trim()) { toast.warning('请输入 Cron 表达式'); return }
  if (!f.taskType) { toast.warning('请选择任务类型'); return }

  editSaving.value = true
  try {
    // 校验表达式
    const vr = await ValidateCronExpr(f.cronExpr.trim())
    if (typeof vr === 'string' && !vr.includes('有效')) {
      toast.error(vr)
      return
    }
    // 校验间隔
    if (!checkInterval(nextRuns.value)) {
      toast.warning('两次执行间隔过短，请设置为至少 60 秒')
      return
    }
    const payload = { ...f, params: buildParams() }
    const res = f.id ? await UpdateCronTask(payload) : await CreateCronTask(payload)
    if (typeof res === 'string' && res.includes('成功')) {
      toast.success(res)
      editVisible.value = false
      fetchPage(1, false)
    } else {
      toast.error(typeof res === 'string' ? res : '保存失败')
    }
  } catch (e) {
    toast.error('保存失败: ' + (e.message || ''))
  } finally {
    editSaving.value = false
  }
}

function fmtTime(s) {
  if (!s) return ''
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) return String(s)
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const statusClass = computed(() => s => ({
  active: 'st--active', paused: 'st--paused', error: 'st--error',
}[s] || 'st--paused'))

const isStockType = computed(() => editForm.value.taskType === 'stock_analysis')
const isMarketType = computed(() => editForm.value.taskType === 'market_analysis')
const isCustomParams = computed(() => !isStockType.value && !isMarketType.value)

onBeforeMount(() => {
  loadOptions()
  fetchPage(1, false)
})
</script>

<template>
  <div class="page">
    <!-- 操作条 -->
    <div class="op-bar">
      <button
        type="button"
        class="filter-trigger"
        :class="{ 'filter-trigger--active': hasFilter }"
        @click="filterVisible = true"
      >
        <MIcon name="search" :size="16" /> <span class="filter-text">{{ hasFilter ? '已筛选' : '搜索任务' }}</span>
      </button>
      <MButton type="primary" size="small" @click="openCreate">新建任务</MButton>
    </div>

    <MPullRefresh class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="list.length" class="task-list">
          <div v-for="task in list" :key="task.id" class="task-card">
            <div class="task-head">
              <div class="task-name-wrap">
                <span class="task-type-tag">{{ taskTypeLabel(task.taskType) }}</span>
                <span class="task-name">{{ task.name }}</span>
              </div>
              <span class="task-status" :class="statusClass(task.status)">{{ task.status }}</span>
            </div>

            <div class="task-info">
              <span class="task-cron"><MIcon name="alarm" :size="16" /> {{ task.cronExpr }}</span>
              <span class="task-runs">运行 {{ task.runCount || 0 }} 次</span>
            </div>
            <div v-if="task.lastRunAt" class="task-last">
              最近：{{ fmtTime(task.lastRunAt) }}
              <span v-if="task.lastRunResult" class="task-result" :class="{ 'task-result--ok': String(task.lastRunResult).startsWith('成功') }">
                {{ task.lastRunResult }}
              </span>
            </div>

            <div class="task-foot">
              <button type="button" class="act-btn" @click="executeNow(task)">执行</button>
              <button type="button" class="act-btn" @click="toggleEnable(task)">{{ task.enable ? '暂停' : '启用' }}</button>
              <button type="button" class="act-btn" @click="openEdit(task)">编辑</button>
              <button type="button" class="act-btn act-btn--danger" @click="removeTask(task)">删除</button>
              <span class="task-enable" :class="{ 'task-enable--on': task.enable }">{{ task.enable ? '已启用' : '已停用' }}</span>
            </div>
          </div>

          <div class="load-more">
            <button v-if="hasMore" type="button" class="load-more-btn" :disabled="loading" @click="loadMore">
              {{ loading ? '加载中...' : '加载更多' }}
            </button>
            <span v-else class="load-more-end">共 {{ total }} 个任务</span>
          </div>
        </div>

        <MLoading v-else-if="loading" text="加载中..." vertical class="page-loading" />
        <MEmpty v-else description="还没有定时任务">
          <MButton type="primary" @click="openCreate">新建任务</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 搜索抽屉 -->
    <MSheet v-model:show="filterVisible" title="搜索任务" height="auto">
      <div class="form">
        <label class="field">
          <span class="field-label">任务名称</span>
          <input v-model="search.name" class="field-input" type="text" placeholder="按名称筛选" />
        </label>
        <div class="field">
          <span class="field-label">任务类型</span>
          <div class="chip-row">
            <button type="button" class="chip" :class="{ 'chip--active': search.taskType === '' }" @click="search.taskType = ''">全部</button>
            <button
              v-for="opt in taskTypeOptions"
              :key="opt.value"
              type="button"
              class="chip"
              :class="{ 'chip--active': search.taskType === opt.value }"
              @click="search.taskType = opt.value"
            >{{ opt.label }}</button>
          </div>
        </div>
        <div class="field">
          <span class="field-label">任务状态</span>
          <div class="chip-row">
            <button
              v-for="opt in STATUS_FILTERS"
              :key="opt.value"
              type="button"
              class="chip"
              :class="{ 'chip--active': search.status === opt.value }"
              @click="search.status = opt.value"
            >{{ opt.label }}</button>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="resetSearch">重置</MButton>
          <MButton type="primary" block @click="applySearch">搜索</MButton>
        </div>
      </template>
    </MSheet>

    <!-- 新建/编辑抽屉 -->
    <MSheet v-model:show="editVisible" :title="editTitle">
      <div class="form">
        <label class="field">
          <span class="field-label">任务名称 <i class="req">*</i></span>
          <input v-model="editForm.name" class="field-input" type="text" placeholder="请输入任务名称" />
        </label>

        <div class="field">
          <span class="field-label">任务类型 <i class="req">*</i></span>
          <div class="chip-row">
            <button
              v-for="opt in taskTypeOptions"
              :key="opt.value"
              type="button"
              class="chip"
              :class="{ 'chip--active': editForm.taskType === opt.value }"
              @click="editForm.taskType = opt.value"
            >{{ opt.label }}</button>
          </div>
        </div>

        <label class="field">
          <span class="field-label">Cron 表达式 <i class="req">*</i></span>
          <input v-model="editForm.cronExpr" class="field-input mono" type="text" placeholder="如：0 0 9 * * ? （秒 分 时 日 月 周）" />
          <div v-if="nextRuns.length" class="next-runs">
            <span class="next-label">未来执行：</span>
            <span v-for="(t, i) in nextRuns" :key="i" class="next-time">{{ fmtTime(t) }}</span>
          </div>
        </label>

        <!-- 股票分析参数 -->
        <template v-if="isStockType">
          <div class="param-group">
            <div class="param-row">
              <span class="param-label">提示词模板</span>
              <select v-model="stockParams.promptId" class="field-select">
                <option :value="null">默认</option>
                <option v-for="o in userPromptOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
              </select>
            </div>
            <div class="param-row">
              <span class="param-label">AI 配置</span>
              <select v-model="stockParams.aiConfigId" class="field-select">
                <option :value="null">请选择</option>
                <option v-for="o in aiConfigOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
              </select>
            </div>
            <div class="param-row">
              <span class="param-label">系统提示词</span>
              <select v-model="stockParams.sysPromptId" class="field-select">
                <option :value="null">默认</option>
                <option v-for="o in sysPromptOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
              </select>
            </div>
            <div class="param-row">
              <span class="param-label">Agent 模式</span>
              <select v-model="stockParams.agentMode" class="field-select">
                <option v-for="m in AGENT_MODES" :key="m.value" :value="m.value">{{ m.label }}</option>
              </select>
            </div>
            <label class="field">
              <span class="field-label">股票代码</span>
              <input v-model="stockParams.stockCode" class="field-input" type="text" placeholder="如：600519" />
            </label>
            <label class="field">
              <span class="field-label">股票名称</span>
              <input v-model="stockParams.stockName" class="field-input" type="text" placeholder="如：贵州茅台" />
            </label>
            <button type="button" class="switch-line" @click="stockParams.thinking = !stockParams.thinking">
              <span>启用思考</span>
              <span class="switch" :class="{ 'switch--on': stockParams.thinking }"><span class="switch-dot" /></span>
            </button>
          </div>
        </template>

        <!-- 市场分析参数 -->
        <template v-else-if="isMarketType">
          <div class="param-group">
            <div class="param-row">
              <span class="param-label">提示词模板</span>
              <select v-model="marketParams.promptId" class="field-select">
                <option :value="null">默认</option>
                <option v-for="o in userPromptOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
              </select>
            </div>
            <div class="param-row">
              <span class="param-label">AI 配置</span>
              <select v-model="marketParams.aiConfigId" class="field-select">
                <option :value="null">请选择</option>
                <option v-for="o in aiConfigOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
              </select>
            </div>
            <div class="param-row">
              <span class="param-label">系统提示词</span>
              <select v-model="marketParams.sysPromptId" class="field-select">
                <option :value="null">默认</option>
                <option v-for="o in sysPromptOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
              </select>
            </div>
            <div class="param-row">
              <span class="param-label">Agent 模式</span>
              <select v-model="marketParams.agentMode" class="field-select">
                <option v-for="m in AGENT_MODES" :key="m.value" :value="m.value">{{ m.label }}</option>
              </select>
            </div>
            <button type="button" class="switch-line" @click="marketParams.thinking = !marketParams.thinking">
              <span>启用思考</span>
              <span class="switch" :class="{ 'switch--on': marketParams.thinking }"><span class="switch-dot" /></span>
            </button>
          </div>
        </template>

        <!-- 其他类型：JSON 参数 -->
        <label v-else class="field">
          <span class="field-label">任务参数</span>
          <textarea v-model="editForm.params" class="field-textarea mono" rows="4" placeholder='JSON 格式，如：{"stock_codes":["600519"]}' />
        </label>

        <label class="field">
          <span class="field-label">任务描述</span>
          <textarea v-model="editForm.description" class="field-textarea" rows="2" placeholder="可选" />
        </label>

        <button type="button" class="switch-line" @click="editForm.enable = !editForm.enable">
          <span>启用任务</span>
          <span class="switch" :class="{ 'switch--on': editForm.enable }"><span class="switch-dot" /></span>
        </button>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="editVisible = false">取消</MButton>
          <MButton type="primary" block :loading="editSaving" @click="saveTask">保存</MButton>
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
  align-items: center;
  gap: var(--m-space-sm);
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.filter-trigger {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--m-space-xs);
  padding: var(--m-space-sm) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  color: var(--m-text-tertiary);
  font: inherit;
  font-size: var(--m-font-sm);
}

.filter-trigger--active {
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
  background: var(--m-color-rise-light);
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

/* 任务卡片 */
.task-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.task-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.task-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-sm);
}

.task-name-wrap {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  min-width: 0;
}

.task-type-tag {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  color: #2080f0;
  background: rgba(32, 128, 240, 0.12);
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
}

.task-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-status {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
}

.st--active { background: var(--m-color-rise-light); color: var(--m-color-rise); }
.st--paused { background: var(--m-divider-color); color: var(--m-text-tertiary); }
.st--error { background: var(--m-color-fall-light); color: var(--m-color-fall); }

.task-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.task-cron {
  font-family: monospace;
}

.task-runs {
  color: var(--m-text-tertiary);
}

.task-last {
  margin-top: var(--m-space-xs);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.task-result {
  color: var(--m-color-fall);
}

.task-result--ok {
  color: var(--m-color-rise);
}

.task-foot {
  display: flex;
  align-items: center;
  gap: var(--m-space-xs);
  margin-top: var(--m-space-sm);
  padding-top: var(--m-space-sm);
  border-top: 1px solid var(--m-divider-color);
}

.act-btn {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.act-btn:active {
  transform: scale(0.95);
}

.act-btn--danger {
  color: var(--m-color-rise);
  border-color: var(--m-color-rise);
}

.task-enable {
  margin-left: auto;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.task-enable--on {
  color: var(--m-color-rise);
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

.field-label {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.req {
  color: var(--m-color-rise);
  font-style: normal;
}

.field-input,
.field-textarea,
.field-select {
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

.mono {
  font-family: monospace;
}

.field-input:focus,
.field-textarea:focus,
.field-select:focus {
  border-color: var(--m-color-rise);
}

.next-runs {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
  margin-top: var(--m-space-xs);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.next-label {
  color: var(--m-color-fall);
}

.next-time {
  font-variant-numeric: tabular-nums;
}

/* 参数组 */
.param-group {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
}

.param-row {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.param-label {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.param-row .field-select {
  background: var(--m-bg-card);
}

/* chip */
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

/* 开关 */
.switch-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: var(--m-space-sm) 0;
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
}

.switch {
  position: relative;
  width: 44px;
  height: 24px;
  border-radius: var(--m-radius-full);
  background: var(--m-border-color);
  transition: background var(--m-duration-fast);
  flex-shrink: 0;
}

.switch--on {
  background: var(--m-color-rise);
}

.switch-dot {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
  transition: transform var(--m-duration-fast);
}

.switch--on .switch-dot {
  transform: translateX(20px);
}

.sheet-footer {
  display: flex;
  gap: var(--m-space-md);
}
</style>
