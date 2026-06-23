<script setup>
import {computed, h, onBeforeMount, onBeforeUnmount, onMounted,onUnmounted, ref,reactive} from 'vue'
import {GetAIResponseResultList, GetConfig, SaveAsMarkdown, ShareAnalysis,DeleteAIResponseResult} from "../api/app";
import {NAvatar, NButton, NEllipsis, NText, useMessage} from "naive-ui";
import {MdEditor, MdPreview} from 'md-editor-v3';
import {useDevice} from "../composables/useDevice";



onBeforeMount(()=> {
  GetConfig().then(result => {
    if (result.darkTheme) {
      editorDataRef.darkTheme = true
    }
  })
})
onMounted(() => {
  query({
    page: 1,
    pageSize: paginationReactive.pageSize,
    order: "desc",
    keyword: paginationReactive.keyword,
    startDate: paginationReactive.range[0],
    endDate: paginationReactive.range[1]
  }).then((data) => {
    console.log( data)
    dataRef.value = data.data
    paginationReactive.page = 1
    paginationReactive.pageCount = data.pageCount
    paginationReactive.itemCount = data.total
    loadingRef.value = false
  })
})
const message = useMessage()
const {isMobile} = useDevice()
const mdPreviewRef = ref(null)
const mdEditorRef = ref(null)
const editorDataRef = reactive({
  show: false,
  loading: false,
  darkTheme: false,
  chatId: "",
  modelName: "",
  CreatedAt: "",
  stockName: "",
  stockCode: "",
  question: "",
  content: "",
})
const dataRef = ref([])
const loadingRef = ref(true)
const columnsRef = ref([
  {
    title: '分析时间',
    key: 'CreatedAt',
    render(row, index) {
      //2026-01-14T22:13:27.2693252+08:00 格式化为常用时间格式
      return row.CreatedAt.substring(0, 19).replace('T', ' ')
    }
  },
  {
    title: '模型名称',
    key: 'modelName'
  },
  {
    title: '分析对象',
    key: 'stockName'
  },
  {
    title: '提示词',
    key: 'question',
    render(row, index) {
      return h(NEllipsis, { tooltip: true ,style: "max-width: 240px;"}, {default: () => h(NText,{type: "info"},{default: () => row.question}),})
    }
  },
  {
    title: '操作',
    render(row, index) {
      return [h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
            type: 'warning', // 橙色按钮
            style: 'font-size: 14px; padding: 0 10px;', // 稍微大一点的按钮
            onClick: () => showReport(row)
          },
          { default: () => '查看分析报告' }
      ),
      h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
            type: 'error', // 橙色按钮
            style: 'font-size: 14px; padding: 0 10px;', // 稍微大一点的按钮
            onClick: () => deleteAIResponseResult(row.ID)
          },
          { default: () => '删除' }
      ),
      ]
    }
  },
])
const paginationReactive = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 12,
  itemCount: 0,
  keyword: "",
  startDate:"",
  range: [
    new Date(new Date().getTime() - 3 * 24 * 60 * 60 * 1000), // 前3天
    new Date() // 当天
  ],
  prefix({ itemCount }) {
    return `${itemCount} 条记录`
  }
})
const theme = computed(() => {
  return editorDataRef.darkTheme ? 'dark' : 'light'
})
function showReport(row) {

  editorDataRef.show = true
  editorDataRef.chatId = row.chatId
  editorDataRef.modelName = row.modelName
  editorDataRef.CreatedAt = row.CreatedAt.substring(0, 19).replace('T', ' ')
  editorDataRef.stockName = row.stockName
  editorDataRef.stockCode = row.stockCode
  editorDataRef.question = row.question
  editorDataRef.content = row.content
  editorDataRef.loading = false
}

function formatCreatedAt(value) {
  return value ? value.substring(0, 19).replace('T', ' ') : '-'
}

function getStockTitle(row) {
  const name = row.stockName || '未命名对象'
  return row.stockCode ? `${name} ${row.stockCode}` : name
}

function query({
                 page,
                 pageSize = 10,
                 order = 'desc',
                 keyword = "",
                 startDate = "",
                 endDate = ""
               }) {
  return new Promise((resolve) => {

    GetAIResponseResultList({
      "page": page,
      "pageSize": pageSize,
      "modelName":keyword,
      "question":keyword,
      "stockName":keyword,
      "stockCode":keyword,
      "startDate":startDate,
      "endDate":endDate
    }).then((res) => {
      const pagedData =res.list
      const total = res.total
      const pageCount =res.totalPages
      resolve({
        pageCount,
        data: pagedData,
        total
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
      order: "desc",
      keyword: paginationReactive.keyword,
      startDate: formatDate(paginationReactive.range[0]),
      endDate: formatDate(paginationReactive.range[1])
    }).then((data) => {
      dataRef.value = data.data
      paginationReactive.page = currentPage
      paginationReactive.pageCount = data.pageCount
      paginationReactive.itemCount = data.total
      loadingRef.value = false
    })
  }
}
function handleSearch() {
  if (!loadingRef.value) {
    loadingRef.value = true
    query({
      page: paginationReactive?.page ?? 1,
      pageSize: paginationReactive.pageSize,
      order: "desc",
      keyword: paginationReactive.keyword,
      startDate: formatDate(paginationReactive.range[0]),
      endDate: formatDate(paginationReactive.range[1])
    }).then((data) => {
      dataRef.value = data.data
      paginationReactive.page = data.page
      paginationReactive.pageCount = data.pageCount
      paginationReactive.itemCount = data.total
      loadingRef.value = false
    })
  }
}
function share(code, name) {
  ShareAnalysis(code, name).then(msg => {
    //message.info(msg)
    notify.info({
      avatar: () =>
          h(NAvatar, {
            size: 'small',
            round: false,
            src: icon.value
          }),
      title: '分享到社区',
      duration: 1000 * 30,
      content: () => {
        return h('div', {
          style: {
            'text-align': 'left',
            'font-size': '14px',
          }
        }, {default: () => msg})
      },
    })
  })
}

function saveAsMarkdown(code,name) {
  SaveAsMarkdown(code, name).then(result => {
    if(result !== ""){
      message.success(result)
    }
  })
}
async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(editorDataRef.content);
    message.success('分析结果已复制到剪切板');
  } catch (err) {
    message.error('复制失败: ' + err);
  }
}
function formatDate(dateString) {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  // const hours = String(date.getHours()).padStart(2, '0')
  // const minutes = String(date.getMinutes()).padStart(2, '0')
  // const seconds = String(date.getSeconds()).padStart(2, '0')
  //return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  return `${year}-${month}-${day}`
}

function deleteAIResponseResult(id){
  DeleteAIResponseResult(id).then(result => {
    if(result !== ""){
      message.success(result)
    }
    handleSearch()
  })
}
</script>

<template>
  <div class="research-report-page">
    <n-input-group class="research-report-search">
      <n-date-picker  v-model:value="paginationReactive.range" type="daterange"   style="width: 50%"/>
      <n-input clearable placeholder="输入关键词搜索" v-model:value="paginationReactive.keyword"/>
      <n-button type="primary" ghost @click="handleSearch"  @input="handleSearch">
        搜索
      </n-button>
    </n-input-group>
    <n-data-table
        v-if="!isMobile"
        class="research-report-table"
        remote
        size="small"
        :columns="columnsRef"
        :data="dataRef"
        :loading="loadingRef"
        :pagination="paginationReactive"
        :row-key="(rowData)=>rowData.ID"
        @update:page="handlePageChange"
        flex-height
        style="height: calc(100vh - 210px);margin-top: 10px"
    />

    <div v-else class="research-report-mobile-list">
      <n-spin :show="loadingRef">
        <div v-if="dataRef.length > 0" class="research-report-mobile-cards">
          <article v-for="row in dataRef" :key="row.ID" class="research-report-mobile-card">
            <div class="research-report-mobile-card__header">
              <div class="research-report-mobile-card__title">{{ getStockTitle(row) }}</div>
              <div class="research-report-mobile-card__time">{{ formatCreatedAt(row.CreatedAt) }}</div>
            </div>
            <div class="research-report-mobile-card__meta">
              <span>{{ row.modelName || '未知模型' }}</span>
              <span v-if="row.chatId" :title="row.chatId">{{ row.chatId }}</span>
            </div>
            <div class="research-report-mobile-card__question">
              {{ row.question || '无提示词内容' }}
            </div>
            <div class="research-report-mobile-card__actions">
              <n-button size="small" type="warning" secondary @click="showReport(row)">查看分析</n-button>
              <n-button size="small" type="error" secondary @click="deleteAIResponseResult(row.ID)">删除</n-button>
            </div>
          </article>
        </div>
        <n-empty v-else-if="!loadingRef" description="暂无AI分析报告" />
      </n-spin>
      <n-pagination
          v-if="paginationReactive.pageCount > 1"
          v-model:page="paginationReactive.page"
          class="research-report-mobile-pagination"
          :page-count="paginationReactive.pageCount"
          :page-size="paginationReactive.pageSize"
          :item-count="paginationReactive.itemCount"
          size="small"
          @update:page="handlePageChange"
      />
    </div>



    <n-modal class="research-ai-modal" transform-origin="center" v-model:show="editorDataRef.show" preset="card" style="width: 800px;max-width: calc(100vw - 32px);"
             :title="'['+editorDataRef.stockName+']AI分析'">
      <n-spin size="small" :show="editorDataRef.loading">
        <MdPreview  class="research-ai-reader" ref="mdPreviewRef" style="height: 540px;max-height: 60vh;text-align: left;overflow-y: auto;"
                   :modelValue="editorDataRef.content" :theme="theme"/>
      </n-spin>
      <template #footer>
        <n-flex class="research-ai-footer" justify="space-between" ref="tipsRef">
          <n-text type="info" v-if="editorDataRef.chatId">
            <n-tag v-if="editorDataRef.modelName" type="warning" round :title="editorDataRef.chatId" :bordered="false">
              {{ editorDataRef.modelName }}
            </n-tag>
            {{ editorDataRef.CreatedAt }}
          </n-text>
          <n-text type="error">*AI分析结果仅供参考，请以实际行情为准。投资需谨慎，风险自担。</n-text>
        </n-flex>
      </template>
      <template #action>
        <n-flex class="research-ai-actions" justify="right">
          <n-button size="tiny" type="success" @click="copyToClipboard">复制到剪切板</n-button>
          <n-button size="tiny" type="primary" @click="saveAsMarkdown(editorDataRef.stockCode,editorDataRef.stockName)">保存为Markdown文件</n-button>
          <n-button size="tiny" type="error" @click="share(editorDataRef.stockCode,editorDataRef.stockName)">分享到项目社区</n-button>
        </n-flex>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
.research-report-page {
  min-width: 0;
}

@media (max-width: 768px) {
  .research-report-page {
    padding: 0 10px calc(var(--mobile-bottom-nav-height) + var(--safe-bottom) + 10px);
  }

  .research-report-search {
    display: grid !important;
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .research-report-search :deep(.n-date-picker),
  .research-report-search :deep(.n-input),
  .research-report-search :deep(.n-button) {
    width: 100% !important;
  }

  .research-report-table {
    height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 170px) !important;
    margin-top: 10px !important;
  }

  .research-report-table :deep(.n-data-table-th),
  .research-report-table :deep(.n-data-table-td) {
    white-space: nowrap;
  }

  .research-report-mobile-list {
    margin-top: 10px;
    min-height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 230px);
  }

  .research-report-mobile-cards {
    display: grid;
    gap: 10px;
  }

  .research-report-mobile-card {
    background: var(--n-card-color, #fff);
    border: 1px solid var(--n-border-color, #efeff5);
    border-radius: 8px;
    box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
    padding: 12px;
  }

  .research-report-mobile-card__header {
    align-items: flex-start;
    display: grid;
    gap: 4px;
    grid-template-columns: minmax(0, 1fr);
  }

  .research-report-mobile-card__title {
    color: var(--n-text-color, #1f2329);
    font-size: 16px;
    font-weight: 700;
    line-height: 1.35;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .research-report-mobile-card__time,
  .research-report-mobile-card__meta {
    color: var(--n-text-color-3, #667085);
    font-size: 12px;
    line-height: 1.35;
  }

  .research-report-mobile-card__meta {
    display: flex;
    gap: 8px;
    margin-top: 8px;
    min-width: 0;
  }

  .research-report-mobile-card__meta span {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .research-report-mobile-card__question {
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    color: var(--n-text-color-2, #344054);
    display: -webkit-box;
    font-size: 14px;
    line-height: 1.5;
    margin-top: 10px;
    overflow: hidden;
  }

  .research-report-mobile-card__actions {
    display: grid;
    gap: 8px;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin-top: 12px;
  }

  .research-report-mobile-card__actions :deep(.n-button) {
    width: 100%;
  }

  .research-report-mobile-pagination {
    justify-content: center;
    margin-top: 12px;
  }

  :deep(.research-ai-modal.n-modal) {
    margin: 0 !important;
    max-width: 100vw !important;
    width: calc(100vw - 12px) !important;
  }

  :deep(.research-ai-modal .n-card) {
    display: flex;
    flex-direction: column;
    max-height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 12px);
  }

  :deep(.research-ai-modal .n-card__content) {
    display: flex;
    flex: 1 1 auto;
    flex-direction: column;
    min-height: 0;
    overflow: hidden;
    padding: 10px 12px;
  }

  :deep(.research-ai-modal .n-card__content > .n-spin) {
    display: flex;
    flex: 1 1 auto;
    flex-direction: column;
    min-height: 0;
  }

  :deep(.research-ai-modal .n-card__content > .n-spin .n-spin-content) {
    display: flex;
    flex: 1 1 auto;
    flex-direction: column;
    min-height: 0;
  }

  :deep(.research-ai-modal .n-card__footer),
  :deep(.research-ai-modal .n-card__action) {
    flex: 0 0 auto;
    padding: 10px 12px;
  }

  /* reader 自适应填满剩余空间，不再用视口减固定像素 */
  .research-ai-reader {
    flex: 1 1 auto;
    height: auto !important;
    max-height: none !important;
    min-height: 200px;
    overflow-y: auto;
  }

  .research-ai-footer {
    align-items: flex-start !important;
    flex-direction: column;
    gap: 6px !important;
    line-height: 1.45;
  }

  .research-ai-actions {
    align-items: stretch !important;
    flex-wrap: wrap;
    gap: 8px !important;
  }

  .research-ai-actions :deep(.n-button) {
    flex: 1 1 100%;
  }
}
</style>
