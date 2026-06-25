<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import MCard from '../../components/base/MCard.vue'

const router = useRouter()

const appInfo = ref({
  name: 'go-stock',
  version: '1.0.0',
  buildTime: '2024-06-20',
  description: '移动端股票分析工具'
})

function handleBack() {
  router.back()
}
</script>

<template>
  <div class="about-page">
    <!-- 顶部导航 -->
    <div class="page-header">
      <button class="back-btn" @click="handleBack">←</button>
      <h1 class="page-title">关于</h1>
      <div class="header-placeholder" />
    </div>

    <!-- 内容 -->
    <div class="about-content">
      <!-- 应用信息 -->
      <MCard>
        <div class="app-logo">📈</div>
        <h2 class="app-name">{{ appInfo.name }}</h2>
        <p class="app-version">版本 {{ appInfo.version }}</p>
        <p class="app-desc">{{ appInfo.description }}</p>
      </MCard>

      <!-- 详细信息 -->
      <MCard>
        <div class="info-item">
          <span class="info-label">构建时间</span>
          <span class="info-value">{{ appInfo.buildTime }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">技术栈</span>
          <span class="info-value">Vue 3 + Vite</span>
        </div>
        <div class="info-item">
          <span class="info-label">开源协议</span>
          <span class="info-value">MIT</span>
        </div>
      </MCard>

      <!-- 链接 -->
      <MCard>
        <div class="link-item">用户协议</div>
        <div class="link-item">隐私政策</div>
        <div class="link-item">问题反馈</div>
      </MCard>
    </div>
  </div>
</template>

<style scoped>
.about-page { display: flex; flex-direction: column; height: 100%; background: var(--m-bg-primary); }
.page-header { display: flex; align-items: center; justify-content: space-between; padding: var(--m-space-md); background: var(--m-bg-card); border-bottom: 1px solid var(--m-divider-color); }
.back-btn { width: var(--m-touch-min); height: var(--m-touch-min); display: flex; align-items: center; justify-content: center; background: transparent; border: none; font-size: 24px; color: var(--m-text-primary); cursor: pointer; }
.back-btn:active { opacity: 0.6; }
.page-title { font-size: var(--m-font-xl); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.header-placeholder { width: var(--m-touch-min); }
.about-content { flex: 1; overflow-y: auto; padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.app-logo { font-size: 64px; text-align: center; margin-bottom: var(--m-space-lg); }
.app-name { font-size: var(--m-font-2xl); font-weight: var(--m-font-weight-bold); text-align: center; margin-bottom: var(--m-space-xs); }
.app-version { text-align: center; color: var(--m-text-tertiary); font-size: var(--m-font-sm); margin-bottom: var(--m-space-md); }
.app-desc { text-align: center; color: var(--m-text-secondary); font-size: var(--m-font-sm); }
.info-item { display: flex; justify-content: space-between; padding: var(--m-space-lg) 0; border-bottom: 1px solid var(--m-divider-color); }
.info-item:last-child { border-bottom: none; }
.info-label { color: var(--m-text-secondary); }
.info-value { color: var(--m-text-primary); font-weight: var(--m-font-weight-medium); }
.link-item { padding: var(--m-space-lg) 0; color: var(--m-text-primary); border-bottom: 1px solid var(--m-divider-color); cursor: pointer; }
.link-item:last-child { border-bottom: none; }
.link-item:active { opacity: 0.6; }
</style>
