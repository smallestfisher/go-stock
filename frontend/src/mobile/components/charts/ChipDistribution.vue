<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { calcChipDistribution } from '../../composables/klineChip'

// 筹码分布（纯前端基于 K 线计算）。
// 输入 K 线数据，canvas 横向柱状图展示各价位持仓占比 + 均成本线 + 获利比例。
const props = defineProps({
  data: { type: Array, default: () => [] },
  // 箱数
  bins: { type: Number, default: 60 },
  height: { type: Number, default: 180 },
})

const canvasRef = ref(null)
const result = ref({ items: [], avgCost: 0, profitRatio: 0, current: 0 })

function compute() {
  if (!props.data || !props.data.length) {
    result.value = { items: [], avgCost: 0, profitRatio: 0, current: 0 }
    return
  }
  result.value = calcChipDistribution(props.data, props.bins)
}

const profitPct = computed(() => (result.value.profitRatio * 100).toFixed(1) + '%')

function cssColor(name, fallback) {
  if (typeof window === 'undefined') return fallback
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback
}

function draw() {
  const canvas = canvasRef.value
  if (!canvas) return
  const items = result.value.items
  if (!items.length) return

  const dpr = window.devicePixelRatio || 1
  const cssW = canvas.clientWidth || 300
  const cssH = props.height
  canvas.width = cssW * dpr
  canvas.height = cssH * dpr
  const ctx = canvas.getContext('2d')
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, cssW, cssH)

  const rise = cssColor('--m-color-rise', '#d03050')
  const fall = cssColor('--m-color-fall', '#18a058')
  const neutral = cssColor('--m-text-tertiary', '#999')
  const cur = result.value.current || 0
  const avg = result.value.avgCost || 0

  const maxRatio = Math.max(...items.map(i => i.ratio), 0.0001)
  const labelW = 44
  const barMaxW = cssW - labelW - 8
  const n = items.length
  const rowH = cssH / n

  for (let i = 0; i < n; i++) {
    const it = items[i]
    const y = i * rowH
    const barW = (it.ratio / maxRatio) * barMaxW
    // 低于当前价（获利盘）用红，高于用绿
    ctx.fillStyle = it.price <= cur ? rise : fall
    ctx.globalAlpha = 0.5
    ctx.fillRect(labelW, y + 0.5, barW, rowH - 1)
    ctx.globalAlpha = 1
  }

  // 均成本线
  const items2 = items
  const minP = result.value.minPrice || items2[0].price
  const maxP = result.value.maxPrice || items2[items2.length - 1].price
  const priceToY = (p) => {
    if (maxP <= minP) return cssH / 2
    return cssH - ((p - minP) / (maxP - minP)) * cssH
  }
  ctx.strokeStyle = neutral
  ctx.lineWidth = 1
  ctx.setLineDash([4, 3])
  ctx.beginPath()
  const avgY = priceToY(avg)
  ctx.moveTo(labelW, avgY)
  ctx.lineTo(cssW, avgY)
  ctx.stroke()
  ctx.setLineDash([])

  // 当前价线
  ctx.strokeStyle = cur >= avg ? rise : fall
  ctx.beginPath()
  const curY = priceToY(cur)
  ctx.moveTo(labelW, curY)
  ctx.lineTo(cssW, curY)
  ctx.stroke()

  // 标签：均成本
  ctx.fillStyle = neutral
  ctx.font = '10px sans-serif'
  ctx.fillText(avg.toFixed(2), 2, avgY - 2)
  ctx.fillText(cur.toFixed(2), 2, curY - 2)
}

function redraw() { nextTick(draw) }

watch(() => props.data, () => { compute(); redraw() }, { deep: true })

onMounted(() => {
  compute()
  redraw()
  window.addEventListener('resize', redraw)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', redraw)
})
</script>

<template>
  <div class="chip-dist">
    <div class="cd-head">
      <span class="cd-title">筹码分布</span>
      <span v-if="result.items.length" class="cd-meta">
        均成本 {{ result.avgCost.toFixed(2) }} · 获利 {{ profitPct }}
      </span>
    </div>
    <canvas
      v-show="result.items.length"
      ref="canvasRef"
      class="cd-canvas"
      :style="{ height: height + 'px' }"
    />
    <div v-if="!result.items.length" class="cd-empty">
      {{ data.length ? '计算中...' : '暂无K线数据' }}
    </div>
  </div>
</template>

<style scoped>
.chip-dist {
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
}

.cd-head {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-sm);
}

.cd-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
}

.cd-meta {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.cd-canvas {
  width: 100%;
  display: block;
}

.cd-empty {
  text-align: center;
  padding: var(--m-space-lg) 0;
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}
</style>
