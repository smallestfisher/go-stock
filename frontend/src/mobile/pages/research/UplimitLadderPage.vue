<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'

// 涨停梯队数据（后续接 GetUplimitHot(start, end) 填充，含连板梯队：首板/2连板/3连板...）
const ladderData = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="ladderData.length" padding="none">
          <div class="ladder-list">
            <div v-for="(item, i) in ladderData" :key="i" class="ladder-item">
              <div class="item-info">
                <div class="item-name">{{ item.stockName }}</div>
                <div class="item-code">{{ item.stockCode }}</div>
              </div>
              <div class="item-right">
                <PercentTag :value="item.changePercent" size="small" />
                <span class="item-ladder">{{ item.ladderCount }}连板</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无涨停数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.ladder-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.ladder-item { display: flex; justify-content: space-between; align-items: center; padding-bottom: var(--m-space-md); border-bottom: 1px solid var(--m-divider-color); }
.ladder-item:last-child { border-bottom: none; }
.item-info { display: flex; flex-direction: column; gap: var(--m-space-xs); }
.item-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.item-code { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.item-right { display: flex; flex-direction: column; align-items: flex-end; gap: var(--m-space-xs); }
.item-ladder { font-size: var(--m-font-xs); color: var(--m-color-rise); font-weight: var(--m-font-weight-medium); }
</style>
