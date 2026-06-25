<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import MPullRefresh from '../components/base/MPullRefresh.vue'
import MDrawer from '../components/base/MDrawer.vue'
import MarketStatusBar from '../components/widgets/MarketStatusBar.vue'
import StockSummaryCard from '../components/cards/StockSummaryCard.vue'
import NewsCard from '../components/cards/NewsCard.vue'
import HotTopicCard from '../components/cards/HotTopicCard.vue'
import AlertCard from '../components/cards/AlertCard.vue'
import IndustryCard from '../components/cards/IndustryCard.vue'
import AiSuggestCard from '../components/cards/AiSuggestCard.vue'
import MCard from '../components/base/MCard.vue'

const router = useRouter()

const drawerVisible = ref(false)

// 模拟数据（后续对接真实API）
const stockData = ref([
  { code: '600519', name: '贵州茅台', price: 1820.50, changePercent: 2.34 },
  { code: '000858', name: '五粮液', price: 156.80, changePercent: -1.23 },
  { code: '00700', name: '腾讯控股', price: 358.20, changePercent: 0.85 },
])

const newsData = ref([
  { title: '央行宣布降准0.5个百分点', time: new Date(Date.now() - 1800000), source: '财联社' },
  { title: 'A股三大指数集体低开，半导体板块领跌', time: new Date(Date.now() - 3600000), source: '证券时报' },
  { title: '外资净流入50亿元，连续五日加仓', time: new Date(Date.now() - 5400000), source: '第一财经' },
])

const topicsData = ref([
  { title: 'AI芯片', heat: '1.2M', changePercent: 5.67, stocks: 23 },
  { title: '新能源汽车', heat: '980K', changePercent: 3.45, stocks: 45 },
  { title: 'ChatGPT概念', heat: '850K', changePercent: -2.12, stocks: 18 },
])

const alertsData = ref([
  { stockName: '寒武纪', stockCode: '688256', type: 'limit_up', changePercent: 10.00, time: new Date() },
  { stockName: '中芯国际', stockCode: '688981', type: 'rapid_rise', changePercent: 7.89, time: new Date(Date.now() - 600000) },
])

const industriesData = ref([
  { name: '半导体', changePercent: 5.23, leadingStock: '寒武纪 +10.00%' },
  { name: '新能源', changePercent: 3.87, leadingStock: '宁德时代 +6.54%' },
  { name: '人工智能', changePercent: 2.95, leadingStock: '科大讯飞 +5.32%' },
  { name: '医药生物', changePercent: 1.45, leadingStock: '恒瑞医药 +3.21%' },
  { name: '白酒', changePercent: -0.89, leadingStock: '贵州茅台 +2.34%' },
])

const aiSuggestion = ref({
  title: '基于你的自选，建议关注半导体板块',
  content: '近期AI芯片需求激增，半导体板块表现强势。你的自选中腾讯控股与多家半导体公司有业务合作，可关注相关概念股。',
  stocks: [
    { code: '688256', name: '寒武纪' },
    { code: '688981', name: '中芯国际' },
  ],
  action: '查看详细分析'
})

const aiLoading = ref(false)
const refreshCount = ref(0)

// 下拉刷新
async function handleRefresh() {
  return new Promise(resolve => {
    setTimeout(() => {
      refreshCount.value++
      // 这里后续对接真实API刷新数据
      resolve()
    }, 1500)
  })
}

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

function handleAiRefresh() {
  aiLoading.value = true
  setTimeout(() => {
    aiLoading.value = false
    // TODO: 调用AI接口
  }, 2000)
}

function handleAiAction() {
  navigateTo('/mobile/research')
}
</script>

<template>
  <div class="home-container">
    <!-- 顶部导航栏 -->
    <div class="home-header">
      <button class="menu-btn" @click="drawerVisible = true">☰</button>
      <h1 class="home-title">go-stock</h1>
      <div class="header-actions">
        <button class="search-btn">🔍</button>
        <button class="notification-btn">🔔</button>
      </div>
    </div>

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

      <!-- AI建议卡片 -->
      <AiSuggestCard
        :suggestion="aiSuggestion"
        :loading="aiLoading"
        @refresh="handleAiRefresh"
        @action-click="handleAiAction"
        @stock-click="handleStockClick"
      />

      <!-- 底部提示 -->
      <div class="home-footer">
        <p class="footer-text">下拉刷新数据 · 已刷新 {{ refreshCount }} 次</p>
        <p class="footer-tip">💡 点击卡片查看详情</p>
      </div>
    </div>
  </MPullRefresh>

  <!-- 侧边抽屉 -->
  <MDrawer v-model:show="drawerVisible" />
</div>
</template>

<style scoped>
.home-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  position: sticky;
  top: 0;
  z-index: var(--m-z-sticky);
}

.menu-btn {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 24px;
  color: var(--m-text-primary);
  cursor: pointer;
}

.menu-btn:active {
  opacity: 0.6;
}

.home-title {
  font-size: var(--m-font-xl);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
}

.header-actions {
  display: flex;
  gap: var(--m-space-sm);
}

.search-btn,
.notification-btn {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 20px;
  color: var(--m-text-primary);
  cursor: pointer;
}

.search-btn:active,
.notification-btn:active {
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
