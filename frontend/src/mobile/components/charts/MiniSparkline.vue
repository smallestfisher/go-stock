<script setup>
import { ref, computed, onMounted, watch } from 'vue'

const props = defineProps({
  // 数据点数组
  data: {
    type: Array,
    default: () => []
  },
  // 宽度
  width: {
    type: Number,
    default: 60
  },
  // 高度
  height: {
    type: Number,
    default: 24
  },
  // 线条颜色（根据涨跌自动计算）
  color: {
    type: String,
    default: ''
  },
  // 填充透明度
  fillOpacity: {
    type: Number,
    default: 0.1
  }
})

const canvasRef = ref(null)

// 计算涨跌状态
const trend = computed(() => {
  if (!props.data.length) return 'flat'
  const first = props.data[0]
  const last = props.data[props.data.length - 1]
  if (last > first) return 'rise'
  if (last < first) return 'fall'
  return 'flat'
})

// 线条颜色
const lineColor = computed(() => {
  if (props.color) return props.color
  const colorMap = {
    rise: '#18a058',
    fall: '#d03050',
    flat: '#666'
  }
  return colorMap[trend.value]
})

// 绘制走势线
function draw() {
  const canvas = canvasRef.value
  if (!canvas || !props.data.length) return

  const ctx = canvas.getContext('2d')
  const dpr = window.devicePixelRatio || 1

  // 设置canvas实际大小（考虑设备像素比）
  canvas.width = props.width * dpr
  canvas.height = props.height * dpr
  ctx.scale(dpr, dpr)

  // 清空画布
  ctx.clearRect(0, 0, props.width, props.height)

  // 计算数据范围
  const values = props.data
  const min = Math.min(...values)
  const max = Math.max(...values)
  const range = max - min || 1

  // 计算坐标点
  const points = values.map((value, index) => {
    const x = (index / (values.length - 1)) * props.width
    const y = props.height - ((value - min) / range) * props.height
    return { x, y }
  })

  // 绘制填充区域
  ctx.beginPath()
  ctx.moveTo(points[0].x, props.height)
  points.forEach(point => {
    ctx.lineTo(point.x, point.y)
  })
  ctx.lineTo(points[points.length - 1].x, props.height)
  ctx.closePath()
  ctx.fillStyle = lineColor.value + Math.round(props.fillOpacity * 255).toString(16).padStart(2, '0')
  ctx.fill()

  // 绘制线条
  ctx.beginPath()
  ctx.moveTo(points[0].x, points[0].y)
  points.forEach(point => {
    ctx.lineTo(point.x, point.y)
  })
  ctx.strokeStyle = lineColor.value
  ctx.lineWidth = 1.5
  ctx.lineJoin = 'round'
  ctx.lineCap = 'round'
  ctx.stroke()
}

// 监听数据变化重绘
watch(() => props.data, draw, { deep: true })
watch(() => [props.width, props.height], draw)

onMounted(() => {
  draw()
})
</script>

<template>
  <canvas
    ref="canvasRef"
    class="mini-sparkline"
    :style="{
      width: `${width}px`,
      height: `${height}px`
    }"
  />
</template>

<style scoped>
.mini-sparkline {
  display: block;
}
</style>
