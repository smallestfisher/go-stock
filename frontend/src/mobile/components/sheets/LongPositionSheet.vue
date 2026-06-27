<script setup>
import { ref, computed, watch } from 'vue'
import MSheet from '../base/MSheet.vue'

// 多单仓位计算器（对齐桌面端 longPositionStats）。
// 输入开仓/止损/止盈/成本，计算风险幅度、目标幅度、盈亏比、风险%与目标%。
const props = defineProps({
  show: { type: Boolean, default: false },
  // 最新收盘价，用于一键填充开仓
  latestClose: { type: [Number, String], default: '' },
})

const emit = defineEmits(['update:show'])

const entryStr = ref('')
const stopStr = ref('')
const takeProfitStr = ref('')
const costStr = ref('')
// 资金（用于算建议仓位）
const capitalStr = ref('')

function toNum(s) {
  const n = Number(s)
  return Number.isFinite(n) ? n : NaN
}

// 填充最新收盘价为开仓价
function fillLatestClose() {
  if (props.latestClose !== '' && props.latestClose != null) {
    entryStr.value = String(props.latestClose)
  }
}

// 清空
function clearAll() {
  entryStr.value = ''
  stopStr.value = ''
  takeProfitStr.value = ''
  costStr.value = ''
  capitalStr.value = ''
}

// 计算结果（对齐桌面端 longPositionStats）
const stats = computed(() => {
  const entry = toNum(entryStr.value)
  const stop = toNum(stopStr.value)
  const tp = toNum(takeProfitStr.value)
  const capital = toNum(capitalStr.value)
  if (!Number.isFinite(entry)) return null

  const risk = Number.isFinite(stop) ? entry - stop : NaN
  const reward = Number.isFinite(tp) ? tp - entry : NaN
  const rr = Number.isFinite(risk) && risk > 0 && Number.isFinite(reward) ? reward / risk : NaN
  const riskPct = Number.isFinite(risk) && entry !== 0 ? (risk / entry) * 100 : NaN
  const rewardPct = Number.isFinite(reward) && entry !== 0 ? (reward / entry) * 100 : NaN

  // 建议仓位：按固定比例风险（单笔亏 2% 本金）反推
  let suggestLots = NaN
  if (Number.isFinite(capital) && capital > 0 && Number.isFinite(risk) && risk > 0) {
    const riskMoney = capital * 0.02
    const lots = Math.floor((riskMoney / risk) / 100) * 100 // A 股按手(100股)取整
    suggestLots = lots > 0 ? lots : 0
  }

  return {
    riskPts: risk,
    rewardPts: reward,
    riskRr: rr,
    riskPct,
    rewardPct,
    suggestLots,
  }
})

// 盈亏比文字
const rrText = computed(() => {
  const s = stats.value
  if (!s || !Number.isFinite(s.riskRr) || s.riskRr <= 0) return '--'
  return `1 : ${s.riskRr.toFixed(2)}`
})

function close() { emit('update:show', false) }
</script>

<template>
  <MSheet
    :show="show"
    title="多单仓位计算"
    height="auto"
    @update:show="close"
  >
    <div class="long-pos">
      <!-- 输入区 -->
      <div class="lp-inputs">
        <label class="lp-field">
          <span>开仓价</span>
          <input v-model="entryStr" type="number" inputmode="decimal" placeholder="开仓">
          <button class="lp-fill" @click="fillLatestClose">现价</button>
        </label>
        <label class="lp-field">
          <span>止损价</span>
          <input v-model="stopStr" type="number" inputmode="decimal" placeholder="止损">
        </label>
        <label class="lp-field">
          <span>止盈价</span>
          <input v-model="takeProfitStr" type="number" inputmode="decimal" placeholder="止盈">
        </label>
        <label class="lp-field">
          <span>成本价</span>
          <input v-model="costStr" type="number" inputmode="decimal" placeholder="成本(可选)">
        </label>
        <label class="lp-field">
          <span>本金(元)</span>
          <input v-model="capitalStr" type="number" inputmode="decimal" placeholder="用于建议仓位">
        </label>
      </div>

      <!-- 计算结果 -->
      <div v-if="stats" class="lp-result">
        <div class="lp-rr">
          <span class="lp-rr-label">盈亏比</span>
          <span class="lp-rr-val" :class="{ 'lp-rr--good': stats.riskRr >= 2 }">{{ rrText }}</span>
        </div>
        <div class="lp-grid">
          <div class="lp-item">
            <span class="lp-item-k">风险幅度</span>
            <span class="lp-item-v c-fall">{{ Number.isFinite(stats.riskPts) ? stats.riskPts.toFixed(2) : '--' }}</span>
          </div>
          <div class="lp-item">
            <span class="lp-item-k">目标幅度</span>
            <span class="lp-item-v c-rise">{{ Number.isFinite(stats.rewardPts) ? stats.rewardPts.toFixed(2) : '--' }}</span>
          </div>
          <div class="lp-item">
            <span class="lp-item-k">风险%</span>
            <span class="lp-item-v c-fall">{{ Number.isFinite(stats.riskPct) ? stats.riskPct.toFixed(2) + '%' : '--' }}</span>
          </div>
          <div class="lp-item">
            <span class="lp-item-k">目标%</span>
            <span class="lp-item-v c-rise">{{ Number.isFinite(stats.rewardPct) ? stats.rewardPct.toFixed(2) + '%' : '--' }}</span>
          </div>
          <div class="lp-item">
            <span class="lp-item-k">建议仓位</span>
            <span class="lp-item-v">{{ Number.isFinite(stats.suggestLots) ? stats.suggestLots + ' 股' : '--' }}</span>
          </div>
        </div>
        <p class="lp-tip">按单笔亏损本金 2% 反推；止损需低于开仓，止盈需高于开仓。</p>
      </div>
      <div v-else class="lp-empty">输入开仓价后显示计算结果</div>

      <button class="lp-clear" @click="clearAll">清空</button>
    </div>
  </MSheet>
</template>

<style scoped>
.long-pos {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.lp-inputs {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.lp-field {
  display: grid;
  grid-template-columns: 64px 1fr auto;
  align-items: center;
  gap: var(--m-space-sm);
}

.lp-field > span {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.lp-field input {
  height: 40px;
  padding: 0 var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-md);
  font-variant-numeric: tabular-nums;
  outline: none;
}

.lp-field input:focus {
  border-color: var(--m-color-rise);
}

.lp-fill {
  height: 32px;
  padding: 0 var(--m-space-sm);
  border: 1px solid var(--m-color-rise);
  border-radius: var(--m-radius-sm);
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
  font-size: var(--m-font-xs);
  white-space: nowrap;
}

.lp-fill:active { opacity: 0.6; }

/* 结果 */
.lp-result {
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
}

.lp-rr {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  padding-bottom: var(--m-space-sm);
  border-bottom: 1px solid var(--m-divider-color);
  margin-bottom: var(--m-space-sm);
}

.lp-rr-label {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.lp-rr-val {
  font-size: var(--m-font-xl);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
}

.lp-rr--good { color: var(--m-color-rise); }

.lp-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--m-space-sm) var(--m-space-md);
}

.lp-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.lp-item-k {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.lp-item-v {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
  color: var(--m-text-primary);
}

.c-rise { color: var(--m-color-rise); }
.c-fall { color: var(--m-color-fall); }

.lp-tip {
  margin-top: var(--m-space-sm);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  line-height: 1.5;
}

.lp-empty {
  text-align: center;
  padding: var(--m-space-lg) 0;
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

.lp-clear {
  height: 40px;
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-card);
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
}

.lp-clear:active { opacity: 0.6; }
</style>
