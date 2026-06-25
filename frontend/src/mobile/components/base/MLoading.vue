<script setup>
defineProps({
  // 加载文本
  text: {
    type: String,
    default: '加载中...'
  },
  // 尺寸
  size: {
    type: String,
    default: 'medium' // 'small' | 'medium' | 'large'
  },
  // 是否垂直布局
  vertical: {
    type: Boolean,
    default: false
  },
  // 颜色
  color: {
    type: String,
    default: '' // 默认使用主题色
  }
})
</script>

<template>
  <div
    class="m-loading"
    :class="{
      [`m-loading--${size}`]: true,
      'm-loading--vertical': vertical
    }"
  >
    <div
      class="m-loading__spinner"
      :style="{ borderTopColor: color || 'var(--m-color-rise)' }"
    />
    <span v-if="text" class="m-loading__text">{{ text }}</span>
  </div>
</template>

<style scoped>
.m-loading {
  display: inline-flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.m-loading--vertical {
  flex-direction: column;
  gap: var(--m-space-md);
}

.m-loading__spinner {
  display: inline-block;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--m-color-rise);
  border-radius: 50%;
  animation: m-spin 0.8s linear infinite;
}

/* 尺寸变体 */
.m-loading--small .m-loading__spinner {
  width: 16px;
  height: 16px;
  border-width: 2px;
}

.m-loading--small .m-loading__text {
  font-size: var(--m-font-sm);
}

.m-loading--medium .m-loading__spinner {
  width: 24px;
  height: 24px;
  border-width: 3px;
}

.m-loading--medium .m-loading__text {
  font-size: var(--m-font-md);
}

.m-loading--large .m-loading__spinner {
  width: 32px;
  height: 32px;
  border-width: 4px;
}

.m-loading--large .m-loading__text {
  font-size: var(--m-font-lg);
}

.m-loading__text {
  color: var(--m-text-secondary);
}

@keyframes m-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
