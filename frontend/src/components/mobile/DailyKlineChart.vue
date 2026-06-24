<script setup>
import {onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {GetStockKLine} from '../../api/app'

const props = defineProps({
    code: {type: String, default: ''},
    name: {type: String, default: ''},
    darkTheme: {type: Boolean, default: false},
    chartHeight: {type: Number, default: 420},
})

const chartRef = ref(null)
let echartsModule = null
let chart = null

// 日K 蜡烛图配色（A 股习惯：红涨绿跌），与 stock.vue handleKLine 保持一致
const upColor = '#ec0000'
const downColor = '#00da3c'

// 迁移自 stock.vue 的 calculateMA
function calculateMA(dayCount, values) {
    var result = []
    for (let i = 0, len = values.length; i < len; i++) {
        if (i < dayCount) {
            result.push('-')
            continue
        }
        var sum = 0
        for (let j = 0; j < dayCount; j++) {
            sum += +values[i - j][1]
        }
        result.push((sum / dayCount).toFixed(2))
    }
    return result
}

async function renderChart() {
    if (!chartRef.value || !props.code) {
        return
    }
    if (!echartsModule) {
        echartsModule = await import('echarts')
    }
    if (chart) {
        chart.dispose()
    }
    chart = echartsModule.init(chartRef.value)
    GetStockKLine(props.code, props.name, 365).then(result => {
        if (!chart) {
            return
        }
        const categoryData = []
        const values = []
        const volumns = []
        for (let i = 0; i < result.length; i++) {
            let resultElement = result[i]
            categoryData.push(resultElement.day)
            let flag = resultElement.close > resultElement.open ? 1 : -1
            values.push([
                resultElement.open,
                resultElement.close,
                resultElement.low,
                resultElement.high
            ])
            volumns.push([i, resultElement.volume / 10000, flag])
        }
        let option = {
            darkMode: props.darkTheme,
            animation: false,
            legend: {
                top: 0,
                left: 'center',
                data: ['日K', 'MA5', 'MA10', 'MA20', 'MA30'],
                textStyle: {
                    color: props.darkTheme ? '#ccc' : '#456'
                },
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
                borderColor: props.darkTheme ? '#456' : '#ccc',
                backgroundColor: props.darkTheme ? '#456' : '#fff',
                padding: 10,
                textStyle: {
                    color: props.darkTheme ? '#ccc' : '#456'
                },
                formatter: function (params) {
                    let volum = params[5].data
                    let ma5 = params[1].data
                    let ma10 = params[2].data
                    let ma20 = params[3].data
                    let ma30 = params[4].data
                    params = params[0]
                    let currentItemData = params.data

                    return params.name + '<br>' +
                        '开盘:' + currentItemData[1] + '<br>' +
                        '收盘:' + currentItemData[2] + '<br>' +
                        '最低:' + currentItemData[3] + '<br>' +
                        '最高:' + currentItemData[4] + '<br>' +
                        '成交量(万手):' + volum[1] + '<br>' +
                        'MA5日均线:' + ma5 + '<br>' +
                        'MA10日均线:' + ma10 + '<br>' +
                        'MA20日均线:' + ma20 + '<br>' +
                        'MA30日均线:' + ma30
                }
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
            visualMap: {
                show: false,
                seriesIndex: 5,
                dimension: 2,
                pieces: [
                    {
                        value: -1,
                        color: downColor
                    },
                    {
                        value: 1,
                        color: upColor
                    }
                ]
            },
            grid: [
                {
                    left: '10%',
                    right: '8%',
                    top: 42,
                    height: '48%',
                },
                {
                    left: '10%',
                    right: '8%',
                    top: '66%',
                    height: '13%'
                }
            ],
            xAxis: [
                {
                    type: 'category',
                    data: categoryData,
                    boundaryGap: false,
                    axisLine: {onZero: false},
                    splitLine: {show: false},
                    min: 'dataMin',
                    max: 'dataMax',
                    axisPointer: {
                        z: 100
                    }
                },
                {
                    type: 'category',
                    gridIndex: 1,
                    data: categoryData,
                    boundaryGap: false,
                    axisLine: {onZero: false},
                    axisTick: {show: false},
                    splitLine: {show: false},
                    axisLabel: {show: false},
                    min: 'dataMin',
                    max: 'dataMax'
                }
            ],
            yAxis: [
                {
                    scale: true,
                    splitArea: {
                        show: true
                    }
                },
                {
                    scale: true,
                    gridIndex: 1,
                    splitNumber: 2,
                    axisLabel: {show: false},
                    axisLine: {show: false},
                    axisTick: {show: false},
                    splitLine: {show: false}
                }
            ],
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
                    bottom: 8,
                    height: 24,
                    start: 86,
                    end: 100
                }
            ],
            series: [
                {
                    name: '日K',
                    type: 'candlestick',
                    data: values,
                    itemStyle: {
                        color: upColor,
                        color0: downColor,
                    },
                    markPoint: {
                        label: {
                            formatter: function (param) {
                                return param != null ? param.value + '' : ''
                            }
                        },
                        data: [
                            {
                                name: '最高',
                                type: 'max',
                                valueDim: 'highest'
                            },
                            {
                                name: '最低',
                                type: 'min',
                                valueDim: 'lowest'
                            },
                            {
                                name: '平均收盘价',
                                type: 'average',
                                valueDim: 'close'
                            }
                        ],
                        tooltip: {
                            formatter: function (param) {
                                return param.name + '<br>' + (param.data.coord || '')
                            }
                        }
                    },
                    markLine: {
                        symbol: ['none', 'none'],
                        data: [
                            [
                                {
                                    name: 'from lowest to highest',
                                    type: 'min',
                                    valueDim: 'lowest',
                                    symbol: 'circle',
                                    symbolSize: 10,
                                    label: {
                                        show: false
                                    },
                                    emphasis: {
                                        label: {
                                            show: false
                                        }
                                    }
                                },
                                {
                                    type: 'max',
                                    valueDim: 'highest',
                                    symbol: 'circle',
                                    symbolSize: 10,
                                    label: {
                                        show: false
                                    },
                                    emphasis: {
                                        label: {
                                            show: false
                                        }
                                    }
                                }
                            ],
                            {
                                name: 'min line on close',
                                type: 'min',
                                valueDim: 'close'
                            },
                            {
                                name: 'max line on close',
                                type: 'max',
                                valueDim: 'close'
                            }
                        ]
                    }
                },
                {
                    name: 'MA5',
                    type: 'line',
                    data: calculateMA(5, values),
                    smooth: true,
                    showSymbol: false,
                    lineStyle: {
                        opacity: 0.6
                    }
                },
                {
                    name: 'MA10',
                    type: 'line',
                    data: calculateMA(10, values),
                    smooth: true,
                    showSymbol: false,
                    lineStyle: {
                        opacity: 0.6
                    }
                },
                {
                    name: 'MA20',
                    type: 'line',
                    data: calculateMA(20, values),
                    smooth: true,
                    showSymbol: false,
                    lineStyle: {
                        opacity: 0.6
                    }
                },
                {
                    name: 'MA30',
                    type: 'line',
                    data: calculateMA(30, values),
                    smooth: true,
                    showSymbol: false,
                    lineStyle: {
                        opacity: 0.6
                    }
                },
                {
                    name: '成交量(手)',
                    type: 'bar',
                    xAxisIndex: 1,
                    yAxisIndex: 1,
                    itemStyle: {
                        color: '#7fbe9e'
                    },
                    data: volumns
                }
            ]
        }
        chart.setOption(option)
        chart.on('click', {seriesName: '日K'}, function (params) {
        })
    })
}

onMounted(() => {
    renderChart()
})

watch(() => props.code, () => {
    renderChart()
})

onBeforeUnmount(() => {
    if (chart) {
        chart.dispose()
        chart = null
    }
})
</script>

<template>
    <div ref="chartRef" :style="{width: '100%', height: chartHeight + 'px'}"></div>
</template>
