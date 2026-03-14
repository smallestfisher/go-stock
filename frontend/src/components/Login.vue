<template>
  <div class="login-wrapper">
    <div class="login-background"></div>
    <div class="login-content">
      <n-card class="login-card" :bordered="false">
        <div class="login-header">
          <div class="login-logo">
            <n-icon size="48" color="#18a058">
              <TrendingUpSharp />
            </n-icon>
          </div>
          <h2 class="login-title">go-stock</h2>
          <p class="login-subtitle">AI 赋能的股票分析工具</p>
        </div>

        <n-form ref="formRef" :model="formValue" :rules="rules" size="large">
          <n-form-item label="管理员账号" path="username">
            <n-input v-model:value="formValue.username" placeholder="请输入用户名" @keyup.enter="handleLogin">
              <template #prefix>
                <n-icon :component="PersonOutline" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item label="访问密码" path="password">
            <n-input
              v-model:value="formValue.password"
              type="password"
              show-password-on="mousedown"
              placeholder="请输入密码"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <n-icon :component="LockClosedOutline" />
              </template>
            </n-input>
          </n-form-item>
          
          <div class="login-options">
            <n-checkbox v-model:checked="rememberMe">记住登录状态</n-checkbox>
          </div>

          <n-button
            type="primary"
            block
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            进入系统
          </n-button>
        </n-form>

        <div class="login-footer">
          <p>© 2026 go-stock Open Source Project</p>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { 
  PersonOutline, 
  LockClosedOutline,
  TrendingUpSharp 
} from '@vicons/ionicons5'

const router = useRouter()
const message = useMessage()
const loading = ref(false)
const rememberMe = ref(true)

const formValue = ref({
  username: '',
  password: ''
})

const rules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' }
}

onMounted(() => {
  // 检查是否已有 token
  if (localStorage.getItem('auth_token')) {
    router.push('/')
  }
})

async function handleLogin() {
  if (!formValue.value.username || !formValue.value.password) {
    message.warning('请输入完整的登录信息')
    return
  }

  loading.value = true
  try {
    const response = await fetch('/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formValue.value)
    })

    const result = await response.json()
    if (response.ok) {
      localStorage.setItem('auth_token', result.token)
      // SSE 不支持 Header，必须通过 Cookie 传 token
      document.cookie = `auth_token=${result.token}; path=/; max-age=86400`;
      message.success('登录成功，正在跳转...')
      
      // 使用 href 强制全页刷新以建立带 Cookie 的 SSE 连接
      setTimeout(() => {
        window.location.href = '/';
      }, 500);
    } else {
      message.error(result.error || '验证失败，请重试')
    }
  } catch (error) {
    message.error('无法连接到服务器')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-wrapper {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  z-index: -1;
}

/* 装饰性光斑 */
.login-background::after {
  content: '';
  position: absolute;
  top: -10%;
  right: -10%;
  width: 40%;
  height: 40%;
  background: radial-gradient(circle, rgba(24, 160, 88, 0.15) 0%, transparent 70%);
  filter: blur(60px);
}

.login-content {
  width: 100%;
  max-width: 420px;
  padding: 20px;
  z-index: 1;
}

.login-card {
  border-radius: 16px;
  padding: 20px 10px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.3), 0 10px 10px -5px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  margin-bottom: 12px;
  display: inline-block;
  background: rgba(24, 160, 88, 0.1);
  padding: 12px;
  border-radius: 50%;
}

.login-title {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  letter-spacing: -0.5px;
}

.login-subtitle {
  margin: 4px 0 0;
  color: #6b7280;
  font-size: 14px;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.login-btn {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  transition: transform 0.2s;
}

.login-btn:active {
  transform: scale(0.98);
}

.login-footer {
  margin-top: 32px;
  text-align: center;
  color: #9ca3af;
  font-size: 12px;
}

:deep(.n-form-item-label) {
  font-weight: 500;
  color: #4b5563;
}

:deep(.n-input) {
  border-radius: 8px;
}
</style>
