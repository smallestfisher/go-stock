<script setup>
import { computed } from 'vue'

const props = defineProps({
  // 新闻数据
  news: {
    type: Object,
    required: true
    // { title, summary, time, source }
  },
  // 是否显示摘要
  showSummary: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

// 格式化时间
const formattedTime = computed(() => {
  if (!props.news.time) return ''

  const now = new Date()
  const newsTime = new Date(props.news.time)
  const diff = Math.floor((now - newsTime) / 1000) // 秒

  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  if (diff < 259200) return `${Math.floor(diff / 86400)}天前`

  // 超过3天显示具体日期
  return newsTime.toLocaleDateString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
})

function handleClick() {
  emit('click', props.news)
}
</script>

<template>
  <div class="news-card" @click="handleClick">
    <!-- 时间标记（左侧竖线） -->
    <div class="news-timeline">
      <div class="timeline-dot" />
      <div class="timeline-line" />
    </div>

    <!-- 内容 -->
    <div class="news-content">
      <div class="news-header">
        <span class="news-time">{{ formattedTime }}</span>
        <span v-if="news.source" class="news-source">{{ news.source }}</span>
      </div>
      <h4 class="news-title">{{ news.title }}</h4>
      <p v-if="showSummary && news.summary" class="news-summary">
        {{ news.summary }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.news-card {
  display: flex;
  gap: var(--m-space-md);
  padding: var(--m-space-md) 0;
  cursor: pointer;
  transition: background var(--m-duration-fast);
}

.news-card:active {
  background: var(--m-bg-primary);
}

.news-timeline {
  position: relative;
  width: 16px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.timeline-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--m-color-rise);
  border: 2px solid var(--m-bg-card);
  box-shadow: 0 0 0 2px var(--m-color-rise-light);
  flex-shrink: 0;
  margin-top: 6px;
}

.timeline-line {
  width: 2px;
  flex: 1;
  background: var(--m-divider-color);
  margin-top: var(--m-space-xs);
}

.news-content {
  flex: 1;
  min-width: 0;
}

.news-header {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  margin-bottom: var(--m-space-xs);
}

.news-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.news-source {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  padding: 2px var(--m-space-xs);
  background: var(--m-bg-primary);
  border-radius: 2px;
}

.news-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  line-height: var(--m-line-height-normal);
  margin-bottom: var(--m-space-xs);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.news-summary {
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>
