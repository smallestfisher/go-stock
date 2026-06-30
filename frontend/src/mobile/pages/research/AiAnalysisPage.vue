<script setup>
import { ref, computed, nextTick, onBeforeMount, onBeforeUnmount } from 'vue'
// md-editor-v3 的轻量预览组件（仅 preview.css，不引入编辑器体积）。移动端首次引入。
import 'md-editor-v3/lib/preview.css'
import { MdPreview } from 'md-editor-v3'
import {
  GetAiConfigs,
  GetPromptTemplates,
  GetAIResponseResult,
  GetAIResponseResultList,
  SummaryStockNews,
  AbortSummaryStockNews,
  SaveAIResponseResult,
  DeleteAIResponseResult,
} from '../../../api/app'
import { EventsOn, EventsOff } from '../../../api/runtime'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MSheet from '../../components/base/MSheet.vue'
import MLoading from '../../components/base/MLoading.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 对齐桌面端 market.vue 的「AI市场资讯总结」(getAiSummary/reAiSummary/summaryStockNews 事件流)
// + researchReport.vue 的历史列表。移动端精简：仅模型选择，提示词/工具/思考用默认。
const SCOPE = '市场资讯'

const aiSummary = ref('')
const aiSummaryTime = ref('')
const modelName = ref('')
const question = ref('')
const chatId = ref('')
// 本次实际分析的题目：供结果区「分析主题」展示，与输入框解耦。
// 这样分析完成后清空输入框，也不会让结果标题丢失。
const analyzedTopic = ref('')
const loading = ref(false)
const analysisStatus = ref('')

// AI 模型配置
const aiConfigs = ref([])
const aiConfigId = ref(null)

// 提示词模板：系统提示词 / 用户提示词（对齐桌面 market.vue:120-123）
const sysPromptOptions = ref([])   // [{ID,name,content}]
const userPromptOptions = ref([])  // [{ID,name,content}]，选中后 content 填入 question
const sysPromptId = ref(null)

// 工具调用 / 思考模式开关（对齐桌面 enableTools/thinkingMode，默认 true）
const enableTools = ref(true)
const thinkingMode = ref(true)

// 用户提示词预设选择（选中即把 content 填入输入框，与桌面 question 下拉同义）
const selectedUserPromptId = ref(null)

// 历史
const history = ref([])
const historyLoading = ref(false)
const detailVisible = ref(false)
const detailItem = ref(null)

const resultScrollRef = ref(null)

// AI 模型配置加载
async function loadAiConfigs() {
  try {
    const res = await GetAiConfigs()
    aiConfigs.value = Array.isArray(res) ? res : []
    if (aiConfigs.value.length && !aiConfigId.value) {
      aiConfigId.value = aiConfigs.value[0].ID
    }
  } catch (e) {
    console.error('加载AI配置失败:', e)
  }
}

// 加载提示词模板，按 type 分系统/用户两组（对齐桌面 market.vue:120-123）
async function loadPromptTemplates() {
  try {
    const res = await GetPromptTemplates('', '')
    const list = Array.isArray(res) ? res : []
    sysPromptOptions.value = list.filter(t => t.type === '模型系统Prompt')
    userPromptOptions.value = list.filter(t => t.type === '模型用户Prompt')
  } catch (e) {
    console.error('加载提示词失败:', e)
  }
}

// 选择用户提示词预设：把 content 填入 question 输入框（与桌面 question 下拉 value-field=content 同义）
function applyUserPrompt(id) {
  const tpl = userPromptOptions.value.find(t => t.ID === id)
  if (tpl && tpl.content) question.value = tpl.content
}

// 清空问题输入与预设选择
function clearQuestion() {
  question.value = ''
  selectedUserPromptId.value = null
}

// 读取最新一条总结（首次进入，不自动生成，省 token）
async function loadLatest() {
  try {
    const result = await GetAIResponseResult(SCOPE)
    if (result && result.content) {
      aiSummary.value = result.content
      // 仅把上次题目存入 analyzedTopic 供结果区展示，不回填输入框；
      // 输入框保持空白，进入即可输入新问题。需要再生成同一题可点「重新总结」。
      analyzedTopic.value = result.question || ''
      modelName.value = result.modelName || ''
      aiSummaryTime.value = fmtTime(result.CreatedAt)
    }
  } catch (e) {
    console.error('读取AI总结失败:', e)
  }
}

// 历史列表
async function loadHistory() {
  historyLoading.value = true
  try {
    const res = await GetAIResponseResultList({ page: 1, pageSize: 20 })
    history.value = (res && Array.isArray(res.list) ? res.list : []).filter(it => it && typeof it === 'object')
  } catch (e) {
    console.error('加载历史失败:', e)
    history.value = []
  } finally {
    historyLoading.value = false
  }
}

// 触发生成（对齐桌面 reAiSummary）
function startSummary() {
  if (loading.value) return
  if (!aiConfigId.value) {
    analysisStatus.value = '请先选择 AI 模型'
    return
  }
  // 输入框为空时（如直接点「重新总结」）沿用上次题目，保证再生成不丢问题；
  // 有新输入则以其为准，并记入 analyzedTopic 供结果区展示。
  const q = question.value.trim() ? question.value : analyzedTopic.value
  analyzedTopic.value = q.trim() || analyzedTopic.value
  aiSummary.value = ''
  loading.value = true
  analysisStatus.value = '正在连接AI服务...'
  // SummaryStockNews(question, aiConfigId, sysPromptId, enableTools, thinkingMode, eventName, stockCode)
  SummaryStockNews(q, aiConfigId.value, sysPromptId.value, enableTools.value, thinkingMode.value, 'summaryStockNews', '')
}

// 当前选中的模型名（生成中展示「正在用什么模型分析什么问题」）
const activeModelName = computed(() => {
  const cfg = aiConfigs.value.find(c => c.ID === aiConfigId.value)
  return cfg ? cfg.name : ''
})
// 当前分析主题摘要（用于生成中/结果区告知用户在分析什么）
const questionPreview = computed(() => {
  // 读 analyzedTopic（本次实际分析的题目），与输入框解耦：
  // 分析完成后清空输入框，结果标题仍保留。
  const q = (analyzedTopic.value || '').trim()
  if (q) return q
  return '市场资讯综合分析'
})

// 中止生成
function stopSummary() {
  AbortSummaryStockNews()
}

// 流式事件：对齐桌面 market.vue summaryStockNews handler
EventsOn('summaryStockNews', async (msg) => {
  if (msg === 'DONE') {
    await SaveAIResponseResult(SCOPE, SCOPE, aiSummary.value, chatId.value, analyzedTopic.value, aiConfigId.value)
    loading.value = false
    analysisStatus.value = '分析完成'
    setTimeout(() => { analysisStatus.value = '' }, 2000)
    loadHistory()
    // 分析完成后清空输入框：本次题目已存入 analyzedTopic 供结果区展示，
    // 输入框回归空白，可直接输入下一个问题，无需手动删除。
    question.value = ''
    selectedUserPromptId.value = null
  } else {
    if (msg.chatId) chatId.value = msg.chatId
    // 后端回传实际分析题目，只更新结果区展示，不再回填输入框
    if (msg.question) analyzedTopic.value = msg.question
    if (msg.content || msg.reasoning_content || msg.extraContent) {
      if (!aiSummary.value) analysisStatus.value = 'AI正在分析中...'
    }
    if (msg.content) aiSummary.value += msg.content
    if (msg.reasoning_content) aiSummary.value += msg.reasoning_content
    if (msg.extraContent) aiSummary.value += msg.extraContent
    if (msg.model) modelName.value = msg.model
    if (msg.time) aiSummaryTime.value = msg.time
    scrollToBottom()
  }
})

function scrollToBottom() {
  nextTick(() => {
    requestAnimationFrame(() => {
      const el = resultScrollRef.value
      if (el) el.scrollTop = el.scrollHeight
    })
  })
}

function fmtTime(s) {
  if (!s) return ''
  return String(s).substring(0, 19).replace('T', ' ')
}

function fmtQuestion(s) {
  if (!s) return SCOPE
  const t = String(s).trim()
  return t.length > 40 ? t.slice(0, 40) + '…' : t
}

// 查看历史详情
function openDetail(item) {
  detailItem.value = item
  detailVisible.value = true
}

// 删除历史
async function deleteDetail(id) {
  try {
    await DeleteAIResponseResult(id)
    history.value = history.value.filter(it => it.ID !== id)
    detailVisible.value = false
  } catch (e) {
    console.error('删除失败:', e)
  }
}

async function handleRefresh() {
  await Promise.all([loadLatest(), loadHistory()])
}

onBeforeMount(() => {
  loadAiConfigs()
  loadPromptTemplates()
  loadLatest()
  loadHistory()
})

onBeforeUnmount(() => {
  EventsOff('summaryStockNews')
})
</script>

<template>
  <div class="page">
    <!-- AI 总结操作区 -->
    <div class="op-bar">
      <div class="op-row">
        <h3 class="op-title">🤖 AI 市场资讯总结</h3>
        <button
          v-if="!loading"
          class="op-btn op-btn--primary"
          type="button"
          @click="startSummary"
        >
          {{ aiSummary ? '重新总结' : '生成总结' }}
        </button>
        <button v-else class="op-btn op-btn--danger" type="button" @click="stopSummary">停止</button>
      </div>

      <!-- 问题输入：决定「生成什么」的核心 -->
      <div class="question-row">
        <textarea
          v-model="question"
          class="question-input"
          rows="2"
          placeholder="输入要分析的问题，例如：总结和分析股票市场新闻中的投资机会"
        />
        <button v-if="question" class="clear-question" type="button" @click="clearQuestion">✕</button>
      </div>

      <!-- 用户提示词预设（选中填入输入框，与桌面 question 下拉同义） -->
      <div v-if="userPromptOptions.length" class="preset-block">
        <span class="preset-label">预设问题</span>
        <div class="chip-scroll">
          <button
            v-for="t in userPromptOptions"
            :key="t.ID"
            type="button"
            class="preset-chip"
            :class="{ 'preset-chip--active': selectedUserPromptId === t.ID }"
            @click="selectedUserPromptId = t.ID; applyUserPrompt(t.ID)"
          >
            {{ t.name }}
          </button>
        </div>
      </div>

      <!-- 模型选择（横向 chip） -->
      <div v-if="aiConfigs.length" class="preset-block">
        <span class="preset-label">模型</span>
        <div class="chip-scroll">
          <button
            v-for="cfg in aiConfigs"
            :key="cfg.ID"
            type="button"
            class="preset-chip"
            :class="{ 'preset-chip--active': aiConfigId === cfg.ID }"
            @click="aiConfigId = cfg.ID"
          >
            {{ cfg.name }}
          </button>
        </div>
      </div>

      <!-- 系统提示词（横向 chip） -->
      <div v-if="sysPromptOptions.length" class="preset-block">
        <span class="preset-label">系统提示词</span>
        <div class="chip-scroll">
          <button
            type="button"
            class="preset-chip"
            :class="{ 'preset-chip--active': sysPromptId === null }"
            @click="sysPromptId = null"
          >
            默认
          </button>
          <button
            v-for="t in sysPromptOptions"
            :key="t.ID"
            type="button"
            class="preset-chip"
            :class="{ 'preset-chip--active': sysPromptId === t.ID }"
            @click="sysPromptId = t.ID"
          >
            {{ t.name }}
          </button>
        </div>
      </div>

      <!-- 工具调用 / 思考模式开关（对齐桌面） -->
      <div class="switch-row">
        <button
          type="button"
          class="toggle-btn"
          :class="{ 'toggle-btn--on': enableTools }"
          @click="enableTools = !enableTools"
        >
          工具调用{{ enableTools ? ' ✓' : '' }}
        </button>
        <button
          type="button"
          class="toggle-btn"
          :class="{ 'toggle-btn--on': thinkingMode }"
          @click="thinkingMode = !thinkingMode"
        >
          思考模式{{ thinkingMode ? ' ✓' : '' }}
        </button>
      </div>
    </div>

    <MPullRefresh class="ai-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <!-- 总结结果 -->
        <div class="result-card">
          <!-- 生成中：告知用户在分析什么 -->
          <div v-if="loading && !aiSummary" class="result-loading">
            <MLoading :text="analysisStatus || 'AI 正在分析...'" vertical />
            <p class="loading-topic">正在分析：<b>{{ questionPreview }}</b></p>
            <p v-if="activeModelName" class="loading-model">模型：{{ activeModelName }}</p>
          </div>
          <template v-else>
            <div v-if="aiSummary" class="result-wrap">
              <div class="result-topic">
                <span class="result-topic-label">分析主题</span>
                <span class="result-topic-text">{{ questionPreview }}</span>
              </div>
              <div ref="resultScrollRef" class="result-body">
                <MdPreview :modelValue="aiSummary" theme="light" />
              </div>
            </div>
            <MEmpty v-else description="暂无 AI 总结，输入问题后点击「生成总结」" />
          </template>
        </div>

        <!-- 状态条 -->
        <div v-if="aiSummary && (modelName || aiSummaryTime)" class="meta-bar">
          <span v-if="modelName" class="meta-model">{{ modelName }}</span>
          <span v-if="aiSummaryTime" class="meta-time">{{ aiSummaryTime }}</span>
          <span v-if="analysisStatus" class="meta-status">{{ analysisStatus }}</span>
        </div>

        <!-- 历史分析 -->
        <div class="history-section">
          <h4 class="section-title">历史分析</h4>
          <div v-if="history.length" class="history-list">
            <div
              v-for="item in history"
              :key="item.ID"
              class="history-item"
              @click="openDetail(item)"
            >
              <div class="history-head">
                <span class="history-name">{{ item.stockName || SCOPE }}</span>
                <span class="history-time">{{ fmtTime(item.CreatedAt) }}</span>
              </div>
              <div class="history-question">{{ fmtQuestion(item.question) }}</div>
              <div class="history-foot">
                <span v-if="item.modelName" class="history-model">{{ item.modelName }}</span>
              </div>
            </div>
          </div>
          <MEmpty v-else-if="!historyLoading" description="暂无历史分析记录" />
        </div>

        <p class="tip">*AI 分析结果仅供参考，请以实际行情为准。投资需谨慎，风险自担。</p>
      </div>
    </MPullRefresh>

    <!-- 历史详情抽屉 -->
    <MSheet v-model:show="detailVisible" :title="detailItem?.stockName || SCOPE">
      <div v-if="detailItem" class="detail-content">
        <div class="detail-meta">
          <span v-if="detailItem.modelName" class="detail-model">{{ detailItem.modelName }}</span>
          <span class="detail-time">{{ fmtTime(detailItem.CreatedAt) }}</span>
        </div>
        <div v-if="detailItem.question" class="detail-question">{{ detailItem.question }}</div>
        <div class="detail-body">
          <MdPreview :modelValue="detailItem.content || ''" theme="light" />
        </div>
        <button
          v-if="detailItem.ID"
          class="detail-delete"
          type="button"
          @click="deleteDetail(detailItem.ID)"
        >
          删除此记录
        </button>
      </div>
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

/* 操作区 */
.op-bar {
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.op-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-md);
}

.op-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.op-btn {
  padding: var(--m-space-xs) var(--m-space-lg);
  border: none;
  border-radius: var(--m-radius-sm);
  font: inherit;
  font-size: var(--m-font-sm);
  color: #fff;
  flex-shrink: 0;
}

.op-btn--primary {
  background: var(--m-color-rise);
}

.op-btn--danger {
  background: var(--m-color-fall);
}

.op-btn:active {
  transform: scale(0.95);
}

/* 问题输入框 */
.question-row {
  position: relative;
}

.question-input {
  width: 100%;
  padding: var(--m-space-sm) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font: inherit;
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-normal);
  resize: vertical;
  outline: none;
}

.question-input:focus {
  border-color: var(--m-color-rise);
}

.clear-question {
  position: absolute;
  top: var(--m-space-xs);
  right: var(--m-space-xs);
  width: 22px;
  height: 22px;
  border: none;
  border-radius: var(--m-radius-full);
  background: var(--m-bg-card);
  color: var(--m-text-tertiary);
  font-size: var(--m-font-xs);
  line-height: 1;
}

/* 预设块（标签 + 横向滚动 chip） */
.preset-block {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.preset-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.chip-scroll {
  display: flex;
  gap: var(--m-space-xs);
  overflow-x: auto;
  scrollbar-width: none;
}

.chip-scroll::-webkit-scrollbar {
  display: none;
}

.preset-chip {
  flex-shrink: 0;
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
  white-space: nowrap;
}

.preset-chip--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

/* 工具/思考开关 */
.switch-row {
  display: flex;
  gap: var(--m-space-sm);
}

.toggle-btn {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.toggle-btn--on {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.ai-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

/* 结果卡片 */
.result-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
}

.result-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--m-space-xs);
  padding: var(--m-space-2xl) 0;
}

.loading-topic {
  margin-top: var(--m-space-sm);
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  text-align: center;
}

.loading-topic b {
  color: var(--m-text-primary);
}

.loading-model {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.result-wrap {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

/* 结果区顶部：告知分析的主题 */
.result-topic {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: var(--m-space-sm);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
}

.result-topic-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.result-topic-text {
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
  line-height: var(--m-line-height-normal);
}

.result-body {
  max-height: 50vh;
  overflow-y: auto;
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-loose);
}

.meta-bar {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  flex-wrap: wrap;
  padding: 0 var(--m-space-xs);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.meta-model {
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.meta-status {
  color: var(--m-color-fall);
}

/* 历史区 */
.history-section {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.section-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.history-item {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
}

.history-item:active {
  background: var(--m-bg-primary);
}

.history-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-md);
}

.history-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.history-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.history-question {
  margin-top: var(--m-space-xs);
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.history-foot {
  margin-top: var(--m-space-xs);
}

.history-model {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  padding: 1px var(--m-space-sm);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
}

.tip {
  font-size: var(--m-font-xs);
  color: var(--m-color-rise);
  line-height: var(--m-line-height-normal);
}

/* 详情抽屉 */
.detail-content {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.detail-model {
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.detail-question {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  line-height: var(--m-line-height-normal);
}

.detail-body {
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-loose);
}

.detail-delete {
  align-self: flex-start;
  padding: var(--m-space-xs) var(--m-space-lg);
  border: 1px solid var(--m-color-fall);
  border-radius: var(--m-radius-sm);
  background: transparent;
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-color-fall);
}
</style>
