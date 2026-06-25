import { ref, onMounted, onBeforeUnmount } from 'vue'

/**
 * 左右滑动切换 Composable
 * @param {Object} options 配置选项
 * @param {Ref} options.containerRef 容器元素引用
 * @param {Function} options.onSwipeLeft 左滑回调
 * @param {Function} options.onSwipeRight 右滑回调
 * @param {Number} options.threshold 触发阈值（px）
 * @returns {Object} 滑动状态
 */
export function useSwipe(options = {}) {
  const {
    containerRef,
    onSwipeLeft,
    onSwipeRight,
    threshold = 50
  } = options

  const startX = ref(0)
  const startY = ref(0)
  const isSwiping = ref(false)

  function handleTouchStart(e) {
    startX.value = e.touches[0].clientX
    startY.value = e.touches[0].clientY
    isSwiping.value = true
  }

  function handleTouchMove(e) {
    if (!isSwiping.value) return

    const currentX = e.touches[0].clientX
    const currentY = e.touches[0].clientY
    const diffX = currentX - startX.value
    const diffY = currentY - startY.value

    // 如果垂直滑动距离大于水平滑动，不处理（让页面滚动）
    if (Math.abs(diffY) > Math.abs(diffX)) {
      isSwiping.value = false
      return
    }

    // 阻止默认行为（避免页面滚动）
    if (Math.abs(diffX) > 10) {
      e.preventDefault()
    }
  }

  function handleTouchEnd(e) {
    if (!isSwiping.value) return

    const endX = e.changedTouches[0].clientX
    const diffX = endX - startX.value

    // 判断是否达到阈值
    if (Math.abs(diffX) >= threshold) {
      if (diffX > 0 && onSwipeRight) {
        onSwipeRight()
      } else if (diffX < 0 && onSwipeLeft) {
        onSwipeLeft()
      }
    }

    isSwiping.value = false
  }

  function bindEvents() {
    const container = containerRef?.value
    if (!container) return

    container.addEventListener('touchstart', handleTouchStart, { passive: true })
    container.addEventListener('touchmove', handleTouchMove, { passive: false })
    container.addEventListener('touchend', handleTouchEnd, { passive: true })
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
    isSwiping,
    bindEvents,
    unbindEvents
  }
}
