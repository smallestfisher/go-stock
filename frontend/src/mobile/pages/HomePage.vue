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
import MCard from '../components/base/MCard.vue'

// 导入API
import {
  GetFollowList,
  Greet,
  GetTelegraphList,
  ReFleshTelegraphList,
  HotTopic,
  GetStockChanges,
  GetIndustryRank,
  GetTodayMarketStatistic
} from '../../api/app'
import { registerFeed, stopFeed } from '../../api/scheduler'
import { EventsOn, EventsOff } from '../../api/runtime'
import { formatHeat } from '../composables/useFormat'

const router = useRouter()

// 真实数据
const stockData = ref([])
const newsData = ref([])
const topicsData = ref([])
const alertsData = ref([])
const industriesData = ref([])
const marketStatistic = ref(null)

const refreshCount = ref(0)
const loading = ref(true)

function pick(item, keys, fallback = undefined) {
  for (const key of keys) {
    if (item && item[key] !== undefined && item[key] !== null && item[key] !== '') {
      return item[key]
    }
  }
  return fallback
}

// 加载自选股票（取前3条），并拉取实时行情填充价格/涨跌/量额
async function loadStockData() {
  try {
    // 0 表示全部分组（与桌面端 stock.vue 的默认值一致）
    const result = await GetFollowList(0)
    if (result && Array.isArray(result)) {
      const list = result.slice(0, 3).map(stock => ({
        code: pick(stock, ['StockCode', 'stockCode', 'code'], ''),
        name: pick(stock, ['Name', 'StockName', 'stockName', 'name'], '未命名股票'),
        price: Number(pick(stock, ['Price', 'price', 'currentPrice', '当前价格'], 0)) || 0,
        changePercent: Number(pick(stock, ['ChangePercent', 'changePercent', 'change_percent'], 0)) || 0,
        changeAmount: Number(pick(stock, ['PriceChange', 'ChangePrice', 'changeAmount', 'priceChange', '涨跌额'], 0)) || 0,
        high: Number(pick(stock, ['High', '今日最高价'], 0)) || 0,
        low: Number(pick(stock, ['Low', '今日最低价'], 0)) || 0,
        volume: Number(pick(stock, ['Volume', 'volume', '成交的股票数'], 0)) || 0,
        turnover: Number(pick(stock, ['Turnover', 'Amount', '成交金额'], 0)) || 0,
        time: pick(stock, ['Time', '时间'], '')
      }))
      stockData.value = list
      // 拉取实时行情覆盖（对齐桌面端 Greet + 自选页 fetchRealtime）
      fetchStockRealtime(list)
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
        stock.volume = Number(rt['成交的股票数']) || stock.volume
        stock.turnover = Number(rt['成交金额']) || stock.turnover
        stock.time = rt['时间'] || stock.time
      }
    } catch (e) {
      // 单只失败不打断
    }
  }))
}

// 加载市场快讯（读后端缓存，用于首次加载）
// 注：后端 GetTelegraphList 不论传何源名都返回全部快讯的混合（按时间倒序），
// 每条数据的 source 字段才是真实来源。取最新 3 条。
async function loadNewsData() {
  try {
    const result = await GetTelegraphList('财联社电报')
    if (result && Array.isArray(result)) {
      newsData.value = result.slice(0, 3).map(item => ({
        title: pick(item, ['content', 'title'], ''),
        time: new Date(pick(item, ['dataTime', 'time', 'createdAt'], Date.now())),
        source: pick(item, ['source', 'media'], '市场快讯')
      }))
    }
  } catch (error) {
    console.error('加载市场快讯失败:', error)
  }
}

// 刷新市场快讯（触发后端抓取最新电报，用于下拉刷新）
// 与 GetTelegraphList 的区别：后者只读缓存，这个会真正拉取最新数据
async function refreshNewsData() {
  try {
    const result = await ReFleshTelegraphList('财联社电报')
    if (result && Array.isArray(result)) {
      newsData.value = result.slice(0, 3).map(item => ({
        title: pick(item, ['content', 'title'], ''),
        time: new Date(pick(item, ['dataTime', 'time', 'createdAt'], Date.now())),
        source: pick(item, ['source', 'media'], '市场快讯')
      }))
    }
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

// 加载市场统计
async function loadMarketStatistic() {
  try {
    const result = await GetTodayMarketStatistic()
    marketStatistic.value = result
  } catch (error) {
    console.error('加载市场统计失败:', error)
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
  loading.value = true
  try {
    await Promise.all([
      loadStockData(),
      loadNewsData(),
      loadTopicsData(),
      loadAlertsData(),
      loadIndustriesData(),
      loadMarketStatistic()
    ])
  } finally {
    loading.value = false
  }
}

// 下拉刷新（快讯用 ReFleshTelegraphList 真正抓取最新，其余读缓存）
async function handleRefresh() {
  refreshCount.value++
  await Promise.all([
    loadStockData(),
    refreshNewsData(),
    loadTopicsData(),
    loadAlertsData(),
    loadIndustriesData(),
    loadMarketStatistic()
  ])
}

// 生命周期：初始化加载和轮询
onBeforeMount(async () => {
  await loadAllData()

  // 注册轮询任务（每10秒刷新一次）
  registerFeed('mobile-home-stocks', {
    fetch: loadStockData,
    intervalMs: 10000
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
    const target = stockData.value.find(s => s.code === code)
    if (target && price > 0) {
      target.price = Number(price) || target.price
      target.changePercent = data.changePercent || 0
      target.changeAmount = Number(data.changePrice ?? data['涨跌额']) || 0
      target.high = Number(data['今日最高价']) || target.high
      target.low = Number(data['今日最低价']) || target.low
      target.time = data['时间'] || target.time
    }
  })

  // 订阅新电报推送（对齐桌面端 market.vue 的 newTelegraph 事件）
  // 后端抓到新财联社电报时通过 SSE 推送，收到后插到列表头部
  EventsOn('newTelegraph', (data) => {
    if (!Array.isArray(data) || !data.length) return
    const newItems = data.map(item => ({
      title: pick(item, ['content', 'title'], ''),
      time: new Date(pick(item, ['dataTime', 'time', 'createdAt'], Date.now())),
      source: pick(item, ['source', 'media'], '市场快讯')
    }))
    // 新电报插到头部，保留前 3 条
    newsData.value = [...newItems, ...newsData.value].slice(0, 3)
  })
})

onBeforeUnmount(() => {
  // 停止轮询
  stopFeed('mobile-home-stocks')
  stopFeed('mobile-home-news')
  stopFeed('mobile-home-topics')
  // 取消实时事件订阅
  EventsOff('stock_price')
  EventsOff('newTelegraph')
})

// 导航跳转
function navigateTo(path) {
  router.push(path)
}

// 卡片事件处理
function handleStockClick(stock) {
  console.log('点击股票:', stock)
  // TODO: 打开股票详情抽屉
}

function handleNewsClick(news) {
  console.log('点击新闻:', news)
  // TODO: 打开新闻详情
}

function handleTopicClick(topic) {
  console.log('点击热点:', topic)
  navigateTo('/mobile/market')
}

function handleAlertClick(alert) {
  console.log('点击异动:', alert)
  navigateTo('/mobile/research')
}

function handleIndustryClick(industry) {
  console.log('点击行业:', industry)
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
          <h3 class="news-title">📰 市场快讯</h3>
          <button class="news-action" @click="navigateTo('/mobile/market')">
            更多 →
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
        <p class="footer-text">下拉刷新数据 · 已刷新 {{ refreshCount }} 次</p>
        <p class="footer-tip">💡 点击卡片查看详情</p>
      </div>
    </div>
  </MPullRefresh>
</div>
</template>

<style scoped>
.home-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

/* 顶部操作按钮（搜索/通知） */
.action-btn {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 18px;
  color: var(--m-text-primary);
  cursor: pointer;
}

.action-btn:active {
  opacity: 0.6;
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
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
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
  text-align: center;
  padding: var(--m-space-xl) 0;
  color: var(--m-text-tertiary);
}

.footer-text {
  font-size: var(--m-font-sm);
  margin-bottom: var(--m-space-xs);
}

.footer-tip {
  font-size: var(--m-font-xs);
}
</style>
