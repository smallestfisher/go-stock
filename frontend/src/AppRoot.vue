<script setup>
import { computed, defineAsyncComponent } from 'vue'
import { useDevice } from './composables/useDevice'

// 桌面端组件（现有的 App.vue）
import DesktopApp from './App.vue'

// 移动端组件（新建的 MobileApp.vue）
const MobileApp = defineAsyncComponent(() => import('./mobile/MobileApp.vue'))

const { isMobile } = useDevice()
</script>

<template>
  <!-- 根据设备类型切换 -->
  <component :is="isMobile ? MobileApp : DesktopApp" />
</template>

<style>
/* 全局样式：移动端和桌面端共用 */
</style>
