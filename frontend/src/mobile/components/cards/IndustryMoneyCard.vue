<script setup>
import { computed } from 'vue'
import PercentTag from '../widgets/PercentTag.vue'
import { formatMoney } from '../../composables/useFormat'

// 行业资金/证监会行业/概念板块 共用卡片（数据来自 GetIndustryMoneyRankSina）。
// 字段：name(板块名) avg_changeratio(涨跌幅,0~1) netamount(净流入,元) ratioamount(净流入率,0~1)
//       ts_name/ts_symbol/ts_changeratio(领涨股) ts_trade(领涨股价) ts_ratioamount(领涨股净流入率)
const props = defineProps({
  item: {
    type: Object,
    required: true
  }
})

const net = computed(() => Number(props.item.netamount))
// 净流入方向：红正绿负（国内习惯）。0/NaN 视为平。
const netClass = computed(() => {
  const n = net.value
  if (!Number.isFinite(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
})

// 比率字段 0~1 → 百分数
const avgChangePct = computed(() => Number(props.item.avg_changeratio) * 100)
const netRatePct = computed(() => Number(props.item.ratioamount) * 100)
const tsChangePct = computed(() => Number(props.item.ts_changeratio) * 100)
</script>

<template>
  <div class="money-card">
    <!-- 第一行：板块名 + 净流入大字 -->
    <div class="row-main">
      <span class="name">{{ item.name || '--' }}</span>
      <span class="net" :class="netClass">{{ formatMoney(net) }}</span>
    </div>

    <!-- 第二行：板块涨跌幅 / 净流入率 -->
    <div class="row-sub">
      <span class="tag-wrap"><i>涨跌</i><PercentTag :value="avgChangePct" size="small" /></span>
      <span class="tag-wrap"><i>净流入率</i><PercentTag :value="netRatePct" size="small" /></span>
    </div>

    <!-- 第三行：领涨股 -->
    <div v-if="item.ts_name" class="row-leader">
      <span class="leader-label">领涨</span>
      <span class="leader-name">{{ item.ts_name }}</span>
      <PercentTag :value="tsChangePct" size="small" />
    </div>
  </div>
</template>

<style scoped>
.money-card {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
  padding: var(--m-space-md);
}

.row-main {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: var(--m-space-md);
}

.name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.net {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
  flex-shrink: 0;
}

.net.m-rise {
  color: var(--m-color-rise);
}

.net.m-fall {
  color: var(--m-color-fall);
}

.net.m-flat {
  color: var(--m-color-gray);
}

.row-sub {
  display: flex;
  gap: var(--m-space-lg);
}

.tag-wrap {
  display: inline-flex;
  align-items: center;
  gap: var(--m-space-xs);
}

.tag-wrap i {
  font-style: normal;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.row-leader {
  display: flex;
  align-items: center;
  gap: var(--m-space-xs);
  font-size: var(--m-font-xs);
}

.leader-label {
  color: var(--m-text-tertiary);
  flex-shrink: 0;
}

.leader-name {
  color: var(--m-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}
</style>
