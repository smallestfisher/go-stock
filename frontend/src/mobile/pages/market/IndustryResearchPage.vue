<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 行业研究报告数据（后续接 IndustryResearchReport 填充）
const research = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="research.length" padding="none">
          <div class="list">
            <div v-for="(item, i) in research" :key="i" class="item">
              <div class="title">{{ item.title }}</div>
              <div class="meta">
                <span>{{ item.institution }}</span>
                <span>{{ item.date }}</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无行业研究" />
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
.title { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-sm); line-height: var(--m-line-height-normal); }
.meta { display: flex; gap: var(--m-space-md); font-size: var(--m-font-xs); color: var(--m-text-secondary); }
</style>
