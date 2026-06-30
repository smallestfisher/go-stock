<script setup>
import { useRouter } from 'vue-router'

defineProps({
  title: {
    type: String,
    default: 'go-stock'
  },
  showBack: {
    type: Boolean,
    default: false
  },
  showMenu: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['menu-click'])

const router = useRouter()

function handleBack() {
  router.back()
}

function handleMenu() {
  emit('menu-click')
}
</script>

<template>
  <div class="page-header">
    <button v-if="showBack" class="header-btn" @click="handleBack">←</button>
    <button v-else-if="showMenu" class="header-btn" @click="handleMenu">☰</button>
    <div v-else class="header-placeholder" />

    <h1 class="header-title">{{ title }}</h1>

    <div class="header-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  /* 顶部安全区内边距：viewport-fit=cover 下避开状态栏/刘海，
     不靠整体上移布局，仅本栏向下让出安全区 */
  padding: var(--m-safe-top) var(--m-space-md) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  position: sticky;
  top: 0;
  z-index: var(--m-z-sticky);
  min-height: calc(56px + var(--m-safe-top));
}

.header-btn {
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
  flex-shrink: 0;
}

.header-btn:active {
  opacity: 0.6;
}

.header-placeholder {
  width: var(--m-touch-min);
  flex-shrink: 0;
}

.header-title {
  flex: 1;
  font-size: var(--m-font-xl);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  text-align: center;
}

.header-actions {
  display: flex;
  gap: var(--m-space-sm);
  flex-shrink: 0;
  min-width: var(--m-touch-min);
  justify-content: flex-end;
}
</style>
