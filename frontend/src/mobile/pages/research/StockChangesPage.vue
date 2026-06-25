<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'

// 数据模式：实时 / 历史（对齐桌面端 stockChangesMonitor 的 viewMode）
const modeTabs = [
  { label: '实时', value: 'realtime' },
  { label: '历史', value: 'history' },
]
const activeMode = ref('realtime')

// 异动类型筛选（对齐桌面端 selectedTypes）
const typeTabs = [
  { label: '全部', value: [0] },
  { label: '涨停', value: [1] },
  { label: '跌停', value: [2] },
  { label: '急涨', value: [3] },
  { label: '急跌', value: [4] },
  { label: '放量', value: [5] },
  { label: '突破', value: [6] },
]
const activeType = ref([0])

// 异动数据（后续接 GetStockChanges / GetAllStockChangesWithPaging 填充）
const changesData = ref([])

async function handleRefresh() {
  // 后续根据 activeMode + activeType 调对应 API
}
</script>

<template>
  <div class="page">
    <div class="filter-bar">
      <MTabs v-model="activeMode" :tabs="modeTabs" />
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="changesData.length" padding="none">
          <div class="change-list">
            <div v-for="(item, i) in changesData" :key="i" class="change-item">
              <div class="item-info">
                <div class="item-name">{{ item.stockName }}</div>
                <div class="item-code">{{ item.stockCode }}</div>
              </div>
              <div class="item-right">
                <PercentTag :value="item.changePercent" size="small" />
                <span class="item-time">{{ item.time }}</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无异动数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; display: flex; flex-direction: column; overflow: hidden; }
.filter-bar { background: var(--m-bg-card); border-bottom: 1px solid var(--m-divider-color); }
.container { padding: var(--m-space-md); }
.change-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.change-item { display: flex; justify-content: space-between; align-items: center; padding-bottom: var(--m-space-md); border-bottom: 1px solid var(--m-divider-color); }
.change-item:last-child { border-bottom: none; }
.item-info { display: flex; flex-direction: column; gap: var(--m-space-xs); }
.item-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.item-code { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.item-right { display: flex; flex-direction: column; align-items: flex-end; gap: var(--m-space-xs); }
.item-time { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
</style>
