<script setup>
import {defineAsyncComponent} from 'vue'

const StockSparkLine = defineAsyncComponent(() => import('../stockSparkLine.vue'))

const props = defineProps({
    // keyed object：全部传 sortedResults，分组传 groupResults
    results: {type: Object, default: () => ({})},
    darkTheme: {type: Boolean, default: false},
})

const emit = defineEmits(['select'])

function rows() {
    return Object.values(props.results || {})
}
</script>

<template>
    <div class="stock-mobile-list" :class="{'stock-mobile-list--dark': darkTheme}">
        <button
            v-for="result in rows()"
            :key="result['股票代码']"
            type="button"
            class="stock-mobile-row"
            :class="'stock-mobile-row--' + result.type"
            @click="emit('select', result)"
        >
            <!-- 左：名称 + 代码 + 盈亏 -->
            <div class="stock-mobile-row__main">
                <div class="stock-mobile-row__name">
                    {{ result['股票名称'] || '--' }}
                </div>
                <div class="stock-mobile-row__meta">
                    <span class="stock-mobile-row__code">{{ result['股票代码'] }}</span>
                    <span v-if="result['时间']" class="stock-mobile-row__time">{{ result['时间'] }}</span>
                </div>
                <div class="stock-mobile-row__profit" v-if="result.costVolume > 0" :class="'text-' + result.type">
                  今日
                  <n-number-animation :duration="1000" :precision="2" :from="0" :to="result.profitAmountToday"/>
                </div>
            </div>
            <!-- 中：迷你走势 -->
            <div class="stock-mobile-row__spark">
                <stock-spark-line
                    :id-suffix="'mobile-' + result['股票代码']"
                    :last-price="Number(result['当前价格'])"
                    :open-price="Number(result['昨日收盘价'])"
                    :stock-code="result['股票代码']"
                    :stock-name="result['股票名称']"
                    :dark-theme="darkTheme"
                />
            </div>
            <!-- 右：价格 + 涨跌幅 -->
            <div class="stock-mobile-row__price">
                <div class="stock-mobile-row__now" :class="'text-' + result.type">
                    <n-number-animation :duration="1000" :precision="2" :from="result['上次当前价格']"
                                        :to="Number(result['当前价格'])"/>
                </div>
                <div class="stock-mobile-row__change" :class="'bg-' + result.type">
                    <n-number-animation :duration="1000" :precision="2" :from="0" :to="result.changePercent"/>%
                </div>
                <div class="stock-mobile-row__after" v-if="result['盘前盘后'] > 0" :class="'text-' + result.type">
                    {{ result['盘前盘后'] }} {{ result['盘前盘后涨跌幅'] }}%
                </div>
            </div>
        </button>
        <div v-if="rows().length === 0" class="stock-mobile-list__empty">暂无自选股票</div>
    </div>
</template>

<style scoped>
.stock-mobile-list {
    background: #f6f7f9;
    padding: 6px 8px calc(var(--safe-bottom) + 12px);
}

.stock-mobile-row {
    align-items: center;
    appearance: none;
    background: #ffffff;
    border: 1px solid #edf0f5;
    border-left: 4px solid #d0d5dd;
    border-radius: 8px;
    box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
    color: #1f2329;
    display: flex;
    font: inherit;
    gap: 10px;
    margin-bottom: 8px;
    min-height: 72px;
    padding: 10px 10px 10px 9px;
    text-align: left;
    width: 100%;
}

.stock-mobile-row:active {
    background: #f8fafc;
}

.stock-mobile-row--error {
    border-left-color: #d03050;
}

.stock-mobile-row--success {
    border-left-color: #18a058;
}

.stock-mobile-row--default {
    border-left-color: #909399;
}

.stock-mobile-row__main {
    display: flex;
    flex: 1 1 96px;
    flex-direction: column;
    gap: 3px;
    min-width: 0;
}

.stock-mobile-row__name {
    font-size: 16px;
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.stock-mobile-row__meta {
    align-items: center;
    color: #667085;
    display: flex;
    font-size: 12px;
    gap: 6px;
    min-width: 0;
}

.stock-mobile-row__code {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.stock-mobile-row__time {
    flex: 0 0 auto;
}

.stock-mobile-row__profit {
    font-size: 12px;
    line-height: 1.2;
}

.stock-mobile-row__spark {
    flex: 0 0 86px;
    height: 34px;
    overflow: hidden;
}

.stock-mobile-row__price {
    align-items: flex-end;
    display: flex;
    flex: 0 0 78px;
    flex-direction: column;
    gap: 3px;
    min-width: 0;
}

.stock-mobile-row__now {
    font-size: 17px;
    font-weight: 700;
}

.stock-mobile-row__change {
    border-radius: 4px;
    color: #fff;
    font-size: 14px;
    font-weight: 700;
    min-width: 70px;
    padding: 3px 8px;
    text-align: center;
}

.stock-mobile-row__after {
    font-size: 11px;
}

.stock-mobile-list__empty {
    color: #98a2b3;
    padding: 40px 0;
    text-align: center;
}

/* 涨跌色：result.type 为 success(涨)/error(跌)/default(平)。
   色块用深一档实色，保证白字在移动端小屏上有足够对比度。 */
.text-success {
    color: #0f7a43;
}

.text-error {
    color: #b91c4b;
}

.text-default {
    color: #344054;
}

.bg-success {
    background: #0f7a43;
}

.bg-error {
    background: #b91c4b;
}

.bg-default {
    background: #6b7280;
}

.stock-mobile-list--dark {
    background: #101014;
}

.stock-mobile-list--dark .stock-mobile-row {
    background: #18181c;
    border-color: #2f2f35;
    box-shadow: none;
    color: #f2f4f7;
}

.stock-mobile-list--dark .stock-mobile-row:active {
    background: #202028;
}

.stock-mobile-list--dark .stock-mobile-row__meta,
.stock-mobile-list--dark .stock-mobile-list__empty {
    color: #98a2b3;
}

.stock-mobile-list--dark .text-default {
    color: #f2f4f7;
}
</style>
