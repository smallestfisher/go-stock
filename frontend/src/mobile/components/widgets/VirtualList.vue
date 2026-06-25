<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'

const props = defineProps({
  // 列表数据
  items: {
    type: Array,
    default: () => []
  },
  // 每项高度（px）
  itemHeight: {
    type: Number,
    default: 80
  },
  // 缓冲区数量（上下各渲染多少项）
  buffer: {
    type: Number,
    default: 5
  }
})

const containerRef = ref(null)
const scrollTop = ref(0)

// 可视区域高度
const viewportHeight = ref(0)

// 计算可视区域可显示的项数
const visibleCount = computed(() => {
  return Math.ceil(viewportHeight.value / props.itemHeight)
})

// 计算开始索引
const startIndex = computed(() => {
  const index = Math.floor(scrollTop.value / props.itemHeight) - props.buffer
  return Math.max(0, index)
})

// 计算结束索引
const endIndex = computed(() => {
  const index = startIndex.value + visibleCount.value + props.buffer * 2
  return Math.min(props.items.length, index)
})

// 可视区域的数据
const visibleItems = computed(() => {
  return props.items.slice(startIndex.value, endIndex.value).map((item, index) => ({
    data: item,
    index: startIndex.value + index
  }))
})

// 容器总高度
const totalHeight = computed(() => {
  return props.items.length * props.itemHeight
})

// 偏移量
const offsetY = computed(() => {
  return startIndex.value * props.itemHeight
})

// 滚动事件处理
function handleScroll(e) {
  scrollTop.value = e.target.scrollTop
}

// 更新视口高度
function updateViewportHeight() {
  if (containerRef.value) {
    viewportHeight.value = containerRef.value.clientHeight
  }
}

// 监听窗口大小变化
let resizeObserver = null

onMounted(() => {
  updateViewportHeight()

  // 使用 ResizeObserver 监听容器大小变化
  if (window.ResizeObserver) {
    resizeObserver = new ResizeObserver(() => {
      updateViewportHeight()
    })
    resizeObserver.observe(containerRef.value)
  }
})

onBeforeUnmount(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
})

// 监听数据变化，重置滚动位置
watch(() => props.items.length, () => {
  if (containerRef.value) {
    containerRef.value.scrollTop = 0
    scrollTop.value = 0
  }
})
</script>

<template>
  <div
    ref="containerRef"
    class="virtual-list"
    @scroll="handleScroll"
  >
    <!-- 占位容器（撑起总高度） -->
    <div class="virtual-list-phantom" :style="{ height: `${totalHeight}px` }">
      <!-- 可视区域内容 -->
      <div
        class="virtual-list-content"
        :style="{ transform: `translateY(${offsetY}px)` }"
      >
        <div
          v-for="item in visibleItems"
          :key="item.index"
          class="virtual-list-item"
          :style="{ height: `${itemHeight}px` }"
        >
          <slot :item="item.data" :index="item.index" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.virtual-list {
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
  position: relative;
}

.virtual-list-phantom {
  position: relative;
}

.virtual-list-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  will-change: transform;
}

.virtual-list-item {
  overflow: hidden;
}
</style>
