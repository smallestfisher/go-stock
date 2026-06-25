<script setup>
import { ref, computed, onMounted, watch } from 'vue'

const props = defineProps({
  // 分时数据
  data: {
    type: Array,
    default: () => []
    // [{ time, price, avgPrice, volume }]
  },
  // 昨收价
  preClose: {
    type: Number,
    required: true
  },
  // 宽度
  width: {
    type: Number,
    default: 375
  },
  // 高度
  height: {
    type: Number,
    default: 300
  }
})

const canvasRef = ref(null)

// 计算价格范围
const priceRange = computed(() => {
  if (!props.data.length) return { min: props.preClose, max: props.preClose }

  const prices = props.data.map(d => d.price)
  const min = Math.min(...prices, props.preClose)
  const max = Math.max(...prices, props.preClose)

  // 以昨收价为中心，计算对称范围
  const diff = Math.max(max - props.preClose, props.preClose - min)
  return {
    min: props.preClose - diff,
    max: props.preClose + diff
  }
})

// 绘制图表
function draw() {
  const canvas = canvasRef.value
  if (!canvas || !props.data.length) return

  const ctx = canvas.getContext('2d')
  const dpr = window.devicePixelRatio || 1

  // 设置canvas实际大小
  canvas.width = props.width * dpr
  canvas.height = props.height * dpr
  ctx.scale(dpr, dpr)

  // 清空画布
  ctx.clearRect(0, 0, props.width, props.height)

  // 绘制参数
  const padding = { top: 10, right: 10, bottom: 30, left: 10 }
  const chartWidth = props.width - padding.left - padding.right
  const chartHeight = props.height - padding.top - padding.bottom

  const { min, max } = priceRange.value
  const range = max - min || 1

  // 绘制昨收价参考线
  const preCloseY = padding.top + ((max - props.preClose) / range) * chartHeight
  ctx.strokeStyle = '#999'
  ctx.lineWidth = 1
  ctx.setLineDash([5, 5])
  ctx.beginPath()
  ctx.moveTo(padding.left, preCloseY)
  ctx.lineTo(props.width - padding.right, preCloseY)
  ctx.stroke()
  ctx.setLineDash([])

  // 计算坐标点
  const points = props.data.map((d, index) => {
    const x = padding.left + (index / (props.data.length - 1)) * chartWidth
    const y = padding.top + ((max - d.price) / range) * chartHeight
    return { x, y, price: d.price }
  })

  const avgPoints = props.data.map((d, index) => {
    const x = padding.left + (index / (props.data.length - 1)) * chartWidth
    const y = padding.top + ((max - d.avgPrice) / range) * chartHeight
    return { x, y }
  })

  // 绘制填充区域
  ctx.beginPath()
  ctx.moveTo(points[0].x, props.height - padding.bottom)
  points.forEach(point => ctx.lineTo(point.x, point.y))
  ctx.lineTo(points[points.length - 1].x, props.height - padding.bottom)
  ctx.closePath()

  const lastPrice = points[points.length - 1].price
  const color = lastPrice >= props.preClose ? '#18a058' : '#d03050'
  const gradient = ctx.createLinearGradient(0, padding.top, 0, props.height - padding.bottom)
  gradient.addColorStop(0, color + '30')
  gradient.addColorStop(1, color + '00')
  ctx.fillStyle = gradient
  ctx.fill()

  // 绘制价格线
  ctx.beginPath()
  ctx.moveTo(points[0].x, points[0].y)
  points.forEach(point => ctx.lineTo(point.x, point.y))
  ctx.strokeStyle = color
  ctx.lineWidth = 1.5
  ctx.stroke()

  // 绘制均价线
  ctx.beginPath()
  ctx.moveTo(avgPoints[0].x, avgPoints[0].y)
  avgPoints.forEach(point => ctx.lineTo(point.x, point.y))
  ctx.strokeStyle = '#f59e0b'
  ctx.lineWidth = 1
  ctx.stroke()
}

watch(() => props.data, draw, { deep: true })
watch(() => [props.width, props.height], draw)

onMounted(() => {
  draw()
})
</script>

<template>
  <div class="fenshi-chart">
    <canvas
      ref="canvasRef"
      :style="{ width: `${width}px`, height: `${height}px` }"
    />
    <div class="chart-placeholder" v-if="!data.length">
      <p>📈 分时图</p>
      <p class="placeholder-tip">暂无数据</p>
    </div>
  </div>
</template>

<style scoped>
.fenshi-chart {
  position: relative;
  width: 100%;
  height: 100%;
}

canvas {
  display: block;
  width: 100%;
  height: 100%;
}

.chart-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--m-text-secondary);
}

.placeholder-tip {
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
  margin-top: var(--m-space-xs);
}
</style>
