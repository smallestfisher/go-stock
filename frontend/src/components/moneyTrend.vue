<script setup lang="ts">
import {onMounted, onUnmounted, ref} from "vue";
import {GetStockMoneyTrendByDay} from "../api/app";
import * as echarts from "echarts";
import {useDevice} from "../composables/useDevice";

const {isMobile} = useDevice()

const {code, name, darkTheme, days, chartHeight} = defineProps({
  code: {
    type: String,
    default: ''
  },
  name: {
    type: String,
    default: ''
  },
  days: {
    type: Number,
    default: 14
  },
  chartHeight: {
    type: Number,
    default: 500
  },
  darkTheme: {
    type: Boolean,
    default: false
  }
})
const LineChartRef = ref(null);
let chartInstance = null
let resizeObs = null
let visibleObs = null

onMounted(
    () => {
      handleLine(code, days)
      if (LineChartRef.value && window.ResizeObserver) {
        resizeObs = new ResizeObserver(() => { if (chartInstance) chartInstance.resize() })
        resizeObs.observe(LineChartRef.value)
      }
      // tab 切换/抽屉展开导致容器 display:none→block 时，ResizeObserver 不一定触发，
      // 用 IntersectionObserver 监听可见性，可见即 resize（解决资金图在 tab 内画错尺寸）
      if (LineChartRef.value && window.IntersectionObserver) {
        visibleObs = new IntersectionObserver((entries) => {
          entries.forEach(e => {
            if (e.isIntersecting && chartInstance) {
              chartInstance.resize()
            }
          })
        })
        visibleObs.observe(LineChartRef.value)
      }
      // 兜底：延迟重画一次，应对容器晚展开
      setTimeout(() => { if (chartInstance) chartInstance.resize() }, 400)
    }
)
onUnmounted(() => {
  if (resizeObs) resizeObs.disconnect()
  if (visibleObs) visibleObs.disconnect()
  if (chartInstance) chartInstance.dispose()
})
const handleLine = (code, days) => {
  GetStockMoneyTrendByDay(code, days).then(result => {
    const mobile = isMobile.value
    const chart = echarts.init(LineChartRef.value);
    chartInstance = chart
    const categoryData = [];
    const netamount_values = [];
    const r0_net_values = [];
    const trades_values = [];
    let volume = []

    let min = 0
    let max = 0
    for (let i = 0; i < result.length; i++) {
      let resultElement = result[i]
      categoryData.push(resultElement.opendate)
      let netamount = (resultElement.netamount / 10000).toFixed(2);
      netamount_values.push(netamount)
      let price = Number(resultElement.trade);
      trades_values.push(price)
      r0_net_values.push((resultElement.r0_net / 10000).toFixed(2))

      if (min === 0 || min > price) {
        min = price
      }
      if (max < price) {
        max = price
      }

      if (i > 0) {
        let b = Number(Number(result[i].netamount) + Number(result[i - 1].netamount)) / 10000
        volume.push(b.toFixed(2))
      } else {
        volume.push((Number(result[i].netamount) / 10000).toFixed(2))
      }

    }
    //console.log("volume", volume)
    const upColor = '#ec0000';
    const downColor = '#00da3c';
    let option = {
      title: {
        show: !mobile,
        text: name,
        left: '20px',
        textStyle: {
          color: darkTheme?'#ccc':'#456'
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          lineStyle: {
            color: '#376df4',
            width: 1,
            opacity: 1
          }
        },
        borderWidth: 2,
        borderColor: darkTheme?'#456':'#ccc',
        backgroundColor: darkTheme?'#456':'#fff',
        padding: mobile ? 6 : 10,
        confine: true,
        textStyle: {
          color: darkTheme?'#ccc':'#456',
          fontSize: mobile ? 11 : 12
        },
      },
      axisPointer: {
        link: [
          {
            xAxisIndex: 'all'
          }
        ],
        label: {
          backgroundColor: '#888'
        }
      },
      legend: {
        show: true,
        data: ['当日净流入', '主力当日净流入','累计净流入',  '股价'],
        selected: {
          '当日净流入': true,
          '主力当日净流入': true,
          '累计净流入': true,
          '股价': true,
        },
        // 移动端：顶部一行 + 可滚动分页，避免 4 个图例项在窄屏重叠
        ...(mobile
          ? { top: 0, left: 0, right: 0, type: 'scroll', pageIconSize: 10 }
          : { top: 'auto', right: 150 }),
        textStyle: {
          color: darkTheme ? 'rgb(253,252,252)' : '#456',
          fontSize: mobile ? 10 : 12
        },
        itemWidth: mobile ? 14 : 25,
        itemHeight: mobile ? 8 : 14,
        itemGap: mobile ? 6 : 10,
      },
      dataZoom: [
        {
          type: 'inside',
          xAxisIndex: [0, 1],
          start: 86,
          end: 100
        },
        {
          show: true,
          xAxisIndex: [0, 1],
          type: 'slider',
          top: '90%',
          height: mobile ? 22 : 16,
          handleSize: mobile ? 22 : undefined,
          start: 86,
          end: 100
        }
      ],
      grid: [
        {
          left: mobile ? '14%' : '8%',
          right: mobile ? '6%' : '8%',
          top: mobile ? '18%' : '8%',
          height: '50%',
        },
        {
          left: mobile ? '14%' : '8%',
          right: mobile ? '6%' : '8%',
          top: '74%',
          height: '15%'
        },
      ],
      xAxis: [
        {
          type: 'category',
          data: categoryData,
          axisPointer: {
            z: 100
          },
          axisLabel: { rotate: mobile ? 45 : 0, fontSize: mobile ? 10 : 12 },
          boundaryGap: false,
          axisLine: { onZero: false },
          splitLine: { show: false },
          min: 'dataMin',
          max: 'dataMax',


        },
        {
          gridIndex: 1,
          type: 'category',
          data: categoryData,
          axisLabel: {
            show: false
          },
        }
      ],
      yAxis: [
        {
          name: mobile ? '' : '当日净流入/万',
          type: 'value',
          axisLine: {
            show: true
          },
          axisLabel: { fontSize: mobile ? 10 : 12 },
          splitLine: {
            show: false
          },
        },
        {
          name: mobile ? '' : '股价',
          type: 'value',
          min: min - 1,
          max: max + 1,
          minInterval: 0.01,
          axisLine: {
            show: true
          },
          axisLabel: { fontSize: mobile ? 10 : 12 },
          splitLine: {
            show: false
          },
        },
        {
          gridIndex: 1,
          name: mobile ? '' : '累计净流入/万',
          type: 'value',
          axisLine: {
            show: true
          },
          axisLabel: { fontSize: mobile ? 10 : 12 },
          splitLine: {
            show: false
          },
        },
      ],
      series: [
        {
          yAxisIndex: 0,
          name: '当日净流入',
          data: netamount_values,
          smooth: false,
          showSymbol: false,
          lineStyle: {
            width: 2
          },
          markPoint: {
            symbol: 'arrow',
            symbolRotate: 90,
            symbolSize: [10, 20],
            symbolOffset: [10, 0],
            itemStyle: {
              color: '#0d7dfc'
            },
            label: {
              position: 'right',
              show: !mobile,
            },
            data: [
              {type: 'max', name: 'Max'},
              {type: 'min', name: 'Min'}
            ]
          },
          markLine: {
            symbol: mobile ? 'none' : undefined,
            label: { show: !mobile },
            data: [
              {
                type: 'average',
                name: 'Average',
                lineStyle: {
                  color: '#0077ff',
                  width: 0.5
                },
              },
            ]
          },
          type: 'line'
        },
        {
          yAxisIndex: 0,
          name: '主力当日净流入',
          data: r0_net_values,
          smooth: false,
          showSymbol: false,
          lineStyle: {
            width: 2
          },
          // markPoint: {
          //   symbol: 'arrow',
          //   symbolRotate: 90,
          //   symbolSize: [10, 20],
          //   symbolOffset: [10, 0],
          //   itemStyle: {
          //     color: '#0d7dfc'
          //   },
          //   label: {
          //     position: 'right',
          //   },
          //   data: [
          //     {type: 'max', name: 'Max'},
          //     {type: 'min', name: 'Min'}
          //   ]
          // },
          // markLine: {
          //   data: [
          //     {
          //       type: 'average',
          //       name: 'Average',
          //       lineStyle: {
          //         color: '#0077ff',
          //         width: 0.5
          //       },
          //     },
          //   ]
          // },
          type: 'bar'
        },
        {
          yAxisIndex: 1,
          name: '股价',
          type: 'line',
          data: trades_values,
          smooth: true,
          showSymbol: false,
          lineStyle: {
            width: 3
          },
          markPoint: {
            symbol: 'arrow',
            symbolRotate: 90,
            symbolSize: mobile ? [8, 16] : [10, 20],
            symbolOffset: mobile ? [0, 0] : [10, 0],
            itemStyle: {
              color: '#f39509'
            },
            label: {
              position: 'right',
              show: !mobile,
            },
            data: [
              {type: 'max', name: 'Max'},
              {type: 'min', name: 'Min'}
            ]
          },
          markLine: {
            symbol: mobile ? 'none' : undefined,
            label: { show: !mobile },
            data: [
              {
                type: 'average',
                name: 'Average',
                lineStyle: {
                  color: '#f39509',
                  width: 0.5
                },
              },
            ]
          },
        },
        {
          type: 'bar',
          xAxisIndex: 1,
          yAxisIndex: 2,
          name: '累计净流入',
          data: volume,
          smooth: true,
          showSymbol: false,
          lineStyle: {
            width: 2
          },
          markPoint: {
            symbol: 'arrow',
            symbolRotate: 90,
            symbolSize: mobile ? [8, 16] : [10, 20],
            symbolOffset: mobile ? [0, 0] : [10, 0],
            // itemStyle: {
            //   color: '#f39509'
            // },
            label: {
              position: 'right',
              show: !mobile,
            },
            data: [
              {type: 'max', name: 'Max'},
              {type: 'min', name: 'Min'}
            ]
          },
        },
      ]
    };
    chart.setOption(option);
  })
}
</script>

<template>
  <div ref="LineChartRef" style="width: 100%;height: auto;" :style="{height:chartHeight+'px'}"></div>
</template>

<style scoped>

</style>
