<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import { formatMoney } from '../../composables/useFormat'

// 概念资金流数据（后续接 GetConceptFundFlowList / GetConceptFundFlowTopList 填充）
const concepts = ref([])

async function handleRefresh() {
  // 后续接真实 API
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
