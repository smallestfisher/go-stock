import { createRouter, createWebHashHistory } from 'vue-router'

// 移动端页面（懒加载）
const HomePage = () => import('./pages/HomePage.vue')
const StockListPage = () => import('./pages/stock/StockListPage.vue')
const MarketPage = () => import('./pages/market/MarketPage.vue')

// TODO: 后续添加的页面
// const KlineAnalysisPage = () => import('./pages/kline/KlineAnalysisPage.vue')
// const ResearchPage = () => import('./pages/research/ResearchPage.vue')
// ...

const routes = [
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
  // TODO: 后续添加更多路由
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
