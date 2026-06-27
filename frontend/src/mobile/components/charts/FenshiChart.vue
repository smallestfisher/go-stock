<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick, shallowRef } from 'vue'

const props = defineProps({
  // 分时数据：来自 GetStockMinutePriceLineData 的 priceData
  // [{ time:'09:30', price, volume(累计), amount(累计) }]
  data: {
    type: Array,
    default: () => []
  },
  // 昨收价（基准线 + 涨跌着色）
  preClose: {
    type: Number,
    default: 0
  },
  // 交易日期 YYYYMMDD（顶部标题展示）
  date: {
    type: String,
    default: ''
  },
  // 图表高度
  height: {
    type: Number,
    default: 300
  }
})

const chartRef = ref(null)
let chart = null
const echartsLib = shallowRef(null)

function cssColor(name, fallback) {
  if (typeof window === 'undefined') return fallback
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback
}

// ---- 数据规范化：累计 volume/amount → 每分钟增量；算累计均价 ----
const norm = computed(() => {
  const arr = props.data || []
  if (!arr.length) {
    return { category: [], price: [], avg: [], volume: [], open: 0, close: 0, high: 0, low: 0 }
  }

  let cumVol = 0
  let cumAmt = 0
  let prevVol = 0
  let prevAmt = 0
  let high = -Infinity
  let low = Infinity

  const category = []
  const price = []
  const avg = []
  const volume = []

  arr.forEach((d, i) => {
    const p = Number(d.price) || 0
    const totalVol = Number(d.volume) || 0
    const totalAmt = Number(d.amount) || 0
    const vol = Math.max(0, totalVol - prevVol)        // 差分得本分钟成交量
    const amtDiff = Math.max(0, totalAmt - prevAmt)
    prevVol = totalVol
    prevAmt = totalAmt
    cumVol += vol
    cumAmt += amtDiff
    if (p > high) high = p
    if (p < low) low = p

    category.push(d.time || '')
    price.push(p)
    // 均价 = 累计成交额 / 累计成交量 / 100（后端量额单位换算，对齐桌面端）
    avg.push(cumVol > 0 ? +(cumAmt / cumVol / 100).toFixed(3) : p)
    // 成交量柱：[索引, 量, 涨跌flag]（flag 用于着色）
    volume.push({ value: vol, itemStyle: { color: p >= props.preClose ? volRiseColor() : volFallColor() } })
  })

  return {
    category,
    price,
    avg,
    volume,
    open: price[0],
    close: price[price.length - 1],
    high: high === -Infinity ? 0 : high,
    low: low === Infinity ? 0 : low
  }
})

function volRiseColor() { return cssColor('--m-color-rise', '#d03050') + 'cc' }
function volFallColor() { return cssColor('--m-color-fall', '#18a058') + 'cc' }

// 价格 Y 轴范围：贴合当日实际最高/最低，各留 1% 余量（对齐桌面端，避免曲线挤在一角）
const yAxisRange = computed(() => {
  const n = norm.value
  if (!n.price.length) {
    return { min: undefined, max: undefined }
  }
  // 把昨收也纳入范围（保证基准虚线在可视区内）
  const hi = Math.max(n.high, props.preClose || n.high)
  const lo = Math.min(n.low, props.preClose || n.low)
  return {
    min: +(lo - lo * 0.01).toFixed(2),
    max: +(hi + hi * 0.01).toFixed(2)
  }
})

function buildOption() {
  const n = norm.value
  const pc = props.preClose
  const rise = cssColor('--m-color-rise', '#d03050')
  const fall = cssColor('--m-color-fall', '#18a058')
  const { min, max } = yAxisRange.value

  // 涨跌幅百分比 axisLabel 格式化器
  const pctFmt = (v) => {
    if (!pc) return v.toFixed(2)
    const pct = ((v - pc) / pc * 100).toFixed(2)
    return `${pct >= 0 ? '+' : ''}${pct}%`
  }

  const dateStr = props.date
    ? `${props.date.slice(0, 4)}-${props.date.slice(4, 6)}-${props.date.slice(6, 8)}`
    : ''
  const pct = pc ? ((n.close - pc) / pc * 100).toFixed(2) : '0.00'
  const subTitle = n.price.length
    ? `${dateStr}  开${n.open}  最新${n.close}  高${n.high}  低${n.low}  ${pct >= 0 ? '+' : ''}${pct}%`
    : ''

  return {
    animation: false,
    title: {
      text: subTitle,
      left: 'center',
      top: 2,
      textStyle: {
        color: cssColor('--m-text-secondary', '#6b7077'),
        fontSize: 10,
        fontWeight: 'normal'
      }
    },
    grid: [
      // 价格区
      { left: 6, right: 48, top: 26, height: '58%', containLabel: false },
      // 成交量区
      { left: 6, right: 48, bottom: 22, height: '18%', containLabel: false }
    ],
    tooltip: {
      trigger: 'axis',
      // 移动端：手指按下/拖动都触发游标（默认 mousemove 在 touch 上也会映射，但显式声明更稳）
      triggerOn: 'mousemove|click',
      alwaysShowContent: false,
      enterable: false,
      axisPointer: { type: 'cross', link: [{ xAxisIndex: 'all' }], label: { backgroundColor: '#6b7077' } },
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: '#e7e3da',
      borderWidth: 1,
      padding: [6, 10],
      textStyle: { color: '#1e2227', fontSize: 11 },
      position: (pt, params, dom, rect, size) => {
        // 跟随手指但不超出画布：水平居中偏移，垂直贴近触点上方
        const [x] = pt
        const boxW = size.contentSize[0]
        const viewW = size.viewSize[0]
        let left = x - boxW / 2
        if (left < 4) left = 4
        if (left + boxW > viewW - 4) left = viewW - boxW - 4
        return [left, 8]
      },
      confine: true,
      formatter: (params) => {
        if (!params || !params.length) return ''
        const time = params[0].axisValue
        let html = `<div style="font-weight:600;margin-bottom:2px">${time}</div>`
        params.forEach(p => {
          if (p.seriesName === '成交量') {
            html += `<div>${p.marker}成交量 <b>${Number(p.value.value ?? p.value).toLocaleString()}</b></div>`
          } else {
            const val = Number(p.value)
            const color = pc ? (val >= pc ? rise : fall) : '#1e2227'
            html += `<div>${p.marker}${p.seriesName} <b style="color:${color}">${val.toFixed(2)}</b></div>`
          }
        })
        return html
      }
    },
    axisPointer: { link: [{ xAxisIndex: 'all' }] },
    xAxis: [
      {
        type: 'category',
        data: n.category,
        boundaryGap: false,
        gridIndex: 0,
        axisLine: { lineStyle: { color: '#e7e3da' } },
        axisTick: { show: false },
        axisLabel: { show: false },
        splitLine: { show: false }
      },
      {
        type: 'category',
        data: n.category,
        boundaryGap: false,
        gridIndex: 1,
        axisLine: { lineStyle: { color: '#e7e3da' } },
        axisTick: { show: false },
        axisLabel: {
          color: '#999',
          fontSize: 9,
          interval: (idx) => idx === 0 || idx === Math.floor(n.category.length / 2) || idx === n.category.length - 1,
          formatter: (v) => v
        },
        splitLine: { show: false }
      }
    ],
    yAxis: [
      {
        type: 'value',
        scale: true,
        gridIndex: 0,
        min, max,
        position: 'right',
        splitNumber: 4,
        axisLine: { show: false },
        axisTick: { show: false },
        // 右轴显示涨跌幅百分比；价格用左侧 tooltip 看
        axisLabel: {
          fontSize: 9,
          margin: 4,
          formatter: (v) => {
            if (!pc) return v.toFixed(2)
            const pp = ((v - pc) / pc * 100)
            const c = pp >= 0 ? rise : fall
            return `{a|${v.toFixed(2)}}`
          },
          rich: { a: { color: '#888', fontSize: 9 } }
        },
        splitLine: { lineStyle: { color: '#f0f0f0', type: 'dashed' } }
      },
      {
        type: 'value',
        scale: true,
        gridIndex: 1,
        position: 'right',
        splitNumber: 2,
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: {
          fontSize: 9,
          margin: 4,
          color: '#999',
          formatter: (v) => v >= 10000 ? (v / 10000).toFixed(0) + '万' : v
        },
        splitLine: { show: false }
      }
    ],
    visualMap: {
      show: false,
      seriesIndex: 0,
      dimension: 1,
      // 用 min/max 保证区间合法：下跌股 open>close 时也不会出现非法区间
      // （否则整个 visualMap 失效、价格线画不出来）
      pieces: [
        { lte: Math.min(n.open, n.close), color: fall },                              // 低位段：绿
        { gt: Math.min(n.open, n.close), lte: Math.max(n.open, n.close), color: '#1651ef' }, // 开收之间：蓝
        { gt: Math.max(n.open, n.close), color: rise }                                // 高位段：红
      ],
      outOfRange: { color: rise }
    },
    series: [
      {
        name: '股价',
        type: 'line',
        data: n.price,
        xAxisIndex: 0,
        yAxisIndex: 0,
        showSymbol: false,
        smooth: false,
        lineStyle: { width: 1.2 },
        markLine: pc ? {
          symbol: 'none',
          silent: true,
          data: [
            { yAxis: pc, lineStyle: { color: '#f0b400', type: 'dashed', width: 1 }, label: { show: false } }
          ]
        } : undefined,
        markPoint: n.price.length ? {
          symbol: 'pin',
          symbolSize: 0,
          label: {
            fontSize: 10,
            color: '#888',
            formatter: (p) => p.value
          },
          data: [
            { type: 'max', name: '最高', label: { position: 'top', color: rise } },
            { type: 'min', name: '最低', label: { position: 'bottom', color: fall } }
          ]
        } : undefined
      },
      {
        name: '均价',
        type: 'line',
        data: n.avg,
        xAxisIndex: 0,
        yAxisIndex: 0,
        showSymbol: false,
        smooth: false,
        lineStyle: { width: 0.8, color: '#f0b400' }
      },
      {
        name: '成交量',
        type: 'bar',
        data: n.volume,
        xAxisIndex: 1,
        yAxisIndex: 1,
        barWidth: '60%'
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
  if (!el) return
  const echarts = await ensureEcharts()
  if (!chart) {
    chart = echarts.init(el)
  }
  chart.setOption(buildOption(), true)
}

function resize() {
  chart && chart.resize()
}

watch(() => [props.data, props.preClose], () => { nextTick(render) }, { deep: true })

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
  <div class="fenshi-chart">
    <div
      ref="chartRef"
      class="fenshi-canvas"
      :style="{ height: `${height}px` }"
    />
    <div class="chart-empty" v-if="!data.length">
      <p>📈 分时图</p>
      <p class="empty-tip">暂无分时数据</p>
    </div>
  </div>
</template>

<style scoped>
.fenshi-chart {
  position: relative;
  width: 100%;
}

.fenshi-canvas {
  width: 100%;
}

.chart-empty {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--m-text-secondary, #6b7077);
  pointer-events: none;
}

.empty-tip {
  font-size: var(--m-font-sm, 12px);
  color: var(--m-text-tertiary, #9ca3af);
  margin-top: var(--m-space-xs, 4px);
}
</style>
