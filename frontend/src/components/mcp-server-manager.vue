<template>
  <div class="mcp-server-page">
  <n-space class="mcp-server-toolbar" vertical style="margin-bottom: 12px">
    <n-space class="mcp-server-toolbar__inner">
      <n-input
        v-model:value="searchKeyword"
        placeholder="搜索服务器名称..."
        style="width: 200px"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <n-icon :component="SearchOutline" />
        </template>
      </n-input>

      <n-select
        v-model:value="filterStatus"
        :options="statusOptions"
        placeholder="服务器状态"
        style="width: 120px"
        clearable
      />

      <n-button type="primary" @click="handleSearch">
        搜索
      </n-button>

      <n-button type="warning" @click="handleCreate">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        新建服务器
      </n-button>
    </n-space>
  </n-space>

  <n-data-table
    class="mcp-server-table desktop-only"
    remote
    size="small"
    :columns="columns"
    :data="serverList"
    :loading="loading"
    :pagination="pagination"
    :row-key="(rowData) => rowData.id"
    :expanded-row-keys="expandedKeys"
    :row-props="rowProps"
    @update:expanded-row-keys="handleExpand"
    @update:page="handlePageChange"
    flex-height
    style="height: calc(100vh - 210px); margin-top: 10px"
  />

  <div class="mcp-server-mobile-list mobile-only">
    <n-spin :show="loading">
      <n-space vertical :size="10">
        <n-card v-for="item in serverList" :key="item.id" class="mcp-server-mobile-card" size="small" :bordered="true">
          <template #header>
            <n-space class="mcp-server-mobile-card__title" align="center" :size="8">
              <n-text strong>{{ item.name }}</n-text>
              <n-tag :type="item.enable ? 'success' : 'error'" size="small">{{ item.enable ? '启用' : '禁用' }}</n-tag>
            </n-space>
          </template>
          <template #header-extra>
            <n-tag :type="item.status === 'available' ? 'success' : item.status === 'unavailable' ? 'error' : 'default'" size="small">
              {{ getStatusLabel(item.status) }}
            </n-tag>
          </template>
          <div class="mcp-server-mobile-card__fields">
            <div>
              <span>URL</span>
              <n-text code>{{ item.url || '-' }}</n-text>
            </div>
            <div>
              <span>工具</span>
              <n-text>{{ item.tools?.length || 0 }} 个</n-text>
            </div>
          </div>
          <div v-if="item.description" class="mcp-server-mobile-card__description">{{ item.description }}</div>
          <div v-if="item.testResult" class="mcp-server-mobile-card__result">
            <n-text :type="item.status === 'available' ? 'success' : 'error'">{{ item.testResult }}</n-text>
          </div>
          <div v-if="item.tools?.length" class="mcp-server-mobile-card__tools">
            <n-tag v-for="tool in item.tools.slice(0, 4)" :key="tool.id" size="tiny" type="info" :bordered="false">
              {{ tool.toolName }}
            </n-tag>
          </div>
          <template #action>
            <n-space class="mcp-server-mobile-card__actions" :size="8">
              <n-button size="small" type="info" @click="handleTest(item)">测试</n-button>
              <n-button size="small" :type="item.enable ? 'warning' : 'info'" @click="handleToggleEnable(item)">
                {{ item.enable ? '禁用' : '启用' }}
              </n-button>
              <n-button size="small" type="primary" @click="handleEdit(item)">编辑</n-button>
              <n-popconfirm @positive-click="handleDelete(item.id)">
                <template #trigger>
                  <n-button size="small" type="error">删除</n-button>
                </template>
                确定要删除服务器 "{{ item.name }}" 吗？
              </n-popconfirm>
            </n-space>
          </template>
        </n-card>
        <n-empty v-if="!loading && serverList.length === 0" description="暂无 MCP 服务器" />
      </n-space>
    </n-spin>
    <n-space justify="center" style="margin-top: 12px">
      <n-pagination
        :page="currentPage"
        :page-count="Math.ceil(total / pageSize) || 1"
        :page-size="pageSize"
        @update:page="handlePageChange"
      />
    </n-space>
  </div>

  <!-- 创建/修改服务器（桌面端） -->
  <n-modal
    v-if="!isMobile"
    class="mcp-server-edit-modal"
    v-model:show="showCreateModal"
    :title="editingServer ? '修改服务器' : '创建新服务器'"
    preset="dialog"
    :style="{ width: '750px' }"
    @close="resetForm"
    :z-index="2000"
    to="body"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="130px"
      require-mark-placement="right-hanging"
    >
      <n-form-item label="服务器名称" path="name">
        <n-input v-model:value="formData.name" placeholder="请输入服务器名称" clearable />
      </n-form-item>

      <n-form-item label="描述" path="description">
        <n-input
          v-model:value="formData.description"
          type="textarea"
          :rows="2"
          placeholder="请输入服务器描述（可选）"
          show-count
          maxlength="500"
        />
      </n-form-item>

      <n-form-item label="URL" path="url">
        <n-input v-model:value="formData.url" placeholder="例如：http://localhost:8080 或 SSE 端点地址" clearable />
      </n-form-item>

      <n-form-item label="环境变量" path="env">
        <n-input
          v-model:value="formData.env"
          type="textarea"
          :rows="3"
          placeholder='JSON 对象格式，例如：{"API_KEY": "your-api-key"}'
          show-count
        />
      </n-form-item>

      <n-form-item label="启用状态" path="enable">
        <n-switch v-model:value="formData.enable" size="large">
          <template #checked>
            <n-icon :component="PlayCircleOutline" />
            启用
          </template>
          <template #unchecked>
            <n-icon :component="StopCircleOutline" />
            禁用
          </template>
        </n-switch>
      </n-form-item>
    </n-form>

    <template #action>
      <n-button @click="showCreateModal = false">取消</n-button>
      <n-button type="primary" @click="handleSubmit" :loading="submitting">
        <template #icon>
          <n-icon :component="CheckmarkCircleOutline" />
        </template>
        {{ editingServer ? '修改服务器' : '创建新服务器' }}
      </n-button>
    </template>
  </n-modal>

  <!-- 创建/修改服务器（移动端底部抽屉） -->
  <BottomSheet v-if="isMobile && showCreateModal" :show="showCreateModal" :title="editingServer ? '修改服务器' : '创建新服务器'" height="80vh" @update:show="(v) => showCreateModal = v" @close="resetForm">
    <div class="mcp-form-sheet">
      <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top" require-mark-placement="right-hanging">
        <n-form-item label="服务器名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入服务器名称" clearable />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input v-model:value="formData.description" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入服务器描述（可选）" show-count maxlength="500" />
        </n-form-item>
        <n-form-item label="URL" path="url">
          <n-input v-model:value="formData.url" placeholder="例如：http://localhost:8080 或 SSE 端点地址" clearable />
        </n-form-item>
        <n-form-item label="环境变量" path="env">
          <n-input v-model:value="formData.env" type="textarea" :autosize="{ minRows: 3, maxRows: 6 }" placeholder='JSON 对象格式，例如：{"API_KEY": "your-api-key"}' show-count />
        </n-form-item>
        <n-form-item label="启用状态" path="enable">
          <n-switch v-model:value="formData.enable" size="large">
            <template #checked><n-icon :component="PlayCircleOutline" /> 启用</template>
            <template #unchecked><n-icon :component="StopCircleOutline" /> 禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <div class="mcp-form-sheet__actions">
        <n-button @click="showCreateModal = false">取消</n-button>
        <n-button type="primary" @click="handleSubmit" :loading="submitting">{{ editingServer ? '修改' : '创建' }}</n-button>
      </div>
    </div>
  </BottomSheet>

  <!-- 工具参数详情（桌面端） -->
  <n-modal
    v-if="!isMobile"
    class="mcp-server-tool-modal"
    v-model:show="showToolDetailModal"
    title="工具参数详情"
    preset="card"
    style="width: 850px"
    :z-index="2000"
  >
    <template v-if="currentTool">
      <n-descriptions bordered label-placement="top" :column="1" style="margin-bottom: 16px" content-style="text-align: left">
        <n-descriptions-item label="工具名称">
          <n-text code>{{ currentTool.toolName }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="描述">{{ currentTool.description || '无描述' }}</n-descriptions-item>
      </n-descriptions>

      <template v-if="parsedParams.length > 0">
        <n-text strong style="margin-bottom: 8px; display: block">参数列表</n-text>
        <n-data-table
          :columns="paramDetailColumns"
          :data="parsedParams"
          size="small"
          bordered
          :pagination="false"
        />
      </template>
      <n-text v-else depth="3">此工具无需参数</n-text>

      <n-collapse style="margin-top: 12px" v-if="currentTool.paramsSchema">
        <n-collapse-item title="原始 JSON Schema" name="raw">
          <VueJsonPretty
            :data="parseJSON(currentTool.paramsSchema)"
            :deep="3"
            show-length
            show-line
            collapsed-on-click-bracket
          />
        </n-collapse-item>
      </n-collapse>
    </template>
  </n-modal>

  <!-- 工具参数详情（移动端底部抽屉） -->
  <BottomSheet v-if="isMobile && showToolDetailModal" :show="showToolDetailModal" title="工具参数详情" height="80vh" @update:show="(v) => showToolDetailModal = v">
    <div v-if="currentTool" class="mcp-tool-sheet">
      <div class="mcp-tool-sheet__name"><n-text code>{{ currentTool.toolName }}</n-text></div>
      <div class="mcp-tool-sheet__desc">{{ currentTool.description || '无描述' }}</div>

      <n-divider style="margin: 10px 0" />
      <n-text strong>参数列表</n-text>
      <div v-if="parsedParams.length > 0" class="mcp-params">
        <div v-for="(p, i) in parsedParams" :key="i" class="mcp-param">
          <div class="mcp-param__name">
            <n-text code>{{ p.name }}</n-text>
            <n-tag v-if="p.required" size="tiny" type="error" :bordered="false">必填</n-tag>
          </div>
          <div class="mcp-param__type">{{ p.type }}</div>
          <div class="mcp-param__desc">{{ p.description || '无描述' }}</div>
        </div>
      </div>
      <n-text v-else depth="3">此工具无需参数</n-text>

      <n-collapse style="margin-top: 12px" v-if="currentTool.paramsSchema">
        <n-collapse-item title="原始 JSON Schema" name="raw">
          <VueJsonPretty
              :data="parseJSON(currentTool.paramsSchema)"
              :deep="3"
              show-length
              show-line
              collapsed-on-click-bracket
          />
        </n-collapse-item>
      </n-collapse>
    </div>
  </BottomSheet>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, h } from 'vue'
import {
  NButton,
  NIcon,
  NTag,
  NSpace,
  NPopconfirm,
  NDescriptions,
  NDescriptionsItem,
  NText,
  NDataTable,
  NCollapse,
  NCollapseItem,
  useMessage
} from 'naive-ui'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import {
  SearchOutline,
  AddOutline,
  PlayOutline,
  PauseOutline,
  TrashOutline,
  CreateOutline,
  PlayCircleOutline,
  StopCircleOutline,
  CheckmarkCircleOutline,
  FlashOutline,
  EyeOutline
} from '@vicons/ionicons5'
import {
  CreateMCPServer,
  UpdateMCPServer,
  DeleteMCPServer,
  GetMCPServerByID,
  GetMCPServerList,
  EnableMCPServer,
  TestMCPServer,
  GetMCPToolsByServerID,
  GetAllMCPTools
} from '../api/app'
import {useDevice} from '../composables/useDevice'
import BottomSheet from './mobile/BottomSheet.vue'

const message = useMessage()
const {isMobile} = useDevice()

const formRef = ref(null)

const loading = ref(false)
const submitting = ref(false)
const showCreateModal = ref(false)
const showToolDetailModal = ref(false)
const editingServer = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const expandedKeys = ref([])
const currentTool = ref(null)

const formData = reactive({
  id: null,
  name: '',
  description: '',
  url: '',
  command: '',
  args: '',
  env: '',
  enable: true,
  status: 'untested'
})

const formRules = {
  name: { required: true, message: '请输入服务器名称', trigger: ['input', 'blur'] },
  command: { required: true, message: '请输入命令', trigger: ['input', 'blur'] }
}

const statusOptions = [
  { label: '可用', value: 'available' },
  { label: '未测试', value: 'untested' },
  { label: '不可用', value: 'unavailable' }
]

const getStatusLabel = (status) => {
  switch (status) {
    case 'available':
      return '可用'
    case 'untested':
      return '未测试'
    case 'unavailable':
      return '不可用'
    default:
      return status
  }
}

const formatJSON = (str) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

const parseJSON = (str) => {
  try {
    return JSON.parse(str)
  } catch {
    return str
  }
}

const handleExpand = (keys) => {
  expandedKeys.value = keys
}

const rowProps = (row) => {
  return {
    style: 'cursor: pointer',
    onClick: (e) => {
      if (e.target.closest('.n-button, .n-popconfirm, .n-switch, a')) return
      const key = row.id
      const idx = expandedKeys.value.indexOf(key)
      if (idx === -1) {
        expandedKeys.value = [...expandedKeys.value, key]
      } else {
        expandedKeys.value = expandedKeys.value.filter(k => k !== key)
      }
    }
  }
}

const handleViewToolDetail = (tool) => {
  currentTool.value = tool
  showToolDetailModal.value = true
}

const parsedParams = computed(() => {
  if (!currentTool.value || !currentTool.value.paramsSchema) return []
  try {
    const schema = JSON.parse(currentTool.value.paramsSchema)
    const props = schema.properties || {}
    const required = schema.required || []
    return Object.entries(props).map(([name, prop]) => ({
      name,
      type: prop.type || '-',
      required: required.includes(name),
      description: prop.description || prop.desc || '-',
      enum: prop.enum ? prop.enum.join(', ') : '',
      default: prop.default !== undefined ? String(prop.default) : ''
    }))
  } catch {
    return []
  }
})

const paramDetailColumns = [
  {
    title: '参数名',
    key: 'name',
    width: 180,
    ellipsis: { tooltip: { style: { maxWidth: '300px' } } }
  },
  {
    title: '类型',
    key: 'type',
    width: 80
  },
  {
    title: '必填',
    key: 'required',
    width: 60,
    render(row) {
      return h(NTag, { type: row.required ? 'error' : 'default', size: 'small' }, {
        default: () => row.required ? '是' : '否'
      })
    }
  },
  {
    title: '描述',
    key: 'description',
    ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } }
  },
  {
    title: '枚举值',
    key: 'enum',
    width: 120,
    ellipsis: { tooltip: { style: { maxWidth: '300px' } } }
  },
  {
    title: '默认值',
    key: 'default',
    width: 80,
    ellipsis: { tooltip: true }
  }
]

const columns = [
  {
    type: 'expand',
    renderExpand: (row) => {
      const tools = row.tools
      if (!tools || tools.length === 0) {
        return h(NText, { depth: 3, style: 'padding: 8px 16px' }, { default: () => '暂无工具信息，请先测试连接以获取工具列表' })
      }

      const toolColumns = [
        {
          title: '工具名称',
          key: 'toolName',
          width: 200,
          ellipsis: { tooltip: true }
        },
        {
          title: '描述',
          key: 'description',
          ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } }
        },
        {
          title: '参数',
          key: 'paramsSchema',
          width: 260,
          render(toolRow) {
            if (!toolRow.paramsSchema) {
              return h(NTag, { type: 'default', size: 'small' }, { default: () => '无参数' })
            }
            try {
              const schema = JSON.parse(toolRow.paramsSchema)
              const props = schema.properties || {}
              const required = schema.required || []
              const paramNames = Object.keys(props)
              if (paramNames.length === 0) {
                return h(NTag, { type: 'default', size: 'small' }, { default: () => '无参数' })
              }
              const tags = paramNames.map(name => {
                const prop = props[name]
                const isReq = required.includes(name)
                return h(NTag, {
                  size: 'small',
                  type: isReq ? 'info' : 'default',
                  style: 'margin: 2px'
                }, {
                  default: () => name + (prop.type ? ':' + prop.type : '') + (isReq ? '*' : '')
                })
              })
              return h('div', { style: 'display: flex; flex-wrap: wrap; align-items: center; gap: 0' }, [
                ...tags,
                h(NButton, {
                  size: 'tiny',
                  type: 'info',
                  quaternary: true,
                  style: 'margin-left: 4px',
                  onClick: () => handleViewToolDetail(toolRow)
                }, {
                  icon: () => h(NIcon, { component: EyeOutline })
                })
              ])
            } catch {
              return h(NButton, {
                size: 'tiny',
                type: 'info',
                quaternary: true,
                onClick: () => handleViewToolDetail(toolRow)
              }, {
                icon: () => h(NIcon, { component: EyeOutline }),
                default: () => '查看参数'
              })
            }
          }
        }
      ]

      return h('div', { style: 'padding: 8px 16px' }, [
        h(NDataTable, {
          columns: toolColumns,
          data: tools,
          size: 'small',
          bordered: false,
          rowKey: (t) => t.id,
          pagination: false
        })
      ])
    }
  },
  {
    title: 'ID',
    key: 'id',
    width: 60,
    ellipsis: { tooltip: true }
  },
  {
    title: '服务器名称',
    key: 'name',
    width: 180,
    ellipsis: { tooltip: true }
  },
  {
    title: '描述',
    key: 'description',
    width: 200,
    ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } }
  },
  {
    title: 'URL',
    key: 'url',
    width: 180,
    ellipsis: { tooltip: true },
    render(row) {
      if (!row.url) return h(NText, { depth: 3 }, { default: () => '-' })
      return h(NText, { code: true, depth: 2 }, { default: () => row.url })
    }
  },
  {
    title: '启用',
    key: 'enable',
    width: 70,
    render(row) {
      return h(NTag, { type: row.enable ? 'success' : 'error' }, {
        default: () => (row.enable ? '是' : '否')
      })
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render(row) {
      const typeMap = {
        available: 'success',
        untested: 'default',
        unavailable: 'error'
      }
      return h(NTag, { type: typeMap[row.status] || 'default' }, {
        default: () => getStatusLabel(row.status)
      })
    }
  },
  {
    title: '测试结果',
    key: 'testResult',
    width: 200,
    ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } },
    render(row) {
      if (!row.testResult) return h(NText, { depth: 3 }, { default: () => '-' })
      const isSuccess = row.status === 'available'
      return h(NText, { type: isSuccess ? 'success' : 'error' }, { default: () => row.testResult })
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 280,
    fixed: 'right',
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(
            NButton,
            {
              size: 'tiny',
              type: 'info',
              onClick: () => handleTest(row)
            },
            {
              icon: () => h(NIcon, { component: FlashOutline }),
              default: () => '测试'
            }
          ),
          h(
            NButton,
            {
              size: 'tiny',
              type: row.enable ? 'warning' : 'info',
              onClick: () => handleToggleEnable(row)
            },
            {
              icon: () => h(NIcon, { component: row.enable ? PauseOutline : PlayOutline }),
              default: () => (row.enable ? '禁用' : '启用')
            }
          ),
          h(
            NButton,
            {
              size: 'tiny',
              type: 'primary',
              onClick: () => handleEdit(row)
            },
            {
              icon: () => h(NIcon, { component: CreateOutline }),
              default: () => '编辑'
            }
          ),
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleDelete(row.id)
            },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: 'tiny',
                    type: 'error'
                  },
                  {
                    icon: () => h(NIcon, { component: TrashOutline }),
                    default: () => '删除'
                  }
                ),
              default: () => `确定要删除服务器 "${row.name}" 吗？`
            }
          )
        ]
      })
    }
  }
]

const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  pageCount: Math.ceil(total.value / pageSize.value) || 1,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  prefix: ({ itemCount }) => `共 ${itemCount} 条`,
  onChange: handlePageChange,
  onUpdatePageSize: handlePageSizeChange
}))

const loadServerList = async () => {
  loading.value = true
  try {
    const query = {
      page: currentPage.value,
      pageSize: pageSize.value,
      name: searchKeyword.value,
      status: filterStatus.value
    }

    const result = await GetMCPServerList(query)
    if (result) {
      const servers = result.data || []
      let allTools = []
      try {
        allTools = await GetAllMCPTools() || []
      } catch (error) {
        console.error('加载工具列表失败:', error)
      }
      const toolsMap = {}
      for (const t of allTools) {
        if (!toolsMap[t.mcpServerId]) toolsMap[t.mcpServerId] = []
        toolsMap[t.mcpServerId].push(t)
      }
      for (const server of servers) {
        server.tools = toolsMap[server.id] || []
      }
      serverList.value = servers
      total.value = result.total || 0
    }
  } catch (error) {
    console.error('加载服务器列表失败:', error)
    message.error('加载服务器列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = async () => {
  currentPage.value = 1
  await loadServerList()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadServerList()
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadServerList()
}

const handleTest = async (row) => {
  try {
    const result = await TestMCPServer(row.id)
    message.success(result)
    await loadServerList()
  } catch (error) {
    message.error('测试失败：' + error.message)
  }
}

const handleToggleEnable = async (row) => {
  try {
    const newEnable = !row.enable
    const result = await EnableMCPServer(row.id, newEnable)
    message.success(result)
    await loadServerList()
  } catch (error) {
    message.error('操作失败：' + error.message)
  }
}

const handleCreate = () => {
  editingServer.value = false
  resetForm()
  showCreateModal.value = true
}

const handleEdit = async (row) => {
  editingServer.value = true
  try {
    const server = await GetMCPServerByID(row.id)
    if (server) {
      resetForm()
      formData.id = server.id
      formData.name = server.name
      formData.description = server.description
      formData.url = server.url
      formData.command = server.command
      formData.args = server.args
      formData.env = server.env
      formData.enable = server.enable
      formData.status = server.status
      showCreateModal.value = true
    }
  } catch (error) {
    message.error('获取服务器详情失败：' + error.message)
  }
}

const handleDelete = async (id) => {
  try {
    const result = await DeleteMCPServer(id)
    message.success(result)
    await loadServerList()
  } catch (error) {
    message.error('删除失败：' + error.message)
  }
}

const handleSubmit = async () => {
  try {
    if (formRef.value) {
      try {
        await formRef.value.validate()
      } catch {
        return
      }
    }

    submitting.value = true
    const submitData = { ...formData }

    let result
    if (formData.id) {
      result = await UpdateMCPServer(submitData)
    } else {
      result = await CreateMCPServer(submitData)
    }

    if (result.includes('成功')) {
      message.success(result)
      showCreateModal.value = false
      await loadServerList()
    } else {
      message.error(result)
    }
  } catch (error) {
    message.error('操作失败：' + error.message)
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
    id: null,
    name: '',
    description: '',
    url: '',
    command: '',
    args: '',
    env: '',
    enable: true,
    status: 'stopped'
  })
  if (formRef.value) {
    formRef.value.restoreValidation()
  }
}

const serverList = ref([])

onMounted(async () => {
  await loadServerList()
})
</script>

<style scoped>
@media (max-width: 768px) {
  .mcp-server-page {
    padding: 0 10px calc(var(--mobile-bottom-nav-height) + var(--safe-bottom) + 10px);
    text-align: left;
  }

  .mcp-server-toolbar,
  .mcp-server-toolbar__inner {
    width: 100%;
  }

  .mcp-server-toolbar__inner {
    display: grid !important;
    gap: 8px !important;
    grid-template-columns: 1fr;
  }

  .mcp-server-toolbar__inner :deep(.n-input),
  .mcp-server-toolbar__inner :deep(.n-select),
  .mcp-server-toolbar__inner :deep(.n-button) {
    width: 100% !important;
  }

  .mcp-server-mobile-list {
    display: block !important;
  }

  .mcp-server-mobile-card {
    text-align: left;
  }

  .mcp-server-mobile-card__title {
    align-items: flex-start !important;
    flex-wrap: wrap !important;
    min-width: 0;
  }

  .mcp-server-mobile-card__fields {
    display: grid;
    gap: 8px;
    grid-template-columns: 1fr;
    margin-bottom: 10px;
  }

  .mcp-server-mobile-card__fields > div {
    background: var(--n-color-embedded, rgba(128, 128, 128, 0.06));
    border-radius: 6px;
    display: grid;
    gap: 3px;
    min-width: 0;
    padding: 8px;
  }

  .mcp-server-mobile-card__fields span {
    color: var(--n-text-color-3);
    font-size: 12px;
  }

  .mcp-server-mobile-card__fields :deep(.n-text),
  .mcp-server-mobile-card__description,
  .mcp-server-mobile-card__result {
    overflow-wrap: anywhere;
  }

  .mcp-server-mobile-card__description,
  .mcp-server-mobile-card__result {
    color: var(--n-text-color-2);
    font-size: 12px;
    line-height: 1.45;
    margin-bottom: 8px;
  }

  .mcp-server-mobile-card__tools {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 4px;
  }

  .mcp-server-mobile-card__actions {
    display: grid !important;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    width: 100%;
  }

  .mcp-server-mobile-card__actions :deep(.n-button) {
    width: 100%;
  }

  :deep(.mcp-server-edit-modal.n-modal),
  :deep(.mcp-server-tool-modal.n-modal) {
    margin: 0 !important;
    max-width: 100vw !important;
    width: calc(100vw - 12px) !important;
  }

  :deep(.mcp-server-edit-modal .n-dialog),
  :deep(.mcp-server-tool-modal .n-card) {
    max-height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 12px);
    overflow: auto;
    width: calc(100vw - 12px) !important;
  }

  :deep(.mcp-server-edit-modal .n-form-item) {
    grid-template-columns: 1fr !important;
  }
}

/* ============ 移动端抽屉（仅在 isMobile 渲染） ============ */
.mcp-form-sheet {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 4px 12px calc(var(--safe-bottom) + 8px);
}

.mcp-form-sheet :deep(.n-form-item) {
  display: block;
}

.mcp-form-sheet :deep(.n-form-item-label) {
  align-items: flex-start;
  display: flex;
  margin-bottom: 6px;
  min-height: auto;
  padding: 0;
}

.mcp-form-sheet__actions {
  display: flex;
  gap: 10px;
  margin-top: 4px;
}

.mcp-form-sheet__actions :deep(.n-button) {
  flex: 1 1 0;
}

.mcp-tool-sheet {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 4px 12px calc(var(--safe-bottom) + 8px);
}

.mcp-tool-sheet__name {
  font-size: 16px;
  font-weight: 700;
}

.mcp-tool-sheet__desc {
  color: var(--n-text-color-2, #666);
  font-size: 13px;
  line-height: 1.55;
}

.mcp-params {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}

.mcp-param {
  background: var(--n-color-target, #f5f7fa);
  border-radius: 8px;
  padding: 8px 10px;
}

.mcp-param__name {
  align-items: center;
  display: flex;
  gap: 8px;
  margin-bottom: 2px;
}

.mcp-param__type {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 12px;
}

.mcp-param__desc {
  color: var(--n-text-color-2, #666);
  font-size: 13px;
  line-height: 1.5;
  margin-top: 2px;
  overflow-wrap: anywhere;
}
</style>
