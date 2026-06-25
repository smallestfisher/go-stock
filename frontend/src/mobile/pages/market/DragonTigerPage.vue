<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import MEmpty from '../../components/base/MEmpty.vue'

const dragonData = ref([
  { stockName: '寒武纪', stockCode: '688256', changePercent: 10.00, turnover: 139620000, reason: '日涨幅偏离值达7%' },
  { stockName: '中芯国际', stockCode: '688981', changePercent: 8.45, turnover: 256780000, reason: '日振幅达15%' },
])

async function handleRefresh() { return new Promise(resolve => setTimeout(resolve, 1500)) }
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
