<script setup>
import { ref, computed, nextTick, onBeforeMount, onBeforeUnmount } from 'vue'
import 'md-editor-v3/lib/preview.css'
import { MdPreview } from 'md-editor-v3'
import { ChatWithAgent, AbortChatWithAgent, GetAiConfigs, GetPromptTemplates } from '../../../api/app'
import { EventsOn, EventsOff } from '../../../api/runtime'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MSheet from '../../components/base/MSheet.vue'
import { toast } from '../../composables/useToast'

// 对齐桌面端 agent-chat.vue 的 ChatWithAgent 流式对话。
// 事件 agent-message 推 {role, content, reasoning_content}，结束标记 content === 'agent-DONE'。
// ChatWithAgent(question, aiConfigId, sysPromptId, memoryMode, memoryCount, thinkingMode, agentMode)

const messages = ref([])        // { role, content, reasoning }
const inputValue = ref('')
const messageListRef = ref(null)
const sending = ref(false)

// 模型 / Agent 模式 / 系统提示词 / 思考开关
const aiConfigs = ref([])
const aiConfigId = ref(null)
const sysPromptOptions = ref([])
const sysPromptId = ref(null)
const thinkingMode = ref(false)
const agentMode = ref('')        // '' 自动 / react 快速 / plan_execute 规划

const AGENT_MODES = [
  { label: '🤖 自动', value: '' },
  { label: '⚡ 快速', value: 'react' },
  { label: '🧠 规划', value: 'plan_execute' },
]

const settingsVisible = ref(false)

const activeModelName = computed(() => {
  const cfg = aiConfigs.value.find(c => c.ID === aiConfigId.value)
  return cfg ? cfg.name : '未选择模型'
})

async function loadAiConfigs() {
  try {
    const res = await GetAiConfigs()
    aiConfigs.value = Array.isArray(res) ? res : []
    if (aiConfigs.value.length && !aiConfigId.value) {
      aiConfigId.value = aiConfigs.value[0].ID
    }
  } catch (e) {
    console.error('加载AI配置失败:', e)
  }
}

async function loadSysPrompts() {
  try {
    const res = await GetPromptTemplates('', '模型系统Prompt')
    sysPromptOptions.value = Array.isArray(res) ? res : []
  } catch (e) {
    console.error('加载系统提示词失败:', e)
  }
}

function scrollToEnd() {
  nextTick(() => {
    requestAnimationFrame(() => {
      const el = messageListRef.value
      if (el) el.scrollTop = el.scrollHeight
    })
  })
}

function sendMessage() {
  const text = inputValue.value.trim()
  if (!text || sending.value) return
  if (!aiConfigId.value) {
    toast.warning('请先选择 AI 模型')
    settingsVisible.value = true
    return
  }

  messages.value.push({ role: 'user', content: text })
  // 占位的助手消息，流式往里追加
  messages.value.push({ role: 'assistant', content: '', reasoning: '' })
  inputValue.value = ''
  sending.value = true
  scrollToEnd()

  // ChatWithAgent(question, aiConfigId, sysPromptId, memoryMode, memoryCount, thinkingMode, agentMode)
  ChatWithAgent(text, aiConfigId.value, sysPromptId.value, false, 0, thinkingMode.value, agentMode.value)
}

function stopChat() {
  AbortChatWithAgent()
  sending.value = false
}

// 流式事件处理（对齐桌面 handleAgentMessage）
function handleAgentMessage(data) {
  if (!data) return
  const last = messages.value[messages.value.length - 1]
  if (!last || last.role !== 'assistant') return

  // 结束标记
  if (data.content === 'agent-DONE') {
    sending.value = false
    scrollToEnd()
    return
  }

  if (data.reasoning_content) {
    last.reasoning = (last.reasoning || '') + data.reasoning_content
  }
  if (data.content) {
    last.content = (last.content || '') + data.content
  }
  scrollToEnd()
}

function clearChat() {
  if (sending.value) {
    toast.warning('请先停止当前对话')
    return
  }
  messages.value = []
}

onBeforeMount(() => {
  loadAiConfigs()
  loadSysPrompts()
  EventsOn('agent-message', handleAgentMessage)
})

onBeforeUnmount(() => {
  EventsOff('agent-message')
  if (sending.value) AbortChatWithAgent()
})
</script>

<template>
  <div class="agent-page">
    <PageHeader title="AI智能体">
      <template #actions>
        <button type="button" class="header-act" @click="settingsVisible = true">⚙️</button>
        <button type="button" class="header-act" @click="clearChat">🗑️</button>
      </template>
    </PageHeader>

    <!-- 当前配置条 -->
    <div class="config-bar" @click="settingsVisible = true">
      <span class="config-model">{{ activeModelName }}</span>
      <span class="config-sep">·</span>
      <span class="config-mode">{{ AGENT_MODES.find(m => m.value === agentMode)?.label }}</span>
      <span v-if="thinkingMode" class="config-think">· 思考</span>
      <span class="config-edit">设置</span>
    </div>

    <div ref="messageListRef" class="agent-content">
      <div v-if="messages.length" class="message-list">
        <div
          v-for="(message, index) in messages"
          :key="index"
          class="message-row"
          :class="`message-row--${message.role}`"
        >
          <div class="message-bubble">
            <!-- 用户消息纯文本 -->
            <template v-if="message.role === 'user'">{{ message.content }}</template>
            <!-- 助手消息：思考过程 + Markdown 正文 -->
            <template v-else>
              <details v-if="message.reasoning" class="reasoning">
                <summary>💭 思考过程</summary>
                <div class="reasoning-body">{{ message.reasoning }}</div>
              </details>
              <MdPreview v-if="message.content" :modelValue="message.content" theme="light" />
              <span v-else-if="sending && index === messages.length - 1" class="typing">
                <span class="dot" /><span class="dot" /><span class="dot" />
              </span>
            </template>
          </div>
        </div>
      </div>
      <MEmpty v-else description="输入你的问题，开始与 AI 智能体对话" />
    </div>

    <div class="input-bar">
      <textarea
        v-model="inputValue"
        class="chat-input"
        rows="1"
        placeholder="输入你的问题..."
        @keyup.enter.exact.prevent="sendMessage"
      />
      <button
        v-if="!sending"
        type="button"
        class="send-btn"
        :disabled="!inputValue.trim()"
        @click="sendMessage"
      >发送</button>
      <button v-else type="button" class="send-btn send-btn--stop" @click="stopChat">停止</button>
    </div>

    <!-- 设置抽屉 -->
    <MSheet v-model:show="settingsVisible" title="对话设置" height="auto">
      <div class="settings">
        <div class="setting-block">
          <span class="setting-label">AI 模型</span>
          <div class="chip-wrap">
            <button
              v-for="cfg in aiConfigs"
              :key="cfg.ID"
              type="button"
              class="chip"
              :class="{ 'chip--active': aiConfigId === cfg.ID }"
              @click="aiConfigId = cfg.ID"
            >{{ cfg.name }}</button>
            <span v-if="!aiConfigs.length" class="setting-empty">未配置 AI 模型，请到设置页添加</span>
          </div>
        </div>

        <div class="setting-block">
          <span class="setting-label">Agent 模式</span>
          <div class="chip-wrap">
            <button
              v-for="m in AGENT_MODES"
              :key="m.value"
              type="button"
              class="chip"
              :class="{ 'chip--active': agentMode === m.value }"
              @click="agentMode = m.value"
            >{{ m.label }}</button>
          </div>
        </div>

        <div v-if="sysPromptOptions.length" class="setting-block">
          <span class="setting-label">系统提示词</span>
          <div class="chip-wrap">
            <button
              type="button"
              class="chip"
              :class="{ 'chip--active': sysPromptId === null }"
              @click="sysPromptId = null"
            >默认</button>
            <button
              v-for="t in sysPromptOptions"
              :key="t.ID"
              type="button"
              class="chip"
              :class="{ 'chip--active': sysPromptId === t.ID }"
              @click="sysPromptId = t.ID"
            >{{ t.name }}</button>
          </div>
        </div>

        <button
          type="button"
          class="switch-line"
          @click="thinkingMode = !thinkingMode"
        >
          <span>思考模式</span>
          <span class="switch" :class="{ 'switch--on': thinkingMode }"><span class="switch-dot" /></span>
        </button>
      </div>
    </MSheet>
  </div>
</template>

<style scoped>
.agent-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.header-act {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  background: transparent;
  border: none;
  font-size: var(--m-font-lg);
}

/* 配置条 */
.config-bar {
  display: flex;
  align-items: center;
  gap: var(--m-space-xs);
  padding: var(--m-space-xs) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.config-model {
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.config-sep,
.config-think {
  color: var(--m-text-tertiary);
}

.config-edit {
  margin-left: auto;
  color: var(--m-color-rise);
}

.agent-content {
  flex: 1;
  min-height: 0;
  padding: var(--m-space-md);
  overflow-y: auto;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.message-row {
  display: flex;
}

.message-row--assistant {
  justify-content: flex-start;
}

.message-row--user {
  justify-content: flex-end;
}

.message-bubble {
  max-width: 88%;
  padding: var(--m-space-sm) var(--m-space-md);
  border-radius: var(--m-radius-md);
  font-size: var(--m-font-md);
  line-height: var(--m-line-height-normal);
  word-break: break-word;
}

.message-row--assistant .message-bubble {
  background: var(--m-bg-card);
  color: var(--m-text-primary);
  box-shadow: var(--m-shadow-sm);
}

.message-row--user .message-bubble {
  background: var(--m-color-rise);
  color: #fff;
}

/* 思考过程 */
.reasoning {
  margin-bottom: var(--m-space-sm);
  font-size: var(--m-font-sm);
}

.reasoning summary {
  color: var(--m-text-tertiary);
  cursor: pointer;
}

.reasoning-body {
  margin-top: var(--m-space-xs);
  padding: var(--m-space-sm);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  white-space: pre-wrap;
}

/* 打字指示 */
.typing {
  display: inline-flex;
  gap: 4px;
  align-items: center;
  height: 20px;
}

.typing .dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--m-text-tertiary);
  animation: m-typing 1s infinite ease-in-out;
}

.typing .dot:nth-child(2) { animation-delay: 0.15s; }
.typing .dot:nth-child(3) { animation-delay: 0.3s; }

@keyframes m-typing {
  0%, 60%, 100% { opacity: 0.3; transform: translateY(0); }
  30% { opacity: 1; transform: translateY(-3px); }
}

/* 输入栏 */
.input-bar {
  display: flex;
  gap: var(--m-space-sm);
  align-items: flex-end;
  padding: var(--m-space-sm) var(--m-space-md);
  padding-bottom: calc(var(--m-space-sm) + var(--m-safe-bottom));
  background: var(--m-bg-card);
  border-top: 1px solid var(--m-divider-color);
}

.chat-input {
  flex: 1;
  min-width: 0;
  max-height: 100px;
  padding: var(--m-space-sm) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-md);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font: inherit;
  font-size: var(--m-font-md);
  line-height: var(--m-line-height-normal);
  resize: none;
  outline: none;
}

.chat-input:focus {
  border-color: var(--m-color-rise);
}

.send-btn {
  flex-shrink: 0;
  height: 40px;
  padding: 0 var(--m-space-lg);
  border: none;
  border-radius: var(--m-radius-md);
  background: var(--m-color-rise);
  color: #fff;
  font: inherit;
  font-size: var(--m-font-md);
}

.send-btn:disabled {
  opacity: 0.5;
}

.send-btn--stop {
  background: var(--m-color-fall);
}

/* 设置抽屉 */
.settings {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-lg);
}

.setting-block {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.setting-label {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.setting-empty {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

.chip-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-sm);
}

.chip {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  font: inherit;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.chip--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.switch-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: var(--m-space-sm) 0;
  background: transparent;
  border: none;
  font: inherit;
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
}

.switch {
  position: relative;
  width: 44px;
  height: 24px;
  border-radius: var(--m-radius-full);
  background: var(--m-border-color);
  transition: background var(--m-duration-fast);
  flex-shrink: 0;
}

.switch--on {
  background: var(--m-color-rise);
}

.switch-dot {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
  transition: transform var(--m-duration-fast);
}

.switch--on .switch-dot {
  transform: translateX(20px);
}
</style>
