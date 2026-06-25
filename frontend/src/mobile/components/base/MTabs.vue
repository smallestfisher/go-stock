<script setup>
import { ref, watch, nextTick } from 'vue'

const props = defineProps({
  // 当前激活的标签
  modelValue: {
    type: [String, Number],
    required: true
  },
  // 标签列表
  tabs: {
    type: Array,
    required: true
    // 格式: [{ label: '标签1', value: '1' }, ...]
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const tabsRef = ref(null)
const activeTabRef = ref(null)

function handleTabClick(tab) {
  if (tab.disabled) return
  emit('update:modelValue', tab.value)
  emit('change', tab.value)
}

// 滚动到激活的标签
function scrollToActiveTab() {
  nextTick(() => {
    if (!tabsRef.value || !activeTabRef.value) return

    const container = tabsRef.value
    const activeTab = activeTabRef.value

    const containerLeft = container.scrollLeft
    const containerWidth = container.offsetWidth
    const tabLeft = activeTab.offsetLeft
    const tabWidth = activeTab.offsetWidth

    // 计算目标滚动位置（居中显示）
    const targetScroll = tabLeft - (containerWidth / 2) + (tabWidth / 2)

    container.scrollTo({
      left: targetScroll,
      behavior: 'smooth'
    })
  })
}

watch(() => props.modelValue, () => {
  scrollToActiveTab()
}, { immediate: true })
</script>

<template>
  <div class="m-tabs">
    <div ref="tabsRef" class="m-tabs__wrapper">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        :ref="el => { if (tab.value === modelValue) activeTabRef = el }"
        class="m-tabs__item"
        :class="{
          'm-tabs__item--active': tab.value === modelValue,
          'm-tabs__item--disabled': tab.disabled
        }"
        @click="handleTabClick(tab)"
      >
        {{ tab.label }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.m-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.m-tabs__wrapper {
  display: flex;
  overflow-x: auto;
  overflow-y: hidden;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none; /* Firefox */
}

.m-tabs__wrapper::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

.m-tabs__item {
  position: relative;
  flex-shrink: 0;
  min-height: var(--m-touch-min);
  padding: 0 var(--m-space-lg);
  background: transparent;
  border: none;
  font-size: var(--m-font-md);
  color: var(--m-text-secondary);
  white-space: nowrap;
  cursor: pointer;
  transition: color var(--m-duration-fast);
}

.m-tabs__item::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 2px;
  background: var(--m-color-rise);
  transform: scaleX(0);
  transition: transform var(--m-duration-normal) var(--m-ease-out);
}

.m-tabs__item--active {
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.m-tabs__item--active::after {
  transform: scaleX(1);
}

.m-tabs__item--disabled {
  color: var(--m-text-disabled);
  cursor: not-allowed;
  opacity: 0.5;
}

.m-tabs__item:active:not(.m-tabs__item--disabled) {
  opacity: 0.7;
}
</style>
