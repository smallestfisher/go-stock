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
    <div class="stock-mobile-list">
        <button
            v-for="result in rows()"
            :key="result['股票代码']"
            type="button"
            class="stock-mobile-row"
            @click="emit('select', result)"
        >
            <!-- 左：名称 + 代码 + 盈亏 -->
            <div class="stock-mobile-row__main">
                <div class="stock-mobile-row__name">
                    {{ result['股票名称'] }}
                    <span class="stock-mobile-row__code">{{ result['股票代码'] }}</span>
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
    padding: 4px 8px calc(var(--safe-bottom) + 12px);
}

.stock-mobile-row {
    align-items: center;
    appearance: none;
    background: var(--n-color, #fff);
    border: 0;
    border-bottom: 1px solid var(--n-border-color, #efeff5);
    color: var(--n-text-color, #333);
    display: flex;
    font: inherit;
    gap: 10px;
    padding: 12px 6px;
    text-align: left;
    width: 100%;
}

.stock-mobile-row:active {
    background: rgba(24, 160, 88, 0.06);
}

.stock-mobile-row__main {
    display: flex;
    flex: 1 1 0;
    flex-direction: column;
    gap: 4px;
    min-width: 0;
}

.stock-mobile-row__name {
    font-size: 16px;
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.stock-mobile-row__code {
    color: var(--n-text-color-3, #999);
    font-size: 12px;
    font-weight: 400;
    margin-left: 6px;
}

.stock-mobile-row__profit {
    font-size: 12px;
}

.stock-mobile-row__spark {
    flex: 0 0 80px;
    height: 36px;
    overflow: hidden;
}

.stock-mobile-row__price {
    align-items: flex-end;
    display: flex;
    flex: 0 0 auto;
    flex-direction: column;
    gap: 3px;
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
    min-width: 64px;
    padding: 3px 8px;
    text-align: center;
}

.stock-mobile-row__after {
    font-size: 11px;
}

.stock-mobile-list__empty {
    color: var(--n-text-color-3, #999);
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
    color: var(--n-text-color, #333);
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
</style>
