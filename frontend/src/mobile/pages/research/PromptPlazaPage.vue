<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MTabs from '../../components/base/MTabs.vue'

// 提示词广场：分类 + 排序 + 关键词（对齐桌面端 promptPlaza.vue）
// 数据后续接提示词广场代理 API 填充
const categories = ref([])
const activeCategory = ref(null)
const activeSort = ref('latest')
const keyword = ref('')

const sortTabs = [
  { label: '最新', value: 'latest' },
  { label: '最热', value: 'hot' },
  { label: '精选', value: 'featured' },
]

// 广场提示词列表
const prompts = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="page">
    <div class="plaza-header">
      <!-- 搜索框 -->
      <input
        v-model="keyword"
        class="search-input"
        type="text"
        placeholder="搜索提示词..."
      >
      <!-- 排序 -->
      <MTabs v-model="activeSort" :tabs="sortTabs" />
    </div>

    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 分类 -->
        <div v-if="categories.length" class="category-list">
          <button
            v-for="cat in categories"
            :key="cat"
            class="category-chip"
            :class="{ 'category-chip--active': activeCategory === cat }"
            @click="activeCategory = cat"
          >
            {{ cat }}
          </button>
        </div>

        <!-- 提示词列表 -->
        <MCard v-if="prompts.length" padding="none">
          <div class="prompt-list">
            <div v-for="(item, i) in prompts" :key="i" class="prompt-item">
              <div class="prompt-title">{{ item.title }}</div>
              <p class="prompt-desc">{{ item.content }}</p>
              <div class="prompt-meta">
                <span>👤 {{ item.author }}</span>
                <span>❤️ {{ item.likes }}</span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无提示词" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; display: flex; flex-direction: column; overflow: hidden; }
.plaza-header { padding: var(--m-space-md); background: var(--m-bg-card); border-bottom: 1px solid var(--m-divider-color); display: flex; flex-direction: column; gap: var(--m-space-sm); }
.search-input { width: 100%; padding: var(--m-space-sm) var(--m-space-md); border: 1px solid var(--m-divider-color); border-radius: var(--m-radius-md); font-size: var(--m-font-sm); background: var(--m-bg-primary); }
.search-input:focus { outline: none; border-color: var(--m-color-rise); }
.container { padding: var(--m-space-md); }
.category-list { display: flex; flex-wrap: wrap; gap: var(--m-space-sm); margin-bottom: var(--m-space-md); }
.category-chip { padding: var(--m-space-xs) var(--m-space-md); background: var(--m-bg-card); border: 1px solid var(--m-divider-color); border-radius: var(--m-radius-full); font-size: var(--m-font-sm); }
.category-chip--active { background: var(--m-color-rise); border-color: var(--m-color-rise); color: #fff; }
.prompt-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-lg); }
.prompt-item { padding-bottom: var(--m-space-lg); border-bottom: 1px solid var(--m-divider-color); }
.prompt-item:last-child { border-bottom: none; }
.prompt-title { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); margin-bottom: var(--m-space-xs); }
.prompt-desc { font-size: var(--m-font-sm); color: var(--m-text-secondary); line-height: var(--m-line-height-normal); margin-bottom: var(--m-space-sm); }
.prompt-meta { display: flex; gap: var(--m-space-lg); font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
</style>
