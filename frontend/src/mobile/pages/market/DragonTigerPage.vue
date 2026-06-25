<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 龙虎榜数据（后续接 GetChangeRank(market, date) 填充）
const dragonData = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="dragonData.length">
          <div class="list">
            <div v-for="item in dragonData" :key="item.stockCode" class="item">
              <div class="info">
                <div class="name">{{ item.stockName }} <span class="code">{{ item.stockCode }}</span></div>
                <div class="reason">{{ item.reason }}</div>
              </div>
              <PercentTag :value="item.changePercent" bold size="large" />
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无龙虎榜数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.item { display: flex; justify-content: space-between; align-items: center; padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.info { flex: 1; }
.name { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-xs); }
.code { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.reason { font-size: var(--m-font-sm); color: var(--m-text-secondary); }
</style>
