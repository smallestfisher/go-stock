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
  /* iOS 独立 PWA(black-translucent + viewport-fit=cover)下，height:100% 会少报
     home indicator 那段物理高度，导致布局视口短于屏幕，fixed dock 的 bottom:0
     锚在视口底而非物理屏底，下方露出 manifest 的 background_color(#fff) 成大白带。
     dvh/vh 在独立模式下映射到完整物理屏高，故用 dvh 打底、vh 兜底老设备。 */
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

/* 内容区：flex 占满 shell 剩余高度，自行滚动。
   dock 是 position:fixed 脱离文档流，这里必须预留 dock 的完整高度，
   否则页面底部内容会滑到 dock 下面被遮挡。
   dock 现已去掉底部安全区留白(高度 = 1px边框 + 52px = 53px)，此处同步。
   顶部安全区由 PageHeader 处理。 */
.mobile-main {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  /* dock 是 fixed 脱离文档流，这里预留其完整高度避免遮挡底部内容。
     dock 总高 = 1px边框 + 52px内容 + 底部安全区(env)，此处同步。 */
  padding-bottom: calc(53px + var(--m-safe-bottom));
}
</style>
