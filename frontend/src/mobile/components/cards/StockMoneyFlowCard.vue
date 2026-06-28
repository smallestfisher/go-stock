<script setup>
import { computed } from 'vue'
import PercentTag from '../widgets/PercentTag.vue'
import { formatMoney } from '../../composables/useFormat'

// 个股资金流向卡片（数据来自 GetMoneyRankSina）。
// props.metric = { label, kind:'money'|'rate', get:(item)=>number }，由页面按当前 tab 传入，
// 卡片据此高亮当前排名维度（净流入额/流出额/净流入率/主力·散户 各指标）。
const props = defineProps({
  item: { type: Object, required: true },
  metric: { type: Object, required: true },
})

// 现价
const trade = computed(() => Number(props.item.trade))
// 涨跌幅（0~1 → 百分数）
const changePct = computed(() => Number(props.item.changeratio) * 100)
const changeDir = computed(() => {
  const n = changePct.value
  if (!Number.isFinite(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
})

// 当前 tab 高亮指标值
const metricVal = computed(() => {
  const n = Number(props.metric?.get?.(props.item))
  return Number.isFinite(n) ? n : NaN
})

// 成交额、换手率（次要）
const amount = computed(() => Number(props.item.amount))
const turnoverPct = computed(() => Number(props.item.turnover) / 100)
</script>

<template>
  <div class="flow-card">
    <!-- 第一行：股票名+代码 / 现价+涨跌幅 -->
    <div class="row-top">
      <div class="stock-left">
        <span class="stock-name">{{ item.name || '--' }}</span>
        <span class="stock-code">{{ item.symbol }}</span>
      </div>
      <div class="price-group">
        <span class="trade" :class="changeDir">{{ Number.isFinite(trade) ? trade.toFixed(2) : '--' }}</span>
        <PercentTag :value="changePct" size="small" />
      </div>
    </div>

    <!-- 第二行：当前 tab 高亮指标（随排名维度变化） -->
    <div class="row-metric">
      <span class="metric-label">{{ metric.label }}</span>
      <span v-if="metric.kind === 'money'" class="metric-val" :class="{
        'm-rise': metricVal > 0,
        'm-fall': metricVal < 0,
        'm-flat': !Number.isFinite(metricVal) || metricVal === 0
      }">
        {{ Number.isFinite(metricVal) ? formatMoney(metricVal) : '--' }}
      </span>
      <PercentTag v-else :value="metricVal" size="medium" />
    </div>

    <!-- 第三行：成交额 + 换手率 -->
    <div class="row-sub">
      <span class="sub-item">成交额 {{ formatMoney(amount) }}</span>
      <span class="sub-item">换手 {{ Number.isFinite(turnoverPct) ? turnoverPct.toFixed(2) + '%' : '--' }}</span>
    </div>
  </div>
</template>

<style scoped>
.flow-card {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
  padding: var(--m-space-md);
}

.row-top {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: var(--m-space-md);
}

.stock-left {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  min-width: 0;
}

.stock-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stock-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  flex-shrink: 0;
}

.price-group {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  flex-shrink: 0;
}

.trade {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
  color: var(--m-text-primary);
}

.trade.m-rise {
  color: var(--m-color-rise);
}

.trade.m-fall {
  color: var(--m-color-fall);
}

/* 高亮指标行 */
.row-metric {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: var(--m-space-md);
  padding: var(--m-space-xs) var(--m-space-sm);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
}

.metric-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.metric-val {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
  color: var(--m-color-gray);
}

.metric-val.m-rise {
  color: var(--m-color-rise);
}

.metric-val.m-fall {
  color: var(--m-color-fall);
}

.metric-val.m-flat {
  color: var(--m-color-gray);
}

.row-sub {
  display: flex;
  gap: var(--m-space-lg);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.sub-item {
  white-space: nowrap;
}
</style>
