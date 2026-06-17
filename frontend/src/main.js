import {createApp} from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import Login from './Login.vue'
import router from './router/router'
import {getToken} from './api/transport.js'
// 引入组件库的少量全局样式变量
import 'tdesign-vue-next/es/style/index.css';

// 启动闸门：探测 /api/health。
//  - 200：服务端开放或令牌有效 → 挂载完整应用
//  - 401：需要鉴权 → 挂载登录页，登录成功后重载进入应用
async function bootstrap() {
  const token = getToken()
  const headers = {}
  if (token) headers["Authorization"] = "Bearer " + token

  let ok = false
  try {
    const r = await fetch("/api/health", {headers})
    ok = r.ok
  } catch (_) {
    ok = false
  }

  const app = ok ? createApp(App) : createApp(Login)
  app.config.errorHandler = (err) => {
    if (err && err.message && err.message.includes('ResizeObserver')) {
      return
    }
    console.error(err)
  }
  if (ok) {
    app.use(router)
  }
  app.use(naive)
  app.mount('#app')
}

window.addEventListener('error', (event) => {
  if (event.message && event.message.includes('ResizeObserver')) {
    event.preventDefault()
    return true
  }
})

window.addEventListener('unhandledrejection', (event) => {
  if (event.reason && event.reason.message && event.reason.message.includes('ResizeObserver')) {
    event.preventDefault()
    return true
  }
})

bootstrap()
