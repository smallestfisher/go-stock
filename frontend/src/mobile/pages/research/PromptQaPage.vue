<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import 'md-editor-v3/lib/preview.css'
import { MdPreview } from 'md-editor-v3'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MButton from '../../components/base/MButton.vue'
import MSheet from '../../components/base/MSheet.vue'
import { toast } from '../../composables/useToast'
import { usePromptPlaza } from '../../composables/usePromptPlaza'

// 对齐桌面端 promptQa.vue：问题列表 + 搜索 + 已解决筛选 + 提问 + 详情(回答/点赞/采纳/删除)
// 复用提示词广场同一套账号体系（promptPlazaToken）

const {
  currentUser,
  isLoggedIn,
  apiGet,
  apiPost,
  apiDelete,
  fetchCurrentUser,
  displayName,
  timeAgo,
  formatTime,
} = usePromptPlaza()

const PAGE_SIZE = 15

const FILTER_OPTIONS = [
  { label: '全部', value: '' },
  { label: '待解决', value: 'false' },
  { label: '已解决', value: 'true' },
]

const keyword = ref('')
const resolvedFilter = ref('')
const loading = ref(false)
const questions = ref([])
const page = ref(1)
const total = ref(0)
const totalPages = ref(1)
// 问答接口可能未部署，单独提示（对齐桌面端 apiAvailable）
const apiAvailable = ref(true)

const hasMore = computed(() => page.value < totalPages.value)

async function loadQuestions(targetPage = 1, append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const params = { page: targetPage, pageSize: PAGE_SIZE }
    if (keyword.value) params.keyword = keyword.value
    if (resolvedFilter.value) params.resolved = resolvedFilter.value
    const data = await apiGet('/questions', params)
    apiAvailable.value = true
    const rows = data.list || []
    questions.value = append ? questions.value.concat(rows) : rows
    total.value = data.total || 0
    totalPages.value = Math.ceil((data.total || 0) / PAGE_SIZE) || 1
    page.value = targetPage
  } catch (e) {
    if (e.message && (e.message.includes('接口返回非JSON') || e.message.includes('404'))) {
      apiAvailable.value = false
    } else {
      toast.error('加载问题列表失败: ' + (e.message || ''))
    }
  } finally {
    loading.value = false
  }
}

function handleRefresh() {
  return loadQuestions(1, false)
}

function loadMore() {
  if (hasMore.value && !loading.value) loadQuestions(page.value + 1, true)
}

function applySearch() {
  loadQuestions(1, false)
}

function selectFilter(value) {
  resolvedFilter.value = value
  loadQuestions(1, false)
}

function requireLogin() {
  if (!isLoggedIn.value) {
    toast.warning('请先在「提示词广场」登录后再操作')
    return false
  }
  return true
}

// ===== 提问 =====
const askVisible = ref(false)
const askForm = ref({ title: '', content: '' })
const askLoading = ref(false)

function openAsk() {
  if (!requireLogin()) return
  askForm.value = { title: '', content: '' }
  askVisible.value = true
}

async function submitAsk() {
  if (!askForm.value.title || !askForm.value.content) {
    toast.warning('请填写标题和内容')
    return
  }
  askLoading.value = true
  try {
    await apiPost('/questions', { title: askForm.value.title, content: askForm.value.content })
    askVisible.value = false
    toast.success('提问成功')
    loadQuestions(1, false)
  } catch (e) {
    toast.error('提问失败: ' + (e.message || ''))
  } finally {
    askLoading.value = false
  }
}

// ===== 详情 + 回答 =====
const detailVisible = ref(false)
const detailQuestion = ref(null)
const answers = ref([])
const newAnswer = ref('')

async function openDetail(id) {
  try {
    const data = await apiGet(`/questions/${id}`)
    detailQuestion.value = data.question
    answers.value = data.answers || []
    newAnswer.value = ''
    detailVisible.value = true
  } catch (e) {
    toast.error('加载问题详情失败: ' + (e.message || ''))
  }
}

const isQuestionOwner = computed(
  () => !!(currentUser.value && detailQuestion.value && detailQuestion.value.userId === currentUser.value.id)
)

async function submitAnswer() {
  if (!requireLogin()) return
  if (!newAnswer.value.trim()) {
    toast.warning('请输入回答内容')
    return
  }
  try {
    const data = await apiPost(`/questions/${detailQuestion.value.id}/answers`, { content: newAnswer.value })
    newAnswer.value = ''
    detailQuestion.value.answersCount = (detailQuestion.value.answersCount || 0) + 1
    answers.value.push(data)
    toast.success('回答成功')
  } catch (e) {
    toast.error('回答失败: ' + (e.message || ''))
  }
}

async function acceptAnswer(answer) {
  if (!requireLogin()) return
  try {
    await apiPost(`/answers/${answer.id}/accept`)
    answers.value.forEach(a => { a.isAccepted = false })
    answer.isAccepted = true
    detailQuestion.value.isResolved = true
    // 同步列表中该问题状态
    const inList = questions.value.find(q => q.id === detailQuestion.value.id)
    if (inList) inList.isResolved = true
    toast.success('已采纳该回答')
  } catch (e) {
    toast.error('采纳失败: ' + (e.message || ''))
  }
}

async function likeAnswer(answer) {
  if (!requireLogin()) return
  try {
    const data = await apiPost(`/answers/${answer.id}/like`)
    answer.isLiked = data.isLiked
    answer.likesCount = data.likesCount
  } catch (e) {
    toast.error('操作失败: ' + (e.message || ''))
  }
}

async function deleteAnswer(answer) {
  if (!confirm('确定要删除这条回答吗？')) return
  try {
    await apiDelete(`/answers/${answer.id}`)
    answers.value = answers.value.filter(a => a.id !== answer.id)
    detailQuestion.value.answersCount = Math.max(0, (detailQuestion.value.answersCount || 1) - 1)
    toast.success('删除成功')
  } catch (e) {
    toast.error('删除失败: ' + (e.message || ''))
  }
}

async function deleteQuestion() {
  if (!confirm('确定要删除这个问题吗？所有回答也会被删除。')) return
  try {
    await apiDelete(`/questions/${detailQuestion.value.id}`)
    detailVisible.value = false
    toast.success('删除成功')
    loadQuestions(1, false)
  } catch (e) {
    toast.error('删除失败: ' + (e.message || ''))
  }
}

onBeforeMount(() => {
  loadQuestions(1, false)
  if (isLoggedIn.value) fetchCurrentUser()
})
</script>

<template>
  <div class="page">
    <!-- 顶部：搜索 + 筛选 + 账号 -->
    <div class="op-bar">
      <div class="search-row">
        <input
          v-model="keyword"
          class="search-input"
          type="text"
          placeholder="搜索问题..."
          @keyup.enter="applySearch"
        />
        <MButton type="primary" size="small" @click="openAsk">❓ 提问</MButton>
      </div>
      <div class="filter-row">
        <div class="chip-scroll">
          <button
            v-for="opt in FILTER_OPTIONS"
            :key="opt.value"
            type="button"
            class="chip"
            :class="{ 'chip--active': resolvedFilter === opt.value }"
            @click="selectFilter(opt.value)"
          >{{ opt.label }}</button>
        </div>
        <span v-if="isLoggedIn" class="user-tag">{{ displayName() }}</span>
        <span v-else class="login-hint">请在「广场」登录</span>
      </div>
    </div>

    <!-- 服务不可用提示 -->
    <div v-if="!apiAvailable" class="notice">
      问答广场接口未部署或服务未启动，请联系管理员部署最新版本的服务端程序。
    </div>

    <MPullRefresh v-else class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="questions.length" class="qa-list">
          <div v-for="item in questions" :key="item.id" class="qa-card" @click="openDetail(item.id)">
            <div class="qa-card-head">
              <span class="qa-tag" :class="item.isResolved ? 'qa-tag--resolved' : 'qa-tag--pending'">
                {{ item.isResolved ? '已解决' : '待解决' }}
              </span>
              <span class="qa-title">{{ item.title }}</span>
            </div>
            <div class="qa-card-foot">
              <span class="qa-author">{{ item.user?.nickname || item.user?.username || '匿名' }} · {{ timeAgo(item.createdAt) }}</span>
              <span class="qa-answers">💬 {{ item.answersCount || 0 }} 回答</span>
            </div>
          </div>

          <div class="load-more">
            <button v-if="hasMore" type="button" class="load-more-btn" :disabled="loading" @click="loadMore">
              {{ loading ? '加载中...' : '加载更多' }}
            </button>
            <span v-else class="load-more-end">共 {{ total }} 个问题</span>
          </div>
        </div>

        <MLoading v-else-if="loading" text="加载中..." vertical class="page-loading" />
        <MEmpty v-else description="暂无问题，快来提问吧">
          <MButton type="primary" @click="openAsk">我要提问</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 提问抽屉 -->
    <MSheet v-model:show="askVisible" title="提问" height="auto">
      <div class="form">
        <label class="field">
          <span class="field-label">问题标题 <i class="req">*</i></span>
          <input v-model="askForm.title" class="field-input" type="text" placeholder="一句话描述你的问题" />
        </label>
        <label class="field">
          <span class="field-label">问题描述 <i class="req">*</i></span>
          <textarea v-model="askForm.content" class="field-textarea" rows="6" placeholder="详细描述你的问题..." />
        </label>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="askVisible = false">取消</MButton>
          <MButton type="primary" block :loading="askLoading" @click="submitAsk">提交问题</MButton>
        </div>
      </template>
    </MSheet>

    <!-- 详情抽屉 -->
    <MSheet v-model:show="detailVisible" :title="detailQuestion?.title || '问题详情'">
      <div v-if="detailQuestion" class="detail">
        <div class="detail-meta">
          <span class="qa-tag" :class="detailQuestion.isResolved ? 'qa-tag--resolved' : 'qa-tag--pending'">
            {{ detailQuestion.isResolved ? '已解决' : '待解决' }}
          </span>
          <span class="detail-author">
            {{ detailQuestion.user?.nickname || detailQuestion.user?.username || '匿名' }} · {{ formatTime(detailQuestion.createdAt) }}
          </span>
          <button
            v-if="isQuestionOwner"
            type="button"
            class="del-question"
            @click="deleteQuestion"
          >🗑️ 删除</button>
        </div>

        <div class="detail-content">
          <MdPreview :modelValue="detailQuestion.content" theme="light" />
        </div>

        <!-- 回答区 -->
        <div class="answer-section">
          <div class="answer-title">{{ answers.length || 0 }} 个回答</div>

          <div class="answer-input-row">
            <textarea v-model="newAnswer" class="field-textarea" rows="3" placeholder="写下你的回答..." />
            <div class="answer-input-foot">
              <MButton type="primary" size="small" @click="submitAnswer">提交回答</MButton>
            </div>
          </div>

          <div v-if="answers.length" class="answer-list">
            <div
              v-for="answer in answers"
              :key="answer.id"
              class="answer-item"
              :class="{ 'answer-item--accepted': answer.isAccepted }"
            >
              <div class="answer-head">
                <span class="answer-user">{{ answer.user?.nickname || answer.user?.username }}</span>
                <span class="answer-time">{{ timeAgo(answer.createdAt) }}</span>
                <span v-if="answer.isAccepted" class="accepted-tag">✅ 已采纳</span>
              </div>
              <div class="answer-content">
                <MdPreview :modelValue="answer.content" theme="light" />
              </div>
              <div class="answer-acts">
                <button type="button" class="mini-btn" :class="{ 'mini-btn--on': answer.isLiked }" @click="likeAnswer(answer)">
                  {{ answer.isLiked ? '❤️' : '🤍' }} {{ answer.likesCount || 0 }}
                </button>
                <button
                  v-if="isQuestionOwner && !detailQuestion.isResolved"
                  type="button"
                  class="mini-btn mini-btn--accept"
                  @click="acceptAnswer(answer)"
                >采纳</button>
                <button
                  v-if="currentUser && answer.userId === currentUser.id"
                  type="button"
                  class="mini-btn mini-btn--danger"
                  @click="deleteAnswer(answer)"
                >删除</button>
              </div>
            </div>
          </div>
          <p v-else class="answer-empty">暂无回答，来写下第一个回答吧</p>
        </div>
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

/* 顶部操作条 */
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

.filter-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.chip-scroll {
  display: flex;
  gap: var(--m-space-xs);
  overflow-x: auto;
  scrollbar-width: none;
  flex: 1;
}

.chip-scroll::-webkit-scrollbar {
  display: none;
}

.chip {
  flex-shrink: 0;
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  white-space: nowrap;
}

.chip--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.user-tag {
  flex-shrink: 0;
  font-size: var(--m-font-sm);
  color: var(--m-color-fall);
  font-weight: var(--m-font-weight-medium);
}

.login-hint {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

/* 服务不可用 */
.notice {
  margin: var(--m-space-md);
  padding: var(--m-space-md);
  background: rgba(240, 160, 32, 0.12);
  border: 1px solid #f0a020;
  border-radius: var(--m-radius-md);
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
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

/* 列表 */
.qa-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.qa-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.qa-card:active {
  background: var(--m-bg-primary);
}

.qa-card-head {
  display: flex;
  align-items: flex-start;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-sm);
}

.qa-tag {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  padding: 2px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  white-space: nowrap;
}

.qa-tag--resolved {
  background: var(--m-color-fall-light);
  color: var(--m-color-fall);
}

.qa-tag--pending {
  background: rgba(240, 160, 32, 0.12);
  color: #f0a020;
}

.qa-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  line-height: var(--m-line-height-normal);
}

.qa-card-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
}

.qa-author {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.qa-answers {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
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

/* 详情 */
.detail {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  flex-wrap: wrap;
}

.detail-author {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.del-question {
  margin-left: auto;
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-color-rise);
  border-radius: var(--m-radius-full);
  background: transparent;
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-color-rise);
}

.detail-content {
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-loose);
  border-top: 1px solid var(--m-divider-color);
  padding-top: var(--m-space-sm);
}

/* 回答区 */
.answer-section {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
  border-top: 1px solid var(--m-divider-color);
  padding-top: var(--m-space-md);
}

.answer-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.answer-input-row {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.answer-input-foot {
  display: flex;
  justify-content: flex-end;
}

.answer-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.answer-item {
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
}

.answer-item--accepted {
  border: 1px solid var(--m-color-fall);
  background: var(--m-color-fall-light);
}

.answer-head {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-xs);
}

.answer-user {
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.answer-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.accepted-tag {
  font-size: var(--m-font-xs);
  color: var(--m-color-fall);
  font-weight: var(--m-font-weight-medium);
}

.answer-content {
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-normal);
}

.answer-acts {
  display: flex;
  gap: var(--m-space-md);
  margin-top: var(--m-space-xs);
}

.mini-btn {
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.mini-btn--on {
  color: var(--m-color-rise);
}

.mini-btn--accept {
  color: var(--m-color-fall);
  font-weight: var(--m-font-weight-medium);
}

.mini-btn--danger {
  color: var(--m-color-rise);
}

.answer-empty {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
  text-align: center;
  padding: var(--m-space-lg) 0;
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
