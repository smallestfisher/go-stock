<script setup>
import { watch, nextTick } from 'vue'

const props = defineProps({
  // 是否显示
  show: {
    type: Boolean,
    default: false
  },
  // 标题
  title: {
    type: String,
    default: ''
  },
  // 高度
  height: {
    type: String,
    default: '90vh'
  },
  // 是否显示关闭按钮
  closable: {
    type: Boolean,
    default: true
  },
  // 是否点击遮罩关闭
  closeOnMask: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:show', 'close'])

function close() {
  emit('update:show', false)
  emit('close')
}

function handleMaskClick() {
  if (props.closeOnMask) {
    close()
  }
}

// 阻止滚动穿透
watch(() => props.show, (newVal) => {
  nextTick(() => {
    if (newVal) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  })
})
</script>

<template>
  <!-- 遮罩层 -->
  <Transition name="m-mask">
    <div v-if="show" class="m-sheet-mask" @click="handleMaskClick" />
  </Transition>

  <!-- 抽屉内容 -->
  <Transition name="m-slide-up">
    <div v-if="show" class="m-sheet" :style="{ height }">
      <!-- 头部 -->
      <div v-if="title || closable" class="m-sheet__header">
        <div class="m-sheet__title">{{ title }}</div>
        <button v-if="closable" class="m-sheet__close" @click="close">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
            <path d="M10 8.586L6.707 5.293a1 1 0 00-1.414 1.414L8.586 10l-3.293 3.293a1 1 0 101.414 1.414L10 11.414l3.293 3.293a1 1 0 001.414-1.414L11.414 10l3.293-3.293a1 1 0 00-1.414-1.414L10 8.586z" />
          </svg>
        </button>
      </div>

      <!-- 内容区 -->
      <div class="m-sheet__body">
        <slot />
      </div>

      <!-- 底部插槽 -->
      <div v-if="$slots.footer" class="m-sheet__footer">
        <slot name="footer" />
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.m-sheet-mask {
  position: fixed;
  inset: 0;
  background: var(--m-bg-overlay);
  z-index: var(--m-z-mask);
}

.m-sheet {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  max-height: 95vh;
  background: var(--m-bg-card);
  border-radius: var(--m-radius-lg) var(--m-radius-lg) 0 0;
  z-index: calc(var(--m-z-mask) + 1);
  display: flex;
  flex-direction: column;
  box-shadow: 0 -4px 16px rgba(0, 0, 0, 0.12);
}

.m-sheet__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-lg) var(--m-space-lg);
  border-bottom: 1px solid var(--m-divider-color);
  flex-shrink: 0;
}

.m-sheet__title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
}

.m-sheet__close {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: var(--m-touch-min);
  min-height: var(--m-touch-min);
  margin: calc(var(--m-space-lg) * -1);
  background: transparent;
  border: none;
  color: var(--m-text-secondary);
  cursor: pointer;
}

.m-sheet__close:active {
  opacity: 0.6;
}

.m-sheet__body {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
  padding: var(--m-space-lg);
}

.m-sheet__footer {
  flex-shrink: 0;
  padding: var(--m-space-lg);
  border-top: 1px solid var(--m-divider-color);
  background: var(--m-bg-card);
}
</style>
