import { createRouter, createWebHashHistory } from 'vue-router'

// 移动端页面（懒加载）
const HomePage = () => import('./pages/HomePage.vue')
const StockListPage = () => import('./pages/stock/StockListPage.vue')
const MarketPage = () => import('./pages/market/MarketPage.vue')
const KlineAnalysisPage = () => import('./pages/kline/KlineAnalysisPage.vue')
const ResearchPage = () => import('./pages/research/ResearchPage.vue')
const FundPage = () => import('./pages/fund/FundPage.vue')
const AgentPage = () => import('./pages/agent/AgentPage.vue')
const SettingsPage = () => import('./pages/settings/SettingsPage.vue')
const AboutPage = () => import('./pages/about/AboutPage.vue')

const routes = [
  {
    path: '/',
    redirect: '/mobile'
  },
  {
    path: '/stock',
    redirect: '/mobile/stock'
  },
  {
    path: '/market',
    redirect: '/mobile/market'
  },
  {
    path: '/kline-analysis',
    redirect: '/mobile/kline'
  },
  {
    path: '/research',
    redirect: '/mobile/research'
  },
  {
    path: '/fund',
    redirect: '/mobile/fund'
  },
  {
    path: '/agent',
    redirect: '/mobile/agent'
  },
  {
    path: '/settings',
    redirect: '/mobile/settings'
  },
  {
    path: '/about',
    redirect: '/mobile/about'
  },
  {
    path: '/mobile',
    name: 'MobileHome',
    component: HomePage,
    meta: { title: 'go-stock' }
  },
  {
    path: '/mobile/stock',
    name: 'MobileStock',
    component: StockListPage,
    meta: { title: '自选' }
  },
  {
    path: '/mobile/market',
    name: 'MobileMarket',
    component: MarketPage,
    meta: { title: '市场' }
  },
  {
    path: '/mobile/kline',
    name: 'MobileKline',
    component: KlineAnalysisPage,
    meta: { title: 'K线分析' }
  },
  {
    path: '/mobile/research',
    name: 'MobileResearch',
    component: ResearchPage,
    meta: { title: '研究中心' }
  },
  {
    path: '/mobile/fund',
    name: 'MobileFund',
    component: FundPage,
    meta: { title: '基金中心' }
  },
  {
    path: '/mobile/agent',
    name: 'MobileAgent',
    component: AgentPage,
    meta: { title: 'AI智能体' }
  },
  {
    path: '/mobile/settings',
    name: 'MobileSettings',
    component: SettingsPage,
    meta: { title: '设置' }
  },
  {
    path: '/mobile/about',
    name: 'MobileAbout',
    component: AboutPage,
    meta: { title: '关于' }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/mobile'
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫：设置页面标题
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} - go-stock`
  }
  next()
})

export default router
