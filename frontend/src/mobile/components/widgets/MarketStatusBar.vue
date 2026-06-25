<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { cnOpen, hkOpen, usOpen } from '../../../api/marketClock'

const marketStatus = ref({
  cn: { open: false, label: 'A股' },
  hk: { open: false, label: '港股' },
  us: { open: false, label: '美股' }
})

function updateStatus() {
  marketStatus.value.cn.open = cnOpen.value
  marketStatus.value.hk.open = hkOpen.value
  marketStatus.value.us.open = usOpen.value
}

let timer = null

onMounted(() => {
  updateStatus()
  // 每分钟更新一次
  timer = setInterval(updateStatus, 60000)
})

onBeforeUnmount(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<template>
  <div class="market-status-bar">
    <div
      v-for="(market, key) in marketStatus"
      :key="key"
      class="market-status-item"
    >
      <span
        class="market-status-dot"
        :class="{ 'market-status-dot--open': market.open }"
      />
      <span class="market-status-label">{{ market.label }}</span>
      <span class="market-status-text">
        {{ market.open ? '交易中' : '休市' }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.market-status-bar {
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
}

.market-status-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-xs);
}

.market-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--m-text-tertiary);
  transition: background var(--m-duration-normal);
}

.market-status-dot--open {
  background: var(--m-color-rise);
  box-shadow: 0 0 4px var(--m-color-rise);
  animation: m-pulse 2s ease-in-out infinite;
}

.market-status-label {
  color: var(--m-text-secondary);
  font-weight: var(--m-font-weight-medium);
}

.market-status-text {
  color: var(--m-text-tertiary);
}

@keyframes m-pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
  }
}
</style>
