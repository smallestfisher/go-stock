<script setup>
import { ref, nextTick } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'

// AI 智能体对话（对齐桌面端 agent-chat.vue 的 ChatWithAgent）
// 消息列表后续接 ChatWithAgent 流式填充，当前先留空
const messages = ref([])
const inputValue = ref('')
const messageListRef = ref(null)
const sending = ref(false)

async function scrollToEnd() {
  await nextTick()
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}

async function sendMessage() {
  const text = inputValue.value.trim()
  if (!text || sending.value) return

  messages.value.push({ role: 'user', content: text })
  inputValue.value = ''
  scrollToEnd()

  // TODO: 接 ChatWithAgent 流式接口，收到内容 push 到 messages
  sending.value = true
  // 占位：接入前暂不回复，避免假数据
  sending.value = false
}
</script>

<template>
  <div class="agent-page">
    <PageHeader title="AI智能体" />

    <div ref="messageListRef" class="agent-content">
      <MCard padding="none" class="chat-card">
        <div v-if="messages.length" class="message-list">
          <div
            v-for="(message, index) in messages"
            :key="index"
            class="message-row"
            :class="`message-row--${message.role}`"
          >
            <div class="message-bubble">
              {{ message.content }}
            </div>
          </div>
        </div>
        <MEmpty v-else description="输入你的问题，开始与 AI 智能体对话" />
      </MCard>
    </div>

    <div class="input-bar">
      <input
        v-model="inputValue"
        class="chat-input"
        type="text"
        placeholder="输入你的问题"
        :disabled="sending"
        @keyup.enter="sendMessage"
      >
      <MButton type="primary" size="small" :disabled="sending" @click="sendMessage">
        发送
      </MButton>
    </div>
  </div>
</template>

<style scoped>
.agent-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.agent-content {
  flex: 1;
  min-height: 0;
  padding: var(--m-space-md);
  overflow-y: auto;
}

.chat-card {
  min-height: 100%;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
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
  max-width: 82%;
  padding: var(--m-space-sm) var(--m-space-md);
  border-radius: var(--m-radius-md);
  font-size: var(--m-font-md);
  line-height: var(--m-line-height-normal);
  word-break: break-word;
}

.message-row--assistant .message-bubble {
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
}

.message-row--user .message-bubble {
  background: var(--m-color-rise);
  color: #fff;
}

.input-bar {
  display: flex;
  gap: var(--m-space-sm);
  align-items: center;
  padding: var(--m-space-sm) var(--m-space-md);
  padding-bottom: calc(var(--m-space-sm) + var(--m-safe-bottom));
  background: var(--m-bg-card);
  border-top: 1px solid var(--m-divider-color);
}

.chat-input {
  flex: 1;
  min-width: 0;
  height: 40px;
  padding: 0 var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-md);
  outline: none;
}

.chat-input:focus {
  border-color: var(--m-color-rise);
}

.chat-input:disabled {
  opacity: 0.6;
}
</style>
