<script setup>
import MCard from '../base/MCard.vue'
import PercentTag from '../widgets/PercentTag.vue'

defineProps({
  // 行业数据
  industries: {
    type: Array,
    default: () => []
    // { name, changePercent, leadingStock }
  },
  // 最多显示数量
  maxShow: {
    type: Number,
    default: 5
  }
})

const emit = defineEmits(['industryClick', 'viewAll'])

function handleIndustryClick(industry) {
  emit('industryClick', industry)
}

function handleViewAll() {
  emit('viewAll')
}
</script>

<template>
  <MCard>
    <!-- 头部 -->
    <div class="industry-header">
      <h3 class="header-title">📈 行业热度 TOP{{ maxShow }}</h3>
      <button class="header-action" @click="handleViewAll">
        全部 →
      </button>
    </div>

    <!-- 行业列表 -->
    <div v-if="industries.length" class="industry-list">
      <div
        v-for="(industry, index) in industries.slice(0, maxShow)"
        :key="index"
        class="industry-item"
        @click="handleIndustryClick(industry)"
      >
        <div class="industry-rank">{{ index + 1 }}</div>
        <div class="industry-info">
          <div class="industry-name">{{ industry.name }}</div>
          <div v-if="industry.leadingStock" class="industry-leading">
            领涨: {{ industry.leadingStock }}
          </div>
        </div>
        <div class="industry-change">
          <PercentTag :value="industry.changePercent" size="large" bold />
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="industry-empty">
      <p>暂无行业数据</p>
    </div>
  </MCard>
</template>

<style scoped>
.industry-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--m-space-lg);
}

.header-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
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

.industry-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.industry-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  cursor: pointer;
  transition: all var(--m-duration-fast);
}

.industry-item:active {
  transform: scale(0.98);
  background: var(--m-divider-color);
}

.industry-rank {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-sm);
  flex-shrink: 0;
}

.industry-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--m-space-xs);
}

.industry-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.industry-leading {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.industry-change {
  flex-shrink: 0;
}

.industry-empty {
  text-align: center;
  padding: var(--m-space-xl) 0;
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
}
</style>
