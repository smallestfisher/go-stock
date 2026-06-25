import {createApp} from 'vue'
import AppRoot from './AppRoot.vue'
import MobileApp from './mobile/MobileApp.vue'
import Login from './Login.vue'
import desktopRouter from './router/router'
import mobileRouter from './mobile/router'
import {getToken} from './api/transport.js'
import {useDevice} from './composables/useDevice'
// 引入组件库的少量全局样式变量
import 'tdesign-vue-next/es/style/index.css';

// 彻底静音浏览器无害的 ResizeObserver 循环提示。
// 该提示是 Chromium 的自我保护（回调里又改了尺寸形成循环，浏览器主动打断），
// 不影响任何功能/数据。现有 error / unhandledrejection 两层已拦截大部分，
// 但较新 Chromium 会把它作为独立异常从 ResizeObserver 内部抛出，
// 故在原型回调外再包一层 try/catch 兜底。
;(function silenceResizeObserverLoop() {
  if (typeof window === 'undefined' || typeof ResizeObserver === 'undefined') return
  const Native = window.ResizeObserver
  if (!Native || Native.__silenced) return
  function Patched(cb) {
    const wrapped = (entries, obs) => {
      try {
        cb(entries, obs)
      } catch (e) {
        // 吞掉 "ResizeObserver loop completed with undelivered notifications."
        if (!(e && typeof e.message === 'string' && e.message.includes('ResizeObserver loop'))) {
          throw e
        }
      }
    }
    return new Native(wrapped)
  }
  Patched.prototype = Native.prototype
  Patched.__silenced = true
  window.ResizeObserver = Patched
})()

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

  const app = ok ? createApp(AppRoot) : createApp(Login)
  app.config.errorHandler = (err) => {
    if (err && err.message && err.message.includes('ResizeObserver')) {
      return
    }
    console.error(err)
  }
  if (ok) {
    // 根据设备类型使用不同的路由和App组件
    const { isMobile } = useDevice()

    if (isMobile.value) {
      // 移动端：使用独立的MobileApp（无桌面端布局）
      const mobileApp = createApp(MobileApp)
      mobileApp.use(mobileRouter)
      mobileApp.mount('#app')

      // 跳转到移动端首页
      if (mobileRouter.currentRoute.value.path === '/') {
        mobileRouter.replace('/mobile')
      }
    } else {
      // 桌面端：使用AppRoot（带侧边栏布局）
      app.use(desktopRouter)
      app.mount('#app')
    }
    return
  }
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
