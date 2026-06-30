<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MDrawer from './MDrawer.vue'

// 底部 dock 栏：4 个高频 tab + 「更多」入口（点开抽屉放次级功能）
// 常驻所有页面，对齐移动端主流交互

const router = useRouter()
const route = useRoute()

// dock 主 tab（对应一级路由）
const dockItems = [
  { label: '首页', icon: '🏠', path: '/mobile' },
  { label: '自选', icon: '⭐', path: '/mobile/stock' },
  { label: '市场', icon: '📊', path: '/mobile/market' },
  { label: '研究', icon: '🔬', path: '/mobile/research' },
]

const moreVisible = ref(false)

// 当前激活的 tab（路径前缀匹配，如 /mobile/stock/xxx 也高亮自选）
const activePath = computed(() => {
  const path = route.path
  // 首页精确匹配，其他用前缀
  if (path === '/mobile') return '/mobile'
  const matched = dockItems.find(item => item.path !== '/mobile' && path.startsWith(item.path))
  return matched ? matched.path : ''
})

function navigate(path) {
  if (route.path !== path) {
    router.push(path)
  }
}

function showMore() {
  moreVisible.value = true
}
</script>

<template>
  <div class="mobile-dock">
    <div class="dock-bar">
      <!-- 4 个主 tab -->
      <button
        v-for="item in dockItems"
        :key="item.path"
        class="dock-item"
        :class="{ 'dock-item--active': activePath === item.path }"
        @click="navigate(item.path)"
      >
        <span class="dock-icon">{{ item.icon }}</span>
        <span class="dock-label">{{ item.label }}</span>
      </button>

      <!-- 更多 -->
      <button class="dock-item" @click="showMore">
        <span class="dock-icon">☰</span>
        <span class="dock-label">更多</span>
      </button>
    </div>

    <!-- 更多抽屉（K线/基金/智能体/设置/关于） -->
    <MDrawer v-model:show="moreVisible" />
  </div>
</template>

<style scoped>
.mobile-dock {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: var(--m-z-fixed);
}

.dock-bar {
  display: flex;
  background: var(--m-bg-card);
  border-top: 1px solid var(--m-divider-color);
  /* 安全区底部留白 + 固定兜底，避免压到 iPhone 底部 home indicator 小白条 */
  padding-bottom: calc(var(--m-safe-bottom) + 6px);
  box-shadow: 0 -1px 8px rgba(0, 0, 0, 0.04);
}

.dock-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 3px;
  padding: var(--m-space-sm) 0 var(--m-space-xs);
  background: transparent;
  border: none;
  color: var(--m-text-tertiary);
  cursor: pointer;
  min-height: 52px;
  transition: color var(--m-duration-fast);
}

.dock-item--active {
  color: var(--m-color-rise);
}

.dock-icon {
  font-size: 24px;
  line-height: 1;
}

.dock-label {
  font-size: 11px;
  line-height: 1;
}

.dock-item:active {
  opacity: 0.7;
}
</style>
