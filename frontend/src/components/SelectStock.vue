<script setup lang="ts">
import {h, onBeforeMount, onMounted, onUnmounted, ref, reactive, computed} from 'vue'
import {SearchStock, GetHotStrategy, OpenURL, Follow, GetFollowList, GetAllCustomStrategies, SaveCustomStrategy, DeleteCustomStrategy, GetConfig} from "../api/app";
import {useMessage, NText, NTag, NButton, NPopconfirm} from 'naive-ui'
import {BookmarkOutline, TrashOutline, CreateOutline, AddOutline} from "@vicons/ionicons5";
import {EventsEmit} from "../api/runtime";
import StockLightweightKlineChart from "./StockLightweightKlineChart.vue";
import {useDevice} from "../composables/useDevice";
import BottomSheet from "./mobile/BottomSheet.vue";

const message = useMessage()
const {isMobile} = useDevice()
const search = ref('')
const columns = ref([])
const dataList = ref([])
const hotStrategy = ref([])
const customStrategies = ref([])
const traceInfo = ref('')
const tableScrollX = ref(2800)
const leftTab = ref('hot')
const showSaveModal = ref(false)
const darkTheme = ref(false)
const klineModalShow = ref(false)
const klineStockCode = ref('')
const klineStockName = ref('')
const saveForm = reactive({
  id: 0,
  name: '',
  query: '',
  description: '',
  sortOrder: 0,
})

const paginationProps = computed(() => ({
  pageSize: 10,
  prefix: ({itemCount}) => h('span', {style: 'margin-right: 8px'}, [
    '共找到 ',
    h(NTag, {type: 'info', bordered: false, size: 'small'}, {default: () => itemCount}),
    ' 只股',
  ]),
}))

function calculateTableWidth(cols) {
  let totalWidth = 0;
  cols.forEach(col => {
    if (col.children && col.children.length > 0) {
      let childrenWidth = 0;
      col.children.forEach(child => {
        childrenWidth += child.width || child.minWidth || 100;
      });
      totalWidth += Math.max(col.width || col.minWidth || 200, childrenWidth);
    } else {
      totalWidth += col.width || col.minWidth || 120;
    }
  });
  totalWidth += 100;
  return Math.max(totalWidth, 1200);
}

function Search() {
  if (!search.value) {
    message.warning('请输入选股指标或者要求')
    return
  }
  const loading = message.loading("正在获取选股数据...", {duration: 0});
  SearchStock(search.value).then(res => {
    loading.destroy()
    if (res.code == 100) {
      traceInfo.value = res.data.traceInfo.showText
      columns.value = res.data.result.columns.filter(item => !item.hiddenNeed && (item.title != "市场码" && item.title != "市场简称")).map(item => {
        if (item.children) {
          return {
            title: item.title + (item.unit ? '[' + item.unit + ']' : ''),
            key: item.key,
            resizable: true,
            minWidth: 200,
            ellipsis: {tooltip: true},
            children: item.children.filter(item => !item.hiddenNeed).map(item => {
              return {
                title: item.dateMsg,
                key: item.key,
                minWidth: 100,
                resizable: true,
                ellipsis: {tooltip: true},
                sorter: (row1, row2) => {
                  if (isNumeric(row1[item.key]) && isNumeric(row2[item.key])) {
                    return row1[item.key] - row2[item.key];
                  } else {
                    return 'default'
                  }
                },
              }
            })
          }
        } else {
          return {
            title: item.title + (item.unit ? '[' + item.unit + ']' : ''),
            key: item.key,
            resizable: true,
            minWidth: 120,
            ellipsis: {tooltip: true},
            sorter: (row1, row2) => {
              if (isNumeric(row1[item.key]) && isNumeric(row2[item.key])) {
                return row1[item.key] - row2[item.key];
              } else {
                return 'default'
              }
            },
          }
        }
      })
      columns.value.push({
        title: '操作',
        key: 'actions',
        width: 130,
        fixed: 'right',
        render: (row) => {
          return h('div', {style: 'display:flex;gap:4px;'}, [
            h(
              NButton,
              {
                size: 'tiny',
                type: 'info',
                onClick: () => showStockKline(row)
              },
              {default: () => 'K线'}
            ),
            h(
              NButton,
              {
                strong: true,
                tertiary: true,
                size: 'small',
                type: 'warning',
                style: 'font-size: 14px; padding: 0 10px;',
                onClick: () => handleFollow(row)
              },
              {default: () => '关注'}
            )
          ])
        }
      });
      dataList.value = res.data.result.dataList
      tableScrollX.value = calculateTableWidth(columns.value);
    } else {
      if (res.msg) {
        message.error(res.msg)
      }
      if (res.message) {
        message.error(res.message)
      }
    }
  }).catch(err => {
    message.error(err)
  })
}

function toEastMoneyCode(stockCode, marketShortName) {
  const m = (marketShortName || '').toUpperCase()
  if (m === 'SH' || m === 'SZ' || m === 'BJ') return stockCode + '.' + m
  if (m === 'HK') return stockCode + '.HK'
  if (m === 'US') return stockCode + '.US'
  if (/^(6|5)/.test(stockCode)) return stockCode + '.SH'
  if (/^(8|4)/.test(stockCode)) return stockCode + '.BJ'
  // 纯字母代码视为美股（如 AAPL）
  if (/^[a-zA-Z]+$/.test(stockCode)) return stockCode.toUpperCase() + '.US'
  return stockCode + '.SZ'
}

function showStockKline(row) {
  const stockCode = row.SECURITY_CODE
  const stockName = row.SECURITY_SHORT_NAME
  const em = toEastMoneyCode(stockCode, row.MARKET_SHORT_NAME)
  if (!em) {
    message.warning('当前代码暂不支持K线图')
    return
  }
  klineStockCode.value = em
  klineStockName.value = stockName || ''
  klineModalShow.value = true
}

function handleFollow(row) {
  let code = row.MARKET_SHORT_NAME.toLowerCase() + row.SECURITY_CODE
  Follow(code).then(result => {
    if (result === "关注成功") {
      message.success(result)
    } else {
      message.error(result)
    }
  });
}

function isNumeric(value) {
  return !isNaN(parseFloat(value)) && isFinite(value);
}

function flattenMobileColumns(list) {
  const result = []
  for (const col of list || []) {
    if (col.key === 'actions') continue
    if (col.children && col.children.length) {
      for (const child of col.children) {
        result.push({
          title: `${col.title}/${child.title}`,
          key: child.key
        })
      }
    } else {
      result.push({
        title: col.title,
        key: col.key
      })
    }
  }
  return result
}

function getMobileResultFields(row) {
  return flattenMobileColumns(columns.value)
    .filter(col => !['SECURITY_CODE', 'SECURITY_SHORT_NAME', 'MARKET_SHORT_NAME', 'SERIAL'].includes(col.key))
    .map(col => ({...col, value: row[col.key]}))
    .filter(item => item.value !== undefined && item.value !== null && item.value !== '')
    .slice(0, 6)
}

onBeforeMount(() => {
  GetConfig().then(result => {
    if (result.darkTheme) darkTheme.value = true
  })
  GetHotStrategy().then(res => {
    if (res.code == 1) {
      hotStrategy.value = res.data
      search.value = hotStrategy.value[0].question
      Search()
    }
  }).catch(err => {
    message.error(err)
  })
  loadCustomStrategies()
})

function loadCustomStrategies() {
  GetAllCustomStrategies().then(res => {
    customStrategies.value = res || []
  }).catch(err => {
    message.error(err)
  })
}

function DoSearch(question) {
  search.value = question
  Search()
}

function openSaveModal(isEdit = false, strategy = null) {
  if (isEdit && strategy) {
    saveForm.id = strategy.id
    saveForm.name = strategy.name
    saveForm.query = strategy.query
    saveForm.description = strategy.description || ''
    saveForm.sortOrder = strategy.sortOrder || 0
  } else {
    saveForm.id = 0
    saveForm.name = ''
    saveForm.query = search.value
    saveForm.description = ''
    saveForm.sortOrder = 0
  }
  showSaveModal.value = true
}

function handleSaveStrategy() {
  if (!saveForm.name.trim()) {
    message.warning('请输入策略名称')
    return
  }
  if (!saveForm.query.trim()) {
    message.warning('请输入选股条件')
    return
  }
  SaveCustomStrategy({
    id: saveForm.id || 0,
    name: saveForm.name,
    query: saveForm.query,
    description: saveForm.description,
    sortOrder: saveForm.sortOrder,
  }).then(res => {
    message.success(res)
    showSaveModal.value = false
    loadCustomStrategies()
  }).catch(err => {
    message.error(err)
  })
}

function handleDeleteStrategy(id) {
  DeleteCustomStrategy(id).then(res => {
    message.success(res)
    loadCustomStrategies()
  }).catch(err => {
    message.error(err)
  })
}

function openCenteredWindow(url, width, height) {
  const left = (window.screen.width - width) / 2;
  const top = (window.screen.height - height) / 2;
  window.open(
      url,
      'centeredWindow',
      `width=${width},height=${height},left=${left},top=${top},location=no,menubar=no,toolbar=no,display=standalone`
  )
}
</script>

<template>
  <div class="select-stock-page">
  <n-grid class="select-stock-layout" :cols="24" style="max-height: calc(100vh - 165px)">
    <n-gi class="select-stock-sidebar" :span="4">
      <n-tabs v-model:value="leftTab" type="segment" size="small" style="margin-bottom: 4px;">
        <n-tab name="hot">热门策略</n-tab>
        <n-tab name="custom">我的策略</n-tab>
      </n-tabs>

      <n-list bordered style="text-align: left;" hoverable clickable v-show="leftTab==='hot'">
        <n-scrollbar style="max-height: calc(100vh - 210px);">
          <n-list-item v-for="item in hotStrategy" :key="item.rank" @click="DoSearch(item.question)">
            <n-ellipsis line-clamp="1" :tooltip="true">
              <n-tag size="small" :bordered="false" type="info">#{{ item.rank }}</n-tag>
              <n-text type="warning">{{ item.question }}</n-text>
              <template #tooltip>
                <div style="text-align: center;max-width: 180px">
                  <n-text type="warning">{{ item.question }}</n-text>
                </div>
              </template>
            </n-ellipsis>
          </n-list-item>
        </n-scrollbar>
      </n-list>

      <div v-show="leftTab==='custom'">
        <n-scrollbar style="max-height: calc(100vh - 250px);">
          <n-list bordered hoverable clickable v-if="customStrategies.length > 0">
            <n-list-item v-for="item in customStrategies" :key="item.id">
              <template #suffix>
                <n-flex :size="2" align="center">
                  <n-button text type="info" size="small" @click.stop="openSaveModal(true, item)">
                    <template #icon><n-icon :component="CreateOutline"/></template>
                  </n-button>
                  <n-popconfirm @positive-click="handleDeleteStrategy(item.id)">
                    <template #trigger>
                      <n-button text type="error" size="small" @click.stop>
                        <template #icon><n-icon :component="TrashOutline"/></template>
                      </n-button>
                    </template>
                    确定删除策略「{{ item.name }}」吗？
                  </n-popconfirm>
                </n-flex>
              </template>
              <div @click="DoSearch(item.query)" style="cursor: pointer;">
                <n-ellipsis line-clamp="1" :tooltip="true">
                  <n-tag size="small" :bordered="false" type="success">
                    <template #icon><n-icon :component="BookmarkOutline" size="12"/></template>
                  </n-tag>
                  <n-text strong>{{ item.name }}</n-text>
                  <template #tooltip>
                    <div style="max-width: 200px">
                      <div><n-text strong>{{ item.name }}</n-text></div>
                      <div v-if="item.description" style="margin-top:2px"><n-text depth="3">{{ item.description }}</n-text></div>
                      <div style="margin-top:2px"><n-text type="warning">{{ item.query }}</n-text></div>
                    </div>
                  </template>
                </n-ellipsis>
                <n-ellipsis line-clamp="1" style="margin-top: 2px;">
                  <n-text depth="3" style="font-size: 12px;">{{ item.query }}</n-text>
                </n-ellipsis>
              </div>
            </n-list-item>
          </n-list>
          <n-empty v-else description="暂无自定义策略" style="margin-top: 40px;"/>
        </n-scrollbar>
        <n-button block dashed type="primary" size="small" @click="openSaveModal(false)" style="margin-top: 4px;">
          <template #icon><n-icon :component="AddOutline"/></template>
          添加策略
        </n-button>
      </div>
    </n-gi>
    <n-gi class="select-stock-main" :span="20">
      <div style="">
        <n-input-group class="select-stock-search" style="text-align: left">
          <n-input :rows="1" clearable v-model:value="search" placeholder="请输入选股指标或者要求" @keyup.enter="Search"/>
          <n-button type="primary" @click="Search">搜索A股</n-button>
          <n-button type="warning" @click="openSaveModal(false)" :disabled="!search">
            <template #icon><n-icon :component="BookmarkOutline" size="16"/></template>
            保存策略
          </n-button>
        </n-input-group>
      </div>
      <div v-if="traceInfo" style="margin: 5px 0">
        <n-ellipsis line-clamp="1" :tooltip="true">
          <n-text type="info" :bordered="false">选股条件：</n-text>
          <n-text type="warning" :bordered="true">{{ traceInfo }}</n-text>
          <template #tooltip>
            <div style="text-align: center;max-width: 580px">
              <n-text type="warning">{{ traceInfo }}</n-text>
            </div>
          </template>
        </n-ellipsis>
      </div>
      <n-data-table
          class="select-stock-table desktop-only"
          :striped="true"
          flex-height
          size="small"
          :columns="columns"
          :data="dataList"
          :pagination="paginationProps"
          :scroll-x="tableScrollX"
          style="height: calc(100vh - 240px)"
          :render-cell="(value, rowData, column) => {
        if(column.key=='SECURITY_CODE'||column.key=='SERIAL'){
          return h(NText, { type: 'info',border: false }, { default: () => `${value}` })
        }
        if (isNumeric(value)) {
          let type='info';
          if (Number(value)<0){
            type='success';
          }
          if(Number(value)>=0&&Number(value)<=5){
            type='warning';
          }
          if (Number(value)>5){
            type='error';
          }
            return h(NText, { type: type }, { default: () => `${value}` })
        }else{
            if(column.key=='SECURITY_SHORT_NAME'){
              return h(NText, { type: 'info',bordered: false ,size:'small',onClick:()=>{
               openCenteredWindow(`https://quote.eastmoney.com/${rowData.MARKET_SHORT_NAME}${rowData.SECURITY_CODE}.html#fullScreenChart`,1240,700)
              }}, { default: () => `${value}` })
            }else{
              return h(NText, { type: 'info' }, { default: () => `${value}` })
            }
          }
      }"
      />
      <div class="select-stock-mobile-list mobile-only">
        <n-space vertical :size="10">
          <n-card v-for="(item, idx) in dataList" :key="item.SECURITY_CODE || idx" class="select-stock-mobile-card" size="small" :bordered="true">
            <template #header>
              <n-space class="select-stock-mobile-card__title" align="center" :size="8">
                <n-text strong>{{ item.SECURITY_SHORT_NAME || item.SECURITY_NAME_ABBR || '-' }}</n-text>
                <n-text depth="3">{{ item.SECURITY_CODE }}</n-text>
                <n-tag v-if="item.MARKET_SHORT_NAME" size="small" type="info">{{ item.MARKET_SHORT_NAME }}</n-tag>
              </n-space>
            </template>
            <div class="select-stock-mobile-card__fields">
              <div v-for="field in getMobileResultFields(item)" :key="field.key">
                <span>{{ field.title }}</span>
                <n-text :type="isNumeric(field.value) ? (Number(field.value) < 0 ? 'success' : Number(field.value) > 5 ? 'error' : 'warning') : 'default'">
                  {{ field.value }}
                </n-text>
              </div>
            </div>
            <template #action>
              <n-space class="select-stock-mobile-card__actions" :size="8">
                <n-button size="small" type="info" @click="showStockKline(item)">K线</n-button>
                <n-button size="small" type="warning" @click="handleFollow(item)">关注</n-button>
              </n-space>
            </template>
          </n-card>
          <n-empty v-if="dataList.length === 0" description="暂无选股结果" />
        </n-space>
      </div>
    </n-gi>
  </n-grid>

  <!-- 保存策略弹窗（桌面端） -->
  <n-modal v-if="!isMobile" class="select-stock-save-modal" v-model:show="showSaveModal" preset="dialog" :title="saveForm.id ? '编辑策略' : '保存策略'" positive-text="保存" negative-text="取消"
           @positive-click="handleSaveStrategy" style="width: 500px;">
    <n-form label-placement="left" label-width="80">
      <n-form-item label="策略名称">
        <n-input v-model:value="saveForm.name" placeholder="请输入策略名称"/>
      </n-form-item>
      <n-form-item label="选股条件">
        <n-input v-model:value="saveForm.query" type="textarea" :rows="3" placeholder="请输入选股条件"/>
      </n-form-item>
      <n-form-item label="策略描述">
        <n-input v-model:value="saveForm.description" type="textarea" :rows="2" placeholder="可选，对策略的简要说明"/>
      </n-form-item>
    </n-form>
  </n-modal>

  <!-- 保存策略（移动端底部抽屉） -->
  <BottomSheet v-else :show="showSaveModal" :title="saveForm.id ? '编辑策略' : '保存策略'" height="68vh" @update:show="(v) => showSaveModal = v">
    <div class="ss-save-sheet">
      <n-form label-placement="top">
        <n-form-item label="策略名称">
          <n-input v-model:value="saveForm.name" placeholder="请输入策略名称"/>
        </n-form-item>
        <n-form-item label="选股条件">
          <n-input v-model:value="saveForm.query" type="textarea" :autosize="{ minRows: 3, maxRows: 6 }" placeholder="请输入选股条件"/>
        </n-form-item>
        <n-form-item label="策略描述">
          <n-input v-model:value="saveForm.description" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="可选，对策略的简要说明"/>
        </n-form-item>
      </n-form>
      <div class="ss-save-sheet__actions">
        <n-button @click="showSaveModal = false">取消</n-button>
        <n-button type="primary" @click="() => { if (!handleSaveStrategy()) showSaveModal = false }">保存</n-button>
      </div>
    </div>
  </BottomSheet>

  <!-- K 线弹窗（桌面端） -->
  <n-modal
    v-if="!isMobile"
    class="select-stock-kline-modal"
    v-model:show="klineModalShow"
    :title="(klineStockName || '') + ' - ' + klineStockCode + ' K线图'"
    preset="card"
    style="width: 1400px;max-width: calc(100vw - 32px);"
    :mask-closable="true"
  >
    <StockLightweightKlineChart
      v-if="klineModalShow && klineStockCode"
      :key="klineStockCode"
      :code="klineStockCode"
      :stock-name="klineStockName"
      :dark-theme="darkTheme"
      :chart-height="460"
    />
  </n-modal>

  <!-- K 线（移动端底部抽屉） -->
  <BottomSheet v-else :show="klineModalShow" :title="(klineStockName || '') + ' · K线'" height="80vh" @update:show="(v) => klineModalShow = v">
    <div class="ss-kline-wrap">
      <StockLightweightKlineChart
          v-if="klineModalShow && klineStockCode"
          :key="klineStockCode"
          :code="klineStockCode"
          :stock-name="klineStockName"
          :dark-theme="darkTheme"
          :chart-height="440"
      />
    </div>
  </BottomSheet>
  </div>
</template>

<style scoped>
@media (max-width: 768px) {
  .select-stock-page {
    padding: 0 10px calc(var(--mobile-bottom-nav-height) + var(--safe-bottom) + 10px);
    text-align: left;
  }

  .select-stock-layout {
    display: grid !important;
    gap: 10px;
    grid-template-columns: minmax(0, 1fr) !important;
    max-height: none !important;
    width: 100% !important;
  }

  .select-stock-sidebar,
  .select-stock-main {
    grid-column: 1 / -1 !important;
    max-width: 100%;
    min-width: 0 !important;
    width: 100% !important;
  }

  .select-stock-layout > :deep(.n-grid-item),
  .select-stock-sidebar,
  .select-stock-main {
    grid-column-end: -1 !important;
    grid-column-start: 1 !important;
  }

  .select-stock-sidebar :deep(.n-list) {
    max-height: 170px;
    overflow: auto;
  }

  .select-stock-sidebar :deep(.n-list-item) {
    padding: 8px 10px;
  }

  .select-stock-sidebar :deep(.n-scrollbar) {
    max-height: 170px !important;
  }

  .select-stock-search {
    display: grid !important;
    gap: 8px;
    grid-template-columns: 1fr;
  }

  .select-stock-search :deep(.n-input),
  .select-stock-search :deep(.n-button) {
    width: 100% !important;
  }

  .select-stock-mobile-list {
    display: block !important;
    margin-top: 10px;
  }

  .select-stock-mobile-card {
    text-align: left;
  }

  .select-stock-mobile-card__title {
    align-items: flex-start !important;
    flex-wrap: wrap !important;
    min-width: 0;
  }

  .select-stock-mobile-card__fields {
    display: grid;
    gap: 8px;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .select-stock-mobile-card__fields > div {
    background: var(--n-color-embedded, rgba(128, 128, 128, 0.06));
    border-radius: 6px;
    display: grid;
    gap: 3px;
    min-width: 0;
    padding: 8px;
  }

  .select-stock-mobile-card__fields span {
    color: var(--n-text-color-3);
    font-size: 12px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .select-stock-mobile-card__fields :deep(.n-text) {
    overflow-wrap: anywhere;
  }

  .select-stock-mobile-card__actions {
    display: grid !important;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    width: 100%;
  }

  .select-stock-mobile-card__actions :deep(.n-button) {
    width: 100%;
  }

  :deep(.select-stock-save-modal.n-modal),
  :deep(.select-stock-kline-modal.n-modal) {
    margin: 0 !important;
    max-width: 100vw !important;
    width: calc(100vw - 12px) !important;
  }

  :deep(.select-stock-save-modal .n-card),
  :deep(.select-stock-kline-modal .n-card) {
    max-height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 12px);
    overflow: auto;
  }
}

/* ============ 移动端抽屉（仅在 isMobile 渲染） ============ */
.ss-save-sheet {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 4px 12px calc(var(--safe-bottom) + 8px);
}

.ss-save-sheet :deep(.n-form-item) {
  display: block;
}

.ss-save-sheet :deep(.n-form-item-label) {
  align-items: flex-start;
  display: flex;
  margin-bottom: 6px;
  min-height: auto;
  padding: 0;
}

.ss-save-sheet__actions {
  display: flex;
  gap: 10px;
  margin-top: 4px;
}

.ss-save-sheet__actions :deep(.n-button) {
  flex: 1 1 0;
}

.ss-kline-wrap {
  padding: 4px 8px 8px;
}
</style>
