<script setup>
import { ref, onBeforeMount } from 'vue'
import { GetStockList, StockNotice } from '../../../api/app'
import { BrowserOpenURL } from '../../../api/runtime'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'

// 对齐桌面端 StockNoticeList.vue：StockNotice(code)（空=全部最新公告）。
// 顶部搜索选股（同个股研报范式）。公告类型按关键词着色（复用桌面 getTypeColor 语义）。
const list = ref([])
const loading = ref(false)

const searchKeyword = ref('')
const searchResults = ref([])
const selected = ref(null)

function pickArray(res) {
  return (Array.isArray(res) ? res : []).filter(it => it && typeof it === 'object')
}

async function loadNotices(code = '') {
  loading.value = true
  try {
    const res = await StockNotice(code)
    list.value = pickArray(res)
  } catch (e) {
    console.error('加载公司公告失败:', e)
    list.value = []
  } finally {
    loading.value = false
  }
}

async function handleSearch() {
  const kw = searchKeyword.value.trim()
  if (!kw) {
    searchResults.value = []
    return
  }
  try {
    const res = await GetStockList(kw)
    searchResults.value = (Array.isArray(res) ? res : []).slice(0, 30)
  } catch {
    searchResults.value = []
  }
}

function selectStock(stock) {
  selected.value = { code: stock.ts_code, name: stock.name }
  searchKeyword.value = ''
  searchResults.value = []
  loadNotices(stock.ts_code)
}

function clearSelected() {
  selected.value = null
  loadNotices('')
}

// 公告类型着色：error(红)/warning(橙)/info(默认)，对齐桌面 getTypeColor
function typeLevel(name) {
  const s = String(name || '')
  if (/质押|冻结|解冻|解押|解禁|异常|减持|增发|重大|季度报告|年度报告|澄清公告|风险|终止|复牌|停牌|退市|破产|清算/.test(s)) return 'error'
  if (/回购|重组|诉讼|仲裁|转让|收购|调研|募集/.test(s)) return 'warning'
  return 'info'
}

// 打开公告详情（对齐桌面 openWin，art_code）
function openNotice(artCode) {
  if (!artCode) return
  BrowserOpenURL(`https://np-cnotice-stock.eastmoney.com/api/content/ann?art_code=${artCode}&client_source=web&page_index=1`)
}

function firstCode(item) {
  return Array.isArray(item.codes) && item.codes.length ? item.codes[0] : null
}

function firstColumn(item) {
  return Array.isArray(item.columns) && item.columns.length ? item.columns[0] : null
}

function fmtDate(s) {
  return s ? String(s).substring(0, 10) : ''
}

async function handleRefresh() {
  await loadNotices(selected.value?.code || '')
}

onBeforeMount(() => {
  loadNotices('')
})
</script>

<template>
  <div class="page">
    <!-- 搜索栏 -->
    <div class="search-area">
      <input
        v-model="searchKeyword"
        class="search-input"
        type="text"
        placeholder="输入股票名称/代码筛选"
        @input="handleSearch"
      />
      <div v-if="searchKeyword && searchResults.length" class="search-results">
        <div
          v-for="stock in searchResults"
          :key="stock.ts_code"
          class="search-item"
          @click="selectStock(stock)"
        >
          <span class="search-name">{{ stock.name }}</span>
          <span class="search-code">{{ stock.ts_code }}</span>
        </div>
      </div>
    </div>

    <!-- 当前筛选标的 -->
    <div v-if="selected" class="selected-bar">
      <span class="selected-text">{{ selected.name }} <span class="selected-code">{{ selected.code }}</span></span>
      <button class="clear-btn" type="button" @click="clearSelected">✕</button>
    </div>

    <MPullRefresh class="notice-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !list.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else>
          <div v-if="list.length" class="notice-list">
            <div v-for="(item, i) in list" :key="item.art_code || i" class="notice-item">
              <!-- 类型标签 + 标题（可点开） -->
              <div class="title-row">
                <span
                  v-if="firstColumn(item)?.column_name"
                  class="type-tag"
                  :class="`type-tag--${typeLevel(firstColumn(item).column_name)}`"
                >
                  {{ firstColumn(item).column_name }}
                </span>
                <a class="title" @click="openNotice(item.art_code)">{{ item.title }}</a>
              </div>

              <!-- 元信息：公司 / 日期 -->
              <div class="meta">
                <span v-if="firstCode(item)?.short_name" class="meta-company">
                  {{ firstCode(item).short_name }}
                </span>
                <span class="meta-date">{{ fmtDate(item.notice_date) }}</span>
              </div>
            </div>
          </div>
          <MEmpty v-else description="暂无公告" />
        </template>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.search-area {
  position: relative;
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.search-input {
  width: 100%;
  height: 36px;
  padding: 0 var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-md);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-sm);
  outline: none;
}

.search-input:focus {
  border-color: var(--m-color-rise);
}

.search-results {
  position: absolute;
  left: var(--m-space-md);
  right: var(--m-space-md);
  top: calc(100% - 1px);
  background: var(--m-bg-elevated);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-md);
  box-shadow: var(--m-shadow-md);
  max-height: 260px;
  overflow-y: auto;
  z-index: var(--m-z-dropdown);
}

.search-item {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: var(--m-space-md);
  padding: var(--m-space-sm) var(--m-space-md);
}

.search-item:not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.search-item:active {
  background: var(--m-bg-primary);
}

.search-name {
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
}

.search-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.selected-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-md);
  padding: var(--m-space-xs) var(--m-space-md);
  background: var(--m-color-rise-light);
  border-bottom: 1px solid var(--m-divider-color);
}

.selected-text {
  font-size: var(--m-font-sm);
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.selected-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-weight: normal;
}

.clear-btn {
  width: 24px;
  height: 24px;
  border: none;
  border-radius: var(--m-radius-full);
  background: var(--m-bg-card);
  color: var(--m-text-secondary);
  font-size: var(--m-font-xs);
  line-height: 1;
}

.notice-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

.notice-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.notice-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.notice-item {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
  padding: var(--m-space-md);
}

.title-row {
  display: flex;
  align-items: flex-start;
  gap: var(--m-space-sm);
}

.type-tag {
  flex-shrink: 0;
  padding: 1px var(--m-space-sm);
  font-size: var(--m-font-xs);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-secondary);
  line-height: 1.6;
  white-space: nowrap;
}

.type-tag--error {
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.type-tag--warning {
  background: rgba(245, 158, 11, 0.12);
  color: #d97706;
}

.title {
  flex: 1;
  min-width: 0;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  line-height: var(--m-line-height-normal);
}

.title:active {
  color: var(--m-color-rise);
}

.meta {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.meta-company {
  color: var(--m-text-secondary);
}

.meta-date {
  margin-left: auto;
  font-variant-numeric: tabular-nums;
}
</style>
