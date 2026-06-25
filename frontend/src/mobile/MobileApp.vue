<script setup>
import { ref, provide, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MobileLayout from './layouts/MobileLayout.vue'

// 引入移动端样式
import './styles/mobile.css'
import './styles/animations.css'

const router = useRouter()
const drawerVisible = ref(false)

// 全局状态注入：抽屉控制
provide('toggleDrawer', () => {
  drawerVisible.value = !drawerVisible.value
})

provide('closeDrawer', () => {
  drawerVisible.value = false
})

provide('openDrawer', () => {
  drawerVisible.value = true
})

onMounted(() => {
  console.log('📱 Mobile App Initialized')
})
</script>

<template>
  <div class="mobile-app">
    <MobileLayout
      :drawer-visible="drawerVisible"
      @update:drawer-visible="drawerVisible = $event"
    >
      <RouterView />
    </MobileLayout>
  </div>
</template>

<style scoped>
.mobile-app {
  width: 100vw;
  height: 100dvh;
  overflow: hidden;
  background: var(--m-bg-primary);
}
</style>
