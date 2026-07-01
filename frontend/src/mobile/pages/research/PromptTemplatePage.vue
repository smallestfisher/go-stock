<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import 'md-editor-v3/lib/preview.css'
import { MdPreview } from 'md-editor-v3'
import {
  GetPromptTemplateList,
  AddPromptTemplate,
  UpdatePromptTemplate,
  DeletePromptTemplate,
} from '../../../api/app'
import { EventsEmit } from '../../../api/runtime'
import { promptPlazaURL, promptPlazaHeaders, parsePromptPlazaResponse } from '../../../api/promptPlaza'
import { RECOMMENDED_PROMPTS } from '../../composables/recommendedPrompts'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MButton from '../../components/base/MButton.vue'
import MSheet from '../../components/base/MSheet.vue'
import MIcon from '../../components/base/MIcon.vue'
import { toast } from '../../composables/useToast'

// 对齐桌面端 promptTemplateList.vue：分页列表 + 搜索 + 新增/编辑 + 删除 + 分享到广场
// 字段：{ ID, name, content, type, CreatedAt, UpdatedAt }，type ∈ {模型系统Prompt, 模型用户Prompt}

const PAGE_SIZE = 12

const TYPE_OPTIONS = [
  { label: '模型系统Prompt', value: '模型系统Prompt' },
  { label: '模型用户Prompt', value: '模型用户Prompt' },
]

const list = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const totalPages = ref(1)

// 搜索条件
const search = ref({ name: '', type: '', content: '' })
const filterVisible = ref(false)
const hasFilter = computed(() => !!(search.value.name || search.value.type || search.value.content))

// 列表是否还有下一页
const hasMore = computed(() => page.value < totalPages.value)

async function fetchPage(targetPage, append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const res = await GetPromptTemplateList({
      page: targetPage,
      pageSize: PAGE_SIZE,
      name: search.value.name,
      type: search.value.type,
      content: search.value.content,
    })
    const rows = res && Array.isArray(res.list) ? res.list : []
    list.value = append ? list.value.concat(rows) : rows
    total.value = (res && res.total) || 0
    totalPages.value = (res && res.totalPages) || 1
    page.value = targetPage
  } catch (e) {
    console.error('加载提示词模板失败:', e)
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
  search.value = { name: '', type: '', content: '' }
  filterVisible.value = false
  fetchPage(1, false)
}

// ===== 新增 / 编辑 =====
const editVisible = ref(false)
const editForm = ref({ ID: 0, name: '', type: '', content: '' })
const editSaving = ref(false)
const previewMode = ref(false)
const editTitle = computed(() => (editForm.value.ID > 0 ? '编辑模板' : '新增模板'))

function openCreate() {
  editForm.value = { ID: 0, name: '', type: '', content: '' }
  previewMode.value = false
  editVisible.value = true
}

function openEdit(item) {
  editForm.value = { ID: item.ID, name: item.name, type: item.type, content: item.content }
  previewMode.value = false
  editVisible.value = true
}

async function saveTemplate() {
  const f = editForm.value
  if (!f.name || !f.type || !f.content) {
    toast.warning('请填写名称、类型和内容')
    return
  }
  editSaving.value = true
  try {
    const apiCall = f.ID > 0 ? UpdatePromptTemplate : AddPromptTemplate
    const res = await apiCall({ ID: f.ID, name: f.name, type: f.type, content: f.content })
    toast.success(typeof res === 'string' ? res : '保存成功')
    editVisible.value = false
    EventsEmit('promptTemplatesChanged')
    fetchPage(1, false)
  } catch (e) {
    console.error('保存模板失败:', e)
    toast.error('保存失败')
  } finally {
    editSaving.value = false
  }
}

async function removeTemplate(item) {
  if (!confirm(`确定删除模板「${item.name}」吗？删除后不可恢复。`)) return
  try {
    const res = await DeletePromptTemplate(item.ID)
    toast.success(typeof res === 'string' ? res : '删除成功')
    EventsEmit('promptTemplatesChanged')
    // 删除后留在当前页，重新拉第一页保持简单一致
    fetchPage(1, false)
  } catch (e) {
    console.error('删除模板失败:', e)
    toast.error('删除失败')
  }
}

// ===== 分享到提示词广场 =====
const shareVisible = ref(false)
const shareForm = ref({ title: '', content: '', description: '', category: '', tags: '', isPublic: true })
const shareSaving = ref(false)

function openShare(item) {
  shareForm.value = {
    title: item.name || '',
    content: item.content || '',
    description: '',
    category: item.type || '',
    tags: '',
    isPublic: true,
  }
  shareVisible.value = true
}

async function submitShare() {
  const f = shareForm.value
  if (!f.title || !f.content) {
    toast.warning('标题和内容不能为空')
    return
  }
  const plazaToken = localStorage.getItem('promptPlazaToken')
  if (!plazaToken) {
    toast.warning('请先在「提示词广场」登录后再分享')
    return
  }
  shareSaving.value = true
  try {
    const resp = await fetch(promptPlazaURL('/prompts'), {
      method: 'POST',
      headers: promptPlazaHeaders(plazaToken),
      body: JSON.stringify({
        title: f.title,
        content: f.content,
        description: f.description,
        category: f.category,
        tags: f.tags,
        isPublic: f.isPublic,
      }),
    })
    await parsePromptPlazaResponse(resp)
    toast.success('分享成功')
    shareVisible.value = false
  } catch (e) {
    console.error('分享失败:', e)
    toast.error('分享失败: ' + (e.message || ''))
  } finally {
    shareSaving.value = false
  }
}

function fmtTime(s) {
  if (!s) return '-'
  return String(s).substring(0, 19).replace('T', ' ')
}

function isSystemType(type) {
  return type === '模型系统Prompt'
}

// ===== 导入推荐模板（内置 A股 长线/短线/选股/诊股 4 套）=====
const importing = ref(false)

async function importRecommended() {
  if (importing.value) return
  importing.value = true
  try {
    // 拉全量现有模板，按 name 去重，避免重复导入
    const existing = new Set()
    // 一次性多取，覆盖绝大多数情况（模板量通常很小）
    const res = await GetPromptTemplateList({ page: 1, pageSize: 200, name: '', type: '', content: '' })
    const rows = res && Array.isArray(res.list) ? res.list : []
    rows.forEach(r => existing.add(r.name))

    const toAdd = RECOMMENDED_PROMPTS.filter(t => !existing.has(t.name))
    if (!toAdd.length) {
      toast.info('推荐模板已全部导入')
      return
    }
    let ok = 0
    for (const t of toAdd) {
      try {
        await AddPromptTemplate({ ID: 0, name: t.name, type: t.type, content: t.content })
        ok++
      } catch (e) {
        console.error('导入模板失败:', t.name, e)
      }
    }
    toast.success(`已导入 ${ok} 套推荐模板`)
    EventsEmit('promptTemplatesChanged')
    fetchPage(1, false)
  } catch (e) {
    console.error('导入推荐模板失败:', e)
    toast.error('导入失败')
  } finally {
    importing.value = false
  }
}

onBeforeMount(() => {
  fetchPage(1, false)
})
</script>

<template>
  <div class="page">
    <!-- 操作条：搜索入口 + 新建 -->
    <div class="op-bar">
      <button
        type="button"
        class="filter-trigger"
        :class="{ 'filter-trigger--active': hasFilter }"
        @click="filterVisible = true"
      >
        <MIcon name="search" :size="16" />
        <span v-if="hasFilter" class="filter-trigger__text">已筛选</span>
        <span v-else class="filter-trigger__text">搜索模板</span>
      </button>
      <MButton type="default" size="small" :loading="importing" @click="importRecommended">导入推荐</MButton>
      <MButton type="primary" size="small" @click="openCreate">新建模板</MButton>
    </div>

    <MPullRefresh class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="list.length" class="template-list">
          <div v-for="item in list" :key="item.ID" class="template-card">
            <div class="card-head">
              <span class="card-name" :class="{ 'card-name--sys': isSystemType(item.type) }">{{ item.name }}</span>
              <span class="card-type" :class="{ 'card-type--sys': isSystemType(item.type) }">{{ item.type }}</span>
            </div>
            <p class="card-content">{{ item.content }}</p>
            <div class="card-foot">
              <span class="card-time">{{ fmtTime(item.CreatedAt) }}</span>
              <div class="card-actions">
                <button type="button" class="act-btn" @click="openEdit(item)">编辑</button>
                <button type="button" class="act-btn" @click="openShare(item)">分享</button>
                <button type="button" class="act-btn act-btn--danger" @click="removeTemplate(item)">删除</button>
              </div>
            </div>
          </div>

          <!-- 加载更多 -->
          <div class="load-more">
            <button v-if="hasMore" type="button" class="load-more-btn" :disabled="loading" @click="loadMore">
              {{ loading ? '加载中...' : '加载更多' }}
            </button>
            <span v-else class="load-more-end">共 {{ total }} 条</span>
          </div>
        </div>

        <MLoading v-else-if="loading" text="加载中..." vertical class="page-loading" />
        <MEmpty v-else description="还没有提示词模板">
          <div class="empty-actions">
            <MButton type="primary" @click="openCreate">新建模板</MButton>
            <MButton type="default" :loading="importing" @click="importRecommended">导入推荐模板</MButton>
          </div>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 搜索抽屉 -->
    <MSheet v-model:show="filterVisible" title="搜索模板" height="auto">
      <div class="form">
        <label class="field">
          <span class="field-label">模板名称</span>
          <input v-model="search.name" class="field-input" type="text" placeholder="按名称筛选" />
        </label>
        <div class="field">
          <span class="field-label">模板类型</span>
          <div class="chip-row">
            <button
              type="button"
              class="chip"
              :class="{ 'chip--active': search.type === '' }"
              @click="search.type = ''"
            >全部</button>
            <button
              v-for="opt in TYPE_OPTIONS"
              :key="opt.value"
              type="button"
              class="chip"
              :class="{ 'chip--active': search.type === opt.value }"
              @click="search.type = opt.value"
            >{{ opt.label }}</button>
          </div>
        </div>
        <label class="field">
          <span class="field-label">内容关键词</span>
          <input v-model="search.content" class="field-input" type="text" placeholder="按内容筛选" />
        </label>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="resetSearch">重置</MButton>
          <MButton type="primary" block @click="applySearch">搜索</MButton>
        </div>
      </template>
    </MSheet>

    <!-- 新增/编辑抽屉 -->
    <MSheet v-model:show="editVisible" :title="editTitle">
      <div class="form">
        <label class="field">
          <span class="field-label">模板名称 <i class="req">*</i></span>
          <input v-model="editForm.name" class="field-input" type="text" placeholder="请输入模板名称" />
        </label>
        <div class="field">
          <span class="field-label">模板类型 <i class="req">*</i></span>
          <div class="chip-row">
            <button
              v-for="opt in TYPE_OPTIONS"
              :key="opt.value"
              type="button"
              class="chip"
              :class="{ 'chip--active': editForm.type === opt.value }"
              @click="editForm.type = opt.value"
            >{{ opt.label }}</button>
          </div>
        </div>
        <div class="field">
          <div class="field-label-row">
            <span class="field-label">模板内容 <i class="req">*</i></span>
            <button type="button" class="link-btn" @click="previewMode = !previewMode">
              {{ previewMode ? '编辑' : '预览' }}
            </button>
          </div>
          <textarea
            v-if="!previewMode"
            v-model="editForm.content"
            class="field-textarea"
            rows="10"
            placeholder="请输入模板内容（支持 Markdown）"
          />
          <div v-else class="preview-box">
            <MdPreview v-if="editForm.content" :modelValue="editForm.content" theme="light" />
            <p v-else class="preview-empty">暂无内容</p>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="editVisible = false">取消</MButton>
          <MButton type="primary" block :loading="editSaving" @click="saveTemplate">保存</MButton>
        </div>
      </template>
    </MSheet>

    <!-- 分享到广场抽屉 -->
    <MSheet v-model:show="shareVisible" title="分享到提示词广场">
      <div class="form">
        <label class="field">
          <span class="field-label">标题 <i class="req">*</i></span>
          <input v-model="shareForm.title" class="field-input" type="text" placeholder="提示词标题" />
        </label>
        <div class="field-grid">
          <label class="field">
            <span class="field-label">分类</span>
            <input v-model="shareForm.category" class="field-input" type="text" placeholder="如：AI编程" />
          </label>
          <label class="field">
            <span class="field-label">标签</span>
            <input v-model="shareForm.tags" class="field-input" type="text" placeholder="逗号分隔" />
          </label>
        </div>
        <label class="field">
          <span class="field-label">描述</span>
          <textarea v-model="shareForm.description" class="field-textarea" rows="2" placeholder="简短描述提示词用途" />
        </label>
        <label class="field">
          <span class="field-label">内容 <i class="req">*</i></span>
          <textarea v-model="shareForm.content" class="field-textarea" rows="6" placeholder="提示词内容" />
        </label>
        <button
          type="button"
          class="switch-line"
          @click="shareForm.isPublic = !shareForm.isPublic"
        >
          <span>公开到广场</span>
          <span class="switch" :class="{ 'switch--on': shareForm.isPublic }">
            <span class="switch-dot" />
          </span>
        </button>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="shareVisible = false">取消</MButton>
          <MButton type="primary" block :loading="shareSaving" @click="submitShare">分享</MButton>
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

/* 列表 */
.template-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.template-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
  box-shadow: var(--m-shadow-sm);
}

.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-sm);
}

.card-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: #2080f0;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-name--sys {
  color: var(--m-color-fall);
}

.card-type {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  color: #2080f0;
  background: rgba(32, 128, 240, 0.12);
  padding: 2px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
}

.card-type--sys {
  color: var(--m-color-fall);
  background: var(--m-color-fall-light);
}

.card-content {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: var(--m-space-sm);
}

.card-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
}

.card-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.card-actions {
  display: flex;
  gap: var(--m-space-xs);
}

.empty-actions {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  margin-top: var(--m-space-md);
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

.field-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.req {
  color: var(--m-color-rise);
  font-style: normal;
}

.link-btn {
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-color-rise);
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

.preview-box {
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  padding: var(--m-space-sm) var(--m-space-md);
  min-height: 120px;
  max-height: 50vh;
  overflow-y: auto;
}

.preview-empty {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

/* chip 选择 */
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

/* 开关行 */
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

/* 抽屉底部按钮 */
.sheet-footer {
  display: flex;
  gap: var(--m-space-md);
}
</style>
