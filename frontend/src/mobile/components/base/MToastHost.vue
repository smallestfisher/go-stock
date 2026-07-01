<script setup>
// 全局 Toast 渲染宿主，挂在 MobileApp 根部。订阅 useToast 的全局队列。
import { toasts } from '../../composables/useToast'
import MIcon from './MIcon.vue'
</script>

<template>
  <Teleport to="body">
    <div class="m-toast-host">
      <TransitionGroup name="m-toast">
        <div
          v-for="t in toasts"
          :key="t.id"
          class="m-toast"
          :class="`m-toast--${t.type}`"
        >
          <span class="m-toast__icon">
            <template v-if="t.type === 'success'"><MIcon name="check" :size="12" /></template>
            <template v-else-if="t.type === 'error'"><MIcon name="close" :size="12" /></template>
            <template v-else-if="t.type === 'warning'">!</template>
            <template v-else>i</template>
          </span>
          <span class="m-toast__text">{{ t.message }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.m-toast-host {
  position: fixed;
  top: calc(var(--m-safe-top) + 12%);
  left: 0;
  right: 0;
  z-index: var(--m-z-toast);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--m-space-sm);
  pointer-events: none;
  padding: 0 var(--m-space-xl);
}

.m-toast {
  display: inline-flex;
  align-items: center;
  gap: var(--m-space-sm);
  max-width: 100%;
  padding: var(--m-space-sm) var(--m-space-lg);
  border-radius: var(--m-radius-full);
  background: rgba(30, 34, 39, 0.92);
  color: #fff;
  font-size: var(--m-font-sm);
  line-height: var(--m-line-height-normal);
  box-shadow: var(--m-shadow-lg);
  backdrop-filter: blur(4px);
}

.m-toast__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: var(--m-radius-full);
  font-size: 11px;
  font-weight: var(--m-font-weight-bold);
  flex-shrink: 0;
}

.m-toast__text {
  word-break: break-all;
}

.m-toast--success .m-toast__icon { background: var(--m-color-fall); }
.m-toast--error .m-toast__icon { background: var(--m-color-rise); }
.m-toast--warning .m-toast__icon { background: #f0a020; color: #1e2227; }
.m-toast--info .m-toast__icon { background: #2080f0; }

.m-toast-enter-active,
.m-toast-leave-active {
  transition: opacity var(--m-duration-normal) var(--m-ease-out),
    transform var(--m-duration-normal) var(--m-ease-out);
}

.m-toast-enter-from,
.m-toast-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}
</style>
