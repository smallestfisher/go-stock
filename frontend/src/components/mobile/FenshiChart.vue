<script setup>
import {onBeforeUnmount, onMounted, ref, watch} from 'vue'
import {GetStockMinutePriceLineData} from '../../api/app'

const props = defineProps({
    code: {type: String, default: ''},
    name: {type: String, default: ''},
    darkTheme: {type: Boolean, default: false},
    chartHeight: {type: Number, default: 420},
})

const chartRef = ref(null)
let echartsModule = null
let chart = null
let refreshTimer = null

async function renderChart() {
    if (!chartRef.value || !props.code) {
        return
    }
    if (!echartsModule) {
        echartsModule = await import('echarts')
    }
    if (!chart) {
        chart = echartsModule.init(chartRef.value)
    }
    GetStockMinutePriceLineData(props.code, props.name).then(result => {
        if (!chart) {
            return
        }
        const priceData = result.priceData
        let category = []
        let price = []
        let openprice = 0
        let closeprice = 0
        let volume = []
        let volumeRate = []
        let min = 0
        let max = 0
        openprice = priceData[0].price
        closeprice = priceData[priceData.length - 1].price
        for (let i = 0; i < priceData.length; i++) {
            category.push(priceData[i].time)
            price.push(priceData[i].price)
            if (min === 0 || min > priceData[i].price) {
                min = priceData[i].price
            }
            if (max < priceData[i].price) {
                max = priceData[i].price
            }
            if (i > 0) {
                let b = priceData[i].volume - priceData[i - 1].volume
                volumeRate.push(((b - volume[i - 1]) / volume[i - 1] * 100).toFixed(2))
                volume.push(b)
            } else {
                volume.push(priceData[i].volume)
                volumeRate.push(0)
            }
        }

        let option = {
            title: {
                subtext: "[" + result.date + "] 开盘:" + openprice + " 最新:" + closeprice + " 最高:" + max + " 最低:" + min,
                left: 'center',
                top: '10',
                textStyle: {
                    color: props.darkTheme ? '#ccc' : '#456'
                }
            },
            legend: {
                data: ['股价', '成交量'],
                textStyle: {
                    color: props.darkTheme ? '#ccc' : '#456'
                },
                right: 50,
            },
            darkMode: props.darkTheme,
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'cross',
                    animation: false,
                    label: {
                        backgroundColor: '#505765'
                    }
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
            xAxis: [
                {
                    type: 'category',
                    data: category,
                    axisLabel: {
                        show: false
                    }
                },
                {
                    gridIndex: 1,
                    type: 'category',
                    data: category,
                },
            ],
            grid: [
                {
                    left: '8%',
                    right: '8%',
                    height: '50%',
                },
                {
                    left: '8%',
                    right: '8%',
                    top: '70%',
                    height: '15%'
                },
            ],
            yAxis: [
                {
                    axisLine: {
                        show: true
                    },
                    splitLine: {
                        show: false
                    },
                    name: "股价",
                    min: (min - min * 0.01).toFixed(2),
                    max: (max + max * 0.01).toFixed(2),
                    minInterval: 0.01,
                    type: 'value'
                },
                {
                    gridIndex: 1,
                    axisLine: {
                        show: true
                    },
                    splitLine: {
                        show: false
                    },
                    name: "成交量",
                    type: 'value',
                },
            ],
            visualMap: {
                type: 'piecewise',
                seriesIndex: 0,
                top: 0,
                left: 10,
                orient: 'horizontal',
                textStyle: {
                    color: props.darkTheme ? '#fff' : '#456'
                },
                pieces: [
                    {
                        text: '低于开盘价',
                        gt: 0,
                        lte: openprice,
                        color: '#31F113',
                        textStyle: {
                            color: props.darkTheme ? '#fff' : '#456'
                        },
                    },
                    {
                        text: '大于开盘价小于收盘价',
                        gt: openprice,
                        lte: closeprice,
                        color: '#1651EF',
                        textStyle: {
                            color: props.darkTheme ? '#fff' : '#456'
                        },
                    },
                    {
                        text: '大于收盘价',
                        gt: closeprice,
                        color: '#AC3B2A',
                        textStyle: {
                            color: props.darkTheme ? '#fff' : '#456'
                        },
                    }
                ],
            },
            series: [
                {
                    name: "股价",
                    data: price,
                    type: 'line',
                    smooth: false,
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
                            color: '#FC290D'
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
                        symbol: 'none',
                        data: [
                            {type: 'average', name: 'Average'},
                            {
                                lineStyle: {
                                    color: '#FFCB00',
                                    width: 0.5
                                },
                                yAxis: openprice,
                                name: '开盘价'
                            },
                            {
                                yAxis: closeprice,
                                symbol: 'none',
                                lineStyle: {
                                    color: 'red',
                                    width: 0.5
                                },
                            }
                        ]
                    },
                },
                {
                    xAxisIndex: 1,
                    yAxisIndex: 1,
                    name: "成交量",
                    data: volume,
                    type: 'bar',
                },

            ]
        };
        chart.setOption(option);
    })
}

function startRefresh() {
    stopRefresh()
    refreshTimer = setInterval(() => {
        renderChart()
    }, 1000 * 10)
}

function stopRefresh() {
    if (refreshTimer) {
        clearInterval(refreshTimer)
        refreshTimer = null
    }
}

onMounted(() => {
    renderChart()
    startRefresh()
})

watch(() => props.code, () => {
    renderChart()
})

onBeforeUnmount(() => {
    stopRefresh()
    if (chart) {
        chart.dispose()
        chart = null
    }
})
</script>

<template>
    <div ref="chartRef" :style="{width: '100%', height: chartHeight + 'px'}"></div>
</template>
