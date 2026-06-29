<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import 'md-editor-v3/lib/preview.css'
import { MdPreview } from 'md-editor-v3'
import { AddPromptTemplate } from '../../../api/app'
import { EventsEmit } from '../../../api/runtime'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MButton from '../../components/base/MButton.vue'
import MSheet from '../../components/base/MSheet.vue'
import { toast } from '../../composables/useToast'
import { usePromptPlaza } from '../../composables/usePromptPlaza'

// 对齐桌面端 promptPlaza.vue：登录/注册 + 分类/排序/搜索 + 列表 + 详情(点赞/收藏/下载/复制/评论) + 发布/编辑/删除 + 排行榜 + 添加到我的模板

const {
  currentUser,
  isLoggedIn,
  apiGet,
  apiPost,
  apiPut,
  apiDelete,
  fetchCurrentUser,
  login,
  register,
  logout,
  displayName,
  rememberedUsername,
  timeAgo,
  formatTime,
  copyText,
} = usePromptPlaza()

const PAGE_SIZE = 12

const SORT_OPTIONS = [
  { label: '🕐 最新', value: 'latest' },
  { label: '🔥 热度', value: 'hot' },
  { label: '❤️ 点赞', value: 'likes' },
  { label: '⭐ 收藏', value: 'favorites' },
  { label: '⬇️ 下载', value: 'downloads' },
  { label: '💬 评论', value: 'comments' },
]

const categories = ref([])
const activeCategory = ref(null)
const activeSort = ref('latest')
const keyword = ref('')
const loading = ref(false)
const prompts = ref([])
const page = ref(1)
const total = ref(0)
const totalPages = ref(1)

const hasMore = computed(() => page.value < totalPages.value)

// 后端可能返回 vip 限制标记，移动端不展示付费门槛，统一抹平
function stripVip(p) {
  if (p && typeof p === 'object') {
    p.needVip = false
    p.vipOnly = false
  }
  return p
}

async function loadCategories() {
  try {
    const data = await apiGet('/prompts/categories')
    categories.value = Array.isArray(data) ? data : []
  } catch (e) {
    console.warn('加载分类失败', e)
  }
}

async function loadPrompts(targetPage = 1, append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const params = { page: targetPage, pageSize: PAGE_SIZE, sort: activeSort.value }
    if (activeCategory.value) params.category = activeCategory.value
    if (keyword.value) params.keyword = keyword.value
    const data = await apiGet('/prompts', params)
    const rows = (data.list || []).map(stripVip)
    prompts.value = append ? prompts.value.concat(rows) : rows
    total.value = data.total || 0
    totalPages.value = Math.ceil((data.total || 0) / (data.pageSize || PAGE_SIZE)) || 1
    page.value = targetPage
  } catch (e) {
    console.error('加载提示词列表失败:', e)
    toast.error('加载失败: ' + (e.message || ''))
  } finally {
    loading.value = false
  }
}

function handleRefresh() {
  return Promise.all([loadCategories(), loadPrompts(1, false)])
}

function loadMore() {
  if (hasMore.value && !loading.value) loadPrompts(page.value + 1, true)
}

function applySearch() {
  loadPrompts(1, false)
}

function selectCategory(cat) {
  activeCategory.value = cat
  loadPrompts(1, false)
}

function selectSort(value) {
  activeSort.value = value
  loadPrompts(1, false)
}

// ===== 登录 / 注册 =====
const loginVisible = ref(false)
const loginTab = ref('login')
const loginForm = ref({ username: rememberedUsername(), password: '', nickname: '' })
const loginLoading = ref(false)

function openLogin() {
  loginTab.value = 'login'
  loginForm.value.username = rememberedUsername()
  loginVisible.value = true
}

async function submitLogin() {
  if (!loginForm.value.username || !loginForm.value.password) {
    toast.warning('请填写用户名和密码')
    return
  }
  loginLoading.value = true
  try {
    if (loginTab.value === 'login') {
      await login(loginForm.value.username, loginForm.value.password)
      toast.success('登录成功')
    } else {
      await register(loginForm.value.username, loginForm.value.password, loginForm.value.nickname)
      toast.success('注册成功')
    }
    loginVisible.value = false
    loginForm.value.password = ''
    loadPrompts(1, false)
  } catch (e) {
    toast.error((loginTab.value === 'login' ? '登录失败: ' : '注册失败: ') + (e.message || ''))
  } finally {
    loginLoading.value = false
  }
}

function handleLogout() {
  if (!confirm('确定要退出登录吗？')) return
  logout()
  toast.success('已退出登录')
  loadPrompts(1, false)
}

function requireLogin() {
  if (!isLoggedIn.value) {
    toast.warning('请先登录')
    openLogin()
    return false
  }
  return true
}

// ===== 详情 + 评论 =====
const detailVisible = ref(false)
const detail = ref(null)
const comments = ref([])
const commentLoading = ref(false)
const newComment = ref('')
const replyTo = ref(null)

async function openDetail(id) {
  try {
    const data = await apiGet(`/prompts/${id}`)
    detail.value = stripVip(data)
    newComment.value = ''
    replyTo.value = null
    detailVisible.value = true
    loadComments(id)
  } catch (e) {
    toast.error('加载详情失败: ' + (e.message || ''))
  }
}

async function loadComments(promptId) {
  commentLoading.value = true
  try {
    const data = await apiGet(`/prompts/${promptId}/comments`, { page: 1, pageSize: 50 })
    comments.value = data.list || []
  } catch (e) {
    console.warn('加载评论失败', e)
  } finally {
    commentLoading.value = false
  }
}

// 列表项与详情同步点赞/收藏的本地计数
function syncCounts(target, patch) {
  Object.assign(target, patch)
  if (detail.value && detail.value.id === target.id) Object.assign(detail.value, patch)
  const inList = prompts.value.find(p => p.id === target.id)
  if (inList && inList !== target) Object.assign(inList, patch)
}

async function toggleLike(item) {
  if (!requireLogin()) return
  try {
    const data = await apiPost(`/prompts/${item.id}/like`)
    syncCounts(item, { isLiked: data.isLiked, likesCount: data.likesCount })
  } catch (e) {
    toast.error('操作失败: ' + (e.message || ''))
  }
}

async function toggleFavorite(item) {
  if (!requireLogin()) return
  try {
    const data = await apiPost(`/prompts/${item.id}/favorite`)
    syncCounts(item, { isFavorited: data.isFavorited, favoritesCount: data.favoritesCount })
  } catch (e) {
    toast.error('操作失败: ' + (e.message || ''))
  }
}

async function downloadPrompt(item) {
  try {
    const data = await apiGet(`/prompts/${item.id}/download`)
    await copyText(data.content)
    syncCounts(item, { downloadsCount: (item.downloadsCount || 0) + 1 })
    toast.success('内容已复制到剪贴板')
  } catch (e) {
    if (item.content) {
      await copyText(item.content)
      toast.success('已复制当前内容')
      return
    }
    toast.error('下载失败: ' + (e.message || ''))
  }
}

async function copyContent(content) {
  try {
    await copyText(content)
    toast.success('已复制到剪贴板')
  } catch (e) {
    toast.error('复制失败')
  }
}

async function addToTemplate(item) {
  try {
    const res = await AddPromptTemplate({ name: item.title, content: item.content, type: '模型系统Prompt' })
    if (res === '添加成功') {
      toast.success('已添加到我的模板')
      EventsEmit('promptTemplatesChanged')
    } else {
      toast.warning(typeof res === 'string' ? res : '添加完成')
    }
  } catch (e) {
    toast.error('添加失败: ' + (e.message || ''))
  }
}

async function submitComment() {
  if (!requireLogin()) return
  if (!newComment.value.trim()) {
    toast.warning('请输入评论内容')
    return
  }
  try {
    const body = { content: newComment.value }
    if (replyTo.value) body.parentId = replyTo.value.id
    await apiPost(`/prompts/${detail.value.id}/comments`, body)
    newComment.value = ''
    replyTo.value = null
    detail.value.commentsCount = (detail.value.commentsCount || 0) + 1
    loadComments(detail.value.id)
    toast.success('评论成功')
  } catch (e) {
    toast.error('评论失败: ' + (e.message || ''))
  }
}

async function deleteComment(commentId) {
  if (!confirm('确定要删除这条评论吗？')) return
  try {
    await apiDelete(`/comments/${commentId}`)
    detail.value.commentsCount = Math.max(0, (detail.value.commentsCount || 1) - 1)
    loadComments(detail.value.id)
    toast.success('删除成功')
  } catch (e) {
    toast.error('删除失败: ' + (e.message || ''))
  }
}

const isOwner = computed(() => !!(currentUser.value && detail.value && detail.value.userId === currentUser.value.id))

// ===== 发布 / 编辑 =====
const editVisible = ref(false)
const editForm = ref({ id: 0, title: '', content: '', description: '', category: '', tags: '', isPublic: true })
const editSaving = ref(false)
const editTitle = computed(() => (editForm.value.id > 0 ? '编辑提示词' : '发布提示词'))

function openCreate() {
  if (!requireLogin()) return
  editForm.value = { id: 0, title: '', content: '', description: '', category: '', tags: '', isPublic: true }
  editVisible.value = true
}

function openEdit(item) {
  editForm.value = {
    id: item.id,
    title: item.title || '',
    content: item.content || '',
    description: item.description || '',
    category: item.category || '',
    tags: item.tags || '',
    isPublic: item.isPublic !== false,
  }
  editVisible.value = true
}

async function submitEdit() {
  const f = editForm.value
  if (!f.title || !f.content) {
    toast.warning('请填写标题和内容')
    return
  }
  editSaving.value = true
  try {
    const body = {
      title: f.title,
      content: f.content,
      description: f.description,
      category: f.category,
      tags: f.tags,
      isPublic: f.isPublic,
    }
    if (f.id > 0) {
      await apiPut(`/prompts/${f.id}`, body)
      toast.success('修改成功')
      detailVisible.value = false
    } else {
      await apiPost('/prompts', body)
      toast.success('发布成功')
    }
    editVisible.value = false
    loadPrompts(1, false)
    loadCategories()
  } catch (e) {
    toast.error((f.id > 0 ? '修改失败: ' : '发布失败: ') + (e.message || ''))
  } finally {
    editSaving.value = false
  }
}

async function deletePrompt(item) {
  if (!confirm('确定要删除这个提示词吗？删除后不可恢复。')) return
  try {
    await apiDelete(`/prompts/${item.id}`)
    detailVisible.value = false
    toast.success('删除成功')
    loadPrompts(1, false)
    loadCategories()
  } catch (e) {
    toast.error('删除失败: ' + (e.message || ''))
  }
}

// ===== 排行榜 =====
const rankingVisible = ref(false)
const rankingType = ref('hot')
const rankingRange = ref('all')
const rankingList = ref([])
const rankingLoading = ref(false)

const RANK_TYPES = [
  { label: '🔥 热度', value: 'hot' },
  { label: '❤️ 点赞', value: 'likes' },
  { label: '⬇️ 下载', value: 'downloads' },
  { label: '⭐ 收藏', value: 'favorites' },
]
const RANK_RANGES = [
  { label: '全部', value: 'all' },
  { label: '今日', value: 'daily' },
  { label: '本周', value: 'weekly' },
  { label: '本月', value: 'monthly' },
]

function openRanking() {
  rankingVisible.value = true
  loadRanking('hot', 'all')
}

async function loadRanking(type, range) {
  rankingType.value = type
  rankingRange.value = range
  rankingLoading.value = true
  try {
    const data = await apiGet('/prompts/ranking', { type, range, limit: 50 })
    rankingList.value = (data.list || []).map(stripVip)
  } catch (e) {
    toast.error('加载排行榜失败: ' + (e.message || ''))
  } finally {
    rankingLoading.value = false
  }
}

function openFromRanking(id) {
  rankingVisible.value = false
  openDetail(id)
}

function splitTags(tags) {
  if (!tags) return []
  return tags.split(',').map(t => t.trim()).filter(Boolean)
}

onBeforeMount(() => {
  loadCategories()
  loadPrompts(1, false)
  if (isLoggedIn.value) fetchCurrentUser()
})
</script>

<template>
  <div class="page">
    <!-- 顶部：搜索 + 账号 -->
    <div class="op-bar">
      <div class="search-row">
        <input
          v-model="keyword"
          class="search-input"
          type="text"
          placeholder="搜索提示词..."
          @keyup.enter="applySearch"
        />
        <button type="button" class="icon-btn" @click="openRanking">🏆</button>
      </div>
      <div class="account-row">
        <MButton type="primary" size="small" @click="openCreate">✏️ 发布</MButton>
        <template v-if="isLoggedIn">
          <span class="user-tag">{{ displayName() }}</span>
          <button type="button" class="link-btn" @click="handleLogout">退出</button>
        </template>
        <button v-else type="button" class="link-btn" @click="openLogin">登录 / 注册</button>
      </div>
      <!-- 排序 chip 横滚 -->
      <div class="chip-scroll">
        <button
          v-for="opt in SORT_OPTIONS"
          :key="opt.value"
          type="button"
          class="chip"
          :class="{ 'chip--active': activeSort === opt.value }"
          @click="selectSort(opt.value)"
        >{{ opt.label }}</button>
      </div>
      <!-- 分类 chip 横滚 -->
      <div v-if="categories.length" class="chip-scroll">
        <button
          type="button"
          class="chip"
          :class="{ 'chip--active': activeCategory === null }"
          @click="selectCategory(null)"
        >全部</button>
        <button
          v-for="cat in categories"
          :key="cat"
          type="button"
          class="chip"
          :class="{ 'chip--active': activeCategory === cat }"
          @click="selectCategory(cat)"
        >{{ cat }}</button>
      </div>
    </div>

    <MPullRefresh class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="prompts.length" class="prompt-list">
          <div v-for="item in prompts" :key="item.id" class="prompt-card" @click="openDetail(item.id)">
            <div class="card-head">
              <span class="card-title">{{ item.title }}</span>
              <span v-if="item.category" class="card-cat">{{ item.category }}</span>
            </div>
            <p class="card-desc">{{ item.summary || item.description || item.content }}</p>
            <div v-if="splitTags(item.tags).length" class="card-tags">
              <span v-for="tag in splitTags(item.tags).slice(0, 3)" :key="tag" class="tag">{{ tag }}</span>
            </div>
            <div class="card-foot">
              <span class="card-author">{{ item.user?.nickname || item.user?.username || '匿名' }} · {{ timeAgo(item.createdAt) }}</span>
            </div>
            <div class="card-stats">
              <span class="stat">👁️ {{ item.viewsCount || 0 }}</span>
              <button type="button" class="stat stat-btn" :class="{ 'stat--on': item.isLiked }" @click.stop="toggleLike(item)">
                {{ item.isLiked ? '❤️' : '🤍' }} {{ item.likesCount || 0 }}
              </button>
              <button type="button" class="stat stat-btn" :class="{ 'stat--warn': item.isFavorited }" @click.stop="toggleFavorite(item)">
                {{ item.isFavorited ? '⭐' : '☆' }} {{ item.favoritesCount || 0 }}
              </button>
              <span class="stat">💬 {{ item.commentsCount || 0 }}</span>
              <span class="stat">⬇️ {{ item.downloadsCount || 0 }}</span>
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
        <MEmpty v-else description="暂无提示词">
          <MButton type="primary" @click="openCreate">发布提示词</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 详情抽屉 -->
    <MSheet v-model:show="detailVisible" :title="detail?.title || '提示词详情'">
      <div v-if="detail" class="detail">
        <div class="detail-meta">
          <span v-if="detail.category" class="card-cat">{{ detail.category }}</span>
          <span class="detail-author">{{ detail.user?.nickname || detail.user?.username || '匿名' }} · {{ formatTime(detail.createdAt) }}</span>
        </div>

        <div v-if="splitTags(detail.tags).length" class="card-tags">
          <span v-for="tag in splitTags(detail.tags)" :key="tag" class="tag">{{ tag }}</span>
        </div>

        <!-- 操作区 -->
        <div class="detail-actions">
          <button type="button" class="act-chip" :class="{ 'act-chip--on': detail.isLiked }" @click="toggleLike(detail)">
            {{ detail.isLiked ? '❤️ 已赞' : '🤍 点赞' }} {{ detail.likesCount || 0 }}
          </button>
          <button type="button" class="act-chip" :class="{ 'act-chip--warn': detail.isFavorited }" @click="toggleFavorite(detail)">
            {{ detail.isFavorited ? '⭐ 已收藏' : '☆ 收藏' }} {{ detail.favoritesCount || 0 }}
          </button>
          <button type="button" class="act-chip" @click="downloadPrompt(detail)">⬇️ 下载 {{ detail.downloadsCount || 0 }}</button>
          <button type="button" class="act-chip" @click="copyContent(detail.content)">📋 复制</button>
          <button type="button" class="act-chip" @click="addToTemplate(detail)">➕ 加到模板</button>
          <button v-if="isOwner" type="button" class="act-chip" @click="openEdit(detail)">✏️ 编辑</button>
          <button v-if="isOwner" type="button" class="act-chip act-chip--danger" @click="deletePrompt(detail)">🗑️ 删除</button>
        </div>

        <!-- 内容 -->
        <div class="detail-content">
          <MdPreview :modelValue="detail.content" theme="light" />
        </div>

        <!-- 评论 -->
        <div class="comment-section">
          <div class="comment-title">评论 ({{ detail.commentsCount || 0 }})</div>
          <div class="comment-input-row">
            <textarea
              v-model="newComment"
              class="field-textarea"
              rows="2"
              :placeholder="replyTo ? `回复 @${replyTo.user?.nickname || replyTo.user?.username}...` : '发表评论...'"
            />
            <div class="comment-input-foot">
              <span v-if="replyTo" class="reply-hint">
                回复 @{{ replyTo.user?.nickname || replyTo.user?.username }}
                <button type="button" class="link-btn" @click="replyTo = null">取消</button>
              </span>
              <span v-else />
              <MButton type="primary" size="small" @click="submitComment">发表</MButton>
            </div>
          </div>

          <MLoading v-if="commentLoading" text="加载评论..." class="comment-loading" />
          <div v-else-if="comments.length" class="comment-list">
            <div v-for="c in comments" :key="c.id" class="comment-item">
              <div class="comment-head">
                <span class="comment-user">{{ c.user?.nickname || c.user?.username }}</span>
                <span class="comment-time">{{ timeAgo(c.createdAt) }}</span>
              </div>
              <p class="comment-content">{{ c.content }}</p>
              <div class="comment-acts">
                <button type="button" class="mini-btn" @click="replyTo = c">回复</button>
                <button
                  v-if="currentUser && c.userId === currentUser.id"
                  type="button"
                  class="mini-btn mini-btn--danger"
                  @click="deleteComment(c.id)"
                >删除</button>
              </div>
            </div>
          </div>
          <p v-else class="comment-empty">暂无评论</p>
        </div>
      </div>
    </MSheet>

    <!-- 发布/编辑抽屉 -->
    <MSheet v-model:show="editVisible" :title="editTitle">
      <div class="form">
        <label class="field">
          <span class="field-label">标题 <i class="req">*</i></span>
          <input v-model="editForm.title" class="field-input" type="text" placeholder="标题" />
        </label>
        <div class="field-grid">
          <label class="field">
            <span class="field-label">分类</span>
            <input v-model="editForm.category" class="field-input" type="text" placeholder="如：AI编程" />
          </label>
          <label class="field">
            <span class="field-label">标签</span>
            <input v-model="editForm.tags" class="field-input" type="text" placeholder="逗号分隔" />
          </label>
        </div>
        <label class="field">
          <span class="field-label">简短描述</span>
          <textarea v-model="editForm.description" class="field-textarea" rows="2" placeholder="简短描述" />
        </label>
        <label class="field">
          <span class="field-label">内容 <i class="req">*</i></span>
          <textarea v-model="editForm.content" class="field-textarea" rows="8" placeholder="提示词内容（支持 Markdown）" />
        </label>
        <button type="button" class="switch-line" @click="editForm.isPublic = !editForm.isPublic">
          <span>公开到广场</span>
          <span class="switch" :class="{ 'switch--on': editForm.isPublic }"><span class="switch-dot" /></span>
        </button>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="editVisible = false">取消</MButton>
          <MButton type="primary" block :loading="editSaving" @click="submitEdit">
            {{ editForm.id > 0 ? '保存' : '发布' }}
          </MButton>
        </div>
      </template>
    </MSheet>

    <!-- 登录/注册抽屉 -->
    <MSheet v-model:show="loginVisible" title="账号" height="auto">
      <div class="login-tabs">
        <button type="button" class="login-tab" :class="{ 'login-tab--active': loginTab === 'login' }" @click="loginTab = 'login'">登录</button>
        <button type="button" class="login-tab" :class="{ 'login-tab--active': loginTab === 'register' }" @click="loginTab = 'register'">注册</button>
      </div>
      <div class="form">
        <label class="field">
          <span class="field-label">用户名</span>
          <input v-model="loginForm.username" class="field-input" type="text" :placeholder="loginTab === 'register' ? '用户名 (3-50字)' : '用户名'" />
        </label>
        <label class="field">
          <span class="field-label">密码</span>
          <input v-model="loginForm.password" class="field-input" type="password" :placeholder="loginTab === 'register' ? '密码 (6字以上)' : '密码'" />
        </label>
        <label v-if="loginTab === 'register'" class="field">
          <span class="field-label">昵称</span>
          <input v-model="loginForm.nickname" class="field-input" type="text" placeholder="昵称 (可选)" />
        </label>
      </div>
      <template #footer>
        <MButton type="primary" block :loading="loginLoading" @click="submitLogin">
          {{ loginTab === 'login' ? '登录' : '注册' }}
        </MButton>
      </template>
    </MSheet>

    <!-- 排行榜抽屉 -->
    <MSheet v-model:show="rankingVisible" title="🏆 排行榜">
      <div class="rank-filters">
        <div class="chip-scroll">
          <button
            v-for="opt in RANK_TYPES"
            :key="opt.value"
            type="button"
            class="chip"
            :class="{ 'chip--active': rankingType === opt.value }"
            @click="loadRanking(opt.value, rankingRange)"
          >{{ opt.label }}</button>
        </div>
        <div class="chip-scroll">
          <button
            v-for="opt in RANK_RANGES"
            :key="opt.value"
            type="button"
            class="chip"
            :class="{ 'chip--active': rankingRange === opt.value }"
            @click="loadRanking(rankingType, opt.value)"
          >{{ opt.label }}</button>
        </div>
      </div>
      <MLoading v-if="rankingLoading" text="加载中..." vertical class="page-loading" />
      <div v-else-if="rankingList.length" class="rank-list">
        <div v-for="item in rankingList" :key="item.id" class="rank-item" @click="openFromRanking(item.id)">
          <span class="rank-no" :class="{ 'rank-no--top': item.rank <= 3 }">{{ item.rank }}</span>
          <div class="rank-body">
            <div class="rank-title">{{ item.title }}</div>
            <div class="rank-stats">
              <span>❤️ {{ item.likesCount || 0 }}</span>
              <span>⬇️ {{ item.downloadsCount || 0 }}</span>
              <span>⭐ {{ item.favoritesCount || 0 }}</span>
              <span v-if="item.hotScore" class="rank-hot">🔥 {{ item.hotScore }}</span>
            </div>
          </div>
        </div>
      </div>
      <MEmpty v-else description="暂无排行数据" />
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

.icon-btn {
  flex-shrink: 0;
  width: var(--m-touch-min);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font-size: var(--m-font-lg);
}

.account-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.user-tag {
  margin-left: auto;
  font-size: var(--m-font-sm);
  color: var(--m-color-fall);
  font-weight: var(--m-font-weight-medium);
}

.account-row .link-btn {
  margin-left: auto;
}

.account-row .user-tag + .link-btn {
  margin-left: 0;
}

.link-btn {
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-color-rise);
}

/* chip 横滚 */
.chip-scroll {
  display: flex;
  gap: var(--m-space-xs);
  overflow-x: auto;
  scrollbar-width: none;
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
.prompt-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.prompt-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.prompt-card:active {
  background: var(--m-bg-primary);
}

.card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-xs);
}

.card-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  flex: 1;
  min-width: 0;
}

.card-cat {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  color: #2080f0;
  background: rgba(32, 128, 240, 0.12);
  padding: 2px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
}

.card-desc {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: var(--m-space-sm);
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
  margin-bottom: var(--m-space-sm);
}

.tag {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  background: var(--m-bg-primary);
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-full);
}

.card-foot {
  margin-bottom: var(--m-space-sm);
}

.card-author {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.card-stats {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  flex-wrap: wrap;
}

.stat {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.stat-btn {
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-xs);
  padding: 0;
}

.stat--on {
  color: var(--m-color-rise);
}

.stat--warn {
  color: #f0a020;
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

.detail-actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
}

.act-chip {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.act-chip--on {
  color: var(--m-color-rise);
  border-color: var(--m-color-rise);
  background: var(--m-color-rise-light);
}

.act-chip--warn {
  color: #f0a020;
  border-color: #f0a020;
}

.act-chip--danger {
  color: var(--m-color-rise);
  border-color: var(--m-color-rise);
}

.detail-content {
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-loose);
  border-top: 1px solid var(--m-divider-color);
  padding-top: var(--m-space-sm);
}

/* 评论 */
.comment-section {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  border-top: 1px solid var(--m-divider-color);
  padding-top: var(--m-space-md);
}

.comment-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.comment-input-row {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.comment-input-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.reply-hint {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.comment-loading {
  padding: var(--m-space-lg) 0;
  justify-content: center;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.comment-item {
  padding-bottom: var(--m-space-md);
  border-bottom: 1px solid var(--m-divider-color);
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-head {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-xs);
}

.comment-user {
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.comment-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.comment-content {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  margin-bottom: var(--m-space-xs);
}

.comment-acts {
  display: flex;
  gap: var(--m-space-md);
}

.mini-btn {
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.mini-btn--danger {
  color: var(--m-color-rise);
}

.comment-empty {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
  text-align: center;
  padding: var(--m-space-lg) 0;
}

/* 表单（与模板页一致） */
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

/* 登录 Tab */
.login-tabs {
  display: flex;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-lg);
}

.login-tab {
  flex: 1;
  padding: var(--m-space-sm);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-md);
  color: var(--m-text-secondary);
}

.login-tab--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

/* 排行榜 */
.rank-filters {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-md);
}

.rank-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.rank-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-sm) 0;
  border-bottom: 1px solid var(--m-divider-color);
}

.rank-item:last-child {
  border-bottom: none;
}

.rank-no {
  flex-shrink: 0;
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  font-variant-numeric: tabular-nums;
}

.rank-no--top {
  background: var(--m-color-rise);
  color: #fff;
}

.rank-body {
  flex: 1;
  min-width: 0;
}

.rank-title {
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-stats {
  display: flex;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  margin-top: 2px;
}

.rank-hot {
  color: #f0a020;
}
</style>
