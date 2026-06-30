import {createApp} from 'vue'
import {getToken} from './api/transport.js'

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
//  - 200：服务端开放或令牌有效 → 挂载移动端应用
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

  if (!ok) {
    const { default: Login } = await import('./Login.vue')
    const app = createConfiguredApp(Login)
    app.mount('#app')
    return
  }

  const [{ default: MobileApp }, { default: mobileRouter }] = await Promise.all([
    import('./mobile/MobileApp.vue'),
    import('./mobile/router')
  ])
  const app = createConfiguredApp(MobileApp)
  app.use(mobileRouter)
  app.mount('#app')

  if (mobileRouter.currentRoute.value.path === '/') {
    mobileRouter.replace('/mobile')
  }
}

function createConfiguredApp(rootComponent) {
  const app = createApp(rootComponent)
  app.config.errorHandler = (err) => {
    if (err && err.message && err.message.includes('ResizeObserver')) {
      return
    }
    console.error(err)
  }
  return app
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
