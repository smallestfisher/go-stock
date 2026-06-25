<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 个股研报数据（后续接 GetTdxCompanyCategoryList/Content 或研报接口填充）
const reports = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
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
