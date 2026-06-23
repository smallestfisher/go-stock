<script setup>
import {h, ref, computed, reactive, onMounted, watch} from "vue";
import {NButton, NText, NFlex, NTag, NDataTable} from "naive-ui";
import {
  FollowFund,
  GetConfig,
  GetFollowedFund,
  GetFundRanking,
  GetFundTop10Holdings,
  OpenURL,
  SearchFundCodes
} from "../api/app";
import {useMessage} from "naive-ui";
import StockLightweightKlineChart from "./StockLightweightKlineChart.vue";
import StockSparkLine from "./stockSparkLine.vue";
import {useDevice} from "../composables/useDevice";
import BottomSheet from "./mobile/BottomSheet.vue";

const message = useMessage()
const {isMobile} = useDevice()

const darkTheme = ref(false)

const marketType = ref('kf')
const rankingFundType = ref('all')
const rankingSortField = ref('jnzf')
const rankingSortOrder = ref('desc')
const rankingLoading = ref(false)
const rankingData = ref([])
const followList = ref([])
const searchKeyword = ref('')
const searchCodes = ref(null)
let searchTimer = null
let skipSortWatch = false

const holdingsModalShow = ref(false)
const holdingsFundCode = ref('')
const holdingsFundName = ref('')
const holdingsData = ref([])
const holdingsLoading = ref(false)

const klineModalShow = ref(false)
const klineStockCode = ref('')
const klineStockName = ref('')

const paginationReactive = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 50,
  itemCount: 0,
  prefix({itemCount}) {
    return `${itemCount} 只基金`
  }
})

const filteredData = computed(() => {
  if (searchCodes.value === null) return rankingData.value
  if (searchCodes.value.length === 0) return []
  return rankingData.value.filter(item => searchCodes.value.includes(item.code))
})

function onSearchKeywordChange(kw) {
  if (searchTimer) clearTimeout(searchTimer)
  const trimmed = kw.trim()
  if (!trimmed) {
    searchCodes.value = null
    return
  }
  searchTimer = setTimeout(() => {
    SearchFundCodes(trimmed).then(result => {
      if (result && result.length > 0) {
        searchCodes.value = result.map(item => item.code)
      } else {
        searchCodes.value = []
      }
    }).catch(() => {
      searchCodes.value = []
    })
  }, 500)
}

const marketTypeOptions = [
  {label: '场外基金', value: 'kf'},
  {label: '场内基金', value: 'fb'},
]

const offExchangeFundTypeOptions = [
  {label: '全部', value: 'all'},
  {label: '股票型', value: 'gp'},
  {label: '混合型', value: 'hh'},
  {label: '债券型', value: 'zq'},
  {label: '指数型', value: 'zs'},
  {label: 'QDII', value: 'qdii'},
  {label: 'FOF', value: 'fof'},
]

const onExchangeFundTypeOptions = [
  {label: '全部', value: 'ct'},
  {label: 'ETF', value: 'etf'},
  {label: 'LOF', value: 'lof'},
]

const fundTypeOptions = computed(() => {
  return marketType.value === 'fb' ? onExchangeFundTypeOptions : offExchangeFundTypeOptions
})

const sortFieldOptions = [
  {label: '今年来涨幅', value: 'jnzf'},
  {label: '日涨幅', value: 'rzdf'},
  {label: '近1周涨幅', value: '1yzf'},
  {label: '近1月涨幅', value: '1mzf'},
  {label: '近3月涨幅', value: '3mzf'},
  {label: '近6月涨幅', value: '6mzf'},
  {label: '近1年涨幅', value: '1nzf'},
  {label: '近2年涨幅', value: '2nzf'},
  {label: '近3年涨幅', value: '3nzf'},
  {label: '成立来涨幅', value: 'clzf'},
  {label: '规模', value: 'gm'},
]

watch(marketType, () => {
  rankingFundType.value = marketType.value === 'fb' ? 'ct' : 'all'
  paginationReactive.page = 1
  searchKeyword.value = ''
  searchCodes.value = null
  fetchFundRanking()
})

watch(rankingFundType, () => {
  paginationReactive.page = 1
  fetchFundRanking()
})

watch(rankingSortField, () => {
  if (skipSortWatch) {
    skipSortWatch = false
    return
  }
  rankingSortOrder.value = 'desc'
  paginationReactive.page = 1
  fetchFundRanking()
})

const keyToSortField = {
  dailyGrowth: 'rzdf',
  weekGrowth: '1yzf',
  monthGrowth: '1mzf',
  threeMonthGrowth: '3mzf',
  sixMonthGrowth: '6mzf',
  yearGrowth: '1nzf',
  threeYearGrowth: '3nzf',
  ytdGrowth: 'jnzf',
  sinceInception: 'clzf',
  scale: 'gm',
  netUnitValue: 'dwjz',
  netAccumulated: 'ljjz',
}

const sortFieldToKey = {}
for (const [k, v] of Object.entries(keyToSortField)) {
  sortFieldToKey[v] = k
}

function getSortOrder(key) {
  if (sortFieldToKey[rankingSortField.value] === key) {
    return rankingSortOrder.value === 'asc' ? 'ascend' : 'descend'
  }
  return false
}

function renderGrowth(val) {
  if (val == null) return '-'
  const color = val > 0 ? '#ef5350' : val < 0 ? '#26a69a' : undefined
  return h(NText, {style: {color}}, () => (val > 0 ? '+' : '') + val.toFixed(2) + '%')
}

const rankingColumns = computed(() => {
  const cols = [
    {title: '代码', key: 'code', width: 90, fixed: 'left'},
    {title: '名称', key: 'name', width: 160, ellipsis: {tooltip: true}},
  ]
  if (marketType.value === 'fb') {
    cols.push({title: '类型', key: 'fundTypeDetail', width: 90, render: (row) => row.fundTypeDetail ? h(NTag, {size: 'tiny', bordered: false, type: 'info'}, () => row.fundTypeDetail) : '-'})
  }
  cols.push(
    {title: '净值日期', key: 'netValueDate', width: 95},
    {title: '单位净值', key: 'netUnitValue', width: 85, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('netUnitValue'), render: (row) => row.netUnitValue?.toFixed(4) ?? '-'},
    {title: '累计净值', key: 'netAccumulated', width: 85, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('netAccumulated'), render: (row) => row.netAccumulated?.toFixed(4) ?? '-'},
    {title: '日涨幅', key: 'dailyGrowth', width: 78, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('dailyGrowth'), render: (row) => renderGrowth(row.dailyGrowth)},
    {title: '近1周', key: 'weekGrowth', width: 72, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('weekGrowth'), render: (row) => renderGrowth(row.weekGrowth)},
    {title: '近1月', key: 'monthGrowth', width: 72, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('monthGrowth'), render: (row) => renderGrowth(row.monthGrowth)},
    {title: '近3月', key: 'threeMonthGrowth', width: 72, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('threeMonthGrowth'), render: (row) => renderGrowth(row.threeMonthGrowth)},
    {title: '近6月', key: 'sixMonthGrowth', width: 72, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('sixMonthGrowth'), render: (row) => renderGrowth(row.sixMonthGrowth)},
    {title: '近1年', key: 'yearGrowth', width: 72, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('yearGrowth'), render: (row) => renderGrowth(row.yearGrowth)},
    {title: '近3年', key: 'threeYearGrowth', width: 72, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('threeYearGrowth'), render: (row) => renderGrowth(row.threeYearGrowth)},
    {title: '今年来', key: 'ytdGrowth', width: 78, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('ytdGrowth'), render: (row) => renderGrowth(row.ytdGrowth)},
    {title: '成立来', key: 'sinceInception', width: 78, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('sinceInception'), render: (row) => renderGrowth(row.sinceInception)},
    {title: '规模(亿)', key: 'scale', width: 78, className: 'nowrap-cell', sorter: true, sortOrder: getSortOrder('scale'), render: (row) => row.scale?.toFixed(2) ?? '-'},
    {title: '成立日期', key: 'establishDate', width: 95},
    {
      title: '操作', key: 'actions', width: 170, fixed: 'right',
      render: (row) => {
        const isFollowed = followList.value.some(f => f.code === row.code)
        return h(NFlex, {size: 4, wrap: false}, () => [
          h(NButton, {
            size: 'tiny',
            type: isFollowed ? 'default' : 'primary',
            disabled: isFollowed,
            onClick: () => rankingFollowFund(row.code)
          }, () => isFollowed ? '已关注' : '关注'),
          h(NButton, {size: 'tiny', type: 'info', onClick: () => showHoldings(row.code, row.name)}, () => '持仓'),
          h(NButton, {size: 'tiny', type: 'warning', onClick: () => search(row.code)}, () => '详情'),
        ])
      }
    },
  )
  return cols
})

function fetchFundRanking() {
  rankingLoading.value = true
  GetFundRanking(marketType.value, rankingFundType.value, rankingSortField.value, rankingSortOrder.value, paginationReactive.page, paginationReactive.pageSize).then(result => {
    if (result) {
      rankingData.value = result.items || []
      paginationReactive.pageCount = result.totalPages || 0
      paginationReactive.itemCount = result.totalCount || 0
    }
  }).catch(() => {
    rankingData.value = []
  }).finally(() => {
    rankingLoading.value = false
  })
}

function resetRankingFilter() {
  rankingFundType.value = marketType.value === 'fb' ? 'ct' : 'all'
  rankingSortField.value = 'jnzf'
  rankingSortOrder.value = 'desc'
  paginationReactive.page = 1
  searchKeyword.value = ''
  searchCodes.value = null
  fetchFundRanking()
}

function handleSorterChange(sorter) {
  if (!sorter || sorter.order === false) {
    rankingSortField.value = 'jnzf'
    rankingSortOrder.value = 'desc'
  } else {
    const field = keyToSortField[sorter.columnKey]
    if (field) {
      skipSortWatch = true
      rankingSortField.value = field
      rankingSortOrder.value = sorter.order === 'ascend' ? 'asc' : 'desc'
    }
  }
  paginationReactive.page = 1
  fetchFundRanking()
}

function handlePageChange(currentPage) {
  if (!rankingLoading.value) {
    paginationReactive.page = currentPage
    fetchFundRanking()
  }
}

function rankingFollowFund(code) {
  FollowFund(code).then(result => {
    if (result) {
      message.success('关注成功')
      loadFollowList()
    }
  })
}

function loadFollowList() {
  GetFollowedFund().then(result => {
    followList.value = result
  })
}

function search(code) {
  setTimeout(() => {
    window.open("https://fund.eastmoney.com/" + code + ".html", "_blank", "noreferrer,width=1000,top=100,left=100,status=no,toolbar=no,location=no,scrollbars=no")
  }, 300)
}

function showHoldings(code, name) {
  holdingsFundCode.value = code
  holdingsFundName.value = name
  holdingsModalShow.value = true
  holdingsLoading.value = true
  holdingsData.value = []
  GetFundTop10Holdings(code).then(result => {
    holdingsData.value = result || []
  }).catch(() => {
    holdingsData.value = []
  }).finally(() => {
    holdingsLoading.value = false
  })
}

const holdingsColumns = [
  {title: '排名', key: 'rank', width: 50},
  {title: '代码', key: 'stockCode', width: 75},
  {title: '名称', key: 'stockName', width: 90, ellipsis: {tooltip: true}},
  {title: '占比(%)', key: 'ratio', width: 70, render: (row) => row.ratio?.toFixed(2) ?? '-'},
  {
    title: '涨跌幅', key: 'changeRate', width: 75,
    render: (row) => {
      const v = row.changeRate
      if (v == null) return '-'
      const color = v > 0 ? '#ef5350' : v < 0 ? '#26a69a' : undefined
      return h(NText, {style: {color}}, () => (v > 0 ? '+' : '') + v.toFixed(2) + '%')
    }
  },
  {
    title: '分时', key: 'sparkline', width: 120,
    render: (row) => {
      if (row.market !== 'A') return '-'
      const lastPrice = row.price || 0
      const openPrice = (row.changeRate != null && row.price != null && row.changeRate !== 0)
        ? row.price / (1 + row.changeRate / 100) : lastPrice
      const prefix = /^(6|5)/.test(row.stockCode) ? 'sh' : 'sz'
      return h(StockSparkLine, {
        stockCode: prefix + row.stockCode,
        stockName: row.stockName,
        lastPrice: lastPrice,
        openPrice: openPrice,
        darkTheme: darkTheme.value,
        idSuffix: '_fh_' + row.stockCode,
      })
    }
  },
  {
    title: '操作', key: 'actions', width: 70,
    render: (row) => h(NButton, {
      size: 'tiny', type: 'info',
      onClick: () => showStockKline(row.stockCode, row.stockName, row.market)
    }, () => 'K线')
  },
]

loadFollowList()
fetchFundRanking()
GetConfig().then(result => {
  if (result.darkTheme) darkTheme.value = true
})

function toEastMoneyCode(stockCode, market) {
  if (market === 'A') {
    if (/^(6|5)/.test(stockCode)) return stockCode + '.SH'
    return stockCode + '.SZ'
  }
  if (market === 'HK') return stockCode + '.HK'
  return stockCode + '.US'
}

function showStockKline(stockCode, stockName, market) {
  klineStockCode.value = toEastMoneyCode(stockCode, market)
  klineStockName.value = stockName
  klineModalShow.value = true
}

// 移动端：把增长率格式化为 {text, type}，供卡片模板使用
function growthText(val) {
  if (val == null) return {text: '-', type: 'default'}
  const sign = val > 0 ? '+' : ''
  return {
    text: sign + val.toFixed(2) + '%',
    type: val > 0 ? 'error' : val < 0 ? 'success' : 'default',
  }
}

// 移动端：是否已关注
function isFollowed(code) {
  return followList.value.some(f => f.code === code)
}

// 移动端：当前关注卡片的目标基金（持仓抽屉复用 desktop 的 holdings* 状态）
function mobileFollow(row) {
  if (isFollowed(row.code)) return
  rankingFollowFund(row.code)
}
</script>

<template>
  <div class="fund-ranking-page">
    <n-flex class="fund-ranking-toolbar" :wrap="false" align="center" :size="12" style="flex-shrink: 0; margin-bottom: 8px;">
      <n-select v-model:value="marketType" :options="marketTypeOptions" style="width: 120px;" size="small"/>
      <n-input v-model:value="searchKeyword" placeholder="基金名称/代码" clearable size="small" style="width: 180px;" @update:value="onSearchKeywordChange"/>
      <n-select v-model:value="rankingFundType" :options="fundTypeOptions" style="width: 120px;" size="small"/>
      <n-select v-model:value="rankingSortField" :options="sortFieldOptions" style="width: 140px;" size="small"/>
      <n-button type="primary" size="small" @click="fetchFundRanking" :loading="rankingLoading">查询</n-button>
      <n-button size="small" @click="resetRankingFilter">重置</n-button>
      <n-text depth="3" v-if="searchCodes !== null" style="font-size: 12px;">搜索到 {{ searchCodes.length }} 只，当前页匹配 {{ filteredData.length }} 只</n-text>
    </n-flex>
    <!-- ===================== 桌面端：宽表格 ===================== -->
    <n-data-table
      v-if="!isMobile"
      class="fund-ranking-table"
      remote
      :columns="rankingColumns"
      :data="filteredData"
      :loading="rankingLoading"
      :bordered="false"
      size="small"
      striped
      :pagination="paginationReactive"
      @update:page="handlePageChange"
      @update:sorter="handleSorterChange"
      :scroll-x="1700"
      flex-height
      style="height: calc(100vh - 210px);margin-top: 10px"
    />

    <!-- ===================== 移动端：基金卡片列表 ===================== -->
    <div v-else class="fr-mobile">
      <n-spin :show="rankingLoading">
        <article v-for="row in filteredData" :key="row.code" class="fr-card">
          <div class="fr-card__head">
            <div class="fr-card__title">
              <span class="fr-card__name">{{ row.name }}</span>
              <n-tag size="tiny" :bordered="false">{{ row.code }}</n-tag>
              <n-tag v-if="row.fundTypeDetail" size="tiny" :bordered="false" type="info">{{ row.fundTypeDetail }}</n-tag>
            </div>
            <span class="fr-card__primary" :class="'fr-card__primary--' + growthText(row.dailyGrowth).type">
              {{ growthText(row.dailyGrowth).text }}
            </span>
          </div>

          <div class="fr-card__core">
            <div class="fr-metric">
              <span class="fr-metric__label">单位净值</span>
              <span class="fr-metric__value">{{ row.netUnitValue != null ? row.netUnitValue.toFixed(4) : '-' }}</span>
            </div>
            <div class="fr-metric">
              <span class="fr-metric__label">累计净值</span>
              <span class="fr-metric__value">{{ row.netAccumulated != null ? row.netAccumulated.toFixed(4) : '-' }}</span>
            </div>
            <div class="fr-metric">
              <span class="fr-metric__label">规模/亿</span>
              <span class="fr-metric__value">{{ row.scale != null ? row.scale.toFixed(2) : '-' }}</span>
            </div>
          </div>

          <!-- 业绩区间：随排序字段高亮当前列 -->
          <div class="fr-card__growth">
            <div v-for="g in [
              {label: '近1周', v: row.weekGrowth},
              {label: '近1月', v: row.monthGrowth},
              {label: '近3月', v: row.threeMonthGrowth},
              {label: '近1年', v: row.yearGrowth},
              {label: '今年来', v: row.ytdGrowth},
              {label: '成立来', v: row.sinceInception},
            ]" :key="g.label" class="fr-growth">
              <span class="fr-growth__label">{{ g.label }}</span>
              <span class="fr-growth__value" :class="'text-' + growthText(g.v).type">{{ growthText(g.v).text }}</span>
            </div>
          </div>

          <div class="fr-card__actions">
            <n-button size="tiny" :type="isFollowed(row.code) ? 'default' : 'primary'" :disabled="isFollowed(row.code)" @click="mobileFollow(row)">
              {{ isFollowed(row.code) ? '已关注' : '关注' }}
            </n-button>
            <n-button size="tiny" type="info" @click="showHoldings(row.code, row.name)">持仓</n-button>
            <n-button size="tiny" type="warning" @click="search(row.code)">详情</n-button>
          </div>
        </article>
        <n-empty v-if="!rankingLoading && filteredData.length === 0" description="暂无基金数据" style="padding: 40px 0" />
      </n-spin>

      <n-flex v-if="paginationReactive.pageCount > 1" justify="center" style="margin-top: 12px">
        <n-pagination
            v-model:page="paginationReactive.page"
            :page-count="paginationReactive.pageCount"
            :page-size="paginationReactive.pageSize"
            size="small"
            @update:page="handlePageChange"
        />
      </n-flex>
    </div>

  <!-- ===================== 桌面端：持仓弹窗 ===================== -->
  <n-modal
    v-if="!isMobile"
    class="fund-ranking-holdings-modal"
    v-model:show="holdingsModalShow"
    :title="holdingsFundName + ' - ' + holdingsFundCode + ' 十大持仓'"
    preset="card"
    style="max-width: 1400px;"
    :mask-closable="true"
  >
    <n-text v-if="holdingsData.length > 0 && holdingsData[0]?.quarter" depth="3" style="font-size: 12px; margin-bottom: 4px; display: inline-block;">
      {{ holdingsData[0].quarter }}
    </n-text>
    <n-data-table
      :columns="holdingsColumns"
      :data="holdingsData"
      :loading="holdingsLoading"
      :pagination="false"
      size="small"
      :bordered="false"
      :max-height="500"
      striped
    />
  </n-modal>

  <!-- ===================== 桌面端：K 线弹窗 ===================== -->
  <n-modal
    v-if="!isMobile"
    class="fund-ranking-kline-modal"
    v-model:show="klineModalShow"
    :title="klineStockName + ' - ' + klineStockCode + ' K线图'"
    preset="card"
    style="max-width: 1400px;"
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

  <!-- ===================== 移动端：持仓底部抽屉 ===================== -->
  <BottomSheet v-else-if="isMobile && holdingsModalShow" :show="holdingsModalShow" :title="holdingsFundName + ' · 十大持仓'" height="80vh" @update:show="(v) => holdingsModalShow = v">
    <div class="fr-holdings">
      <n-text v-if="holdingsData.length > 0 && holdingsData[0]?.quarter" depth="3" class="fr-holdings__quarter">{{ holdingsData[0].quarter }}</n-text>
      <n-spin :show="holdingsLoading">
        <button
            v-for="h in holdingsData" :key="h.stockCode"
            type="button" class="fr-holding" @click="showStockKline(h.stockCode, h.stockName, h.market)"
        >
          <div class="fr-holding__rank">{{ h.rank }}</div>
          <div class="fr-holding__main">
            <div class="fr-holding__name">{{ h.stockName }}</div>
            <div class="fr-holding__code">{{ h.stockCode }}</div>
          </div>
          <div class="fr-holding__ratio">{{ h.ratio != null ? h.ratio.toFixed(2) : '-' }}%</div>
          <div class="fr-holding__change" :class="'text-' + growthText(h.changeRate).type">
            {{ growthText(h.changeRate).text }}
          </div>
          <span class="fr-holding__arrow">›</span>
        </button>
        <n-empty v-if="!holdingsLoading && holdingsData.length === 0" description="暂无持仓数据" size="small" />
      </n-spin>
    </div>
  </BottomSheet>

  <!-- ===================== 移动端：K 线底部抽屉 ===================== -->
  <BottomSheet v-if="isMobile && klineModalShow" :show="klineModalShow" :title="klineStockName + ' · K线'" height="78vh" @update:show="(v) => klineModalShow = v">
    <div class="fr-kline-wrap">
      <StockLightweightKlineChart
          v-if="klineModalShow && klineStockCode"
          :key="klineStockCode"
          :code="klineStockCode"
          :stock-name="klineStockName"
          :dark-theme="darkTheme"
          :chart-height="420"
      />
    </div>
  </BottomSheet>
  </div>
</template>

<style scoped>
.fund-ranking-page {
  min-width: 0;
}

:deep(.nowrap-cell) {
  white-space: nowrap;
}

@media (max-width: 768px) {
  .fund-ranking-page {
    padding-bottom: calc(var(--mobile-bottom-nav-height) + var(--safe-bottom) + 8px);
  }

  .fund-ranking-toolbar {
    display: grid !important;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 8px !important;
    align-items: stretch !important;
  }

  .fund-ranking-toolbar :deep(.n-select),
  .fund-ranking-toolbar :deep(.n-input),
  .fund-ranking-toolbar :deep(.n-button) {
    width: 100% !important;
  }

  .fund-ranking-toolbar :deep(.n-text) {
    grid-column: 1 / -1;
    line-height: 1.35;
  }

  .fund-ranking-table {
    height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 190px) !important;
    margin-top: 8px !important;
  }

  .fund-ranking-table :deep(.n-data-table-th),
  .fund-ranking-table :deep(.n-data-table-td) {
    white-space: nowrap;
  }

  :deep(.fund-ranking-holdings-modal.n-modal),
  :deep(.fund-ranking-kline-modal.n-modal) {
    margin: 0 !important;
    max-width: 100vw !important;
    width: calc(100vw - 12px) !important;
  }

  :deep(.fund-ranking-holdings-modal .n-card),
  :deep(.fund-ranking-kline-modal .n-card) {
    display: flex;
    flex-direction: column;
    max-height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 12px);
  }

  :deep(.fund-ranking-holdings-modal .n-card__content),
  :deep(.fund-ranking-kline-modal .n-card__content) {
    display: flex;
    flex: 1 1 auto;
    flex-direction: column;
    min-height: 0;
    overflow: auto;
    padding: 10px 12px;
  }

  @media (max-width: 420px) {
    .fund-ranking-toolbar {
      grid-template-columns: minmax(0, 1fr);
    }
  }
}

/* ============ 移动端基金卡片（仅在 isMobile 渲染） ============ */
.fr-mobile {
  display: flex;
  flex-direction: column;
}

.fr-card {
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 10px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 8px;
  padding: 11px 12px;
}

.fr-card__head {
  align-items: flex-start;
  display: flex;
  gap: 8px;
  justify-content: space-between;
}

.fr-card__title {
  align-items: center;
  display: flex;
  flex: 1 1 auto;
  flex-wrap: wrap;
  gap: 6px;
  min-width: 0;
}

.fr-card__name {
  font-size: 16px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fr-card__primary {
  border-radius: 4px;
  color: #fff;
  flex: 0 0 auto;
  font-size: 15px;
  font-weight: 700;
  padding: 2px 8px;
}

.fr-card__primary--error {
  background: #d03050;
}

.fr-card__primary--success {
  background: #18a058;
}

.fr-card__primary--default {
  background: #909399;
}

.fr-card__core {
  display: grid;
  gap: 8px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.fr-metric {
  align-items: flex-start;
  background: var(--n-color-target, #f5f7fa);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 7px 9px;
}

.fr-metric__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.fr-metric__value {
  font-size: 15px;
  font-weight: 700;
}

.fr-card__growth {
  border-top: 1px dashed var(--n-divider-color, #eef0f4);
  display: grid;
  gap: 6px 8px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  padding-top: 8px;
}

.fr-growth {
  align-items: center;
  display: flex;
  justify-content: space-between;
}

.fr-growth__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.fr-growth__value {
  font-size: 12px;
  font-weight: 700;
}

.fr-card__actions {
  display: flex;
  gap: 8px;
}

.fr-card__actions :deep(.n-button) {
  flex: 1 1 0;
}

/* 持仓抽屉行 */
.fr-holdings {
  padding: 4px 8px calc(var(--safe-bottom) + 8px);
}

.fr-holdings__quarter {
  display: block;
  font-size: 12px;
  margin-bottom: 8px;
}

.fr-holding {
  align-items: center;
  appearance: none;
  background: var(--n-color, #fff);
  border: 0;
  border-bottom: 1px solid var(--n-border-color, #eef0f4);
  color: inherit;
  display: flex;
  font: inherit;
  gap: 10px;
  padding: 9px 8px;
  text-align: left;
  width: 100%;
}

.fr-holding:active {
  background: var(--n-color-hover, #f8fafc);
}

.fr-holding__rank {
  color: var(--n-text-color-3, #98a2b3);
  flex: 0 0 22px;
  font-size: 13px;
  font-weight: 700;
  text-align: center;
}

.fr-holding__main {
  display: flex;
  flex: 1 1 auto;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.fr-holding__name {
  font-size: 14px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fr-holding__code {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.fr-holding__ratio {
  font-size: 13px;
  font-weight: 600;
  text-align: right;
  width: 56px;
}

.fr-holding__change {
  font-size: 13px;
  font-weight: 700;
  text-align: right;
  width: 64px;
}

.fr-holding__arrow {
  color: var(--n-text-color-3, #c0c4cc);
  flex: 0 0 auto;
  font-size: 18px;
}

.fr-kline-wrap {
  padding: 4px 8px 8px;
}

/* 涨跌色 */
.text-error {
  color: #d03050;
}

.text-success {
  color: #18a058;
}

.text-default {
  color: var(--n-text-color, #333);
}
</style>
