<script setup>
import { ref, computed, onMounted, watch } from 'vue'

const props = defineProps({
  // K线数据
  data: {
    type: Array,
    default: () => []
    // [{ time, open, close, high, low, volume }]
  },
  // 宽度
  width: {
    type: Number,
    default: 375
  },
  // 高度
  height: {
    type: Number,
    default: 400
  },
  // 显示的MA线
  maLines: {
    type: Array,
    default: () => [5, 10, 20, 30]
  }
})

const canvasRef = ref(null)

// 计算MA均线
function calculateMA(data, period) {
  const result = []
  for (let i = 0; i < data.length; i++) {
    if (i < period - 1) {
      result.push(null)
      continue
    }
    let sum = 0
    for (let j = 0; j < period; j++) {
      sum += data[i - j].close
    }
    result.push(sum / period)
  }
  return result
}

// 计算价格范围
const priceRange = computed(() => {
  if (!props.data.length) return { min: 0, max: 0 }

  const highs = props.data.map(d => d.high)
  const lows = props.data.map(d => d.low)
  return {
    min: Math.min(...lows),
    max: Math.max(...highs)
  }
})

// 绘制图表
function draw() {
  const canvas = canvasRef.value
  if (!canvas || !props.data.length) return

  const ctx = canvas.getContext('2d')
  const dpr = window.devicePixelRatio || 1

  canvas.width = props.width * dpr
  canvas.height = props.height * dpr
  ctx.scale(dpr, dpr)

  ctx.clearRect(0, 0, props.width, props.height)

  const padding = { top: 20, right: 10, bottom: 40, left: 10 }
  const chartWidth = props.width - padding.left - padding.right
  const chartHeight = props.height - padding.top - padding.bottom

  const { min, max } = priceRange.value
  const range = max - min || 1

  // 每根K线的宽度
  const candleWidth = Math.max(2, chartWidth / props.data.length - 1)

  // 绘制K线
  props.data.forEach((d, index) => {
    const x = padding.left + (index / props.data.length) * chartWidth
    const openY = padding.top + ((max - d.open) / range) * chartHeight
    const closeY = padding.top + ((max - d.close) / range) * chartHeight
    const highY = padding.top + ((max - d.high) / range) * chartHeight
    const lowY = padding.top + ((max - d.low) / range) * chartHeight

    const isRise = d.close >= d.open
    ctx.fillStyle = isRise ? '#18a058' : '#d03050'
    ctx.strokeStyle = isRise ? '#18a058' : '#d03050'

    // 绘制上下影线
    ctx.beginPath()
    ctx.moveTo(x + candleWidth / 2, highY)
    ctx.lineTo(x + candleWidth / 2, lowY)
    ctx.stroke()

    // 绘制实体
    const bodyHeight = Math.abs(closeY - openY)
    const bodyY = Math.min(openY, closeY)
    ctx.fillRect(x, bodyY, candleWidth, Math.max(bodyHeight, 1))
  })

  // 绘制MA线
  const maColors = ['#f59e0b', '#8b5cf6', '#3b82f6', '#10b981']
  props.maLines.forEach((period, idx) => {
    const maData = calculateMA(props.data, period)
    ctx.strokeStyle = maColors[idx] || '#999'
    ctx.lineWidth = 1
    ctx.beginPath()

    let started = false
    maData.forEach((ma, index) => {
      if (ma === null) return
      const x = padding.left + (index / props.data.length) * chartWidth + candleWidth / 2
      const y = padding.top + ((max - ma) / range) * chartHeight

      if (!started) {
        ctx.moveTo(x, y)
        started = true
      } else {
        ctx.lineTo(x, y)
      }
    })
    ctx.stroke()
  })
}

watch(() => props.data, draw, { deep: true })
watch(() => [props.width, props.height], draw)

onMounted(() => {
  draw()
})
</script>

<template>
  <div class="kline-chart">
    <canvas
      ref="canvasRef"
      :style="{ width: `${width}px`, height: `${height}px` }"
    />
    <div class="chart-placeholder" v-if="!data.length">
      <p>📊 K线图</p>
      <p class="placeholder-tip">暂无数据</p>
    </div>
  </div>
</template>

<style scoped>
.kline-chart {
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
