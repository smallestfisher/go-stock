<script setup>
defineProps({
  // 空状态图片（可选）
  image: {
    type: String,
    default: ''
  },
  // 空状态描述
  description: {
    type: String,
    default: '暂无数据'
  },
  // 图片尺寸
  imageSize: {
    type: String,
    default: '120px'
  }
})
</script>

<template>
  <div class="m-empty">
    <!-- 自定义图片插槽 -->
    <div v-if="$slots.image" class="m-empty__image">
      <slot name="image" />
    </div>
    <!-- 默认图片 -->
    <div v-else-if="image" class="m-empty__image">
      <img :src="image" :style="{ width: imageSize, height: imageSize }" alt="empty" />
    </div>
    <!-- 默认图标 -->
    <div v-else class="m-empty__icon" :style="{ width: imageSize, height: imageSize }">
      <svg viewBox="0 0 64 41" fill="currentColor">
        <g transform="translate(0 1)" fill="none" fill-rule="evenodd">
          <ellipse fill="#f5f5f5" cx="32" cy="33" rx="32" ry="7" />
          <g fill-rule="nonzero" stroke="#d9d9d9">
            <path d="M55 12.76L44.854 1.258C44.367.474 43.656 0 42.907 0H21.093c-.749 0-1.46.474-1.947 1.257L9 12.761V22h46v-9.24z" />
            <path d="M41.613 15.931c0-1.605.994-2.93 2.227-2.931H55v18.137C55 33.26 53.68 35 52.05 35h-40.1C10.32 35 9 33.259 9 31.137V13h11.16c1.233 0 2.227 1.323 2.227 2.928v.022c0 1.605 1.005 2.901 2.237 2.901h14.752c1.232 0 2.237-1.308 2.237-2.913v-.007z" fill="#fafafa" />
          </g>
        </g>
      </svg>
    </div>

    <!-- 描述文字 -->
    <div class="m-empty__description">
      <slot name="description">
        {{ description }}
      </slot>
    </div>

    <!-- 操作按钮插槽 -->
    <div v-if="$slots.default" class="m-empty__actions">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.m-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--m-space-2xl) var(--m-space-xl);
  min-height: 200px;
}

.m-empty__image,
.m-empty__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--m-space-lg);
}

.m-empty__image img {
  display: block;
  object-fit: contain;
}

.m-empty__icon {
  color: #d9d9d9;
  opacity: 0.6;
}

.m-empty__description {
  font-size: var(--m-font-md);
  color: var(--m-text-secondary);
  text-align: center;
  margin-bottom: var(--m-space-lg);
  line-height: var(--m-line-height-normal);
}

.m-empty__actions {
  margin-top: var(--m-space-md);
}
</style>
