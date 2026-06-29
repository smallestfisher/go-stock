<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import {
  GetMCPServerList,
  GetMCPServerByID,
  GetMCPToolsByServerID,
  CreateMCPServer,
  UpdateMCPServer,
  DeleteMCPServer,
  EnableMCPServer,
  TestMCPServer,
} from '../../../api/app'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MButton from '../../components/base/MButton.vue'
import MSheet from '../../components/base/MSheet.vue'
import { toast } from '../../composables/useToast'

// 对齐桌面端 mcp-server-manager.vue：列表 + 搜索 + 状态筛选 + 测试 + 启用切换 + 新建/编辑 + 删除 + 工具列表
// GetMCPServerList(query{page,pageSize,name,status}) → {data:[{id,name,description,url,enable,status,testResult}], total}
// 表单字段：{id,name,description,url,env,enable}

const PAGE_SIZE = 20

const STATUS_OPTIONS = [
  { label: '全部', value: '' },
  { label: '可用', value: 'available' },
  { label: '未测试', value: 'untested' },
  { label: '不可用', value: 'unavailable' },
]

const STATUS_LABEL = {
  available: '可用',
  untested: '未测试',
  unavailable: '不可用',
}

const servers = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const keyword = ref('')
const filterStatus = ref('')
const busyId = ref(null)        // 正在测试/切换的服务 id

const totalPages = computed(() => Math.ceil(total.value / PAGE_SIZE) || 1)
const hasMore = computed(() => page.value < totalPages.value)

async function fetchPage(targetPage, append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const res = await GetMCPServerList({
      page: targetPage,
      pageSize: PAGE_SIZE,
      name: keyword.value,
      status: filterStatus.value,
    })
    const rows = res && Array.isArray(res.data) ? res.data : []
    servers.value = append ? servers.value.concat(rows) : rows
    total.value = (res && res.total) || 0
    page.value = targetPage
  } catch (e) {
    console.error('加载 MCP 服务失败:', e)
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
  fetchPage(1, false)
}

function selectStatus(value) {
  filterStatus.value = value
  fetchPage(1, false)
}

function statusLabel(s) {
  return STATUS_LABEL[s] || s || '未知'
}

// 测试连接
async function testServer(server) {
  if (busyId.value) return
  busyId.value = server.id
  try {
    const res = await TestMCPServer(server.id)
    toast.success(typeof res === 'string' ? res : '测试完成')
    fetchPage(page.value, false)
  } catch (e) {
    toast.error('测试失败: ' + (e.message || ''))
  } finally {
    busyId.value = null
  }
}

// 启用 / 禁用
async function toggleEnable(server) {
  if (busyId.value) return
  busyId.value = server.id
  try {
    const next = !server.enable
    const res = await EnableMCPServer(server.id, next)
    server.enable = next
    toast.success(typeof res === 'string' ? res : (next ? '已启用' : '已禁用'))
  } catch (e) {
    toast.error('操作失败: ' + (e.message || ''))
  } finally {
    busyId.value = null
  }
}

// ===== 工具列表 =====
const toolsVisible = ref(false)
const toolsServer = ref(null)
const tools = ref([])
const toolsLoading = ref(false)

async function openTools(server) {
  toolsServer.value = server
  toolsVisible.value = true
  toolsLoading.value = true
  tools.value = []
  try {
    const res = await GetMCPToolsByServerID(server.id)
    tools.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error('加载工具列表失败:', e)
    toast.error('加载工具失败')
  } finally {
    toolsLoading.value = false
  }
}

// 解析工具参数 schema 为名称列表
function toolParams(tool) {
  if (!tool.paramsSchema) return []
  try {
    const schema = JSON.parse(tool.paramsSchema)
    const props = schema.properties || {}
    const required = schema.required || []
    return Object.entries(props).map(([name, prop]) => ({
      name,
      type: prop.type || '-',
      required: required.includes(name),
      description: prop.description || prop.desc || '',
    }))
  } catch {
    return []
  }
}

// ===== 新建 / 编辑 =====
const editVisible = ref(false)
const editForm = ref({ id: null, name: '', description: '', url: '', env: '', enable: true })
const editSaving = ref(false)
const editTitle = computed(() => (editForm.value.id ? '编辑服务' : '新建服务'))

function openCreate() {
  editForm.value = { id: null, name: '', description: '', url: '', env: '', enable: true }
  editVisible.value = true
}

async function openEdit(server) {
  try {
    const detail = await GetMCPServerByID(server.id)
    const s = detail || server
    editForm.value = {
      id: s.id,
      name: s.name || '',
      description: s.description || '',
      url: s.url || '',
      env: s.env || '',
      enable: s.enable !== false,
    }
    editVisible.value = true
  } catch (e) {
    toast.error('获取服务详情失败: ' + (e.message || ''))
  }
}

async function saveServer() {
  const f = editForm.value
  if (!f.name.trim()) {
    toast.warning('请输入服务器名称')
    return
  }
  if (f.env && f.env.trim()) {
    try {
      JSON.parse(f.env)
    } catch {
      toast.warning('环境变量需为合法 JSON')
      return
    }
  }
  editSaving.value = true
  try {
    const payload = {
      id: f.id || null,
      name: f.name,
      description: f.description,
      url: f.url,
      env: f.env,
      enable: f.enable,
    }
    const res = f.id ? await UpdateMCPServer(payload) : await CreateMCPServer(payload)
    if (typeof res === 'string' && !res.includes('成功')) {
      toast.error(res)
    } else {
      toast.success(typeof res === 'string' ? res : '保存成功')
      editVisible.value = false
      fetchPage(1, false)
    }
  } catch (e) {
    toast.error('保存失败: ' + (e.message || ''))
  } finally {
    editSaving.value = false
  }
}

async function removeServer(server) {
  if (!confirm(`确定删除服务「${server.name}」吗？`)) return
  try {
    const res = await DeleteMCPServer(server.id)
    toast.success(typeof res === 'string' ? res : '删除成功')
    fetchPage(1, false)
  } catch (e) {
    toast.error('删除失败: ' + (e.message || ''))
  }
}

onBeforeMount(() => {
  fetchPage(1, false)
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
          placeholder="搜索服务器名称..."
          @keyup.enter="applySearch"
        />
        <MButton type="primary" size="small" @click="openCreate">新建</MButton>
      </div>
      <div class="chip-scroll">
        <button
          v-for="opt in STATUS_OPTIONS"
          :key="opt.value"
          type="button"
          class="chip"
          :class="{ 'chip--active': filterStatus === opt.value }"
          @click="selectStatus(opt.value)"
        >{{ opt.label }}</button>
      </div>
    </div>

    <MPullRefresh class="refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="servers.length" class="server-list">
          <div v-for="server in servers" :key="server.id" class="server-card">
            <div class="card-head">
              <span class="server-name">{{ server.name }}</span>
              <span class="status-tag" :class="`status-tag--${server.status || 'untested'}`">{{ statusLabel(server.status) }}</span>
            </div>
            <div v-if="server.url" class="server-url">{{ server.url }}</div>
            <p v-if="server.description" class="server-desc">{{ server.description }}</p>
            <p v-if="server.testResult" class="server-result" :class="{ 'server-result--ok': server.status === 'available' }">
              {{ server.testResult }}
            </p>
            <div class="card-foot">
              <span class="enable-tag" :class="{ 'enable-tag--on': server.enable }">
                {{ server.enable ? '已启用' : '已禁用' }}
              </span>
              <div class="card-actions">
                <button type="button" class="act-btn" :disabled="busyId === server.id" @click="testServer(server)">
                  {{ busyId === server.id ? '...' : '测试' }}
                </button>
                <button type="button" class="act-btn" @click="openTools(server)">工具</button>
                <button type="button" class="act-btn" :disabled="busyId === server.id" @click="toggleEnable(server)">
                  {{ server.enable ? '禁用' : '启用' }}
                </button>
                <button type="button" class="act-btn" @click="openEdit(server)">编辑</button>
                <button type="button" class="act-btn act-btn--danger" @click="removeServer(server)">删除</button>
              </div>
            </div>
          </div>

          <div class="load-more">
            <button v-if="hasMore" type="button" class="load-more-btn" :disabled="loading" @click="loadMore">
              {{ loading ? '加载中...' : '加载更多' }}
            </button>
            <span v-else class="load-more-end">共 {{ total }} 个服务</span>
          </div>
        </div>

        <MLoading v-else-if="loading" text="加载中..." vertical class="page-loading" />
        <MEmpty v-else description="还没有配置 MCP 服务">
          <MButton type="primary" @click="openCreate">添加服务</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 新建/编辑抽屉 -->
    <MSheet v-model:show="editVisible" :title="editTitle">
      <div class="form">
        <label class="field">
          <span class="field-label">服务器名称 <i class="req">*</i></span>
          <input v-model="editForm.name" class="field-input" type="text" placeholder="请输入服务器名称" />
        </label>
        <label class="field">
          <span class="field-label">URL</span>
          <input v-model="editForm.url" class="field-input" type="text" placeholder="如：http://localhost:8080 或 SSE 端点" />
        </label>
        <label class="field">
          <span class="field-label">描述</span>
          <textarea v-model="editForm.description" class="field-textarea" rows="2" placeholder="服务器描述（可选）" />
        </label>
        <label class="field">
          <span class="field-label">环境变量</span>
          <textarea v-model="editForm.env" class="field-textarea" rows="3" placeholder='JSON 对象，如：{"API_KEY": "xxx"}' />
        </label>
        <button type="button" class="switch-line" @click="editForm.enable = !editForm.enable">
          <span>启用服务</span>
          <span class="switch" :class="{ 'switch--on': editForm.enable }"><span class="switch-dot" /></span>
        </button>
      </div>
      <template #footer>
        <div class="sheet-footer">
          <MButton type="default" block @click="editVisible = false">取消</MButton>
          <MButton type="primary" block :loading="editSaving" @click="saveServer">保存</MButton>
        </div>
      </template>
    </MSheet>

    <!-- 工具列表抽屉 -->
    <MSheet v-model:show="toolsVisible" :title="`${toolsServer?.name || '服务'} · 工具`">
      <MLoading v-if="toolsLoading" text="加载中..." vertical class="page-loading" />
      <div v-else-if="tools.length" class="tool-list">
        <div v-for="tool in tools" :key="tool.id" class="tool-item">
          <div class="tool-name">{{ tool.toolName }}</div>
          <p v-if="tool.description" class="tool-desc">{{ tool.description }}</p>
          <div v-if="toolParams(tool).length" class="tool-params">
            <span v-for="p in toolParams(tool)" :key="p.name" class="param-tag" :class="{ 'param-tag--req': p.required }">
              {{ p.name }}<i v-if="p.required">*</i>
            </span>
          </div>
          <span v-else class="tool-noparam">无参数</span>
        </div>
      </div>
      <MEmpty v-else description="暂无工具信息，请先测试连接以获取工具列表" />
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
.server-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.server-card {
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
  margin-bottom: var(--m-space-xs);
}

.server-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-tag {
  flex-shrink: 0;
  font-size: var(--m-font-xs);
  padding: 2px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
}

.status-tag--available {
  background: var(--m-color-fall-light);
  color: var(--m-color-fall);
}

.status-tag--untested {
  background: var(--m-divider-color);
  color: var(--m-text-tertiary);
}

.status-tag--unavailable {
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.server-url {
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
  word-break: break-all;
  margin-bottom: var(--m-space-xs);
}

.server-desc {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  margin-bottom: var(--m-space-xs);
}

.server-result {
  font-size: var(--m-font-xs);
  color: var(--m-color-rise);
  line-height: var(--m-line-height-normal);
  margin-bottom: var(--m-space-sm);
  word-break: break-all;
}

.server-result--ok {
  color: var(--m-color-fall);
}

.card-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-sm);
  flex-wrap: wrap;
}

.enable-tag {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.enable-tag--on {
  color: var(--m-color-fall);
}

.card-actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
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

.act-btn:disabled {
  opacity: 0.5;
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

/* 工具列表 */
.tool-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.tool-item {
  padding-bottom: var(--m-space-md);
  border-bottom: 1px solid var(--m-divider-color);
}

.tool-item:last-child {
  border-bottom: none;
}

.tool-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  margin-bottom: var(--m-space-xs);
}

.tool-desc {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  margin-bottom: var(--m-space-sm);
}

.tool-params {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
}

.param-tag {
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
  background: var(--m-bg-primary);
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
}

.param-tag--req {
  color: var(--m-color-rise);
  background: var(--m-color-rise-light);
}

.param-tag i {
  font-style: normal;
}

.tool-noparam {
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
