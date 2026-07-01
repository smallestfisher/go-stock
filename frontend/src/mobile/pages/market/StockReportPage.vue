<script setup>
import { ref, onBeforeMount } from 'vue'
import { GetStockList, StockResearchReport } from '../../../api/app'
import { BrowserOpenURL } from '../../../api/runtime'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MIcon from '../../components/base/MIcon.vue'

// 对齐桌面端 StockResearchReportList.vue：StockResearchReport(code)（空=最近7天全部研报）。
// 顶部搜索选股（复用 KlineAnalysisPage 范式），选中查该股研报，默认展示全部最新。
const list = ref([])
const loading = ref(false)

// 搜索选股
const searchKeyword = ref('')
const searchResults = ref([])
const selected = ref(null) // { code, name }，null=查全部

function pickArray(res) {
  return (Array.isArray(res) ? res : []).filter(it => it && typeof it === 'object')
}

async function loadReports(code = '') {
  loading.value = true
  try {
    const res = await StockResearchReport(code)
    list.value = pickArray(res)
  } catch (e) {
    console.error('加载个股研报失败:', e)
    list.value = []
  } finally {
    loading.value = false
  }
}

// 搜索股票（远程过滤）
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

// 选中搜索结果：查该股研报
function selectStock(stock) {
  selected.value = { code: stock.ts_code, name: stock.name }
  searchKeyword.value = ''
  searchResults.value = []
  loadReports(stock.ts_code)
}

// 清除选中 → 回到全部
function clearSelected() {
  selected.value = null
  loadReports('')
}

// 评级变动：0调高(红)/1调低/2首次/3维持/4无变化
function ratingChangeName(v) {
  return { 0: '调高', 1: '调低', 2: '首次', 3: '维持', 4: '无变化' }[Number(v)] || ''
}

// 打开研报 PDF（对齐桌面 openWin）
function openPdf(infoCode) {
  if (!infoCode) return
  BrowserOpenURL(`https://pdf.dfcfw.com/pdf/H3_${infoCode}_1.pdf?1749744888000.pdf`)
}

function fmtDate(s) {
  return s ? String(s).substring(0, 10) : ''
}

async function handleRefresh() {
  await loadReports(selected.value?.code || '')
}

onBeforeMount(() => {
  loadReports('')
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

    <!-- 当前筛选标的（可清除） -->
    <div v-if="selected" class="selected-bar">
      <span class="selected-text">{{ selected.name }} <span class="selected-code">{{ selected.code }}</span></span>
      <button class="clear-btn" type="button" @click="clearSelected"><MIcon name="close" :size="14" /></button>
    </div>

    <MPullRefresh class="report-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !list.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else>
          <div v-if="list.length" class="report-list">
            <div v-for="(item, i) in list" :key="item.infoCode || i" class="report-item">
              <!-- 标题（可点开 PDF） -->
              <a class="title" @click="openPdf(item.infoCode)">{{ item.title }}</a>

              <!-- 评级 + 评级变动 -->
              <div class="rating-row">
                <span v-if="item.emRatingName" class="rating" :class="{ 'rating--up': item.emRatingName === '增持' }">
                  {{ item.emRatingName }}
                </span>
                <span v-if="item.sRatingName" class="rating-org">{{ item.sRatingName }}</span>
                <span v-if="ratingChangeName(item.ratingChange)" class="rating-change" :class="{ 'rating-change--up': Number(item.ratingChange) === 0 }">
                  {{ ratingChangeName(item.ratingChange) }}
                </span>
              </div>

              <!-- 元信息：行业 / 分析师+机构 / 日期 -->
              <div class="meta">
                <span v-if="item.indvInduName" class="meta-tag">{{ item.indvInduName }}</span>
                <span v-if="item.researcher || item.orgSName" class="meta-text">
                  {{ [item.orgSName, item.researcher].filter(Boolean).join(' · ') }}
                </span>
                <span class="meta-date">{{ fmtDate(item.publishDate) }}</span>
              </div>
            </div>
          </div>
          <MEmpty v-else description="暂无研报" />
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

/* 搜索栏 */
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

/* 当前选中标的条 */
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

.report-refresh {
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

/* 研报列表 */
.report-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.report-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.report-item {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
  padding: var(--m-space-md);
}

.title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  line-height: var(--m-line-height-normal);
}

.title:active {
  color: var(--m-color-rise);
}

.rating-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  flex-wrap: wrap;
}

.rating {
  padding: 1px var(--m-space-sm);
  font-size: var(--m-font-xs);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-secondary);
}

.rating--up {
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.rating-org {
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.rating-change {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.rating-change--up {
  color: var(--m-color-rise);
}

.meta {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  flex-wrap: wrap;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.meta-tag {
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-secondary);
}

.meta-text {
  color: var(--m-text-secondary);
}

.meta-date {
  margin-left: auto;
  font-variant-numeric: tabular-nums;
}
</style>
