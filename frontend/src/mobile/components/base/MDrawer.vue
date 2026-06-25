<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:show'])

const router = useRouter()

const menuItems = [
  { label: '首页', icon: '🏠', path: '/mobile' },
  { label: '自选股票', icon: '⭐', path: '/mobile/stock' },
  { label: '市场行情', icon: '📊', path: '/mobile/market' },
  { label: 'K线分析', icon: '📈', path: '/mobile/kline' },
  { label: '研究中心', icon: '🔬', path: '/mobile/research' },
  { label: '设置', icon: '⚙️', path: '/mobile/settings' },
  { label: '关于', icon: 'ℹ️', path: '/mobile/about' },
]

function handleClose() {
  emit('update:show', false)
}

function handleNavigate(path) {
  router.push(path)
  handleClose()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="show" class="drawer-overlay" @click="handleClose">
        <div class="drawer-content" @click.stop>
          <!-- 头部 -->
          <div class="drawer-header">
            <h2 class="drawer-title">go-stock</h2>
            <button class="drawer-close" @click="handleClose">✕</button>
          </div>

          <!-- 菜单列表 -->
          <div class="drawer-menu">
            <div
              v-for="item in menuItems"
              :key="item.path"
              class="menu-item"
              @click="handleNavigate(item.path)"
            >
              <span class="menu-icon">{{ item.icon }}</span>
              <span class="menu-label">{{ item.label }}</span>
              <span class="menu-arrow">→</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.drawer-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  display: flex;
}

.drawer-content {
  width: 280px;
  height: 100%;
  background: var(--m-bg-card);
  display: flex;
  flex-direction: column;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-xl) var(--m-space-lg);
  border-bottom: 1px solid var(--m-divider-color);
}

.drawer-title {
  font-size: var(--m-font-2xl);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
}

.drawer-close {
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

.drawer-close:active {
  opacity: 0.6;
}

.drawer-menu {
  flex: 1;
  overflow-y: auto;
  padding: var(--m-space-md) 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-lg) var(--m-space-xl);
  color: var(--m-text-primary);
  cursor: pointer;
  transition: background var(--m-duration-fast);
}

.menu-item:active {
  background: var(--m-bg-primary);
}

.menu-icon {
  font-size: 24px;
  width: 32px;
  text-align: center;
}

.menu-label {
  flex: 1;
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
}

.menu-arrow {
  color: var(--m-text-tertiary);
  font-size: var(--m-font-lg);
}

/* 动画 */
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity var(--m-duration-normal);
}

.drawer-enter-active .drawer-content,
.drawer-leave-active .drawer-content {
  transition: transform var(--m-duration-normal);
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-from .drawer-content,
.drawer-leave-to .drawer-content {
  transform: translateX(-100%);
}
</style>
