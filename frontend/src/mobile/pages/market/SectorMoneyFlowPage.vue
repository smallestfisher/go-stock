<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

const sectors = ref([
  { name: '半导体板块', netInflow: 125000000, changePercent: 5.23 },
  { name: '新能源汽车', netInflow: 89000000, changePercent: 3.87 },
  { name: '白酒板块', netInflow: -45000000, changePercent: -1.23 },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}

function formatMoney(value) {
  if (value >= 100000000) return `${(value / 100000000).toFixed(2)}亿`
  if (value >= 10000) return `${(value / 10000).toFixed(2)}万`
  return value
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="sectors.length">
          <div class="list">
            <div v-for="(item, i) in sectors" :key="i" class="item">
              <div class="name">{{ item.name }}</div>
              <div class="flow" :class="item.netInflow > 0 ? 'm-rise' : 'm-fall'">
                {{ item.netInflow > 0 ? '+' : '' }}{{ formatMoney(item.netInflow) }}
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.item { display: flex; justify-content: space-between; padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.name { font-weight: var(--m-font-weight-medium); }
.flow { font-weight: var(--m-font-weight-bold); }
</style>
