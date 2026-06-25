<script setup>
defineProps({
  // 是否可点击
  clickable: {
    type: Boolean,
    default: false
  },
  // 内边距大小
  padding: {
    type: String,
    default: 'default', // 'none' | 'small' | 'default' | 'large'
  }
})

const emit = defineEmits(['click'])

function handleClick() {
  emit('click')
}
</script>

<template>
  <div
    class="m-card"
    :class="{
      'm-card--clickable': clickable,
      [`m-card--padding-${padding}`]: true
    }"
    @click="clickable && handleClick()"
  >
    <slot />
  </div>
</template>

<style scoped>
.m-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  box-shadow: var(--m-shadow-sm);
  transition: all var(--m-duration-fast);
}

.m-card--clickable {
  cursor: pointer;
}

.m-card--clickable:active {
  transform: scale(0.98);
  box-shadow: var(--m-shadow-md);
}

/* 内边距变体 */
.m-card--padding-none {
  padding: 0;
}

.m-card--padding-small {
  padding: var(--m-space-sm);
}

.m-card--padding-default {
  padding: var(--m-card-padding);
}

.m-card--padding-large {
  padding: var(--m-space-xl);
}
</style>
