<script setup>
import { computed } from 'vue'
import MIcon from '../base/MIcon.vue'
import PercentTag from '../widgets/PercentTag.vue'
import PriceTag from '../widgets/PriceTag.vue'

// 单条热门股票卡片。对齐桌面端 HotStockList.vue 的字段：
// name/code/percent(涨跌幅)/current(股价)/value(热度)/increment(热度变化)/rank_change(排名变化)。
const props = defineProps({
  item: {
    type: Object,
    required: true,
  },
  // 榜单名次（从 1 开始）
  rank: {
    type: Number,
    default: 0,
  },
})

// 热度数值偏大，做千分位更易读。
const heatText = computed(() => {
  const v = Number(props.item.value)
  if (!Number.isFinite(v)) return props.item.value ?? '-'
  return v.toLocaleString('en-US')
})

const rankChange = computed(() => Number(props.item.rank_change) || 0)

// 涨红跌绿（与桌面端一致）；0 为平。
function signed(n) {
  const v = Number(n)
  if (!Number.isFinite(v) || v === 0) return '0'
  return v > 0 ? `+${v}` : `${v}`
}
function deltaClass(n) {
  const v = Number(n)
  if (!Number.isFinite(v) || v === 0) return 'hot-delta--flat'
  return v > 0 ? 'hot-delta--up' : 'hot-delta--down'
}
</script>

<template>
  <div class="hot-stock">
    <!-- 名次 -->
    <span class="hot-rank" :class="{ 'hot-rank--top': rank <= 3 }">{{ rank }}</span>

    <!-- 名称 + 代码 -->
    <div class="hot-main">
      <span class="hot-name">{{ item.name }}</span>
      <span class="hot-code">{{ item.code }}</span>
    </div>

    <!-- 价格 + 涨跌幅 -->
    <div class="hot-quote">
      <PriceTag :price="item.current" :change="item.percent" size="small" bold />
      <PercentTag :value="item.percent" size="small" />
    </div>

    <!-- 热度 + 变化 -->
    <div class="hot-heat">
      <span class="hot-heat__value"><MIcon name="fire" :size="14" /> {{ heatText }}</span>
      <div class="hot-heat__delta">
        <span class="hot-delta" :class="deltaClass(item.increment)">
          热{{ signed(item.increment) }}
        </span>
        <span v-if="rankChange !== 0" class="hot-delta" :class="deltaClass(rankChange)">
          名{{ signed(rankChange) }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hot-stock {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  padding: var(--m-space-md) var(--m-space-lg);
}

.hot-rank {
  flex-shrink: 0;
  width: 20px;
  text-align: center;
  font-size: var(--m-font-sm);
  font-variant-numeric: tabular-nums;
  color: var(--m-text-tertiary);
  font-weight: var(--m-font-weight-medium);
}

.hot-rank--top {
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-bold);
}

.hot-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.hot-name {
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hot-code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.hot-quote {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1px;
  min-width: 64px;
}

.hot-heat {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1px;
  min-width: 72px;
}

.hot-heat__value {
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
  font-variant-numeric: tabular-nums;
}

.hot-heat__delta {
  display: flex;
  gap: var(--m-space-xs);
}

.hot-delta {
  font-size: var(--m-font-xs);
  font-variant-numeric: tabular-nums;
}

.hot-delta--up {
  color: var(--m-color-rise);
}

.hot-delta--down {
  color: var(--m-color-fall);
}

.hot-delta--flat {
  color: var(--m-text-tertiary);
}
</style>
