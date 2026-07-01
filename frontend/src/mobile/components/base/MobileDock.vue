<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MDrawer from './MDrawer.vue'
import MIcon from './MIcon.vue'

// 底部 dock 栏：4 个高频 tab + 「更多」入口（点开抽屉放次级功能）
// 常驻所有页面，对齐移动端主流交互

const router = useRouter()
const route = useRoute()

// dock 主 tab（对应一级路由）。icon 为 MIcon 图标名（矢量，吃 currentColor）
const dockItems = [
  { label: '首页', icon: 'home', path: '/mobile' },
  { label: '自选', icon: 'star', path: '/mobile/stock' },
  { label: '市场', icon: 'chart', path: '/mobile/market' },
  { label: '研究', icon: 'research', path: '/mobile/research' },
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
        <MIcon class="dock-icon" :name="item.icon" :size="24" />
        <span class="dock-label">{{ item.label }}</span>
      </button>

      <!-- 更多 -->
      <button
        class="dock-item"
        :class="{ 'dock-item--active': moreVisible }"
        @click="showMore"
      >
        <MIcon class="dock-icon" name="menu" :size="24" />
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
  /* 根因修复：black-translucent + viewport-fit=cover 下，高度链是 height:100%，
     视口底不含 home indicator 安全区，dock 的 bottom:0 锚到视口底而非屏幕物理底，
     二者之间那条 ~34px 安全区会露出 body 的浅灰背景(看着像"空白带/不是纯白")。
     这里让 fixed 容器自身铺 dock 同色背景、并向下吃掉安全区，
     把那条带子变成 dock 白色背景的自然延伸，而非露底色的空带。 */
  background: var(--m-bg-card);
  padding-bottom: var(--m-safe-bottom);
}

.dock-bar {
  display: flex;
  background: var(--m-bg-card);
  border-top: 1px solid var(--m-divider-color);
  /* dock 可点内容保持在安全区之上：安全区由外层 .mobile-dock 的
     padding-bottom 让出，这里不再额外留白，避免图标落入手势区。 */
  padding-bottom: 0;
  box-shadow: 0 -1px 8px rgba(0, 0, 0, 0.04);
}

.dock-item {
  position: relative;
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

/* 选中态顶部指示条：emoji 图标无视 color，靠这条 + 去色对比来强化高亮。 */
.dock-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  width: 24px;
  height: 3px;
  border-radius: 0 0 var(--m-radius-full) var(--m-radius-full);
  background: var(--m-color-rise);
  transform: translateX(-50%) scaleX(0);
  transition: transform var(--m-duration-fast);
}

.dock-item--active {
  color: var(--m-color-rise);
}

.dock-item--active::before {
  transform: translateX(-50%) scaleX(1);
}

.dock-icon {
  /* 矢量图标吃 currentColor：颜色随 .dock-item 的 color 走，
     选中态由下方 --active 的 color 直接驱动，无需 emoji 去色 hack。 */
  transition: transform var(--m-duration-fast);
}

.dock-item--active .dock-icon {
  transform: scale(1.1);
}

.dock-label {
  font-size: 11px;
  line-height: 1;
}

.dock-item--active .dock-label {
  font-weight: var(--m-font-weight-medium);
}

.dock-item:active {
  opacity: 0.7;
}
</style>
