<script setup>
import { ref, onBeforeMount, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import MPullRefresh from '../components/base/MPullRefresh.vue'
import PageHeader from '../components/widgets/PageHeader.vue'
import MarketStatusBar from '../components/widgets/MarketStatusBar.vue'
import StockSummaryCard from '../components/cards/StockSummaryCard.vue'
import NewsCard from '../components/cards/NewsCard.vue'
import HotTopicCard from '../components/cards/HotTopicCard.vue'
import AlertCard from '../components/cards/AlertCard.vue'
import IndustryCard from '../components/cards/IndustryCard.vue'
import StockDetailSheet from '../components/sheets/StockDetailSheet.vue'
import MCard from '../components/base/MCard.vue'
import MIcon from '../components/base/MIcon.vue'

// 导入API
import {
  GetFollowList,
  Greet,
  GetTelegraphList,
  ReFleshTelegraphList,
  HotTopic,
  GetStockChanges,
  GetIndustryRank
} from '../../api/app'
import { registerFeed, stopFeed } from '../../api/scheduler'
import { anyOpen } from '../../api/marketClock'
import { EventsOn, EventsOff } from '../../api/runtime'
import { formatHeat } from '../composables/useFormat'

const router = useRouter()

// 真实数据
const stockData = ref([])
const newsData = ref([])
const topicsData = ref([])
const alertsData = ref([])
const industriesData = ref([])

const NEWS_SOURCES = ['财联社电报', '新浪财经', '外媒']
const HOME_NEWS_LIMIT = 3

// 个股详情抽屉
const detailVisible = ref(false)
const selectedStock = ref(null)

function pick(item, keys, fallback = undefined) {
  for (const key of keys) {
    if (item && item[key] !== undefined && item[key] !== null && item[key] !== '') {
      return item[key]
    }
  }
  return fallback
}

function stockCodeKey(code) {
  return String(code || '').replace(/[^a-z0-9]/gi, '').toLowerCase()
}

function quoteNumber(item, keys, fallback = 0, allowZero = false) {
  const raw = pick(item, keys, undefined)
  const n = Number(raw)
  if (!Number.isFinite(n)) return fallback
  if (!allowZero && n === 0) return fallback
  return n
}

function sameStockCodes(a, b) {
  if (!Array.isArray(a) || !Array.isArray(b) || a.length !== b.length) return false
  return a.every((stock, index) => stockCodeKey(stock.code) === stockCodeKey(b[index]?.code))
}

function commitStockData(next) {
  if (sameStockCodes(next, stockData.value)) {
    stockData.value.forEach((stock, index) => {
      Object.assign(stock, next[index])
    })
    return
  }
  stockData.value = next
}

function normalizeFollowStock(stock, previous = {}) {
  const code = pick(stock, ['StockCode', 'stockCode', 'code'], '')
  return {
    code,
    name: pick(stock, ['Name', 'StockName', 'stockName', 'name'], previous.name || '未命名股票'),
    price: quoteNumber(stock, ['Price', 'price', 'currentPrice', '当前价格'], previous.price || 0),
    changePercent: quoteNumber(stock, ['ChangePercent', 'changePercent', 'change_percent'], previous.changePercent || 0),
    changeAmount: quoteNumber(stock, ['PriceChange', 'ChangePrice', 'changeAmount', 'priceChange', '涨跌额'], previous.changeAmount || 0),
    high: quoteNumber(stock, ['High', '今日最高价'], previous.high || 0),
    low: quoteNumber(stock, ['Low', '今日最低价'], previous.low || 0),
    open: quoteNumber(stock, ['Open', '今日开盘价'], previous.open || 0),
    preClose: quoteNumber(stock, ['PreClose', '昨日收盘价'], previous.preClose || 0),
    volume: quoteNumber(stock, ['成交的股票数'], previous.volume || 0),
    turnover: quoteNumber(stock, ['Turnover', 'Amount', '成交金额'], previous.turnover || 0),
    quoteDate: pick(stock, ['日期', 'Date', 'date'], previous.quoteDate || ''),
    time: pick(stock, ['时间'], previous.time || ''),
    // 持仓成本/数量（GetFollowList 快照带回，实时盈亏由 fetchStockRealtime 覆盖）
    costPrice: quoteNumber(stock, ['CostPrice', 'costPrice'], previous.costPrice || 0),
    costVolume: quoteNumber(stock, ['Volume', 'costVolume'], previous.costVolume || 0),
    profit: quoteNumber(stock, ['profit'], previous.profit || 0),
    profitAmount: quoteNumber(stock, ['profitAmount'], previous.profitAmount || 0),
    profitToday: quoteNumber(stock, ['profitAmountToday'], previous.profitToday || 0)
  }
}

function newsTimestamp(news) {
  const value = pick(news, ['time', 'dataTime', 'createdAt'], Date.now())
  if (value instanceof Date) return value.getTime()
  if (typeof value === 'number') return value < 1e12 ? value * 1000 : value
  const timestamp = new Date(value).getTime()
  return Number.isFinite(timestamp) ? timestamp : Date.now()
}

function newsKey(news) {
  return pick(news, ['id', 'ID', 'url'], `${news.source || ''}-${newsTimestamp(news)}-${news.title || ''}`)
}

function normalizeNewsItem(item) {
  const title = pick(item, ['content', 'title'], '')
  const source = pick(item, ['source', 'media'], '市场快讯')
  const timestamp = newsTimestamp(item)
  return {
    id: pick(item, ['ID', 'id'], undefined),
    title,
    time: new Date(timestamp),
    source,
    url: pick(item, ['url', 'URL'], '')
  }
}

function selectHomeNews(newsLists) {
  const seen = new Set()
  const sorted = newsLists
    .flat()
    .map(normalizeNewsItem)
    .filter(item => item.title)
    .sort((a, b) => newsTimestamp(b) - newsTimestamp(a))

  const unique = []
  for (const item of sorted) {
    const key = newsKey(item)
    if (seen.has(key)) continue
    seen.add(key)
    unique.push(item)
  }

  const selected = []
  const selectedKeys = new Set()
  for (const source of NEWS_SOURCES) {
    const item = unique.find(news => news.source === source)
    if (!item) continue
    const key = newsKey(item)
    selected.push(item)
    selectedKeys.add(key)
  }

  for (const item of unique) {
    if (selected.length >= HOME_NEWS_LIMIT) break
    const key = newsKey(item)
    if (selectedKeys.has(key)) continue
    selected.push(item)
    selectedKeys.add(key)
  }

  return selected
    .sort((a, b) => newsTimestamp(b) - newsTimestamp(a))
    .slice(0, HOME_NEWS_LIMIT)
}

async function getNewsBySource() {
  const results = await Promise.all(NEWS_SOURCES.map(async (source) => {
    try {
      const result = await GetTelegraphList(source)
      return Array.isArray(result) ? result : []
    } catch (error) {
      console.error(`加载${source}失败:`, error)
      return []
    }
  }))
  return selectHomeNews(results)
}

function prependNewsItems(items) {
  if (!Array.isArray(items) || !items.length) return
  newsData.value = selectHomeNews([items, newsData.value])
}

// 加载自选股票（取前3条），并拉取实时行情填充价格/涨跌/量额
async function loadStockData() {
  try {
    // 0 表示全部分组（与桌面端 stock.vue 的默认值一致）
    const result = await GetFollowList(0)
    if (result && Array.isArray(result)) {
      const previousByCode = new Map(stockData.value.map(stock => [stockCodeKey(stock.code), stock]))
      const list = result
        .slice(0, 3)
        .map(stock => {
          const code = pick(stock, ['StockCode', 'stockCode', 'code'], '')
          return normalizeFollowStock(stock, previousByCode.get(stockCodeKey(code)))
        })
        .filter(stock => stock.code)

      if (!stockData.value.length) {
        commitStockData(list)
        await fetchStockRealtime(stockData.value)
        return
      }

      await fetchStockRealtime(list)
      commitStockData(list)
    }
  } catch (error) {
    console.error('加载自选股票失败:', error)
  }
}

// 拉取实时行情增量更新（对齐桌面端 stock.vue 的 Greet + updateData）
async function fetchStockRealtime(list) {
  if (!Array.isArray(list) || !list.length) return
  await Promise.all(list.map(async (stock) => {
    if (!stock.code) return
    try {
      const rt = await Greet(stock.code)
      if (!rt) return
      const price = Number(rt['当前价格'] || rt['卖一报价']) || 0
      if (price > 0) {
        stock.price = price
        stock.changePercent = Number(rt.changePercent) || 0
        stock.changeAmount = Number(rt.changePrice ?? rt['涨跌额']) || 0
        stock.high = Number(rt['今日最高价']) || stock.high
        stock.low = Number(rt['今日最低价']) || stock.low
        stock.open = Number(rt['今日开盘价']) || stock.open
        stock.preClose = Number(rt['昨日收盘价']) || stock.preClose
        stock.volume = Number(rt['成交的股票数']) || stock.volume
        stock.turnover = Number(rt['成交金额']) || stock.turnover
        stock.quoteDate = rt['日期'] || stock.quoteDate
        stock.time = rt['时间'] || stock.time
        // 持仓盈亏（对齐桌面端 profit/profitAmount/profitAmountToday）
        stock.costPrice = Number(rt.costPrice) || 0
        stock.costVolume = Number(rt.costVolume) || 0
        stock.profit = Number(rt.profit) || 0
        stock.profitAmount = Number(rt.profitAmount) || 0
        stock.profitToday = Number(rt.profitAmountToday) || 0
      }
    } catch (e) {
      // 单只失败不打断
    }
  }))
}

// 加载市场快讯（读后端缓存，用于首次加载）
async function loadNewsData() {
  try {
    newsData.value = await getNewsBySource()
  } catch (error) {
    console.error('加载市场快讯失败:', error)
  }
}

// 刷新市场快讯：只触发一次后端三源抓取，再按来源读取缓存合并。
async function refreshNewsData() {
  try {
    await ReFleshTelegraphList('')
    newsData.value = await getNewsBySource()
  } catch (error) {
    console.error('刷新市场快讯失败:', error)
  }
}

// 加载热点话题
async function loadTopicsData() {
  try {
    const result = await HotTopic(10)
    if (result && Array.isArray(result)) {
      topicsData.value = result.slice(0, 3).map(item => ({
        title: pick(item, ['nickname', 'title', 'name'], '热点话题'),
        heat: formatHeat(pick(item, ['clickNumber', 'heat', 'hot'], 0)),
        // HotTopic 接口无涨跌幅字段，不传 changePercent，卡片自动隐藏该列
        stocks: Array.isArray(item.stock_list) ? item.stock_list.length : 0
      }))
    }
  } catch (error) {
    console.error('加载热点话题失败:', error)
  }
}

// 加载异动监控
async function loadAlertsData() {
  try {
    const result = await GetStockChanges([0], 1, 5)
    if (result && result.data && result.data.length > 0) {
      alertsData.value = result.data.slice(0, 2).map(item => ({
        stockName: pick(item, ['stockName', 'StockName', 'name'], '未知股票'),
        stockCode: pick(item, ['stockCode', 'StockCode', 'code'], ''),
        type: mapChangeType(pick(item, ['changeType', 'ChangeType', 'type'], 0)),
        changePercent: pick(item, ['changePercent', 'ChangePercent'], 0),
        time: new Date(pick(item, ['time', 'dataTime', 'createdAt'], Date.now()))
      }))
    } else {
      // 如果没有异动数据，清空数组
      alertsData.value = []
    }
  } catch (error) {
    console.error('加载异动监控失败:', error)
  }
}

// 加载行业排名
async function loadIndustriesData() {
  try {
    const result = await GetIndustryRank(0, 5)
    if (result && Array.isArray(result)) {
      industriesData.value = result.slice(0, 5).map(item => ({
        name: pick(item, ['bd_name', 'name'], '未知行业'),
        changePercent: parseFloat(pick(item, ['bd_zdf', 'changePercent'], 0)) || 0,
        leadingStock: formatLeadingStock(item)
      }))
    }
  } catch (error) {
    console.error('加载行业排名失败:', error)
  }
}

// 格式化领涨股展示
function formatLeadingStock(item) {
  const name = pick(item, ['nzg_name', 'leadingStock'], '')
  const change = pick(item, ['nzg_zdf', 'leadingChangePercent'], '')
  if (!name) return '--'
  const n = Number(change)
  if (!Number.isFinite(n)) return name
  return `${name} ${n > 0 ? '+' : ''}${n}%`
}

// 映射异动类型
function mapChangeType(type) {
  const typeMap = {
    1: 'limit_up',
    2: 'limit_down',
    3: 'rapid_rise',
    4: 'rapid_fall',
    5: 'high_volume',
    6: 'breakthrough'
  }
  return typeMap[type] || 'normal'
}

// 加载所有数据
async function loadAllData() {
  await Promise.all([
    loadStockData(),
    loadNewsData(),
    loadTopicsData(),
    loadAlertsData(),
    loadIndustriesData()
  ])
}

// 下拉刷新（快讯用 ReFleshTelegraphList 真正抓取最新，其余读缓存）
async function handleRefresh() {
  await Promise.all([
    loadStockData(),
    refreshNewsData(),
    loadTopicsData(),
    loadAlertsData(),
    loadIndustriesData()
  ])
}

// 生命周期：初始化加载和轮询
onBeforeMount(async () => {
  await loadAllData()

  // 注册轮询任务（每10秒刷新一次）
  registerFeed('mobile-home-stocks', {
    fetch: loadStockData,
    intervalMs: 10000,
    activeWhen: () => anyOpen.value
  })

  registerFeed('mobile-home-news', {
    fetch: refreshNewsData,
    intervalMs: 30000 // 快讯30秒抓取一次最新
  })

  registerFeed('mobile-home-topics', {
    fetch: loadTopicsData,
    intervalMs: 60000 // 热点1分钟刷新一次
  })

  // 订阅实时价格推送（对齐桌面端 stock.vue 的 stock_price 事件）
  // 后端 MonitorStockPrices 定时通过 SSE 推送，收到后增量更新自选列表的价格
  EventsOn('stock_price', (data) => {
    if (!data) return
    const code = data['股票代码']
    const price = data['当前价格'] || data['卖一报价']
    if (!code) return
    const target = stockData.value.find(s => stockCodeKey(s.code) === stockCodeKey(code))
    if (target && price > 0) {
      target.price = Number(price) || target.price
      target.changePercent = data.changePercent || 0
      target.changeAmount = Number(data.changePrice ?? data['涨跌额']) || 0
      target.open = Number(data['今日开盘价']) || target.open
      target.preClose = Number(data['昨日收盘价']) || target.preClose
      target.high = Number(data['今日最高价']) || target.high
      target.low = Number(data['今日最低价']) || target.low
      target.volume = Number(data['成交的股票数']) || target.volume
      target.turnover = Number(data['成交金额']) || target.turnover
      target.quoteDate = data['日期'] || target.quoteDate
      target.time = data['时间'] || target.time
    }
  })

  // 订阅三类快讯推送（对齐桌面端 market.vue）
  EventsOn('newTelegraph', prependNewsItems)
  EventsOn('newSinaNews', prependNewsItems)
  EventsOn('tradingViewNews', prependNewsItems)
})

onBeforeUnmount(() => {
  // 停止轮询
  stopFeed('mobile-home-stocks')
  stopFeed('mobile-home-news')
  stopFeed('mobile-home-topics')
  // 取消实时事件订阅
  EventsOff('stock_price')
  EventsOff('newTelegraph')
  EventsOff('newSinaNews')
  EventsOff('tradingViewNews')
})

// 导航跳转
function navigateTo(path) {
  router.push(path)
}

// 卡片事件处理
function handleStockClick(stock) {
  // 打开个股详情抽屉（StockDetailSheet 内部按 code 拉分时/K线/盘口/资金）
  selectedStock.value = stock
  detailVisible.value = true
}

function handleNewsClick(news) {
  // 电报无独立详情页：有链接则外开，否则进入市场页看完整快讯
  const url = pick(news, ['url', 'URL'], '')
  if (url) {
    window.open(url, '_blank')
  } else {
    navigateTo('/mobile/market')
  }
}

function handleTopicClick(topic) {
  navigateTo('/mobile/market')
}

function handleAlertClick(alert) {
  navigateTo('/mobile/research')
}

function handleIndustryClick(industry) {
  navigateTo('/mobile/market')
}
</script>

<template>
  <div class="home-container">
    <!-- 顶部导航栏（只保留标题） -->
    <PageHeader title="go-stock" />

    <!-- 下拉刷新内容区 -->
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="home-page">
      <!-- 市场状态条 -->
      <MarketStatusBar />

      <!-- 自选概览卡片 -->
      <StockSummaryCard
        :stocks="stockData"
        :max-show="3"
        @view-all="navigateTo('/mobile/stock')"
        @stock-click="handleStockClick"
      />

      <!-- 市场快讯卡片 -->
      <MCard>
        <div class="news-header">
          <h3 class="news-title">市场快讯</h3>
          <button class="news-action" @click="navigateTo('/mobile/market')">
            更多 <MIcon name="arrow-right" :size="14" />
          </button>
        </div>
        <div class="news-list">
          <NewsCard
            v-for="(news, index) in newsData"
            :key="index"
            :news="news"
            @click="handleNewsClick"
          />
        </div>
      </MCard>

      <!-- 实时热点卡片 -->
      <HotTopicCard
        :topics="topicsData"
        @topic-click="handleTopicClick"
        @view-more="navigateTo('/mobile/market')"
      />

      <!-- 异动监控卡片 -->
      <AlertCard
        :alerts="alertsData"
        @alert-click="handleAlertClick"
        @view-all="navigateTo('/mobile/research')"
      />

      <!-- 行业热度卡片 -->
      <IndustryCard
        :industries="industriesData"
        :max-show="5"
        @industry-click="handleIndustryClick"
        @view-all="navigateTo('/mobile/market')"
      />

      <!-- 底部提示 -->
      <div class="home-footer">
        <span class="footer-line" />
        <p class="footer-text">行情数据仅供参考 · 下拉刷新</p>
      </div>
    </div>
  </MPullRefresh>

    <!-- 个股详情抽屉 -->
    <StockDetailSheet v-model:show="detailVisible" :stock="selectedStock" />
</div>
</template>

<style scoped>
.home-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background:
    radial-gradient(120% 60% at 50% 0%, rgba(208, 48, 80, 0.06) 0%, rgba(208, 48, 80, 0) 60%),
    var(--m-bg-primary);
}

.home-page {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
  padding-bottom: var(--m-space-2xl);
}

.news-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--m-space-lg);
}

.news-title {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  letter-spacing: 0.2px;
}

.news-title::before {
  content: '';
  width: 3px;
  height: 15px;
  border-radius: var(--m-radius-full);
  background: var(--m-color-rise);
}

.news-action {
  padding: var(--m-space-xs) var(--m-space-sm);
  background: transparent;
  border: none;
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
  cursor: pointer;
}

.news-action:active {
  opacity: 0.6;
}

.news-list {
  display: flex;
  flex-direction: column;
}

.home-footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-xl) 0 var(--m-space-lg);
  color: var(--m-text-tertiary);
}

.footer-line {
  width: 32px;
  height: 2px;
  border-radius: var(--m-radius-full);
  background: var(--m-divider-color);
}

.footer-text {
  font-size: var(--m-font-xs);
  letter-spacing: 0.3px;
}
</style>
