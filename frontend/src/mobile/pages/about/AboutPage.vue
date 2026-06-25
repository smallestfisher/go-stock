<script setup>
import { ref, onBeforeMount } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MCard from '../../components/base/MCard.vue'
import MLoading from '../../components/base/MLoading.vue'
import { GetVersionInfo } from '../../../api/app'

// 应用信息（后续接 GetVersionInfo 填充，对齐桌面端 about.vue）
// 返回字段：{ version, content(更新日志), icon, alipay, wxpay, wxgzh }
const versionInfo = ref(null)
const updateLog = ref('')
const icon = ref('')
const loading = ref(false)

async function loadInfo() {
  loading.value = true
  try {
    const res = await GetVersionInfo()
    versionInfo.value = res
    updateLog.value = res?.content || ''
    icon.value = res?.icon || ''
  } catch (error) {
    console.error('加载版本信息失败:', error)
  } finally {
    loading.value = false
  }
}

onBeforeMount(loadInfo)
</script>

<template>
  <div class="about-page">
    <!-- 顶部导航 -->
    <PageHeader title="关于" />

    <!-- 内容 -->
    <div class="about-content">
      <MLoading v-if="loading" text="加载中..." />

      <template v-else>
        <!-- 应用信息 -->
        <MCard>
          <img v-if="icon" :src="icon" class="app-logo" alt="logo">
          <div v-else class="app-logo-emoji">📈</div>
          <h2 class="app-name">go-stock</h2>
          <p v-if="versionInfo?.version" class="app-version">
            版本 {{ versionInfo.version }}
          </p>
          <p class="app-desc">基于大语言模型的 AI 赋能股票分析工具，支持 A股、港股、美股</p>
        </MCard>

        <!-- 更新日志 -->
        <MCard v-if="updateLog">
          <h3 class="section-title">更新日志</h3>
          <div class="update-log" v-html="updateLog" />
        </MCard>

        <!-- 链接 -->
        <MCard>
          <div class="info-item">
            <span class="info-label">开源协议</span>
            <span class="info-value">GPLv3</span>
          </div>
          <div class="info-item">
            <span class="info-label">项目地址</span>
            <span class="info-value">GitHub</span>
          </div>
        </MCard>
      </template>
    </div>
  </div>
</template>

<style scoped>
.about-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.about-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.app-logo {
  width: 80px;
  height: 80px;
  display: block;
  margin: 0 auto var(--m-space-lg);
  border-radius: var(--m-radius-lg);
}

.app-logo-emoji {
  font-size: 64px;
  text-align: center;
  margin-bottom: var(--m-space-lg);
}

.app-name {
  font-size: var(--m-font-2xl);
  font-weight: var(--m-font-weight-bold);
  text-align: center;
  margin-bottom: var(--m-space-xs);
}

.app-version {
  text-align: center;
  color: var(--m-text-tertiary);
  font-size: var(--m-font-sm);
  margin-bottom: var(--m-space-md);
}

.app-desc {
  text-align: center;
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-normal);
}

.section-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  margin-bottom: var(--m-space-md);
}

.update-log {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-loose);
  max-height: 400px;
  overflow-y: auto;
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
