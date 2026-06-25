import { ref, onMounted, onBeforeUnmount } from 'vue'

/**
 * 下拉刷新 Composable
 * @param {Object} options 配置选项
 * @param {Function} options.onRefresh 刷新回调函数
 * @param {Ref} options.containerRef 容器元素引用
 * @param {Boolean} options.disabled 是否禁用
 * @param {Number} options.distance 触发刷新的距离（px）
 * @returns {Object} 下拉刷新状态和方法
 */
export function usePullRefresh(options = {}) {
  const {
    onRefresh,
    containerRef,
    disabled = false,
    distance = 60
  } = options

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
    return containerRef?.value?.scrollTop || 0
  }

  function handleTouchStart(e) {
    if (disabled || state.value === 'refreshing') return
    if (getScrollTop() > 0) return

    startY.value = e.touches[0].clientY
    isTouching.value = true
  }

  function handleTouchMove(e) {
    if (!isTouching.value || disabled) return

    const currentY = e.touches[0].clientY
    const moveDistance = currentY - startY.value

    // 只在顶部且下拉时处理
    if (getScrollTop() === 0 && moveDistance > 0) {
      e.preventDefault()

      // 阻尼效果：距离越大，移动越慢
      pullDistance.value = Math.pow(moveDistance, 0.8)

      if (pullDistance.value >= distance) {
        state.value = 'loosing'
      } else {
        state.value = 'pulling'
      }
    }
  }

  async function handleTouchEnd() {
    if (!isTouching.value || disabled) return

    isTouching.value = false

    if (state.value === 'loosing') {
      // 触发刷新
      state.value = 'refreshing'
      pullDistance.value = distance

      try {
        if (onRefresh && typeof onRefresh === 'function') {
          await onRefresh()
        }
        state.value = 'success'

        // 成功提示显示 500ms 后恢复
        setTimeout(() => {
          reset()
        }, 500)
      } catch (error) {
        console.error('Refresh failed:', error)
        reset()
      }
    } else {
      // 未达到刷新距离，回弹
      reset()
    }
  }

  function reset() {
    state.value = 'normal'
    pullDistance.value = 0
  }

  function bindEvents() {
    const container = containerRef?.value
    if (!container) return

    container.addEventListener('touchstart', handleTouchStart, { passive: false })
    container.addEventListener('touchmove', handleTouchMove, { passive: false })
    container.addEventListener('touchend', handleTouchEnd)
  }

  function unbindEvents() {
    const container = containerRef?.value
    if (!container) return

    container.removeEventListener('touchstart', handleTouchStart)
    container.removeEventListener('touchmove', handleTouchMove)
    container.removeEventListener('touchend', handleTouchEnd)
  }

  onMounted(() => {
    bindEvents()
  })

  onBeforeUnmount(() => {
    unbindEvents()
  })

  return {
    state,
    pullDistance,
    isTouching,
    statusText,
    reset,
    bindEvents,
    unbindEvents
  }
}
