<script setup>
import { computed } from 'vue'

const props = defineProps({
  // 价格
  price: {
    type: [Number, String],
    required: true
  },
  // 变化值（用于判断涨跌）
  change: {
    type: [Number, String],
    default: 0
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
  // 前缀符号
  prefix: {
    type: String,
    default: ''
  }
})

// 涨跌状态
const status = computed(() => {
  const changeNum = Number(props.change)
  if (changeNum > 0) return 'rise'
  if (changeNum < 0) return 'fall'
  return 'flat'
})

// 格式化价格
const formattedPrice = computed(() => {
  const num = Number(props.price)
  if (isNaN(num)) return props.price

  // 保留2位小数
  return num.toFixed(2)
})
</script>

<template>
  <span
    class="price-tag"
    :class="{
      [`price-tag--${status}`]: true,
      [`price-tag--${size}`]: true,
      'price-tag--bold': bold
    }"
  >
    <span v-if="prefix" class="price-tag__prefix">{{ prefix }}</span>
    <span class="price-tag__value">{{ formattedPrice }}</span>
  </span>
</template>

<style scoped>
.price-tag {
  display: inline-flex;
  align-items: baseline;
  font-variant-numeric: tabular-nums;
}

/* 涨跌色 */
.price-tag--rise {
  color: var(--m-color-rise);
}

.price-tag--fall {
  color: var(--m-color-fall);
}

.price-tag--flat {
  color: var(--m-color-gray);
}

/* 尺寸 */
.price-tag--small {
  font-size: var(--m-font-sm);
}

.price-tag--medium {
  font-size: var(--m-font-md);
}

.price-tag--large {
  font-size: var(--m-font-lg);
}

/* 加粗 */
.price-tag--bold {
  font-weight: var(--m-font-weight-bold);
}

.price-tag__prefix {
  margin-right: 2px;
  opacity: 0.8;
}
</style>
