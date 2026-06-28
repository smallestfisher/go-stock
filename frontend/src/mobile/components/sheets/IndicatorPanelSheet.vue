<script setup>
import MSheet from '../base/MSheet.vue'

// 指标开关面板（对齐桌面端侧边栏 5 类分组，完整覆盖桌面端 43 项指标）。
//   主图叠加：绘制在 K 线主图（价格域）
//   副图独立：绘制在独立副图区（自身数值域）
//   主图（量类）：VWAP 等基于价格的量价指标也叠在主图
const props = defineProps({
  show: { type: Boolean, default: false },
  indicators: { type: Array, default: () => ['MA', 'VOL'] },
})

const emit = defineEmits(['update:show', 'update:indicators'])

// 5 类分组，与桌面端侧边栏完全一致。
// sub: main=主图叠加, sub=独立副图
const groups = [
  {
    name: '趋势', icon: '📈', color: 'rise',
    items: [
      { code: 'MA', label: 'MA', sub: '主图' },
      { code: 'EMA', label: 'EMA', sub: '主图' },
      { code: 'BOLL', label: 'BOLL', sub: '主图' },
      { code: 'SAR', label: 'SAR', sub: '主图' },
      { code: 'SUPERTREND', label: 'STrend', sub: '主图' },
      { code: 'DONCHIAN', label: 'Donch', sub: '主图' },
      { code: 'ICHIMOKU', label: 'Ichi', sub: '主图' },
      { code: 'KAMA', label: 'KAMA', sub: '主图' },
      { code: 'HULL', label: 'Hull', sub: '主图' },
      { code: 'DEMA', label: 'DEMA', sub: '主图' },
      { code: 'TEMA', label: 'TEMA', sub: '主图' },
      { code: 'AROON', label: 'Aroon', sub: '副图' },
      { code: 'ALLIGATOR', label: 'Gator', sub: '主图' },
    ],
  },
  {
    name: '动量', icon: '💫', color: 'momentum',
    items: [
      { code: 'MACD', label: 'MACD', sub: '副图' },
      { code: 'KDJ', label: 'KDJ', sub: '副图' },
      { code: 'RSI', label: 'RSI', sub: '副图' },
      { code: 'CCI', label: 'CCI', sub: '副图' },
      { code: 'WR', label: 'W%R', sub: '副图' },
      { code: 'STOCHRSI', label: 'StochRSI', sub: '副图' },
      { code: 'CMO', label: 'CMO', sub: '副图' },
      { code: 'TRIX', label: 'TRIX', sub: '副图' },
      { code: 'ROC', label: 'ROC', sub: '副图' },
      { code: 'COPPOCK', label: 'Coppock', sub: '副图' },
      { code: 'SMI', label: 'SMI', sub: '副图' },
      { code: 'AO', label: 'AO', sub: '副图' },
    ],
  },
  {
    name: '量价', icon: '📊', color: 'vol',
    items: [
      { code: 'VOL', label: '成交量', sub: '副图' },
      { code: 'OBV', label: 'OBV', sub: '副图' },
      { code: 'MFI', label: 'MFI', sub: '副图' },
      { code: 'CMF', label: 'CMF', sub: '副图' },
      { code: 'AD', label: 'A/D', sub: '副图' },
      { code: 'FI', label: 'FI', sub: '副图' },
      { code: 'CHAIKINOSC', label: 'Chaikin', sub: '副图' },
      { code: 'VWAP', label: 'VWAP', sub: '主图' },
      { code: 'VWAPBANDS', label: 'VWAP带', sub: '主图' },
    ],
  },
  {
    name: '波动', icon: '🎢', color: 'wave',
    items: [
      { code: 'ATR', label: 'ATR', sub: '副图' },
      { code: 'KELTNER', label: 'Kelt', sub: '主图' },
      { code: 'TTM', label: 'TTM', sub: '副图' },
      { code: 'ZIGZAG', label: 'ZigZag', sub: '主图' },
      { code: 'AVGAMP', label: '均幅', sub: '副图' },
      { code: 'MASSINDEX', label: 'Mass', sub: '副图' },
      { code: 'ULCER', label: 'Ulcer', sub: '副图' },
      { code: 'SATS', label: 'SATS', sub: '副图' },
    ],
  },
  {
    name: '强度', icon: '📏', color: 'strength',
    items: [
      { code: 'ADX', label: 'ADX/DMI', sub: '副图' },
      { code: 'CHOP', label: 'CHOP', sub: '副图' },
      { code: 'ELDERRAY', label: 'ElderRay', sub: '副图' },
      { code: 'PIVOT', label: 'Pivot', sub: '主图' },
      { code: 'SIGNALRATIO', label: '信号比例', sub: '副图' },
    ],
  },
]

function isActive(code) {
  return props.indicators.includes(code)
}

function toggle(code) {
  const set = new Set(props.indicators)
  if (set.has(code)) set.delete(code)
  else set.add(code)
  emit('update:indicators', Array.from(set))
}

function close() { emit('update:show', false) }
</script>

<template>
  <MSheet
    :show="show"
    title="指标叠加"
    height="75vh"
    @update:show="close"
  >
    <div class="ind-panel">
      <div v-for="g in groups" :key="g.name" class="ind-group">
        <div class="ind-group__head" :class="`ind-group__head--${g.color}`">
          <span class="ind-group__icon">{{ g.icon }}</span>
          <span class="ind-group__name">{{ g.name }}</span>
        </div>
        <div class="ind-group__items">
          <button
            v-for="it in g.items"
            :key="it.code"
            type="button"
            class="ind-chip"
            :class="{ 'ind-chip--active': isActive(it.code) }"
            @click="toggle(it.code)"
          >
            <span class="ind-chip__label">{{ it.label }}</span>
          </button>
        </div>
      </div>

      <p class="ind-tip">主图指标叠加在 K 线上；副图指标显示在独立区域。勾选越多副图，主图区域越小。</p>
    </div>
  </MSheet>
</template>

<style scoped>
.ind-panel {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.ind-group {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

/* 分组标题：紧凑，左侧色条 */
.ind-group__head {
  display: inline-flex;
  align-items: center;
  gap: var(--m-space-xs);
  align-self: flex-start;
  padding: 2px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
  font-weight: var(--m-font-weight-bold);
  border-left: 3px solid;
}

.ind-group__head--rise { color: #ef4444; background: rgba(239,68,68,0.08); border-color: #ef4444; }
.ind-group__head--momentum { color: #3b82f6; background: rgba(59,130,246,0.08); border-color: #3b82f6; }
.ind-group__head--vol { color: #10b981; background: rgba(16,185,129,0.08); border-color: #10b981; }
.ind-group__head--wave { color: #f59e0b; background: rgba(245,158,11,0.08); border-color: #f59e0b; }
.ind-group__head--strength { color: #8b5cf6; background: rgba(139,92,246,0.08); border-color: #8b5cf6; }

/* 指标按钮：紧凑 chip，统一字体，不再撑满高度 */
.ind-group__items {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
}

.ind-chip {
  display: inline-flex;
  align-items: center;
  height: 28px;
  padding: 0 var(--m-space-sm);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-xs);
  font-weight: var(--m-font-weight-normal);
  line-height: 1;
  white-space: nowrap;
}

.ind-chip--active {
  background: var(--m-color-rise);
  border-color: var(--m-color-rise);
  color: #fff;
  font-weight: var(--m-font-weight-medium);
}

.ind-chip:active { transform: scale(0.95); }

.ind-tip {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  line-height: 1.5;
  text-align: center;
}
</style>
