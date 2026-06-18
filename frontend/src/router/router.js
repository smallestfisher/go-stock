import {createRouter, createWebHashHistory} from 'vue-router'

const stockView = () => import('../components/stock.vue')
const settingsView = () => import('../components/settings.vue')
const aboutView = () => import('../components/about.vue')
const fundView = () => import('../components/fund.vue')
const marketView = () => import('../components/market.vue')
const agentChat = () => import('../components/agent-chat.vue')
const research = () => import('../components/researchIndex.vue')
const cronTaskManager = () => import('../components/cron-task-manager.vue')
const mcpServerManager = () => import('../components/mcp-server-manager.vue')
const klineAnalysis = () => import('../components/kline-analysis.vue')

const routes = [
    { path: '/', component: stockView,name: 'stock'},
    { path: '/fund', component: fundView,name: 'fund' },
    { path: '/settings', component: settingsView,name: 'settings' },
    { path: '/about', component: aboutView,name: 'about' },
    { path: '/market', component: marketView,name: 'market' },
    { path: '/agent', component: agentChat,name: 'agent' },
    { path: '/research', component: research,name: 'research' },
    { path: '/cron-tasks', component: cronTaskManager,name: 'cronTasks' },
    { path: '/mcp-servers', component: mcpServerManager,name: 'mcpServers' },
    { path: '/kline-analysis', component: klineAnalysis,name: 'klineAnalysis' },

]

const router = createRouter({
    //history: createWebHistory(),
    history: createWebHashHistory(),
    routes,
})

export default router
