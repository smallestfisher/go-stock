<script setup>
import { ref, computed } from 'vue'

// 单条电报/快讯卡片。对齐桌面端 newsList.vue 的字段语义，
// 但用移动端时间轴 + 标签风格呈现：
//  - 有标题的（多为「外媒」「重要电报」）默认折叠正文，点标题展开；
//  - 无标题的（多为「财联社/新浪」滚动快讯）直接显示正文。
const props = defineProps({
  item: {
    type: Object,
    required: true,
    // { time, dataTime, title, content, isRed, subjects, stocks, url, sentimentResult }
  },
})

const expanded = ref(false)

// 时刻：后端给的是 "HH:MM:SS" 或 "HH:MM"，移动端只取到分钟更紧凑。
const timeText = computed(() => String(props.item.time || '').slice(0, 5))

const hasTitle = computed(() => !!props.item.title)
const showContent = computed(() => !hasTitle.value || expanded.value)

// 情绪：与桌面端一致——看涨=红(涨)，看跌=绿(跌)。
const sentiment = computed(() => props.item.sentimentResult || '')

function toggle() {
  if (hasTitle.value) expanded.value = !expanded.value
}
</script>

<template>
  <div class="tg-item">
    <!-- 左侧时间轴 -->
    <div class="tg-rail">
      <span class="tg-dot" :class="{ 'tg-dot--red': item.isRed }" />
      <span class="tg-line" />
    </div>

    <!-- 右侧内容 -->
    <div class="tg-body">
      <div class="tg-time" :class="{ 'tg-time--red': item.isRed }">{{ timeText }}</div>

      <!-- 标题（可点击展开正文） -->
      <div
        v-if="hasTitle"
        class="tg-title"
        :class="{ 'tg-title--red': item.isRed }"
        @click="toggle"
      >
        <span>{{ item.title }}</span>
        <span class="tg-arrow" :class="{ 'tg-arrow--open': expanded }">›</span>
      </div>

      <!-- 正文 -->
      <p v-if="showContent && item.content" class="tg-content">{{ item.content }}</p>

      <!-- 标签行 -->
      <div
        v-if="(item.subjects && item.subjects.length) || (item.stocks && item.stocks.length) || sentiment || item.url"
        class="tg-tags"
      >
        <span
          v-for="sub in item.subjects"
          :key="'s' + sub"
          class="tg-tag tg-tag--subject"
        >{{ sub }}</span>
        <span
          v-for="st in item.stocks"
          :key="'k' + st"
          class="tg-tag tg-tag--stock"
        >{{ st }}</span>
        <span
          v-if="sentiment"
          class="tg-tag"
          :class="sentiment === '看涨' ? 'tg-tag--up' : sentiment === '看跌' ? 'tg-tag--down' : 'tg-tag--subject'"
        >{{ sentiment }}</span>
        <a
          v-if="item.url"
          class="tg-tag tg-tag--link"
          :href="item.url"
          target="_blank"
          rel="noopener"
          @click.stop
        >查看原文</a>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tg-item {
  display: flex;
  gap: var(--m-space-md);
  padding: var(--m-space-md) var(--m-space-lg);
}

/* 时间轴 */
.tg-rail {
  position: relative;
  width: 10px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.tg-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--m-text-tertiary);
  margin-top: 6px;
  flex-shrink: 0;
}

.tg-dot--red {
  background: var(--m-color-rise);
  box-shadow: 0 0 0 2px var(--m-color-rise-light);
}

.tg-line {
  width: 2px;
  flex: 1;
  background: var(--m-divider-color);
  margin-top: var(--m-space-xs);
}

/* 内容 */
.tg-body {
  flex: 1;
  min-width: 0;
  padding-bottom: var(--m-space-xs);
}

.tg-time {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
  margin-bottom: 2px;
}

.tg-time--red {
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.tg-title {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--m-space-sm);
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  line-height: var(--m-line-height-normal);
  cursor: pointer;
}

.tg-title--red {
  color: var(--m-color-rise);
}

.tg-arrow {
  flex-shrink: 0;
  color: var(--m-text-tertiary);
  transform: rotate(90deg);
  transition: transform var(--m-duration-fast);
}

.tg-arrow--open {
  transform: rotate(-90deg);
}

.tg-content {
  margin: var(--m-space-xs) 0 0;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-normal);
  word-break: break-word;
}

/* 标签 */
.tg-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
  margin-top: var(--m-space-sm);
}

.tg-tag {
  display: inline-flex;
  align-items: center;
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
  line-height: 1.6;
  white-space: nowrap;
}

.tg-tag--subject {
  background: var(--m-bg-primary);
  color: var(--m-text-secondary);
  border: 1px solid var(--m-divider-color);
}

.tg-tag--stock {
  background: rgba(240, 160, 32, 0.12);
  color: #c77800;
}

.tg-tag--up {
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.tg-tag--down {
  background: var(--m-color-fall-light);
  color: var(--m-color-fall);
}

.tg-tag--link {
  background: rgba(240, 160, 32, 0.12);
  color: #c77800;
  text-decoration: none;
}

.tg-tag--link:active {
  opacity: 0.7;
}
</style>
