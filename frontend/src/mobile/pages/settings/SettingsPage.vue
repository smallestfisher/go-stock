<script setup>
import { ref, reactive, onBeforeMount } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'
import { GetConfig, UpdateConfig } from '../../../api/app'

// 对齐桌面端 settings.vue 的配置项，用 GetConfig 读取 / UpdateConfig 保存
// 配置模型见 api/models.ts 的 SettingConfig / AIConfig

const groupTabs = [
  { label: '基本', value: 'general' },
  { label: 'AI诊股', value: 'ai' },
  { label: '通知', value: 'push' },
  { label: '数据源', value: 'source' },
]
const activeTab = ref('general')

const loading = ref(false)
const saving = ref(false)
const savedTip = ref('')

// 配置表单（字段对齐后端 GetConfig 返回结构）
const config = reactive({
  ID: 1,
  // 基本
  refreshInterval: 1,
  darkTheme: false,
  browserPath: '',
  enableFund: false,
  enableNews: true,
  enableDanmu: false,
  enableAgent: false,
  updateBasicInfoOnStart: true,
  // AI诊股
  openAiEnable: false,
  prompt: '',
  questionTemplate: '{{stockName}}分析和总结',
  crawlTimeOut: 60,
  kDays: 30,
  aiConfigs: [],
  // 通知
  dingPushEnable: false,
  dingRobot: '',
  enablePushNews: true,
  enableOnlyPushRedNews: false,
  // 数据源/代理
  httpProxyEnabled: false,
  httpProxy: '',
  tushareToken: '',
  iwencaiApiKey: '',
  emApiKey: '',
  qgqpBId: '',
  promptPlazaApiBase: '',
})

async function loadConfig() {
  loading.value = true
  try {
    const res = await GetConfig()
    Object.assign(config, res)
  } catch (error) {
    console.error('加载配置失败:', error)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const res = await UpdateConfig(config)
    savedTip.value = res === '保存成功！' ? '✓ 已保存' : ('保存失败：' + res)
    setTimeout(() => { savedTip.value = '' }, 2000)
  } catch (error) {
    console.error('保存配置失败:', error)
    savedTip.value = '保存失败'
    setTimeout(() => { savedTip.value = '' }, 2000)
  } finally {
    saving.value = false
  }
}

// AI 配置列表操作
function addAiConfig() {
  config.aiConfigs.push({
    name: '',
    baseUrl: 'https://api.deepseek.com',
    apiKey: '',
    modelName: 'deepseek-chat',
    temperature: 0.1,
    maxTokens: 8192,
    timeOut: 6000,
    httpProxy: '',
    httpProxyEnabled: false,
    thinking: true,
    clientProfile: '',
    clientVersion: '',
  })
}

function removeAiConfig(index) {
  config.aiConfigs.splice(index, 1)
}

onBeforeMount(loadConfig)
</script>

<template>
  <div class="settings-page">
    <!-- 顶部导航（只保留标题） -->
    <PageHeader title="设置" />

    <!-- 分组切换 -->
    <div class="group-tabs">
      <button
        v-for="tab in groupTabs"
        :key="tab.value"
        class="group-btn"
        :class="{ 'group-btn--active': activeTab === tab.value }"
        @click="activeTab = tab.value"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 保存提示 -->
    <Transition name="tip">
      <div v-if="savedTip" class="save-tip">{{ savedTip }}</div>
    </Transition>

    <!-- 内容区 -->
    <div class="settings-content">
      <MEmpty v-if="loading" description="加载配置中..." />

      <!-- 基本设置 -->
      <MCard v-else-if="activeTab === 'general'">
        <div class="switch-item" @click="config.darkTheme = !config.darkTheme">
          <span class="item-label">暗黑主题</span>
          <span class="switch" :class="{ 'switch--on': config.darkTheme }" />
        </div>
        <div class="switch-item" @click="config.enableNews = !config.enableNews">
          <span class="item-label">滚动快讯</span>
          <span class="switch" :class="{ 'switch--on': config.enableNews }" />
        </div>
        <div class="switch-item" @click="config.enableDanmu = !config.enableDanmu">
          <span class="item-label">弹幕功能</span>
          <span class="switch" :class="{ 'switch--on': config.enableDanmu }" />
        </div>
        <div class="switch-item" @click="config.enableFund = !config.enableFund">
          <span class="item-label">基金功能</span>
          <span class="switch" :class="{ 'switch--on': config.enableFund }" />
        </div>
        <div class="switch-item" @click="config.enableAgent = !config.enableAgent">
          <span class="item-label">AI智能体</span>
          <span class="switch" :class="{ 'switch--on': config.enableAgent }" />
        </div>
        <div class="switch-item" @click="config.updateBasicInfoOnStart = !config.updateBasicInfoOnStart">
          <span class="item-label">启动时更新基础信息</span>
          <span class="switch" :class="{ 'switch--on': config.updateBasicInfoOnStart }" />
        </div>
        <div class="input-item">
          <label class="item-label">数据刷新间隔(秒)</label>
          <input v-model.number="config.refreshInterval" type="number" class="text-input" min="1">
        </div>
        <div class="input-item">
          <label class="item-label">服务端 Chromium 路径</label>
          <input v-model="config.browserPath" type="text" class="text-input" placeholder="留空自动检测">
        </div>
      </MCard>

      <!-- AI诊股设置 -->
      <template v-else-if="activeTab === 'ai'">
        <MCard>
          <div class="switch-item" @click="config.openAiEnable = !config.openAiEnable">
            <span class="item-label">AI诊股</span>
            <span class="switch" :class="{ 'switch--on': config.openAiEnable }" />
          </div>
        </MCard>

        <MCard v-if="config.openAiEnable">
          <h3 class="card-title">AI 配置</h3>
          <MEmpty v-if="!config.aiConfigs.length" description="还没有 AI 配置">
            <MButton type="primary" size="small" @click="addAiConfig">添加配置</MButton>
          </MEmpty>
          <div v-else class="ai-config-list">
            <div v-for="(ai, i) in config.aiConfigs" :key="i" class="ai-config-item">
              <div class="ai-config-head">
                <span class="ai-config-name">{{ ai.name || `配置 ${i + 1}` }}</span>
                <button class="del-btn" @click="removeAiConfig(i)">删除</button>
              </div>
              <div class="input-item">
                <label class="item-label-sm">配置名称</label>
                <input v-model="ai.name" type="text" class="text-input">
              </div>
              <div class="input-item">
                <label class="item-label-sm">接口地址</label>
                <input v-model="ai.baseUrl" type="text" class="text-input">
              </div>
              <div class="input-item">
                <label class="item-label-sm">API Key</label>
                <input v-model="ai.apiKey" type="password" class="text-input">
              </div>
              <div class="input-item">
                <label class="item-label-sm">模型名称</label>
                <input v-model="ai.modelName" type="text" class="text-input">
              </div>
              <div class="ai-row">
                <div class="input-item input-item--half">
                  <label class="item-label-sm">Temperature</label>
                  <input v-model.number="ai.temperature" type="number" step="0.1" class="text-input">
                </div>
                <div class="input-item input-item--half">
                  <label class="item-label-sm">MaxTokens</label>
                  <input v-model.number="ai.maxTokens" type="number" class="text-input">
                </div>
              </div>
            </div>
            <MButton type="text" size="small" block @click="addAiConfig">+ 添加配置</MButton>
          </div>
        </MCard>

        <MCard v-if="config.openAiEnable">
          <h3 class="card-title">诊股参数</h3>
          <div class="input-item">
            <label class="item-label">默认个股分析提示词</label>
            <textarea v-model="config.prompt" class="textarea-input" rows="3" />
          </div>
          <div class="input-item">
            <label class="item-label">分析提问模板</label>
            <input v-model="config.questionTemplate" type="text" class="text-input">
          </div>
          <div class="ai-row">
            <div class="input-item input-item--half">
              <label class="item-label-sm">爬虫超时(秒)</label>
              <input v-model.number="config.crawlTimeOut" type="number" class="text-input">
            </div>
            <div class="input-item input-item--half">
              <label class="item-label-sm">日K线天数</label>
              <input v-model.number="config.kDays" type="number" class="text-input">
            </div>
          </div>
        </MCard>
      </template>

      <!-- 通知设置 -->
      <MCard v-else-if="activeTab === 'push'">
        <div class="switch-item" @click="config.enablePushNews = !config.enablePushNews">
          <span class="item-label">市场资讯提醒</span>
          <span class="switch" :class="{ 'switch--on': config.enablePushNews }" />
        </div>
        <div v-if="config.enablePushNews" class="switch-item" @click="config.enableOnlyPushRedNews = !config.enableOnlyPushRedNews">
          <span class="item-label">仅提醒红字/关注个股</span>
          <span class="switch" :class="{ 'switch--on': config.enableOnlyPushRedNews }" />
        </div>
        <div class="divider" />
        <div class="switch-item" @click="config.dingPushEnable = !config.dingPushEnable">
          <span class="item-label">钉钉推送</span>
          <span class="switch" :class="{ 'switch--on': config.dingPushEnable }" />
        </div>
        <div v-if="config.dingPushEnable" class="input-item">
          <label class="item-label">钉钉机器人接口地址</label>
          <input v-model="config.dingRobot" type="text" class="text-input" placeholder="https://oapi.dingtalk.com/...">
        </div>
      </MCard>

      <!-- 数据源设置 -->
      <MCard v-else-if="activeTab === 'source'">
        <div class="switch-item" @click="config.httpProxyEnabled = !config.httpProxyEnabled">
          <span class="item-label">HTTP 代理</span>
          <span class="switch" :class="{ 'switch--on': config.httpProxyEnabled }" />
        </div>
        <div v-if="config.httpProxyEnabled" class="input-item">
          <label class="item-label">代理地址</label>
          <input v-model="config.httpProxy" type="text" class="text-input" placeholder="http://127.0.0.1:7890">
        </div>
        <div class="divider" />
        <div class="input-item">
          <label class="item-label">Tushare Token</label>
          <input v-model="config.tushareToken" type="password" class="text-input">
        </div>
        <div class="input-item">
          <label class="item-label">问财 API 密钥</label>
          <input v-model="config.iwencaiApiKey" type="password" class="text-input">
        </div>
        <div class="input-item">
          <label class="item-label">东财 AI 密钥</label>
          <input v-model="config.emApiKey" type="password" class="text-input">
        </div>
        <div class="input-item">
          <label class="item-label">东财唯一标识</label>
          <input v-model="config.qgqpBId" type="text" class="text-input">
        </div>
        <div class="input-item">
          <label class="item-label">提示词广场地址</label>
          <input v-model="config.promptPlazaApiBase" type="text" class="text-input">
        </div>
      </MCard>
    </div>

    <!-- 底部保存栏（替代顶部保存按钮） -->
    <div class="save-bar">
      <button class="save-btn" :disabled="saving" @click="handleSave">
        {{ saving ? '保存中...' : '保存设置' }}
      </button>
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

/* 底部保存栏 */
.save-bar {
  flex-shrink: 0;
  padding: var(--m-space-md);
  padding-bottom: calc(var(--m-space-md) + var(--m-safe-bottom));
  background: var(--m-bg-card);
  border-top: 1px solid var(--m-divider-color);
}

.save-btn {
  width: 100%;
  height: 44px;
  background: var(--m-color-rise);
  color: #fff;
  border: none;
  border-radius: var(--m-radius-md);
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  cursor: pointer;
}

.save-btn:disabled {
  opacity: 0.6;
}

.save-btn:active {
  opacity: 0.8;
}

.group-tabs {
  display: flex;
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.group-btn {
  flex: 1;
  padding: var(--m-space-md) 0;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: var(--m-font-md);
  color: var(--m-text-secondary);
  cursor: pointer;
}

.group-btn--active {
  color: var(--m-color-rise);
  border-bottom-color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.save-tip {
  position: fixed;
  top: 70px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
  padding: var(--m-space-sm) var(--m-space-lg);
  border-radius: var(--m-radius-md);
  font-size: var(--m-font-sm);
  z-index: var(--m-z-toast);
}

.tip-enter-active,
.tip-leave-active {
  transition: opacity var(--m-duration-fast);
}

.tip-enter-from,
.tip-leave-to {
  opacity: 0;
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.card-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  margin-bottom: var(--m-space-md);
  color: var(--m-text-primary);
}

.switch-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--m-space-md) 0;
  border-bottom: 1px solid var(--m-divider-color);
  cursor: pointer;
}

.switch-item:last-child {
  border-bottom: none;
}

.item-label {
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
}

.item-label-sm {
  display: block;
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
  margin-bottom: var(--m-space-xs);
}

.switch {
  width: 44px;
  height: 26px;
  background: var(--m-text-tertiary);
  border-radius: var(--m-radius-full);
  position: relative;
  transition: background var(--m-duration-normal);
  flex-shrink: 0;
}

.switch::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 22px;
  height: 22px;
  background: #fff;
  border-radius: 50%;
  transition: transform var(--m-duration-normal);
}

.switch--on {
  background: var(--m-color-rise);
}

.switch--on::after {
  transform: translateX(18px);
}

.input-item {
  padding: var(--m-space-md) 0;
  border-bottom: 1px solid var(--m-divider-color);
}

.input-item:last-child {
  border-bottom: none;
}

.text-input {
  width: 100%;
  height: 40px;
  padding: 0 var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-sm);
  outline: none;
}

.text-input:focus {
  border-color: var(--m-color-rise);
}

.textarea-input {
  width: 100%;
  padding: var(--m-space-sm) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-sm);
  outline: none;
  resize: vertical;
  font-family: inherit;
}

.textarea-input:focus {
  border-color: var(--m-color-rise);
}

.ai-row {
  display: flex;
  gap: var(--m-space-md);
}

.input-item--half {
  flex: 1;
  border-bottom: none;
  padding: 0;
}

.ai-config-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.ai-config-item {
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
}

.ai-config-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--m-space-md);
}

.ai-config-name {
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.del-btn {
  padding: var(--m-space-xs) var(--m-space-sm);
  background: transparent;
  border: 1px solid var(--m-color-fall);
  color: var(--m-color-fall);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
  cursor: pointer;
}

.del-btn:active {
  opacity: 0.6;
}

.ai-config-item .input-item {
  border-bottom: none;
  padding: var(--m-space-xs) 0;
}

.divider {
  height: 1px;
  background: var(--m-divider-color);
  margin: var(--m-space-md) 0;
}
</style>
