<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import MIcon from './MIcon.vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:show'])

const router = useRouter()

// 「更多」抽屉：放次级功能（dock 已常驻首页/自选/市场/研究）。icon 为 MIcon 图标名
const menuItems = [
  { label: 'K线分析', icon: 'kline', path: '/mobile/kline' },
  { label: '基金中心', icon: 'fund', path: '/mobile/fund' },
  { label: 'AI智能体', icon: 'robot', path: '/mobile/agent' },
  { label: '设置', icon: 'settings', path: '/mobile/settings' },
  { label: '关于', icon: 'info', path: '/mobile/about' },
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
            <h2 class="drawer-title">更多功能</h2>
            <button class="drawer-close" @click="handleClose">
              <MIcon name="close" :size="22" />
            </button>
          </div>

          <!-- 菜单列表 -->
          <div class="drawer-menu">
            <div
              v-for="item in menuItems"
              :key="item.path"
              class="menu-item"
              @click="handleNavigate(item.path)"
            >
              <MIcon class="menu-icon" :name="item.icon" :size="22" />
              <span class="menu-label">{{ item.label }}</span>
              <MIcon class="menu-arrow" name="arrow-right" :size="18" />
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
  /* 矢量图标：固定占位宽度让文字对齐，颜色随主色，图标本身尺寸由 size 属性定 */
  width: 32px;
  display: flex;
  justify-content: center;
  color: var(--m-color-rise);
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
