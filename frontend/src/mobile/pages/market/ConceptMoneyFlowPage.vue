<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

const concepts = ref([
  { name: 'ChatGPT概念', netInflow: 95000000, changePercent: 4.56 },
  { name: 'AI芯片', netInflow: 78000000, changePercent: 3.45 },
  { name: '储能概念', netInflow: -32000000, changePercent: -1.89 },
])

async function handleRefresh() { return new Promise(resolve => setTimeout(resolve, 1500)) }
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
        <MCard v-if="concepts.length">
          <div class="list">
            <div v-for="(item, i) in concepts" :key="i" class="item">
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
