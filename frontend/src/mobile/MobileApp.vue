<script setup>
// 移动端独立入口，无桌面端布局
// 导入移动端设计系统和全局工具类
import './styles/mobile.css'
import MobileDock from './components/base/MobileDock.vue'
</script>

<template>
  <div class="mobile-shell">
    <div class="mobile-main">
      <router-view />
    </div>
    <!-- 底部 dock 栏：常驻所有页面（4主tab + 更多） -->
    <MobileDock />
  </div>
</template>

<style>
/* 移动端全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  width: 100%;
  height: 100%;
  overflow: hidden;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>

<style scoped>
.mobile-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 内容区高度 = 视口 - 底部dock高度(约56px+安全区) */
.mobile-main {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  /* 底部留出 dock 空间，避免内容被遮挡 */
  padding-bottom: calc(56px + var(--m-safe-bottom));
}
</style>
