<script setup>
import {computed, defineAsyncComponent, ref} from 'vue'
import {toEastMoneyCode} from '../../utils/stockCode'
import BottomSheet from './BottomSheet.vue'
import FenshiChart from './FenshiChart.vue'
import DailyKlineChart from './DailyKlineChart.vue'

const MoneyTrend = defineAsyncComponent(() => import('../moneyTrend.vue'))
const StockLightweightKlineChart = defineAsyncComponent(() => import('../StockLightweightKlineChart.vue'))

const props = defineProps({
    show: {type: Boolean, default: false},
    result: {type: Object, default: null},
    darkTheme: {type: Boolean, default: false},
    openAiEnable: {type: Boolean, default: false},
    groupList: {type: Array, default: () => []},
    currentGroupId: {type: Number, default: 0},
    enableDanmu: {type: Boolean, default: false},
    // 绑定 stock.vue 的 currentStockTradingPrice（含 stockCode/entryPrice/...）
    tradingPrice: {type: Object, default: () => ({})},
})

const emit = defineEmits([
    'update:show',
    'ai', 'set-cost', 'detail', 'notice', 'report',
    'set-group', 'remove-group', 'unfollow', 'danmu',
    'update:costPrice', 'update:longEntryPrice', 'update:longStopLossPrice', 'update:longTakeProfitPrice',
])

// 东方财富格式代码，供轻量K线组件使用
const emCode = computed(() => props.result ? toEastMoneyCode(props.result['股票代码']) : '')
const hasDepth = computed(() => props.result && props.result['买一报价'] > 0)
const activeTab = ref('fenshi')

function close() {
    emit('update:show', false)
}

// 5 档盘口数据
const orderLevels = computed(() => {
    if (!props.result) return []
    const r = props.result
    return [
        {label: '买一', price: r['买一报价'], vol: r['买一申报']},
        {label: '卖一', price: r['卖一报价'], vol: r['卖一申报']},
        {label: '买二', price: r['买二报价'], vol: r['买二申报']},
        {label: '卖二', price: r['卖二报价'], vol: r['卖二申报']},
        {label: '买三', price: r['买三报价'], vol: r['买三申报']},
        {label: '卖三', price: r['卖三报价'], vol: r['卖三申报']},
        {label: '买四', price: r['买四报价'], vol: r['买四申报']},
        {label: '卖四', price: r['卖四报价'], vol: r['卖四申报']},
        {label: '买五', price: r['买五报价'], vol: r['买五申报']},
        {label: '卖五', price: r['卖五报价'], vol: r['卖五申报']},
    ].filter(l => l.price > 0)
})
</script>

<template>
    <BottomSheet :show="show" :title="result ? result['股票名称'] : ''" height="90vh" @update:show="(v) => emit('update:show', v)">
        <div v-if="result" class="stock-detail">
            <!-- 头部：价格 + 涨跌 + 高低开收 -->
            <div class="stock-detail__header">
                <div class="stock-detail__title">
                    <span class="stock-detail__name">{{ result['股票名称'] }}</span>
                    <n-tag size="small" :bordered="false">{{ result['股票代码'] }}</n-tag>
                </div>
                <div class="stock-detail__price-line">
                    <span class="stock-detail__price" :class="'text-' + result.type">
                        <n-number-animation :duration="1000" :precision="2" :from="result['上次当前价格']"
                                            :to="Number(result['当前价格'])"/>
                    </span>
                    <span class="stock-detail__change" :class="'bg-' + result.type">
                        <n-number-animation :duration="1000" :precision="2" :from="0" :to="result.changePercent"/>%
                    </span>
                    <span class="stock-detail__after" v-if="result['盘前盘后'] > 0" :class="'text-' + result.type">
                        {{ result['盘前盘后'] }} {{ result['盘前盘后涨跌幅'] }}%
                    </span>
                </div>
                <div class="stock-detail__profit" v-if="result.costVolume > 0">
                    今日盈亏：
                    <n-text :type="result.type">
                        <n-number-animation :duration="1000" :precision="2" :from="0" :to="result.profitAmountToday"/>
                    </n-text>
                </div>
                <div class="stock-detail__ohl">
                    <span>最高 {{ result['今日最高价'] }} {{ result.highRate }}%</span>
                    <span>最低 {{ result['今日最低价'] }} {{ result.lowRate }}%</span>
                    <span>昨收 {{ result['昨日收盘价'] }}</span>
                    <span>今开 {{ result['今日开盘价'] }}</span>
                </div>
            </div>

            <!-- 4 个 Tab -->
            <n-tabs v-model:value="activeTab" type="line" animated size="small" display-directive="show" class="stock-detail__tabs">
                <n-tab-pane name="fenshi" tab="分时">
                    <FenshiChart
                        :key="'fenshi-' + result['股票代码']"
                        :code="result['股票代码']"
                        :name="result['股票名称']"
                        :dark-theme="darkTheme"
                        :chart-height="380"
                    />
                </n-tab-pane>
                <n-tab-pane name="kline" tab="K线">
                    <stock-lightweight-kline-chart
                        v-if="emCode"
                        :key="'lw-' + emCode"
                        :code="emCode"
                        :stock-name="result['股票名称']"
                        :dark-theme="darkTheme"
                        :chart-height="380"
                        :long-entry-price="tradingPrice.entryPrice"
                        :long-stop-loss-price="tradingPrice.stopLossPrice"
                        :long-take-profit-price="tradingPrice.takeProfitPrice"
                        :cost-price="tradingPrice.costPrice"
                        @update:longEntryPrice="(v) => emit('update:longEntryPrice', v)"
                        @update:longStopLossPrice="(v) => emit('update:longStopLossPrice', v)"
                        @update:longTakeProfitPrice="(v) => emit('update:longTakeProfitPrice', v)"
                        @update:costPrice="(v) => emit('update:costPrice', v)"
                    />
                    <n-empty v-else description="当前代码暂不支持K线图" style="padding: 40px 0" />
                </n-tab-pane>
                <n-tab-pane name="order" tab="盘口">
                    <div v-if="hasDepth" class="stock-detail__order">
                        <div v-for="l in orderLevels" :key="l.label" class="stock-detail__order-item">
                            <span class="stock-detail__order-label">{{ l.label }}</span>
                            <span class="stock-detail__order-price">{{ l.price }}</span>
                            <span class="stock-detail__order-vol">{{ l.vol }}</span>
                        </div>
                    </div>
                    <n-empty v-else description="暂无盘口数据" style="padding: 40px 0" />
                </n-tab-pane>
                <n-tab-pane name="daily" tab="日K">
                    <DailyKlineChart
                        :key="'daily-' + result['股票代码']"
                        :code="result['股票代码']"
                        :name="result['股票名称']"
                        :dark-theme="darkTheme"
                        :chart-height="380"
                    />
                </n-tab-pane>
                <n-tab-pane name="money" tab="资金">
                    <money-trend
                        :key="'money-' + result['股票代码']"
                        :code="result['股票代码']"
                        :name="result['股票名称']"
                        :days="360"
                        :dark-theme="darkTheme"
                        :chart-height="380"
                    />
                </n-tab-pane>
            </n-tabs>

            <!-- 底部动作栏 -->
            <div class="stock-detail__actions">
                <div class="stock-detail__actions-row">
                    <n-button size="small" type="warning" secondary @click="emit('set-cost', result)">成本</n-button>
                    <n-button size="small" type="primary" secondary v-if="openAiEnable" @click="emit('ai', result)">AI分析</n-button>
                    <n-dropdown trigger="click" :options="groupList" key-field="ID" label-field="name"
                                @select="(groupId) => emit('set-group', {groupId, result})">
                        <n-button size="small" secondary>设置分组</n-button>
                    </n-dropdown>
                    <n-button size="small" type="error" secondary @click="emit('unfollow', result)">取消关注</n-button>
                </div>
                <div class="stock-detail__actions-row stock-detail__actions-row--sub">
                    <n-button size="tiny" tertiary @click="emit('detail', result)">详情</n-button>
                    <n-button size="tiny" tertiary v-if="hasDepth" @click="emit('notice', result['股票代码'])">公告</n-button>
                    <n-button size="tiny" tertiary v-if="hasDepth" @click="emit('report', result['股票代码'])">研报</n-button>
                    <n-button size="tiny" tertiary v-if="currentGroupId > 0" @click="emit('remove-group', result)">移出分组</n-button>
                    <n-button size="tiny" tertiary v-if="enableDanmu" @click="emit('danmu', result['股票名称'])">弹幕</n-button>
                </div>
            </div>
        </div>
    </BottomSheet>
</template>

<style scoped>
.stock-detail {
    display: flex;
    flex-direction: column;
    height: 100%;
}

.stock-detail__header {
    padding: 4px 12px 12px;
}

.stock-detail__title {
    align-items: center;
    display: flex;
    gap: 8px;
}

.stock-detail__name {
    font-size: 18px;
    font-weight: 700;
}

.stock-detail__price-line {
    align-items: center;
    display: flex;
    gap: 10px;
    margin-top: 6px;
}

.stock-detail__price {
    font-size: 28px;
    font-weight: 800;
}

.stock-detail__change {
    border-radius: 4px;
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    padding: 2px 8px;
}

.stock-detail__after {
    font-size: 12px;
}

.stock-detail__profit {
    font-size: 13px;
    margin-top: 4px;
}

.stock-detail__ohl {
    color: var(--n-text-color-3, #999);
    display: flex;
    flex-wrap: wrap;
    font-size: 12px;
    gap: 6px 14px;
    margin-top: 8px;
}

.stock-detail__tabs {
    display: flex;
    flex-direction: column;
    flex: 1 1 auto;
    min-height: 0;
    overflow: hidden;
    padding: 0 8px;
}

.stock-detail__tabs :deep(.n-tabs-pane-wrapper) {
    flex: 1 1 auto;
    min-height: 0;
    overflow: auto;
}

.stock-detail__tabs :deep(.n-tab-pane) {
    box-sizing: border-box;
    min-height: 0;
    overflow: auto;
    padding: 8px 0 12px;
}

.stock-detail__order {
    display: grid;
    gap: 6px;
    grid-template-columns: 1fr 1fr;
    padding: 8px 12px;
}

.stock-detail__order-item {
    display: flex;
    font-size: 13px;
    gap: 6px;
    justify-content: space-between;
}

.stock-detail__order-label {
    color: var(--n-text-color-3, #999);
}

.stock-detail__order-price {
    color: var(--n-text-color, #333);
    font-weight: 600;
}

.stock-detail__order-vol {
    color: var(--n-text-color-3, #999);
}

.stock-detail__actions {
    border-top: 1px solid var(--n-border-color, #efeff5);
    padding: 10px 12px calc(var(--safe-bottom) + 8px);
}

.stock-detail__actions-row {
    display: flex;
    gap: 8px;
    justify-content: space-between;
}

.stock-detail__actions-row--sub {
    margin-top: 8px;
}

.text-success {
    color: #18a058;
}

.text-error {
    color: #d03050;
}

.text-default {
    color: var(--n-text-color, #333);
}

.bg-success {
    background: #18a058;
}

.bg-error {
    background: #d03050;
}

.bg-default {
    background: #909399;
}
</style>
