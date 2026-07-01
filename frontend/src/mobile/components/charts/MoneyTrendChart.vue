<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick, shallowRef } from 'vue'
import MIcon from '../base/MIcon.vue'

const props = defineProps({
  // 资金趋势数据：来自 GetStockMoneyTrendByDay
  // [{ opendate, netamount(当日净流入), r0_net(主力净流入), trade(股价) }]
  data: {
    type: Array,
    default: () => []
  },
  // 图表高度
  height: {
    type: Number,
    default: 380
  }
})

const chartRef = ref(null)
let chart = null
const echartsLib = shallowRef(null)

function cssColor(name, fallback) {
  if (typeof window === 'undefined') return fallback
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback
}

// ---- 数据规范化：字符串/数值 → 数值，单位换算为「万」 ----
const norm = computed(() => {
  const arr = props.data || []
  const category = []
  const netAmount = []   // 当日净流入（万）
  const r0Net = []       // 主力当日净流入（万）
  const trade = []       // 股价
  const cumulate = []    // 累计净流入（万）— 对齐桌面端：相邻两日 netamount 之和
  let priceMin = 0
  let priceMax = 0
  for (let i = 0; i < arr.length; i++) {
    const d = arr[i]
    const na = Number(d.netamount) || 0
    const r0 = Number(d.r0_net) || 0
    const price = Number(d.trade) || 0
    category.push(d.opendate || '')
    netAmount.push(+(na / 10000).toFixed(2))
    r0Net.push(+(r0 / 10000).toFixed(2))
    trade.push(price)
    if (i > 0) {
      const prev = Number(arr[i - 1].netamount) || 0
      cumulate.push(+((na + prev) / 10000).toFixed(2))
    } else {
      cumulate.push(+(na / 10000).toFixed(2))
    }
    if (price > 0) {
      if (priceMin === 0 || price < priceMin) priceMin = price
      if (price > priceMax) priceMax = price
    }
  }
  return { category, netAmount, r0Net, trade, cumulate, priceMin, priceMax }
})

// 默认视窗：数据多时只显示后段，数据少时全显示
const zoomStart = computed(() => {
  const n = norm.value.category.length
  if (n <= 60) return 0
  return Math.max(0, Math.round((1 - 60 / n) * 100))
})

function buildOption() {
  const n = norm.value
  const rise = cssColor('--m-color-rise', '#d03050')
  const fall = cssColor('--m-color-fall', '#18a058')
  const textColor = cssColor('--m-text-secondary', '#666')
  const gridColor = cssColor('--m-divider-color', '#eee')

  const netColor = '#3b6df4'   // 当日净流入（线）
  const r0Color = fall          // 主力当日净流入（柱，绿）
  const priceColor = '#f39509'  // 股价（线，橙）
  const cumColor = rise         // 累计净流入（柱，红）

  return {
    animation: false,
    legend: {
      data: ['当日净流入', '主力净流入', '股价', '累计净流入'],
      top: 2,
      left: 'center',
      itemWidth: 14,
      itemHeight: 8,
      itemGap: 8,
      textStyle: { color: textColor, fontSize: 10 }
    },
    tooltip: {
      trigger: 'axis',
      triggerOn: 'mousemove|click',
      axisPointer: { type: 'cross', lineStyle: { color: '#888', width: 1, opacity: 0.8 } },
      backgroundColor: 'rgba(255,255,255,0.96)',
      borderColor: gridColor,
      borderWidth: 1,
      padding: 8,
      confine: true,
      textStyle: { color: '#333', fontSize: 11 },
      formatter: (params) => {
        if (!params || !params.length) return ''
        let html = `<div style="font-weight:600;margin-bottom:4px">${params[0].axisValue}</div>`
        for (const p of params) {
          if (p.value == null || p.value === '-') continue
          const unit = p.seriesName === '股价' ? '' : '万'
          html += `<div><span style="color:${p.color}">${p.seriesName}</span> ${p.value}${unit}</div>`
        }
        return html
      }
    },
    axisPointer: {
      link: [{ xAxisIndex: 'all' }],
      label: { backgroundColor: '#6b7077' }
    },
    grid: [
      { left: 8, right: 44, top: 28, height: '52%' },
      { left: 8, right: 44, top: '72%', height: '16%' }
    ],
    xAxis: [
      {
        type: 'category',
        data: n.category,
        boundaryGap: false,
        axisLine: { onZero: false, lineStyle: { color: gridColor } },
        axisLabel: { fontSize: 9, color: textColor },
        splitLine: { show: false },
        min: 'dataMin',
        max: 'dataMax',
        axisPointer: { z: 100 }
      },
      {
        type: 'category',
        gridIndex: 1,
        data: n.category,
        boundaryGap: false,
        axisLine: { onZero: false, lineStyle: { color: gridColor } },
        axisTick: { show: false },
        axisLabel: { show: false },
        splitLine: { show: false },
        min: 'dataMin',
        max: 'dataMax'
      }
    ],
    yAxis: [
      {
        // 净流入（万）— 左轴
        type: 'value',
        position: 'left',
        axisLabel: { show: false },
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: gridColor, type: 'dashed' } }
      },
      {
        // 股价 — 右轴
        type: 'value',
        position: 'right',
        scale: true,
        min: n.priceMin > 0 ? +(n.priceMin - 0.5).toFixed(2) : undefined,
        max: n.priceMax > 0 ? +(n.priceMax + 0.5).toFixed(2) : undefined,
        axisLabel: { fontSize: 9, color: priceColor, margin: 4 },
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { show: false }
      },
      {
        // 累计净流入（万）— 副图
        type: 'value',
        gridIndex: 1,
        splitNumber: 2,
        position: 'right',
        axisLabel: { show: false },
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { show: false }
      }
    ],
    dataZoom: [
      { type: 'inside', xAxisIndex: [0, 1], start: zoomStart.value, end: 100 },
      {
        show: true,
        type: 'slider',
        xAxisIndex: [0, 1],
        bottom: 2,
        height: 16,
        start: zoomStart.value,
        end: 100,
        handleSize: '120%',
        fillerColor: 'rgba(135,158,222,0.18)',
        borderColor: gridColor,
        textStyle: { fontSize: 9, color: textColor }
      }
    ],
    series: [
      {
        name: '主力净流入',
        type: 'bar',
        yAxisIndex: 0,
        data: n.r0Net,
        itemStyle: { color: r0Color, opacity: 0.65 }
      },
      {
        name: '当日净流入',
        type: 'line',
        yAxisIndex: 0,
        data: n.netAmount,
        smooth: false,
        showSymbol: false,
        lineStyle: { width: 1.5, color: netColor },
        itemStyle: { color: netColor },
        markPoint: {
          symbol: 'pin',
          symbolSize: 36,
          label: { fontSize: 9, color: '#fff' },
          itemStyle: { color: netColor },
          data: [
            { type: 'max', name: '最大' },
            { type: 'min', name: '最小' }
          ]
        }
      },
      {
        name: '股价',
        type: 'line',
        yAxisIndex: 1,
        data: n.trade,
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 2, color: priceColor },
        itemStyle: { color: priceColor }
      },
      {
        name: '累计净流入',
        type: 'bar',
        xAxisIndex: 1,
        yAxisIndex: 2,
        data: n.cumulate,
        itemStyle: { color: cumColor, opacity: 0.65 }
      }
    ]
  }
}

async function ensureEcharts() {
  if (echartsLib.value) return echartsLib.value
  const mod = await import('echarts')
  echartsLib.value = mod
  return mod
}

async function render() {
  const el = chartRef.value
  if (!el || !norm.value.category.length) return
  const echarts = await ensureEcharts()
  if (!chart) {
    chart = echarts.init(el)
  }
  chart.setOption(buildOption(), true)
  // tab 内容晚展开导致初次尺寸为 0，延迟 resize 兜底
  nextTick(() => chart && chart.resize())
}

function resize() {
  chart && chart.resize()
}

watch(() => props.data, () => { nextTick(render) }, { deep: true })

onMounted(() => {
  nextTick(render)
  window.addEventListener('resize', resize)
  setTimeout(resize, 400)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resize)
  if (chart) { chart.dispose(); chart = null }
})
</script>

<template>
  <div class="money-chart">
    <div
      ref="chartRef"
      class="money-canvas"
      :style="{ height: `${height}px` }"
    />
    <div v-if="!data.length" class="chart-empty">
      <p><MIcon name="money" :size="16" /> 暂无资金数据</p>
    </div>
  </div>
</template>

<style scoped>
.money-chart {
  position: relative;
  width: 100%;
}

.money-canvas {
  width: 100%;
}

.chart-empty {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--m-text-tertiary);
  font-size: var(--m-font-sm);
}
</style>
