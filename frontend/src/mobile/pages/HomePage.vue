<script setup>
import { ref, onMounted } from 'vue'

const greeting = ref('Hello Mobile!')
const deviceInfo = ref({
  width: 0,
  height: 0,
  userAgent: ''
})

onMounted(() => {
  deviceInfo.value = {
    width: window.innerWidth,
    height: window.innerHeight,
    userAgent: navigator.userAgent
  }
})
</script>

<template>
  <div class="home-page">
    <div class="welcome-card">
      <h1 class="welcome-title">{{ greeting }}</h1>
      <p class="welcome-subtitle">移动端架构已就绪 🎉</p>

      <div class="device-info">
        <div class="info-item">
          <span class="info-label">设备宽度:</span>
          <span class="info-value">{{ deviceInfo.width }}px</span>
        </div>
        <div class="info-item">
          <span class="info-label">设备高度:</span>
          <span class="info-value">{{ deviceInfo.height }}px</span>
        </div>
      </div>

      <div class="feature-list">
        <div class="feature-item">✅ 基础架构搭建完成</div>
        <div class="feature-item">✅ 侧边抽屉导航就绪</div>
        <div class="feature-item">✅ 设计系统已建立</div>
        <div class="feature-item">✅ 路由系统已配置</div>
      </div>

      <div class="next-steps">
        <p class="next-title">接下来：</p>
        <ul class="next-list">
          <li>开发基础组件库</li>
          <li>实现首页信息流</li>
          <li>完成自选股票模块</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-page {
  min-height: 100%;
  padding: var(--m-space-xl);
  background: var(--m-bg-primary);
}

.welcome-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-lg);
  padding: var(--m-space-xl);
  box-shadow: var(--m-shadow-md);
}

.welcome-title {
  font-size: var(--m-font-2xl);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-color-rise);
  margin-bottom: var(--m-space-sm);
  text-align: center;
}

.welcome-subtitle {
  font-size: var(--m-font-lg);
  color: var(--m-text-secondary);
  text-align: center;
  margin-bottom: var(--m-space-xl);
}

.device-info {
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  padding: var(--m-space-md);
  margin-bottom: var(--m-space-lg);
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: var(--m-space-sm) 0;
  font-size: var(--m-font-sm);
}

.info-label {
  color: var(--m-text-secondary);
}

.info-value {
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
}

.feature-list {
  margin-bottom: var(--m-space-xl);
}

.feature-item {
  padding: var(--m-space-md);
  margin-bottom: var(--m-space-sm);
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
}

.next-steps {
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  padding: var(--m-space-lg);
}

.next-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  margin-bottom: var(--m-space-md);
}

.next-list {
  margin-left: var(--m-space-lg);
  color: var(--m-text-secondary);
  font-size: var(--m-font-md);
  line-height: var(--m-line-height-loose);
}
</style>
