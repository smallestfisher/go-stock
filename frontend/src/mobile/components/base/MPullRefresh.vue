<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  // 下拉刷新回调
  onRefresh: {
    type: Function,
    required: true
  },
  // 是否禁用
  disabled: {
    type: Boolean,
    default: false
  },
  // 触发下拉刷新的距离（px）
  distance: {
    type: Number,
    default: 60
  }
})

const rootRef = ref(null)
const trackRef = ref(null)

const state = ref('normal') // 'normal' | 'pulling' | 'loosing' | 'refreshing' | 'success'
const pullDistance = ref(0)
const startY = ref(0)
const isTouching = ref(false)

const statusText = {
  normal: '下拉刷新',
  pulling: '下拉刷新',
  loosing: '释放刷新',
  refreshing: '刷新中...',
  success: '刷新成功'
}

function getScrollTop() {
  return trackRef.value?.scrollTop || 0
}

function handleTouchStart(e) {
  if (props.disabled || state.value === 'refreshing') return
  if (getScrollTop() > 0) return

  startY.value = e.touches[0].clientY
  isTouching.value = true
}

function handleTouchMove(e) {
  if (!isTouching.value || props.disabled) return

  const currentY = e.touches[0].clientY
  const distance = currentY - startY.value

  // 只在顶部且下拉时处理
  if (getScrollTop() === 0 && distance > 0) {
    e.preventDefault()

    // 阻尼效果：距离越大，移动越慢
    pullDistance.value = Math.pow(distance, 0.8)

    if (pullDistance.value >= props.distance) {
      state.value = 'loosing'
    } else {
      state.value = 'pulling'
    }
  }
}

async function handleTouchEnd() {
  if (!isTouching.value || props.disabled) return

  isTouching.value = false

  if (state.value === 'loosing') {
    // 触发刷新
    state.value = 'refreshing'
    pullDistance.value = props.distance

    try {
      await props.onRefresh()
      state.value = 'success'

      // 成功提示显示 500ms 后恢复
      setTimeout(() => {
        state.value = 'normal'
        pullDistance.value = 0
      }, 500)
    } catch (error) {
      console.error('Refresh failed:', error)
      state.value = 'normal'
      pullDistance.value = 0
    }
  } else {
    // 未达到刷新距离，回弹
    state.value = 'normal'
    pullDistance.value = 0
  }
}

onMounted(() => {
  const track = trackRef.value
  if (!track) return

  track.addEventListener('touchstart', handleTouchStart, { passive: false })
  track.addEventListener('touchmove', handleTouchMove, { passive: false })
  track.addEventListener('touchend', handleTouchEnd)
})

onBeforeUnmount(() => {
  const track = trackRef.value
  if (!track) return

  track.removeEventListener('touchstart', handleTouchStart)
  track.removeEventListener('touchmove', handleTouchMove)
  track.removeEventListener('touchend', handleTouchEnd)
})
</script>

<template>
  <div ref="rootRef" class="m-pull-refresh">
    <!-- 刷新指示器 -->
    <div
      class="m-pull-refresh__indicator"
      :style="{
        transform: `translateY(${Math.min(pullDistance, distance)}px)`,
        transition: isTouching ? 'none' : 'transform 0.3s'
      }"
    >
      <div class="m-pull-refresh__status">
        <!-- 加载动画 -->
        <div
          v-if="state === 'refreshing'"
          class="m-pull-refresh__spinner"
        />
        <!-- 箭头 -->
        <div
          v-else
          class="m-pull-refresh__arrow"
          :class="{ 'm-pull-refresh__arrow--up': state === 'loosing' }"
        >
          ↓
        </div>
        <span class="m-pull-refresh__text">{{ statusText[state] }}</span>
      </div>
    </div>

    <!-- 内容区 -->
    <div
      ref="trackRef"
      class="m-pull-refresh__track"
      :style="{
        transform: `translateY(${Math.min(pullDistance, distance)}px)`,
        transition: isTouching ? 'none' : 'transform 0.3s'
      }"
    >
      <slot />
    </div>
  </div>
</template>

<style scoped>
.m-pull-refresh {
  position: relative;
  overflow: hidden;
  height: 100%;
}

.m-pull-refresh__indicator {
  position: absolute;
  top: -60px;
  left: 0;
  right: 0;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.m-pull-refresh__status {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.m-pull-refresh__arrow {
  font-size: 20px;
  color: var(--m-color-rise);
  transition: transform var(--m-duration-normal);
}

.m-pull-refresh__arrow--up {
  transform: rotate(180deg);
}

.m-pull-refresh__spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-top-color: var(--m-color-rise);
  border-radius: 50%;
  animation: m-spin 0.8s linear infinite;
}

.m-pull-refresh__text {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.m-pull-refresh__track {
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
}

@keyframes m-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
