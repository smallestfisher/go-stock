<script setup>
import { computed } from 'vue'

// 指标信号汇总（对齐桌面端「指标信号汇总 共 N 项」）。
// 输入 summarizeSignals 的结果，渲染：占比条 + 四类计数凡例 + 指标标签云。
const props = defineProps({
  // summarizeSignals(...) 的返回值；null 表示无数据
  summary: {
    type: Object,
    default: null,
  },
})

// 看多=红(涨)、看空=绿(跌)、震荡=黄、中性=灰（与桌面端及涨红跌绿一致）
const SIGNAL_META = {
  bullish: { label: '看多', cls: 'sig--bullish' },
  bearish: { label: '看空', cls: 'sig--bearish' },
  oscillating: { label: '震荡', cls: 'sig--osc' },
  neutral: { label: '中性', cls: 'sig--neutral' },
}

const legend = computed(() => {
  const s = props.summary
  if (!s) return []
  return [
    { key: 'bullish', label: '看多', count: s.bullish, pct: s.bullishPct },
    { key: 'bearish', label: '看空', count: s.bearish, pct: s.bearishPct },
    { key: 'oscillating', label: '震荡', count: s.oscillating, pct: s.oscillatingPct },
    { key: 'neutral', label: '中性', count: s.neutral, pct: s.neutralPct },
  ]
})

function sigClass(signal) {
  return SIGNAL_META[signal]?.cls || 'sig--neutral'
}
</script>

<template>
  <div v-if="summary" class="signal-summary">
    <div class="ss-head">
      <span class="ss-title">指标信号汇总</span>
      <span class="ss-total">共 {{ summary.total }} 项</span>
    </div>

    <!-- 占比条 -->
    <div class="ss-bar">
      <div class="ss-bar__seg sig-bg--bullish" :style="{ width: summary.bullishPct + '%' }" />
      <div class="ss-bar__seg sig-bg--bearish" :style="{ width: summary.bearishPct + '%' }" />
      <div class="ss-bar__seg sig-bg--osc" :style="{ width: summary.oscillatingPct + '%' }" />
      <div class="ss-bar__seg sig-bg--neutral" :style="{ width: summary.neutralPct + '%' }" />
    </div>

    <!-- 凡例 -->
    <div class="ss-legend">
      <span v-for="l in legend" :key="l.key" class="ss-legend__item">
        <span class="ss-dot" :class="sigClass(l.key)" />
        {{ l.label }} {{ l.count }} ({{ l.pct }}%)
      </span>
    </div>

    <!-- 指标标签云 -->
    <div class="ss-tags">
      <span
        v-for="(s, i) in summary.signals"
        :key="s.name + i"
        class="ss-tag"
        :class="sigClass(s.signal)"
      >{{ s.name }}</span>
    </div>
  </div>
</template>

<style scoped>
.signal-summary {
  padding: var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
}

.ss-head {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-sm);
}

.ss-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
}

.ss-total {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

/* 占比条 */
.ss-bar {
  display: flex;
  height: 8px;
  border-radius: var(--m-radius-full);
  overflow: hidden;
  background: var(--m-bg-primary);
}

.ss-bar__seg {
  height: 100%;
  transition: width var(--m-duration-normal) var(--m-ease-out);
}

/* 凡例 */
.ss-legend {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-md);
  margin-top: var(--m-space-sm);
  font-size: var(--m-font-xs);
  color: var(--m-text-secondary);
}

.ss-legend__item {
  display: inline-flex;
  align-items: center;
  gap: var(--m-space-xs);
  font-variant-numeric: tabular-nums;
}

.ss-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

/* 标签云 */
.ss-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
  margin-top: var(--m-space-md);
}

.ss-tag {
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
  line-height: 1.7;
  background: var(--m-bg-primary);
}

/* 信号配色：看多红 / 看空绿 / 震荡黄 / 中性灰 */
.sig--bullish { color: var(--m-color-rise); }
.sig--bearish { color: var(--m-color-fall); }
.sig--osc { color: #e6a23c; }
.sig--neutral { color: var(--m-text-tertiary); }

.sig-bg--bullish { background: var(--m-color-rise); }
.sig-bg--bearish { background: var(--m-color-fall); }
.sig-bg--osc { background: #e6a23c; }
.sig-bg--neutral { background: var(--m-text-disabled); }

.ss-dot.sig--bullish { background: var(--m-color-rise); }
.ss-dot.sig--bearish { background: var(--m-color-fall); }
.ss-dot.sig--osc { background: #e6a23c; }
.ss-dot.sig--neutral { background: var(--m-text-disabled); }

/* 标签按信号着色：浅底 + 同色字 */
.ss-tag.sig--bullish { background: var(--m-color-rise-light); color: var(--m-color-rise); }
.ss-tag.sig--bearish { background: var(--m-color-fall-light); color: var(--m-color-fall); }
.ss-tag.sig--osc { background: rgba(230, 162, 60, 0.14); color: #c77800; }
.ss-tag.sig--neutral { background: var(--m-bg-primary); color: var(--m-text-tertiary); }
</style>
