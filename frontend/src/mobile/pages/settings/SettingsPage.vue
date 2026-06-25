<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import MCard from '../../components/base/MCard.vue'
import MButton from '../../components/base/MButton.vue'
import MTabs from '../../components/base/MTabs.vue'

const router = useRouter()

const settingTabs = [
  { label: '通用', value: 'general' },
  { label: '通知', value: 'notification' },
  { label: '账号', value: 'account' },
]

const activeTab = ref('general')

// 通用设置
const generalSettings = ref({
  theme: 'auto', // 'light' | 'dark' | 'auto'
  language: 'zh-CN',
  fontSize: 'medium', // 'small' | 'medium' | 'large'
})

// 通知设置
const notificationSettings = ref({
  priceAlert: true,
  newsAlert: true,
  aiAnalysis: false,
})

// 账号信息
const accountInfo = ref({
  username: '股票投资者',
  phone: '138****8888',
  email: 'user@example.com',
})

function handleBack() {
  router.back()
}

function handleThemeChange(theme) {
  generalSettings.value.theme = theme
  // TODO: 实际切换主题
}

function handleLogout() {
  console.log('退出登录')
  // TODO: 退出登录逻辑
}
</script>

<template>
  <div class="settings-page">
    <!-- 顶部导航 -->
    <div class="page-header">
      <button class="back-btn" @click="handleBack">←</button>
      <h1 class="page-title">设置</h1>
      <div class="header-placeholder" />
    </div>

    <!-- Tab切换 -->
    <div class="settings-tabs">
      <MTabs v-model="activeTab" :tabs="settingTabs" />
    </div>

    <!-- 内容区 -->
    <div class="settings-content">
      <!-- 通用设置 -->
      <div v-if="activeTab === 'general'" class="settings-section">
        <MCard>
          <div class="setting-item">
            <div class="setting-label">主题</div>
            <div class="theme-options">
              <button
                v-for="theme in ['light', 'dark', 'auto']"
                :key="theme"
                class="theme-btn"
                :class="{ 'theme-btn--active': generalSettings.theme === theme }"
                @click="handleThemeChange(theme)"
              >
                {{ { light: '浅色', dark: '深色', auto: '跟随系统' }[theme] }}
              </button>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-label">字体大小</div>
            <div class="font-options">
              <button
                v-for="size in ['small', 'medium', 'large']"
                :key="size"
                class="font-btn"
                :class="{ 'font-btn--active': generalSettings.fontSize === size }"
                @click="generalSettings.fontSize = size"
              >
                {{ { small: '小', medium: '中', large: '大' }[size] }}
              </button>
            </div>
          </div>
        </MCard>
      </div>

      <!-- 通知设置 -->
      <div v-if="activeTab === 'notification'" class="settings-section">
        <MCard>
          <div class="setting-item" @click="notificationSettings.priceAlert = !notificationSettings.priceAlert">
            <div class="setting-label">价格提醒</div>
            <div class="setting-switch" :class="{ 'setting-switch--on': notificationSettings.priceAlert }" />
          </div>

          <div class="setting-item" @click="notificationSettings.newsAlert = !notificationSettings.newsAlert">
            <div class="setting-label">新闻提醒</div>
            <div class="setting-switch" :class="{ 'setting-switch--on': notificationSettings.newsAlert }" />
          </div>

          <div class="setting-item" @click="notificationSettings.aiAnalysis = !notificationSettings.aiAnalysis">
            <div class="setting-label">AI分析推送</div>
            <div class="setting-switch" :class="{ 'setting-switch--on': notificationSettings.aiAnalysis }" />
          </div>
        </MCard>
      </div>

      <!-- 账号设置 -->
      <div v-if="activeTab === 'account'" class="settings-section">
        <MCard>
          <div class="account-info">
            <div class="info-item">
              <span class="info-label">用户名</span>
              <span class="info-value">{{ accountInfo.username }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">手机号</span>
              <span class="info-value">{{ accountInfo.phone }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">邮箱</span>
              <span class="info-value">{{ accountInfo.email }}</span>
            </div>
          </div>

          <MButton type="danger" block @click="handleLogout">退出登录</MButton>
        </MCard>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.back-btn {
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
}

.back-btn:active {
  opacity: 0.6;
}

.page-title {
  font-size: var(--m-font-xl);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.header-placeholder {
  width: var(--m-touch-min);
}

.settings-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--m-space-md);
}

.settings-section {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--m-space-lg) 0;
  border-bottom: 1px solid var(--m-divider-color);
  cursor: pointer;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-label {
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
}

.theme-options,
.font-options {
  display: flex;
  gap: var(--m-space-sm);
}

.theme-btn,
.font-btn {
  padding: var(--m-space-xs) var(--m-space-md);
  background: var(--m-bg-primary);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.theme-btn--active,
.font-btn--active {
  background: var(--m-color-rise);
  color: white;
  border-color: var(--m-color-rise);
}

.setting-switch {
  width: 48px;
  height: 28px;
  background: var(--m-text-tertiary);
  border-radius: var(--m-radius-full);
  position: relative;
  transition: background var(--m-duration-normal);
}

.setting-switch::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 24px;
  height: 24px;
  background: white;
  border-radius: 50%;
  transition: transform var(--m-duration-normal);
}

.setting-switch--on {
  background: var(--m-color-rise);
}

.setting-switch--on::after {
  transform: translateX(20px);
}

.account-info {
  margin-bottom: var(--m-space-xl);
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: var(--m-space-lg) 0;
  border-bottom: 1px solid var(--m-divider-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: var(--m-text-secondary);
}

.info-value {
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
}
</style>
