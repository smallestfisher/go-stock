<script setup>
import MCard from '../base/MCard.vue'
import MButton from '../base/MButton.vue'

defineProps({
  // AI建议数据
  suggestion: {
    type: Object,
    default: () => null
    // { title, content, stocks, action }
  },
  // 是否加载中
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['refresh', 'actionClick', 'stockClick'])

function handleRefresh() {
  emit('refresh')
}

function handleAction() {
  emit('actionClick')
}

function handleStockClick(stock) {
  emit('stockClick', stock)
}
</script>

<template>
  <MCard>
    <!-- 头部 -->
    <div class="ai-suggest-header">
      <h3 class="header-title">🤖 AI 建议</h3>
      <button
        class="header-refresh"
        :disabled="loading"
        @click="handleRefresh"
      >
        🔄
      </button>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="ai-suggest-loading">
      <div class="loading-spinner" />
      <p>AI 正在分析...</p>
    </div>

    <!-- 建议内容 -->
    <div v-else-if="suggestion" class="ai-suggest-content">
      <h4 class="suggest-title">{{ suggestion.title }}</h4>
      <p class="suggest-content">{{ suggestion.content }}</p>

      <!-- 相关股票 -->
      <div v-if="suggestion.stocks && suggestion.stocks.length" class="suggest-stocks">
        <div class="stocks-label">相关股票:</div>
        <div class="stocks-list">
          <span
            v-for="stock in suggestion.stocks"
            :key="stock.code"
            class="stock-tag"
            @click="handleStockClick(stock)"
          >
            {{ stock.name }}
          </span>
        </div>
      </div>

      <!-- 操作按钮 -->
      <MButton
        v-if="suggestion.action"
        type="primary"
        size="small"
        block
        @click="handleAction"
      >
        {{ suggestion.action }}
      </MButton>
    </div>

    <!-- 空状态 -->
    <div v-else class="ai-suggest-empty">
      <p>暂无 AI 建议</p>
      <MButton type="text" size="small" @click="handleRefresh">
        点击生成
      </MButton>
    </div>
  </MCard>
</template>

<style scoped>
.ai-suggest-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--m-space-lg);
}

.header-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.header-refresh {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 18px;
  cursor: pointer;
  transition: transform var(--m-duration-normal);
}

.header-refresh:active:not(:disabled) {
  transform: rotate(180deg);
}

.header-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.ai-suggest-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-xl) 0;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--m-color-rise);
  border-radius: 50%;
  animation: m-spin 0.8s linear infinite;
}

.ai-suggest-loading p {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.ai-suggest-content {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.suggest-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  line-height: var(--m-line-height-normal);
}

.suggest-content {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-loose);
}

.suggest-stocks {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
}

.stocks-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-weight: var(--m-font-weight-medium);
}

.stocks-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-sm);
}

.stock-tag {
  padding: var(--m-space-xs) var(--m-space-md);
  background: var(--m-bg-card);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-sm);
  color: var(--m-text-primary);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.stock-tag:active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
}

.ai-suggest-empty {
  text-align: center;
  padding: var(--m-space-xl) 0;
  color: var(--m-text-secondary);
}

.ai-suggest-empty p {
  margin-bottom: var(--m-space-md);
}

@keyframes m-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
