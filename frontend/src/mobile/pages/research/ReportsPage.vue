<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'

const reports = ref([
  { title: '半导体行业深度报告：AI算力需求持续爆发', institution: '中金公司', date: '2024-06-20', rating: '买入', views: 1234 },
  { title: '新能源汽车行业周报：销量持续超预期', institution: '华泰证券', date: '2024-06-19', rating: '强推', views: 2345 },
  { title: '贵州茅台：业绩稳健增长，维持买入评级', institution: '中信证券', date: '2024-06-18', rating: '买入', views: 3456 },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard padding="none">
          <div class="report-list">
            <div v-for="(item, i) in reports" :key="i" class="report-item">
              <div class="report-title">{{ item.title }}</div>
              <div class="report-meta">
                <span>{{ item.institution }}</span>
                <span>{{ item.date }}</span>
                <span class="report-rating m-rise">{{ item.rating }}</span>
              </div>
              <div class="report-footer">
                <span class="report-views">👁 {{ item.views }}</span>
              </div>
            </div>
          </div>
        </MCard>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.report-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-xl); }
.report-item { padding-bottom: var(--m-space-xl); border-bottom: 1px solid var(--m-divider-color); }
.report-item:last-child { border-bottom: none; }
.report-title { font-weight: var(--m-font-weight-medium); line-height: var(--m-line-height-normal); margin-bottom: var(--m-space-sm); }
.report-meta { display: flex; gap: var(--m-space-md); font-size: var(--m-font-xs); color: var(--m-text-secondary); margin-bottom: var(--m-space-sm); }
.report-rating { font-weight: var(--m-font-weight-medium); }
.report-footer { display: flex; justify-content: flex-end; }
.report-views { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
</style>
