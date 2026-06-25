<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { MenuOutline, SearchOutline, NotificationsOutline } from '@vicons/ionicons5'

const route = useRoute()
const emit = defineEmits(['menuClick', 'searchClick', 'notificationClick'])

// 根据路由显示不同标题
const title = computed(() => {
  const titleMap = {
    '/': 'go-stock',
    '/mobile': 'go-stock',
    '/mobile/stock': '自选',
    '/mobile/market': '市场',
    '/mobile/kline': 'K线',
    '/mobile/research': '研究',
    '/mobile/fund': '基金',
    '/mobile/agent': 'AI智能体',
    '/mobile/settings': '设置',
    '/mobile/about': '关于',
    '/mobile/search': '搜索'
  }
  return titleMap[route.path] || 'go-stock'
})

function handleMenuClick() {
  emit('menuClick')
}

function handleSearchClick() {
  emit('searchClick')
}

function handleNotificationClick() {
  emit('notificationClick')
}
</script>

<template>
  <header class="mobile-header">
    <button class="header-btn" @click="handleMenuClick" aria-label="打开菜单">
      <n-icon size="24">
        <MenuOutline />
      </n-icon>
    </button>

    <h1 class="header-title">{{ title }}</h1>

    <div class="header-actions">
      <button class="header-btn" @click="handleSearchClick" aria-label="搜索">
        <n-icon size="22">
          <SearchOutline />
        </n-icon>
      </button>
      <button class="header-btn" @click="handleNotificationClick" aria-label="通知">
        <n-icon size="22">
          <NotificationsOutline />
        </n-icon>
      </button>
    </div>
  </header>
</template>

<style scoped>
.mobile-header {
  display: flex;
  align-items: center;
  height: calc(var(--m-header-height) + var(--m-safe-top));
  padding-top: var(--m-safe-top);
  padding-inline: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-border-color);
  position: sticky;
  top: 0;
  z-index: var(--m-z-fixed);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.header-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: var(--m-touch-min);
  min-height: var(--m-touch-min);
  background: transparent;
  border: none;
  color: var(--m-text-primary);
  cursor: pointer;
  transition: opacity var(--m-duration-fast);
}

.header-btn:active {
  opacity: 0.6;
}

.header-title {
  flex: 1;
  margin: 0 var(--m-space-md);
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header-actions {
  display: flex;
  gap: var(--m-space-xs);
}
</style>
