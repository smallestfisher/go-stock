<script setup>
defineProps({
  // 按钮类型
  type: {
    type: String,
    default: 'default', // 'default' | 'primary' | 'success' | 'danger' | 'text'
  },
  // 按钮尺寸
  size: {
    type: String,
    default: 'medium', // 'small' | 'medium' | 'large'
  },
  // 是否禁用
  disabled: {
    type: Boolean,
    default: false
  },
  // 是否块级（占满宽度）
  block: {
    type: Boolean,
    default: false
  },
  // 是否圆角
  round: {
    type: Boolean,
    default: false
  },
  // 是否加载中
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

function handleClick(e) {
  if (!props.disabled && !props.loading) {
    emit('click', e)
  }
}
</script>

<template>
  <button
    class="m-button"
    :class="{
      [`m-button--${type}`]: true,
      [`m-button--${size}`]: true,
      'm-button--block': block,
      'm-button--round': round,
      'm-button--disabled': disabled,
      'm-button--loading': loading,
    }"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <span v-if="loading" class="m-button__loading">
      <span class="m-button__spinner"></span>
    </span>
    <span class="m-button__content">
      <slot />
    </span>
  </button>
</template>

<style scoped>
.m-button {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: var(--m-touch-min);
  padding: 0 var(--m-space-lg);
  border: none;
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  line-height: 1;
  cursor: pointer;
  user-select: none;
  transition: all var(--m-duration-fast);
  -webkit-tap-highlight-color: transparent;
}

.m-button:active:not(.m-button--disabled):not(.m-button--loading) {
  transform: scale(0.96);
}

/* 类型变体 */
.m-button--default {
  background: var(--m-bg-card);
  color: var(--m-text-primary);
  border: 1px solid var(--m-border-color);
}

.m-button--default:active:not(.m-button--disabled) {
  background: var(--m-bg-primary);
}

.m-button--primary {
  background: var(--m-color-rise);
  color: white;
}

.m-button--primary:active:not(.m-button--disabled) {
  background: var(--m-color-rise-hover);
}

.m-button--success {
  background: var(--m-color-rise);
  color: white;
}

.m-button--danger {
  background: var(--m-color-fall);
  color: white;
}

.m-button--text {
  background: transparent;
  color: var(--m-color-rise);
  padding: 0 var(--m-space-md);
}

.m-button--text:active:not(.m-button--disabled) {
  background: var(--m-color-rise-light);
}

/* 尺寸变体 */
.m-button--small {
  min-height: 32px;
  padding: 0 var(--m-space-md);
  font-size: var(--m-font-sm);
}

.m-button--medium {
  min-height: var(--m-touch-min);
  padding: 0 var(--m-space-lg);
  font-size: var(--m-font-md);
}

.m-button--large {
  min-height: 52px;
  padding: 0 var(--m-space-xl);
  font-size: var(--m-font-lg);
}

/* 块级 */
.m-button--block {
  display: flex;
  width: 100%;
}

/* 圆角 */
.m-button--round {
  border-radius: var(--m-radius-full);
}

/* 禁用状态 */
.m-button--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 加载状态 */
.m-button--loading {
  cursor: default;
}

.m-button__loading {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.m-button__spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: m-spin 0.6s linear infinite;
}

.m-button--loading .m-button__content {
  opacity: 0;
}

@keyframes m-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
