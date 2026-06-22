<script setup>
import {computed, h, onBeforeMount, onMounted, ref, reactive} from 'vue'
import {
  GetPromptTemplateList,
  GetConfig,
  AddPromptTemplate,
  DeletePromptTemplate,
  UpdatePromptTemplate
} from "../api/app";
import { EventsEmit } from "../api/runtime";
import {parsePromptPlazaResponse, promptPlazaHeaders, promptPlazaURL} from "../api/promptPlaza";
import {NButton, NInput, NTag, NText, NSwitch, useMessage, useNotification,useDialog, NModal, NCard, NForm, NFormItem, NSpace, NPopover} from "naive-ui";
import { MdEditor, MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

const notify = useNotification()
const message = useMessage()
const dialog = useDialog()
const editorDataRef = reactive({
  darkTheme: false
})
const editorTheme = ref('light')

onBeforeMount(() => {
  GetConfig().then(result => {
    if (result.darkTheme) {
      editorDataRef.darkTheme = true
      editorTheme.value = 'dark'
    }
  })
})

onMounted(() => {
  query({
    page: 1,
    pageSize: paginationReactive.pageSize
  }).then((data) => {
    dataRef.value = data.data
    paginationReactive.page = 1
    paginationReactive.pageCount = data.totalPages
    paginationReactive.itemCount = data.total
    loadingRef.value = false
  })
})

const dataRef = ref([])
const loadingRef = ref(true)

const columnsRef = ref([
  {
    title: '模板名称',
    key: 'name',
    render(row) {
      if (row.type === '模型系统Prompt') {
        return h(NText, { type: "success" }, { default: () => row.name })
      }else{
        return h(NText, { type: "info" }, { default: () => row.name })
      }
    }
  },
  {
    title: '模板类型',
    key: 'type',
    render(row) {
      if (row.type === '模型系统Prompt') {
        return h(NTag, { type: "success" }, { default: () => row.type })
      }else{
        return h(NTag, { type: "info" }, { default: () => row.type })
      }
    }
  },
  {
    title: '创建时间',
    key: 'CreatedAt',
    render(row) {
      return row.CreatedAt.substring(0, 19).replace('T', ' ')
    }
  },
  {
    title: '更新时间',
    key: 'UpdatedAt',
    render(row) {
      return row.UpdatedAt.substring(0, 19).replace('T', ' ')
    }
  },
  {
    title: '模板内容',
    key: 'content',
    width: 200,
    render(row) {
      return h(NPopover, {
        trigger: 'hover',
        placement: 'left',
        showArrow: true,
        style: 'max-width: 800px; max-height: 400px; overflow: hidden',
        scrollable: true
      }, {
        trigger: () => h('span', {
          style: 'display: inline-block; max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; cursor: pointer;'
        }, row.content),
        default: () => h(MdPreview, {
          style:'text-align: left;',
          modelValue: row.content,
          theme: editorTheme.value
        })
      })
    }
  },
  {
    title: '操作',
    width: 260,
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-right: 5px',
            onClick: () => showEditModal(row)
          },
          { default: () => '编辑' }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'info',
            style: 'margin-right: 5px',
            onClick: () => showShareModal(row)
          },
          { default: () => '分享' }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            onClick: () => deletePromptTemplate(row.ID)
          },
          { default: () => '删除' }
        )
      ]
    }
  }
])

const paginationReactive = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 12,
  itemCount: 0,
  prefix({ itemCount }) {
    return `${itemCount} 条记录`
  }
})

const modalDataRef = reactive({
  visible: false,
  isEdit: false,
  formData: {
    ID: 0,
    name: '',
    type: '',
    content: ''
  }
})

const shareDataRef = reactive({
  visible: false,
  title: '',
  content: '',
  description: '',
  category: '',
  tags: '',
  isPublic: true,
  loading: false
})

function query({ page, pageSize = 10, name = "", type = "", content = "" }) {
  return new Promise((resolve) => {
    GetPromptTemplateList({
      "page": page,
      "pageSize": pageSize,
      "name": name,
      "type": type,
      "content": content
    }).then((res) => {
      resolve({
        page,
        data: res.list,
        total: res.total,
        totalPages: res.totalPages
      })
    })
  })
}

function handlePageChange(currentPage) {
  if (!loadingRef.value) {
    loadingRef.value = true
    query({
      page: currentPage,
      pageSize: paginationReactive.pageSize,
      name: searchFormRef.name,
      type: searchFormRef.type,
      content: searchFormRef.content
    }).then((data) => {
      dataRef.value = data.data
      paginationReactive.page = currentPage
      paginationReactive.pageCount = data.totalPages
      paginationReactive.itemCount = data.total
      loadingRef.value = false
    })
  }
}
const promptTypeOptions = [
  {label: "模型系统Prompt", value: '模型系统Prompt'},
  {label: "模型用户Prompt", value: '模型用户Prompt'},]
const searchFormRef = reactive({
  name: "",
  type: null,
  content: ""
})

function handleSearch() {
  if (!loadingRef.value) {
    loadingRef.value = true
    query({
      page: 1,
      pageSize: paginationReactive.pageSize,
      name: searchFormRef.name,
      type: searchFormRef.type,
      content: searchFormRef.content
    }).then((data) => {
      dataRef.value = data.data
      paginationReactive.page = data.page
      paginationReactive.pageCount = data.totalPages
      paginationReactive.itemCount = data.total
      loadingRef.value = false
    })
  }
}

function formatTemplateTime(timeStr) {
  if (!timeStr) return '-'
  return String(timeStr).substring(0, 19).replace('T', ' ')
}

function showAddModal() {
  modalDataRef.isEdit = false
  modalDataRef.formData = {
    ID: 0,
    name: '',
    type: '',
    content: ''
  }
  modalDataRef.visible = true
}

function showEditModal(row) {
  modalDataRef.isEdit = true
  modalDataRef.formData = {
    ID: row.ID,
    name: row.name,
    type: row.type,
    content: row.content
  }
  modalDataRef.visible = true
}

function savePromptTemplate() {
  if (!modalDataRef.formData.name || !modalDataRef.formData.type || !modalDataRef.formData.content) {
    message.warning('请填写完整信息' )
    return
  }

  const apiCall = modalDataRef.isEdit ? UpdatePromptTemplate : AddPromptTemplate
  apiCall(modalDataRef.formData).then((res) => {
    message.info( res )
    modalDataRef.visible = false
    handleSearch()
    EventsEmit('promptTemplatesChanged')
  })
}

function deletePromptTemplate(id) {

  dialog.warning({
    title: '提示',
    content: '确定要删除这个模板吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      DeletePromptTemplate(id).then((res) => {
        message.info( res )
        handleSearch()
        EventsEmit('promptTemplatesChanged')
      })
    }
  })
}

async function showShareModal(row) {
  shareDataRef.title = row.name || ''
  shareDataRef.content = row.content || ''
  shareDataRef.description = ''
  shareDataRef.category = row.type || ''
  shareDataRef.tags = ''
  shareDataRef.isPublic = true
  shareDataRef.visible = true
}

async function handleShare() {
  if (!shareDataRef.title || !shareDataRef.content) {
    message.warning('标题和内容不能为空')
    return
  }
  const token = localStorage.getItem('promptPlazaToken')
  if (!token) {
    message.warning('请先在"提示词广场"登录后再分享')
    return
  }
  shareDataRef.loading = true
  try {
    const resp = await fetch(promptPlazaURL('/prompts'), {
      method: 'POST',
      headers: promptPlazaHeaders(token),
      body: JSON.stringify({
        title: shareDataRef.title,
        content: shareDataRef.content,
        description: shareDataRef.description,
        category: shareDataRef.category,
        tags: shareDataRef.tags,
        isPublic: shareDataRef.isPublic
      })
    })
    await parsePromptPlazaResponse(resp)
    message.success('分享成功！')
    shareDataRef.visible = false
  } catch (e) {
    message.error('分享失败: ' + e.message)
  } finally {
    shareDataRef.loading = false
  }
}
</script>

<template>
  <div class="prompt-template-page">
    <!-- 搜索区域 -->
    <n-space class="prompt-template-search" vertical style="margin-bottom: 16px">
      <n-space class="prompt-template-search__inner">
        <n-input v-model:value="searchFormRef.name" placeholder="模板名称" clearable @keyup.enter="handleSearch" />
        <n-select class="prompt-template-type-select" v-model:value="searchFormRef.type" :options="promptTypeOptions" placeholder="请选择提示词类型" clearable/>
        <n-input v-model:value="searchFormRef.content" placeholder="内容关键词" clearable @keyup.enter="handleSearch" />
        <n-button type="success" @click="handleSearch">搜索</n-button>
        <n-button type="warning" @click="showAddModal">新增模板</n-button>
      </n-space>
    </n-space>

    <!-- 数据表格 -->
    <n-data-table
      class="prompt-template-table desktop-only"
      remote
      size="small"
      :columns="columnsRef"
      :data="dataRef"
      :loading="loadingRef"
      :pagination="paginationReactive"
      :row-key="(rowData) => rowData.ID"
      @update:page="handlePageChange"
      flex-height
      style="height: calc(100vh - 250px)"
    />

    <div class="prompt-template-mobile-list mobile-only">
      <n-spin :show="loadingRef">
        <n-space vertical :size="10">
          <n-card
            v-for="item in dataRef"
            :key="item.ID"
            class="prompt-template-mobile-card"
            size="small"
            :bordered="true"
          >
            <template #header>
              <n-space class="prompt-template-mobile-card__title" align="center" :size="8">
                <n-text strong>{{ item.name }}</n-text>
                <n-tag :type="item.type === '模型系统Prompt' ? 'success' : 'info'" size="small">
                  {{ item.type }}
                </n-tag>
              </n-space>
            </template>
            <div class="prompt-template-mobile-card__content">
              {{ item.content }}
            </div>
            <div class="prompt-template-mobile-card__meta">
              <span>创建 {{ formatTemplateTime(item.CreatedAt) }}</span>
              <span>更新 {{ formatTemplateTime(item.UpdatedAt) }}</span>
            </div>
            <template #action>
              <n-space class="prompt-template-mobile-card__actions" :size="8">
                <n-button size="small" type="primary" @click="showEditModal(item)">编辑</n-button>
                <n-button size="small" type="info" @click="showShareModal(item)">分享</n-button>
                <n-button size="small" type="error" @click="deletePromptTemplate(item.ID)">删除</n-button>
              </n-space>
            </template>
          </n-card>
          <n-empty v-if="!loadingRef && dataRef.length === 0" description="暂无提示词模板" />
        </n-space>
      </n-spin>
      <n-space justify="center" style="margin-top: 12px">
        <n-pagination
          v-model:page="paginationReactive.page"
          :page-count="paginationReactive.pageCount"
          :page-size="paginationReactive.pageSize"
          @update:page="handlePageChange"
        />
      </n-space>
    </div>

    <!-- 编辑/新增模态框 -->
    <n-modal class="prompt-template-edit-modal" v-model:show="modalDataRef.visible" preset="card" style="width: 1100px;text-align: left" :title="modalDataRef.formData.ID>0?'修改':'新增'+'Prompt模板'">
      <n-form :model="modalDataRef.formData" label-placement="left" label-width="80">
        <n-form-item label="模板名称" required>
          <n-input v-model:value="modalDataRef.formData.name" placeholder="请输入模板名称" />
        </n-form-item>
        <n-form-item label="模板类型" required>
          <n-select v-model:value="modalDataRef.formData.type" :options="promptTypeOptions" placeholder="请选择提示词类型"/>
        </n-form-item>
        <n-form-item label="模板内容" required>
          <MdEditor
            v-model="modalDataRef.formData.content"
            style="height: 400px"
            :theme="editorTheme"
            :preview="true"
            :toolbarsExclude="['github', 'htmlPreview', 'catalog', 'save']"
            placeholder="请输入模板内容"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="modalDataRef.visible = false">取消</n-button>
          <n-button type="primary" @click="savePromptTemplate">保存</n-button>
        </n-space>
      </template>
    </n-modal>

    <n-modal class="prompt-template-share-modal" v-model:show="shareDataRef.visible" preset="card" style="width: 700px;text-align: left" title="分享到提示词广场">
      <n-form :model="shareDataRef" label-placement="left" label-width="80">
        <n-form-item label="标题" required>
          <n-input v-model:value="shareDataRef.title" placeholder="提示词标题" />
        </n-form-item>
        <n-space :size="8">
          <n-form-item label="分类" label-placement="left" style="width: 300px">
            <n-input v-model:value="shareDataRef.category" placeholder="如: AI编程, 数据分析" />
          </n-form-item>
          <n-form-item label="标签" label-placement="left" style="width: 300px">
            <n-input v-model:value="shareDataRef.tags" placeholder="逗号分隔" />
          </n-form-item>
        </n-space>
        <n-form-item label="描述">
          <n-input v-model:value="shareDataRef.description" type="textarea" :rows="2" placeholder="简短描述提示词用途" />
        </n-form-item>
        <n-form-item label="内容" required>
          <n-input v-model:value="shareDataRef.content" type="textarea" :rows="6" placeholder="提示词内容" />
        </n-form-item>
        <n-form-item label="公开">
          <n-space align="center">
            <n-switch v-model:value="shareDataRef.isPublic" />
          </n-space>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="shareDataRef.visible = false">取消</n-button>
          <n-button type="primary" :loading="shareDataRef.loading" @click="handleShare">分享</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
.prompt-template-type-select {
  width: 200px;
}

:deep(.md-editor) {
  text-align: left;
}
:deep(.n-popover .md-editor-preview) {
  padding: 8px 12px;
}
:deep(.n-popover .md-editor-preview-wrapper) {
  padding: 0;
}

@media (max-width: 768px) {
  .prompt-template-page {
    padding: 0 10px calc(var(--mobile-bottom-nav-height) + var(--safe-bottom) + 10px);
    text-align: left;
  }

  .prompt-template-search,
  .prompt-template-search__inner {
    width: 100%;
  }

  .prompt-template-search__inner {
    display: grid !important;
    gap: 8px !important;
    grid-template-columns: 1fr;
  }

  .prompt-template-search__inner :deep(.n-input),
  .prompt-template-search__inner :deep(.n-select),
  .prompt-template-search__inner :deep(.n-button),
  .prompt-template-type-select {
    width: 100% !important;
  }

  .prompt-template-mobile-list {
    display: block !important;
  }

  .prompt-template-mobile-card {
    text-align: left;
  }

  .prompt-template-mobile-card__title {
    align-items: flex-start !important;
    flex-wrap: wrap !important;
    min-width: 0;
  }

  .prompt-template-mobile-card__title :deep(.n-text) {
    line-height: 1.35;
    overflow-wrap: anywhere;
  }

  .prompt-template-mobile-card__content {
    color: var(--n-text-color-2);
    display: -webkit-box;
    font-size: 13px;
    line-height: 1.55;
    margin-bottom: 10px;
    max-height: 86px;
    overflow: hidden;
    overflow-wrap: anywhere;
    white-space: pre-wrap;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 4;
  }

  .prompt-template-mobile-card__meta {
    color: var(--n-text-color-3);
    display: grid;
    font-size: 12px;
    gap: 4px;
  }

  .prompt-template-mobile-card__actions {
    display: grid !important;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    width: 100%;
  }

  .prompt-template-mobile-card__actions :deep(.n-button) {
    width: 100%;
  }

  :deep(.prompt-template-edit-modal.n-modal),
  :deep(.prompt-template-share-modal.n-modal) {
    margin: 0 !important;
    max-width: 100vw !important;
    width: calc(100vw - 12px) !important;
  }

  :deep(.prompt-template-edit-modal .n-card),
  :deep(.prompt-template-share-modal .n-card) {
    display: flex;
    flex-direction: column;
    max-height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 12px);
  }

  :deep(.prompt-template-edit-modal .n-card__content),
  :deep(.prompt-template-share-modal .n-card__content),
  :deep(.prompt-template-edit-modal .n-card__footer),
  :deep(.prompt-template-share-modal .n-card__footer) {
    padding: 10px 12px;
  }

  :deep(.prompt-template-edit-modal .n-card__content) {
    display: flex;
    flex: 1 1 auto;
    flex-direction: column;
    min-height: 0;
    overflow: hidden;
  }

  :deep(.prompt-template-edit-modal .n-form-item),
  :deep(.prompt-template-share-modal .n-form-item) {
    grid-template-columns: 1fr !important;
  }

  /* md-editor 自适应填满剩余空间，不再用视口减固定像素 */
  :deep(.prompt-template-edit-modal .md-editor) {
    flex: 1 1 auto;
    height: auto !important;
    min-height: 260px;
  }

  :deep(.prompt-template-share-modal .n-space) {
    flex-wrap: wrap !important;
  }

  :deep(.prompt-template-share-modal .n-space > .n-form-item) {
    width: 100% !important;
  }
}
</style>
