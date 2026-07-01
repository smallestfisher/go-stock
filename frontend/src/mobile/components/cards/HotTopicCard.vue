<script setup>
import MCard from '../base/MCard.vue'
import MIcon from '../base/MIcon.vue'
import PercentTag from '../widgets/PercentTag.vue'

defineProps({
  // 热点数据
  topics: {
    type: Array,
    default: () => []
    // { title, heat, changePercent, stocks }
  }
})

const emit = defineEmits(['topicClick', 'viewMore'])

function handleTopicClick(topic) {
  emit('topicClick', topic)
}

function handleViewMore() {
  emit('viewMore')
}
</script>

<template>
  <MCard>
    <!-- 头部 -->
    <div class="hot-topic-header">
      <h3 class="header-title">实时热点</h3>
      <button class="header-action" @click="handleViewMore">
        更多 <MIcon name="arrow-right" :size="14" />
      </button>
    </div>

    <!-- 热点列表 -->
    <div v-if="topics.length" class="topic-list">
      <div
        v-for="(topic, index) in topics"
        :key="index"
        class="topic-item"
        @click="handleTopicClick(topic)"
      >
        <div class="topic-rank" :class="`topic-rank--${index + 1}`">
          {{ index + 1 }}
        </div>
        <div class="topic-info">
          <div class="topic-title">{{ topic.title }}</div>
          <div class="topic-meta">
            <span class="topic-heat"><MIcon name="fire" :size="14" /> {{ topic.heat }}</span>
            <span v-if="topic.stocks" class="topic-stocks">
              {{ topic.stocks }} 只相关
            </span>
          </div>
        </div>
        <div v-if="topic.changePercent != null" class="topic-change">
          <PercentTag :value="topic.changePercent" size="small" />
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="topic-empty">
      <p>暂无热点数据</p>
    </div>
  </MCard>
</template>

<style scoped>
.hot-topic-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--m-space-lg);
}

.header-title {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  letter-spacing: 0.2px;
}

.header-title::before {
  content: '';
  width: 3px;
  height: 15px;
  border-radius: var(--m-radius-full);
  background: var(--m-color-rise);
}

.header-action {
  padding: var(--m-space-xs) var(--m-space-sm);
  background: transparent;
  border: none;
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
  cursor: pointer;
}

.header-action:active {
  opacity: 0.6;
}

.topic-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.topic-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.topic-item:active {
  transform: scale(0.98);
  background: var(--m-divider-color);
}

.topic-rank {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-tertiary);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-sm);
  flex-shrink: 0;
}

/* 前三名特殊样式 */
.topic-rank--1 {
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  color: white;
}

.topic-rank--2 {
  background: linear-gradient(135deg, #c0c0c0, #e8e8e8);
  color: white;
}

.topic-rank--3 {
  background: linear-gradient(135deg, #cd7f32, #daa520);
  color: white;
}

.topic-info {
  flex: 1;
  min-width: 0;
}

.topic-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  margin-bottom: var(--m-space-xs);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.topic-meta {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.topic-heat {
  color: var(--m-color-fall);
  font-weight: var(--m-font-weight-medium);
}

.topic-change {
  flex-shrink: 0;
}

.topic-empty {
  text-align: center;
  padding: var(--m-space-xl) 0;
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
}
</style>
