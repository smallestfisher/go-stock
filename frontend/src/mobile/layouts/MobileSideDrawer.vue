<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  StarOutline,
  NewspaperOutline,
  AnalyticsOutline,
  FlaskOutline,
  SparklesOutline,
  RocketOutline,
  SettingsOutline,
  InformationOutline
} from '@vicons/ionicons5'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible'])
const router = useRouter()

// 主功能菜单
const menuItems = [
  { icon: StarOutline, label: '我的自选', path: '/mobile/stock' },
  { icon: NewspaperOutline, label: '市场行情', path: '/mobile/market' },
  { icon: AnalyticsOutline, label: 'K线分析', path: '/mobile/kline' },
  { icon: FlaskOutline, label: '研究中心', path: '/mobile/research' },
  { icon: SparklesOutline, label: '基金中心', path: '/mobile/fund' },
  { icon: RocketOutline, label: 'AI智能体', path: '/mobile/agent' },
]

// 设置菜单
const settingItems = [
  { icon: SettingsOutline, label: '设置', path: '/mobile/settings' },
  { icon: InformationOutline, label: '关于', path: '/mobile/about' },
]

function navigate(path) {
  router.push(path)
  emit('update:visible', false)
}

function close() {
  emit('update:visible', false)
}
</script>

<template>
  <!-- 遮罩层 -->
  <Transition name="m-mask">
    <div v-if="visible" class="drawer-mask" @click="close" />
  </Transition>

  <!-- 抽屉内容 -->
  <Transition name="m-slide-left">
    <aside v-if="visible" class="drawer-content">
      <!-- 用户信息区 -->
      <div class="drawer-header">
        <div class="user-avatar">
          <n-icon size="32">
            <StarOutline />
          </n-icon>
        </div>
        <div class="user-info">
          <div class="user-name">股票投资者</div>
          <div class="user-profit">
            今日盈亏: <span class="m-rise">+¥1,234</span>
          </div>
        </div>
      </div>

      <!-- 主菜单 -->
      <nav class="drawer-nav">
        <button
          v-for="item in menuItems"
          :key="item.path"
          class="nav-item"
          @click="navigate(item.path)"
        >
          <n-icon size="20" class="nav-icon">
            <component :is="item.icon" />
          </n-icon>
          <span class="nav-label">{{ item.label }}</span>
        </button>
      </nav>

      <div class="drawer-divider" />

      <!-- 设置菜单 -->
      <nav class="drawer-nav">
        <button
          v-for="item in settingItems"
          :key="item.path"
          class="nav-item"
          @click="navigate(item.path)"
        >
          <n-icon size="20" class="nav-icon">
            <component :is="item.icon" />
          </n-icon>
          <span class="nav-label">{{ item.label }}</span>
        </button>
      </nav>

      <!-- 底部版本信息 -->
      <div class="drawer-footer">
        <div class="version-info">go-stock v1.0.0</div>
      </div>
    </aside>
  </Transition>
</template>

<style scoped>
.drawer-mask {
  position: fixed;
  inset: 0;
  background: var(--m-bg-overlay);
  z-index: var(--m-z-mask);
}

.drawer-content {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 280px;
  max-width: 80vw;
  background: var(--m-bg-card);
  z-index: calc(var(--m-z-mask) + 1);
  display: flex;
  flex-direction: column;
  box-shadow: var(--m-shadow-xl);
}

.drawer-header {
  padding: calc(var(--m-safe-top) + var(--m-space-xl)) var(--m-space-xl) var(--m-space-xl);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.user-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--m-color-rise), var(--m-color-rise-hover));
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--m-space-md);
  color: white;
}

.user-name {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  margin-bottom: var(--m-space-xs);
  color: var(--m-text-primary);
}

.user-profit {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.drawer-nav {
  flex: 1;
  padding: var(--m-space-md) 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  width: 100%;
  min-height: var(--m-touch-min);
  padding: var(--m-space-md) var(--m-space-xl);
  background: transparent;
  border: none;
  color: var(--m-text-primary);
  cursor: pointer;
  transition: background var(--m-duration-fast);
  text-align: left;
}

.nav-item:active {
  background: var(--m-bg-primary);
}

.nav-icon {
  margin-right: var(--m-space-md);
  color: var(--m-text-secondary);
  flex-shrink: 0;
}

.nav-label {
  font-size: var(--m-font-md);
  flex: 1;
}

.drawer-divider {
  height: 1px;
  background: var(--m-divider-color);
  margin: 0 var(--m-space-xl);
}

.drawer-footer {
  padding: var(--m-space-lg) var(--m-space-xl);
  padding-bottom: calc(var(--m-space-lg) + var(--m-safe-bottom));
  border-top: 1px solid var(--m-divider-color);
}

.version-info {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  text-align: center;
}
</style>
