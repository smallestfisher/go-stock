<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

const reports = ref([
  { title: '贵州茅台：业绩稳健增长，维持买入评级', author: '中信证券', date: '2024-06-20', rating: '买入' },
  { title: '宁德时代：海外业务加速拓展', author: '招商证券', date: '2024-06-19', rating: '强烈推荐' },
])

async function handleRefresh() { return new Promise(resolve => setTimeout(resolve, 1500)) }
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="reports.length" padding="none">
          <div class="list">
            <div v-for="(item, i) in reports" :key="i" class="item">
              <div class="title">{{ item.title }}</div>
              <div class="meta">
                <span>{{ item.author }}</span>
                <span>{{ item.date }}</span>
                <span class="rating m-rise">{{ item.rating }}</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无研报" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-lg); }
.item { padding-bottom: var(--m-space-lg); border-bottom: 1px solid var(--m-divider-color); }
.item:last-child { border-bottom: none; }
.title { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-sm); }
.meta { display: flex; gap: var(--m-space-md); font-size: var(--m-font-xs); color: var(--m-text-secondary); }
.rating { font-weight: var(--m-font-weight-medium); }
</style>
