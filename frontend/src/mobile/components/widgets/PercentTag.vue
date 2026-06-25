<script setup>
import { computed } from 'vue'

const props = defineProps({
  // 百分比值
  value: {
    type: [Number, String],
    required: true
  },
  // 是否显示符号（+/-）
  showSign: {
    type: Boolean,
    default: true
  },
  // 是否加粗
  bold: {
    type: Boolean,
    default: false
  },
  // 字体大小
  size: {
    type: String,
    default: 'medium' // 'small' | 'medium' | 'large'
  },
  // 是否显示百分号
  showPercent: {
    type: Boolean,
    default: true
  }
})

// 涨跌状态
const status = computed(() => {
  const num = Number(props.value)
  if (num > 0) return 'rise'
  if (num < 0) return 'fall'
  return 'flat'
})

// 格式化百分比
const formattedValue = computed(() => {
  const num = Number(props.value)
  if (isNaN(num)) return props.value

  // 保留2位小数
  const formatted = Math.abs(num).toFixed(2)

  // 添加符号
  let result = formatted
  if (props.showSign && num !== 0) {
    result = (num > 0 ? '+' : '-') + result
  }

  // 添加百分号
  if (props.showPercent) {
    result += '%'
  }

  return result
})
</script>

<template>
  <span
    class="percent-tag"
    :class="{
      [`percent-tag--${status}`]: true,
      [`percent-tag--${size}`]: true,
      'percent-tag--bold': bold
    }"
  >
    {{ formattedValue }}
  </span>
</template>

<style scoped>
.percent-tag {
  display: inline-block;
  font-variant-numeric: tabular-nums;
}

/* 涨跌色 */
.percent-tag--rise {
  color: var(--m-color-rise);
}

.percent-tag--fall {
  color: var(--m-color-fall);
}

.percent-tag--flat {
  color: var(--m-color-gray);
}

/* 尺寸 */
.percent-tag--small {
  font-size: var(--m-font-sm);
}

.percent-tag--medium {
  font-size: var(--m-font-md);
}

.percent-tag--large {
  font-size: var(--m-font-lg);
}

/* 加粗 */
.percent-tag--bold {
  font-weight: var(--m-font-weight-bold);
}
</style>
