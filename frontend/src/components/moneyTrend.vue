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

onMounted(
    () => {
      handleLine(code, days)
      if (LineChartRef.value && window.ResizeObserver) {
        resizeObs = new ResizeObserver(() => { if (chartInstance) chartInstance.resize() })
        resizeObs.observe(LineChartRef.value)
      }
    }
)
onUnmounted(() => {
  if (resizeObs) resizeObs.disconnect()
  if (chartInstance) chartInstance.dispose()
})
const handleLine = (code, days) => {
  GetStockMoneyTrendByDay(code, days).then(result => {
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
        padding: isMobile ? 6 : 10,
        confine: true,
        textStyle: {
          color: darkTheme?'#ccc':'#456',
          fontSize: isMobile ? 11 : 12
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
        textStyle: {
          color: darkTheme ? 'rgb(253,252,252)' : '#456',
          fontSize: isMobile ? 10 : 12
        },
        right: isMobile ? 8 : 150,
        itemWidth: isMobile ? 14 : 25,
        itemHeight: isMobile ? 8 : 14,
        itemGap: isMobile ? 6 : 10,
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
          height: isMobile ? 22 : 16,
          handleSize: isMobile ? 22 : undefined,
          start: 86,
          end: 100
        }
      ],
      grid: [
        {
          left: isMobile ? '14%' : '8%',
          right: isMobile ? '6%' : '8%',
          top: isMobile ? '18%' : '8%',
          height: '50%',
        },
        {
          left: isMobile ? '14%' : '8%',
          right: isMobile ? '6%' : '8%',
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
          axisLabel: { rotate: isMobile ? 45 : 0, fontSize: isMobile ? 10 : 12 },
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
          name: '当日净流入/万',
          type: 'value',
          axisLine: {
            show: true
          },
          splitLine: {
            show: false
          },
        },
        {
          name: '股价',
          type: 'value',
          min: min - 1,
          max: max + 1,
          minInterval: 0.01,
          axisLine: {
            show: true
          },
          splitLine: {
            show: false
          },
        },
        {
          gridIndex: 1,
          name: '累计净流入/万',
          type: 'value',
          axisLine: {
            show: true
          },
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
            },
            data: [
              {type: 'max', name: 'Max'},
              {type: 'min', name: 'Min'}
            ]
          },
          markLine: {
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
            symbolSize: [10, 20],
            symbolOffset: [10, 0],
            itemStyle: {
              color: '#f39509'
            },
            label: {
              position: 'right',
            },
            data: [
              {type: 'max', name: 'Max'},
              {type: 'min', name: 'Min'}
            ]
          },
          markLine: {
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
            symbolSize: [10, 20],
            symbolOffset: [10, 0],
            // itemStyle: {
            //   color: '#f39509'
            // },
            label: {
              position: 'right',
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