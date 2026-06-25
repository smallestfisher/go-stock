<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import NewsCard from '../../components/cards/NewsCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 新闻列表（后续接 GetTelegraphList 填充）
const newsData = ref([])

const refreshCount = ref(0)
const loadingMore = ref(false)

// 下拉刷新（后续接真实 API）
async function handleRefresh() {
  refreshCount.value++
}

// 加载更多（后续接分页 API）
function handleLoadMore() {
  if (loadingMore.value) return
  loadingMore.value = true
  setTimeout(() => {
    loadingMore.value = false
  }, 1000)
}

// 点击新闻
function handleNewsClick(news) {
  console.log('点击新闻:', news)
  // TODO: 打开新闻详情
}
</script>

<template>
  <div class="news-list-page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="news-list-container">
        <MCard v-if="newsData.length" padding="none">
          <div class="news-list">
            <NewsCard
              v-for="(news, index) in newsData"
              :key="index"
              :news="news"
              show-summary
              @click="handleNewsClick"
            />
          </div>

          <!-- 加载更多 -->
          <div class="load-more" @click="handleLoadMore">
            <span v-if="!loadingMore">加载更多...</span>
            <span v-else>加载中...</span>
          </div>
        </MCard>

        <MEmpty v-else description="暂无快讯" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.news-list-page {
  height: 100%;
  overflow: hidden;
}

.news-list-container {
  padding: var(--m-space-md);
}

.news-list {
  padding: var(--m-space-md);
}

.load-more {
  text-align: center;
  padding: var(--m-space-lg);
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
  cursor: pointer;
  border-top: 1px solid var(--m-divider-color);
}

.load-more:active {
  opacity: 0.6;
}
</style>
