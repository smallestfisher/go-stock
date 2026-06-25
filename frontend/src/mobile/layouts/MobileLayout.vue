<script setup>
import { ref } from 'vue'
import MobileHeader from './MobileHeader.vue'
import MobileSideDrawer from './MobileSideDrawer.vue'

const props = defineProps({
  drawerVisible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:drawerVisible'])

function handleMenuClick() {
  emit('update:drawerVisible', true)
}

function handleDrawerUpdate(visible) {
  emit('update:drawerVisible', visible)
}
</script>

<template>
  <div class="mobile-layout">
    <!-- 顶部固定栏 -->
    <MobileHeader @menu-click="handleMenuClick" />

    <!-- 侧边抽屉 -->
    <MobileSideDrawer
      :visible="drawerVisible"
      @update:visible="handleDrawerUpdate"
    />

    <!-- 内容区（路由视图） -->
    <main class="mobile-content mobile-scrollable">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.mobile-layout {
  display: flex;
  flex-direction: column;
  height: 100dvh;
  width: 100vw;
  overflow: hidden;
}

.mobile-content {
  flex: 1;
  position: relative;
  overflow-y: auto;
  overflow-x: hidden;
  background: var(--m-bg-primary);
}
</style>
