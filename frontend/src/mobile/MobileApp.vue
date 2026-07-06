<script setup>
// 移动端独立入口，无桌面端布局
// 导入移动端设计系统和全局工具类
import './styles/mobile.css'
import MobileDock from './components/base/MobileDock.vue'
import MToastHost from './components/base/MToastHost.vue'
</script>

<template>
  <div class="mobile-shell">
    <div class="mobile-main">
      <router-view />
    </div>
    <!-- 底部 dock 栏：常驻所有页面（4主tab + 更多） -->
    <MobileDock />
    <!-- 全局 toast 宿主 -->
    <MToastHost />
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
  /* iOS 独立 PWA(black-translucent + viewport-fit=cover)下 height:100% 少报
     home indicator 那段物理高度，布局视口短于屏幕，fixed dock 的 bottom:0 锚在
     视口底而非物理屏底，下方露出 manifest background_color(#fff) 成大白带。
     dvh 映射到完整物理屏高，整体内容下移贴到屏幕底；vh 兜底老设备。 */
  height: 100vh;
  height: 100dvh;
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

/* 内容区高度 = 视口 - 底部dock高度(约56px+安全区)，避免内容被遮挡。
   顶部安全区由 PageHeader 处理。 */
.mobile-main {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  /* 底部留出 dock 空间，避免内容被遮挡 */
  padding-bottom: 56px;
}
</style>
