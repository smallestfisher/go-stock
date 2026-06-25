<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'

// 交易日志统计（对齐桌面端 TradingRecordManager.vue）
// 后续接 GetTradingRecordStatistics 填充
const statistics = ref(null)

// 交易记录列表（后续接 GetTradingRecordList 填充）
// 字段参考：{ stockName, stockCode, type(买/卖), price, volume, time }
const records = ref([])

async function handleRefresh() {
  // 后续接真实 API
}

function handleAdd() {
  // 后续打开新增抽屉，接 AddTradingRecord
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 统计概览 -->
        <MCard v-if="statistics">
          <h3 class="section-title">交易统计</h3>
          <div class="stat-grid">
            <div class="stat-item">
              <div class="stat-value m-rise">{{ statistics.totalProfit || 0 }}</div>
              <div class="stat-label">总盈亏</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ statistics.totalCount || 0 }}</div>
              <div class="stat-label">交易次数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value m-rise">{{ statistics.winRate || 0 }}%</div>
              <div class="stat-label">胜率</div>
            </div>
          </div>
        </MCard>

        <!-- 交易记录 -->
        <MCard v-if="records.length" padding="none">
          <h3 class="section-title list-title">交易记录</h3>
          <div class="record-list">
            <div v-for="(item, i) in records" :key="i" class="record-item">
              <div class="record-info">
                <div class="record-name">{{ item.stockName }}</div>
                <div class="record-code">{{ item.stockCode }}</div>
              </div>
              <div class="record-detail">
                <span class="record-type" :class="item.type === 'buy' ? 'm-rise' : 'm-fall'">
                  {{ item.type === 'buy' ? '买入' : '卖出' }}
                </span>
                <span class="record-price">¥{{ item.price }}</span>
                <span class="record-volume">{{ item.volume }}股</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="还没有交易记录">
          <MButton type="primary" @click="handleAdd">添加记录</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.section-title { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-md); }
.list-title { padding: var(--m-space-md) var(--m-space-md) 0; }
.stat-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: var(--m-space-md); }
.stat-item { text-align: center; }
.stat-value { font-size: var(--m-font-xl); font-weight: var(--m-font-weight-bold); color: var(--m-text-primary); font-variant-numeric: tabular-nums; }
.stat-label { font-size: var(--m-font-xs); color: var(--m-text-tertiary); margin-top: var(--m-space-xs); }
.record-list { padding: 0 var(--m-space-md) var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.record-item { display: flex; justify-content: space-between; align-items: center; padding-bottom: var(--m-space-md); border-bottom: 1px solid var(--m-divider-color); }
.record-item:last-child { border-bottom: none; }
.record-info { display: flex; flex-direction: column; gap: 2px; }
.record-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.record-code { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.record-detail { display: flex; flex-direction: column; align-items: flex-end; gap: 2px; font-size: var(--m-font-sm); }
.record-type { font-weight: var(--m-font-weight-medium); }
.record-price { color: var(--m-text-primary); font-variant-numeric: tabular-nums; }
.record-volume { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
</style>
