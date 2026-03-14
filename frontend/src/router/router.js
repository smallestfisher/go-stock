import {createMemoryHistory, createRouter, createWebHashHistory, createWebHistory} from 'vue-router'

import stockView from '../components/stock.vue'
import settingsView from '../components/settings.vue'
import aboutView from "../components/about.vue";
import fundView from "../components/fund.vue";
import marketView from "../components/market.vue";
import agentChat from "../components/agent-chat.vue"
import research from "../components/researchIndex.vue";
import cronTaskManager from "../components/cron-task-manager.vue"
import loginView from "../components/Login.vue"

const routes = [
    { path: '/login', component: loginView, name: 'login' },
    { path: '/', component: stockView,name: 'stock'},
    { path: '/fund', component: fundView,name: 'fund' },
    { path: '/settings', component: settingsView,name: 'settings' },
    { path: '/about', component: aboutView,name: 'about' },
    { path: '/market', component: marketView,name: 'market' },
    { path: '/agent', component: agentChat,name: 'agent' },
    { path: '/research', component: research,name: 'research' },
    { path: '/cron-tasks', component: cronTaskManager,name: 'cronTasks' },

]

const router = createRouter({
    //history: createWebHistory(),
    history: createWebHashHistory(),
    routes,
})

// Authentication Guard
router.beforeEach((to, from, next) => {
    const isLoginPath = to.name === 'login';
    const token = localStorage.getItem('auth_token');

    // 如果是登录页，直接放行
    if (isLoginPath) {
        next();
        return;
    }

    // 如果没有 token，强制跳到登录页
    if (!token) {
        next({ name: 'login' });
    } else {
        next();
    }
});

export default router