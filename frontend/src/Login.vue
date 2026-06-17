<template>
  <div class="login-wrap">
    <n-config-provider :theme="darkTheme">
      <n-card class="login-card" :bordered="false" size="large">
        <div class="title">📈 go-stock</div>
        <div class="subtitle">远程访问 · 请输入访问令牌</div>
        <n-space vertical :size="12">
          <n-input
            v-model:value="token"
            type="password"
            show-password-on="click"
            placeholder="访问令牌 (GO_STOCK_TOKEN)"
            @keyup.enter="submit"
          />
          <n-button type="primary" block :loading="loading" @click="submit">登录</n-button>
          <div v-if="errMsg" class="err">{{ errMsg }}</div>
        </n-space>
      </n-card>
    </n-config-provider>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { darkTheme } from 'naive-ui'
import { setToken } from './api/transport.js'

const token = ref('')
const loading = ref(false)
const errMsg = ref('')

async function submit() {
  const t = token.value.trim()
  if (!t) {
    errMsg.value = '请输入访问令牌'
    return
  }
  loading.value = true
  errMsg.value = ''
  // 先存令牌，再用它探测 /api/health 校验
  setToken(t)
  try {
    const r = await fetch('/api/health', { headers: { Authorization: 'Bearer ' + t } })
    if (r.ok) {
      location.reload()
    } else {
      errMsg.value = '令牌无效，请重试'
      loading.value = false
    }
  } catch (e) {
    errMsg.value = '无法连接服务器'
    loading.value = false
  }
}
</script>

<style scoped>
.login-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #18181c;
}
.login-card {
  width: 360px;
  background: #23232a;
}
.title {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 4px;
  color: #fff;
}
.subtitle {
  font-size: 13px;
  color: #9aa0a6;
  margin-bottom: 18px;
}
.err {
  color: #e06c75;
  font-size: 13px;
}
</style>
