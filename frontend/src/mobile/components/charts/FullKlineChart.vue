<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick, shallowRef } from 'vue'

const props = defineProps({
  // K线数据：来自 GetStockKLine，字段为字符串
  // [{ day, open, close, high, low, volume, amount }]
  data: {
    type: Array,
    default: () => []
  },
  // 图表高度
  height: {
    type: Number,
    default: 360
  }
})

const chartRef = ref(null)
let chart = null
const echartsLib = shallowRef(null)

function cssColor(name, fallback) {
  if (typeof window === 'undefined') return fallback
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback
}

// ---- 数据规范化：字符串字段 → 数值，拆出 category/values/volumes ----
const norm = computed(() => {
  const arr = props.data || []
  const category = []
  const values = []   // [open, close, low, high]
  const volumes = []  // [index, 成交量(万手), flag]
  for (let i = 0; i < arr.length; i++) {
    const d = arr[i]
    const open = Number(d.open) || 0
    const close = Number(d.close) || 0
    const low = Number(d.low) || 0
    const high = Number(d.high) || 0
    const volume = Number(d.volume) || 0
    category.push(d.day || '')
    values.push([open, close, low, high])
    // 涨跌 flag：收盘 >= 开盘 为涨(1)，否则跌(-1)，给成交量着色用
    volumes.push([i, +(volume / 10000).toFixed(2), close >= open ? 1 : -1])
  }
  return { category, values, volumes }
})

// MA 均线：基于收盘价(values[i][1])
function calcMA(period, values) {
  const result = []
  for (let i = 0; i < values.length; i++) {
    if (i < period - 1) {
      result.push('-')
      continue
    }
    let sum = 0
    for (let j = 0; j < period; j++) {
      sum += values[i - j][1]
    }
    result.push(+(sum / period).toFixed(2))
  }
  return result
}

// 默认视窗：数据多时只显示后段（对齐桌面端 start:86），数据少时全显示
const zoomStart = computed(() => {
  const n = norm.value.category.length
  if (n <= 60) return 0
  // 大约显示最近 60 根
  return Math.max(0, Math.round((1 - 60 / n) * 100))
})

function buildOption() {
  const n = norm.value
  const rise = cssColor('--m-color-rise', '#d03050')
  const fall = cssColor('--m-color-fall', '#18a058')
  const textColor = cssColor('--m-text-secondary', '#666')
  const gridColor = cssColor('--m-divider-color', '#eee')

  return {
    animation: false,
    legend: {
      data: ['日K', 'MA5', 'MA10', 'MA20', 'MA30'],
      bottom: 2,
      left: 'center',
      itemWidth: 14,
      itemHeight: 8,
      itemGap: 8,
      textStyle: { color: textColor, fontSize: 10 }
    },
    tooltip: {
      trigger: 'axis',
      // 移动端用 click 触发：避免 touchmove 同时驱动十字光标与 dataZoom 平移，导致滑动卡顿/抢手势
      triggerOn: 'click',
      axisPointer: { type: 'cross', lineStyle: { color: '#888', width: 1, opacity: 0.8 } },
      backgroundColor: 'rgba(255,255,255,0.96)',
      borderColor: gridColor,
      borderWidth: 1,
      padding: 8,
      textStyle: { color: '#333', fontSize: 11 },
      formatter: (params) => {
        if (!params || !params.length) return ''
        // 找出 K线 与成交量项
        const k = params.find(p => p.seriesName === '日K')
        const vol = params.find(p => p.seriesName === '成交量')
        const mas = params.filter(p => p.seriesName.startsWith('MA'))
        let html = `<div style="font-weight:600;margin-bottom:4px">${params[0].axisValue}</div>`
        if (k && Array.isArray(k.data)) {
          // candlestick data: [index?, open, close, low, high] → echarts 传入的是 [open,close,low,high]
          const d = k.data
          // d[0] 是 dataIndex（echarts 会前置），实际 OHLC 从 d[1] 起；兼容两种
          const o = d.length === 5 ? d[1] : d[0]
          const c = d.length === 5 ? d[2] : d[1]
          const l = d.length === 5 ? d[3] : d[2]
          const h = d.length === 5 ? d[4] : d[3]
          const up = c >= o
          const col = up ? rise : fall
          html += `<div style="color:${col}">开 ${o}　收 ${c}</div>`
          html += `<div style="color:${col}">低 ${l}　高 ${h}</div>`
        }
        if (vol && Array.isArray(vol.data)) {
          html += `<div>量 ${vol.data[1]}万手</div>`
        }
        if (mas.length) {
          html += '<div style="margin-top:2px">'
          for (const m of mas) {
            if (m.data != null && m.data !== '-') {
              html += `<span style="color:${m.color};margin-right:6px">${m.seriesName} ${m.data}</span>`
            }
          }
          html += '</div>'
        }
        return html
      }
    },
    axisPointer: {
      link: [{ xAxisIndex: 'all' }],
      label: { backgroundColor: '#6b7077' }
    },
    visualMap: {
      show: false,
      seriesIndex: 5,
      dimension: 2,
      pieces: [
        { value: 1, color: rise },
        { value: -1, color: fall }
      ]
    },
    grid: [
      { left: 8, right: 52, top: 12, height: '60%' },
      { left: 8, right: 52, top: '74%', height: '14%' }
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
        scale: true,
        position: 'right',
        axisLabel: { fontSize: 9, color: textColor, margin: 4 },
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: gridColor, type: 'dashed' } }
      },
      {
        scale: true,
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
        bottom: 24,
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
        name: '日K',
        type: 'candlestick',
        data: n.values,
        itemStyle: {
          color: rise,        // 阳线（涨）
          color0: fall,       // 阴线（跌）
          borderColor: rise,
          borderColor0: fall
        },
        markPoint: {
          symbol: 'pin',
          symbolSize: 38,
          label: { fontSize: 9, color: '#fff' },
          data: [
            { name: '最高', type: 'max', valueDim: 'highest', itemStyle: { color: rise } },
            { name: '最低', type: 'min', valueDim: 'lowest', itemStyle: { color: fall } }
          ]
        }
      },
      { name: 'MA5', type: 'line', data: calcMA(5, n.values), smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#3b82f6' } },
      { name: 'MA10', type: 'line', data: calcMA(10, n.values), smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#22c55e' } },
      { name: 'MA20', type: 'line', data: calcMA(20, n.values), smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#f59e0b' } },
      { name: 'MA30', type: 'line', data: calcMA(30, n.values), smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#ef4444' } },
      {
        name: '成交量',
        type: 'bar',
        xAxisIndex: 1,
        yAxisIndex: 1,
        data: n.volumes
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
}

function resize() {
  chart && chart.resize()
}

watch(() => props.data, () => { nextTick(render) }, { deep: true })

onMounted(() => {
  nextTick(render)
  window.addEventListener('resize', resize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resize)
  if (chart) { chart.dispose(); chart = null }
})
</script>

<template>
  <div class="kline-chart">
    <div
      ref="chartRef"
      class="kline-canvas"
      :style="{ height: `${height}px` }"
    />
    <div v-if="!data.length" class="chart-empty">
      <p>📊 暂无K线数据</p>
    </div>
  </div>
</template>

<style scoped>
.kline-chart {
  position: relative;
  width: 100%;
}

.kline-canvas {
  width: 100%;
  /* 图表区触摸手势交给 echarts，避免单指横滑平移K线被外层 sheet 垂直滚动抢走 */
  touch-action: none;
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
