<template>
  <div class="login-wrap">
    <div class="login-card">
      <div class="title">📈 go-stock</div>
      <div class="subtitle">远程访问 · 请输入访问令牌</div>
      <input
        v-model="token"
        class="login-input"
        type="password"
        placeholder="访问令牌 (GO_STOCK_TOKEN)"
        @keyup.enter="submit"
      />
      <button class="login-btn" type="button" :disabled="loading" @click="submit">
        {{ loading ? '登录中...' : '登录' }}
      </button>
      <div v-if="errMsg" class="err">{{ errMsg }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
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
  padding: 16px;
}

.login-card {
  width: 100%;
  max-width: 360px;
  background: #23232a;
  border-radius: 12px;
  padding: 28px 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
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

.login-input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  margin-bottom: 12px;
  border: 1px solid #3a3a42;
  border-radius: 8px;
  background: #18181c;
  color: #fff;
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
}

.login-input:focus {
  border-color: #18a058;
}

.login-input::placeholder {
  color: #6a6f76;
}

.login-btn {
  width: 100%;
  height: 44px;
  border: none;
  border-radius: 8px;
  background: #18a058;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
}

.login-btn:active:not(:disabled) {
  background: #158048;
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.err {
  color: #e06c75;
  font-size: 13px;
  margin-top: 10px;
}
</style>
