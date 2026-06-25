<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MTabs from '../../components/base/MTabs.vue'

// 问答广场（对齐桌面端 promptQa.vue）：关键词搜索 + 已解决筛选
const keyword = ref('')
const filterTabs = [
  { label: '全部', value: '' },
  { label: '待回答', value: 'unresolved' },
  { label: '已解决', value: 'resolved' },
]
const resolvedFilter = ref('')

// 问题列表（后续接问答广场 API 填充）
const questions = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="page">
    <div class="qa-header">
      <input
        v-model="keyword"
        class="search-input"
        type="text"
        placeholder="搜索问题..."
      >
      <MTabs v-model="resolvedFilter" :tabs="filterTabs" />
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="questions.length" padding="none">
          <div class="qa-list">
            <div v-for="(item, i) in questions" :key="i" class="qa-item">
              <div class="qa-title">
                <span v-if="item.resolved" class="qa-tag qa-tag--resolved">已解决</span>
                {{ item.title }}
              </div>
              <p v-if="item.content" class="qa-content">{{ item.content }}</p>
              <div class="qa-meta">
                <span>👤 {{ item.author }}</span>
                <span>💬 {{ item.answerCount || 0 }} 回答</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无问题" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; display: flex; flex-direction: column; overflow: hidden; }
.qa-header { padding: var(--m-space-md); background: var(--m-bg-card); border-bottom: 1px solid var(--m-divider-color); display: flex; flex-direction: column; gap: var(--m-space-sm); }
.search-input { width: 100%; padding: var(--m-space-sm) var(--m-space-md); border: 1px solid var(--m-divider-color); border-radius: var(--m-radius-md); font-size: var(--m-font-sm); background: var(--m-bg-primary); }
.search-input:focus { outline: none; border-color: var(--m-color-rise); }
.container { padding: var(--m-space-md); }
.qa-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-lg); }
.qa-item { padding-bottom: var(--m-space-lg); border-bottom: 1px solid var(--m-divider-color); }
.qa-item:last-child { border-bottom: none; }
.qa-title { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); margin-bottom: var(--m-space-xs); }
.qa-tag { font-size: var(--m-font-xs); padding: 2px var(--m-space-sm); border-radius: var(--m-radius-sm); margin-right: var(--m-space-sm); }
.qa-tag--resolved { background: var(--m-color-rise-light); color: var(--m-color-rise); }
.qa-content { font-size: var(--m-font-sm); color: var(--m-text-secondary); line-height: var(--m-line-height-normal); margin-bottom: var(--m-space-sm); }
.qa-meta { display: flex; gap: var(--m-space-lg); font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
</style>
