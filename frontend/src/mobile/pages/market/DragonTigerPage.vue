<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import { LongTigerRank } from '../../../api/app'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import { formatMoney } from '../../composables/useFormat'

// 对齐桌面端 LongTigerRankList.vue：LongTigerRank(date) 返回当日龙虎榜。
// 当天数据通常收盘后约1小时更新；某日无数据时自动往前找最近有数据的一天（最多回退7天）。
const MAX_BACKTRACK = 7

const todayStr = (() => {
  const d = new Date()
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
})()

const list = ref([])
const loading = ref(false)
const searchDate = ref(todayStr) // 当前查看的日期（可能被自动回退修改）

function pickArray(res) {
  return (Array.isArray(res) ? res : []).filter(it => it && typeof it === 'object')
}

// SECUCODE "000001.SZ" → "sz000001"（行情代码，便于后续点击查K线）
function marketCode(item) {
  try {
    const [code, market] = String(item.SECUCODE).split('.')
    return (market || '').toLowerCase() + (code || '')
  } catch {
    return item.SECUCODE || ''
  }
}

// 净买额方向：红正绿负
function netClass(v) {
  const n = Number(v)
  if (!Number.isFinite(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
}

// 偏移日期：days 为正往后、负往前，返回 YYYY-MM-DD
function shiftDate(dateStr, days) {
  const d = new Date(dateStr)
  d.setDate(d.getDate() + days)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 取数（带空数据回退）。target 为目标日期；back 记录已回退次数。
async function fetchDate(target, back = 0) {
  const res = pickArray(await LongTigerRank(target))
  if (res.length) {
    list.value = res
    searchDate.value = target
    return true
  }
  // 空数据：未超回退上限则往前一天再试
  if (back < MAX_BACKTRACK) {
    return fetchDate(shiftDate(target, -1), back + 1)
  }
  list.value = []
  searchDate.value = target
  return false
}

async function loadData() {
  loading.value = true
  try {
    await fetchDate(searchDate.value)
  } catch (e) {
    console.error('加载龙虎榜失败:', e)
    list.value = []
  } finally {
    loading.value = false
  }
}

// 切换日期（用户手动点 ‹/›）：重置到该日期再取数（仍允许空数据回退）
function changeDate(delta) {
  searchDate.value = shiftDate(searchDate.value, delta)
  loadData()
}

async function handleRefresh() {
  await loadData()
}

onBeforeMount(() => {
  loadData()
})
</script>

<template>
  <div class="page">
    <!-- 日期切换条 -->
    <div class="date-bar">
      <button class="date-btn" type="button" @click="changeDate(-1)">‹</button>
      <span class="date-text">{{ searchDate }}</span>
      <button class="date-btn" type="button" @click="changeDate(1)">›</button>
    </div>

    <MPullRefresh class="dragon-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !list.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>
        <template v-else>
          <div v-if="list.length" class="dragon-list">
            <div v-for="(item, i) in list" :key="(item.SECUCODE || '') + i" class="dragon-item">
              <!-- 第一行：股票名+代码 / 收盘价+涨跌幅 -->
              <div class="row-top">
                <div class="stock-left">
                  <span class="stock-name">{{ item.SECURITY_NAME_ABBR }}</span>
                  <span class="stock-code">{{ marketCode(item) }}</span>
                </div>
                <div class="price-group">
                  <span class="close" :class="netClass(item.CHANGE_RATE)">{{ Number(item.CLOSE_PRICE).toFixed(2) }}</span>
                  <PercentTag :value="Number(item.CHANGE_RATE)" size="small" />
                </div>
              </div>

              <!-- 第二行：龙虎榜净买额（主指标） + 换手率 -->
              <div class="row-net">
                <span class="net-label">龙虎榜净买额</span>
                <span class="net-val" :class="netClass(item.BILLBOARD_NET_AMT)">
                  {{ formatMoney(item.BILLBOARD_NET_AMT) }}
                </span>
                <span class="turnover">换手 {{ Number(item.TURNOVERRATE || 0).toFixed(2) }}%</span>
              </div>

              <!-- 第三行：上榜原因 -->
              <div v-if="item.EXPLANATION" class="row-reason">
                <span class="reason-tag">{{ item.EXPLANATION }}</span>
              </div>
            </div>
          </div>
          <MEmpty v-else description="该日期暂无龙虎榜数据" />
        </template>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.date-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--m-space-xl);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.date-btn {
  width: 36px;
  height: 36px;
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-lg);
  line-height: 1;
}

.date-btn:active {
  transform: scale(0.92);
}

.date-text {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
  min-width: 110px;
  text-align: center;
}

.dragon-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

/* 卡片列表：分隔线 */
.dragon-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.dragon-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.dragon-item {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
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

.close {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
}

.close.m-rise {
  color: var(--m-color-rise);
}

.close.m-fall {
  color: var(--m-color-fall);
}

/* 净买额行 */
.row-net {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
}

.net-label {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.net-val {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
  color: var(--m-color-gray);
}

.net-val.m-rise {
  color: var(--m-color-rise);
}

.net-val.m-fall {
  color: var(--m-color-fall);
}

.turnover {
  margin-left: auto;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

/* 上榜原因 */
.row-reason {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
}

.reason-tag {
  padding: 1px var(--m-space-sm);
  font-size: var(--m-font-xs);
  color: var(--m-color-rise);
  background: var(--m-color-rise-light);
  border-radius: var(--m-radius-sm);
}
</style>
