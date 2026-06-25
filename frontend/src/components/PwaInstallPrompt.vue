<script setup>
import {ref, onMounted} from 'vue'
import {PhoneLandscapeOutline} from '@vicons/ionicons5'

const deferredPrompt = ref(null)
const showPrompt = ref(false)

onMounted(() => {
  // 检查是否已安装
  if (window.matchMedia('(display-mode: standalone)').matches) {
    return // 已安装，不显示提示
  }

  // 检查本地存储（用户是否已拒绝）
  const dismissed = localStorage.getItem('pwa-install-dismissed')
  if (dismissed) return

  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault()
    deferredPrompt.value = e

    // 延迟 3 秒显示，避免首次访问立即打扰
    setTimeout(() => {
      showPrompt.value = true
    }, 3000)
  })
})

async function install() {
  if (!deferredPrompt.value) return

  deferredPrompt.value.prompt()
  const {outcome} = await deferredPrompt.value.userChoice

  if (outcome === 'accepted') {
    console.log('PWA 已安装')
  }

  showPrompt.value = false
  deferredPrompt.value = null
}

function dismiss() {
  showPrompt.value = false
  localStorage.setItem('pwa-install-dismissed', 'true')
}
</script>

<template>
  <n-alert
    v-if="showPrompt"
    type="info"
    closable
    @close="dismiss"
    class="pwa-install-alert"
  >
    <template #header>
      <n-space align="center" :size="8">
        <n-icon size="20">
          <PhoneLandscapeOutline />
        </n-icon>
        <span>安装 go-stock</span>
      </n-space>
    </template>
    添加到主屏幕，获得类原生 App 体验
    <template #action>
      <n-button size="small" type="primary" @click="install">
        立即安装
      </n-button>
    </template>
  </n-alert>
</template>

<style scoped>
.pwa-install-alert {
  position: fixed;
  bottom: calc(var(--mobile-bottom-nav-height) + 12px + env(safe-area-inset-bottom));
  left: 12px;
  right: 12px;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

@media (min-width: 769px) {
  .pwa-install-alert {
    bottom: calc(var(--desktop-bottom-menu-height) + 12px);
    left: auto;
    right: 24px;
    max-width: 400px;
  }
}
</style>
