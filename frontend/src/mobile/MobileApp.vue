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

/* 内容区：flex 占满 shell 剩余高度，自行滚动。
   dock 是 position:fixed 脱离文档流，这里必须预留 dock 的完整高度，
   否则页面底部内容会滑到 dock 下面被遮挡。
   dock 现已去掉底部安全区留白(高度 = 1px边框 + 52px = 53px)，此处同步。
   顶部安全区由 PageHeader 处理。 */
.mobile-main {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  padding-bottom: 53px;
  /* 背景与 dock 栏同色：dock 是 fixed 盖在底部 53px 上，若 dock 实际高度
     与 53px 有 1~2px 出入，这里露出的是同色，避免出现一条异色细缝。 */
  background: var(--m-bg-card);
}
</style>
