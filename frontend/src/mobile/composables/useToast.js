/**
 * 移动端全局轻量 Toast
 * 移动端无 naive-ui 的 useMessage/useDialog，统一用本模块替代原生 alert 做即时反馈。
 * 单例：MobileApp.vue 挂载 <MToastHost/> 订阅 toasts，任意页面 import { toast } 调用。
 */
import { reactive } from 'vue'

// 全局响应式 toast 队列，MToastHost 渲染它
export const toasts = reactive([])

let seed = 0

/**
 * 弹出一条 toast
 * @param {string} message 文本
 * @param {'info'|'success'|'error'|'warning'} type 类型
 * @param {number} duration 毫秒，默认 2000
 */
export function showToast(message, type = 'info', duration = 2000) {
  if (!message) return
  const id = ++seed
  toasts.push({ id, message: String(message), type })
  setTimeout(() => {
    const i = toasts.findIndex(t => t.id === id)
    if (i > -1) toasts.splice(i, 1)
  }, duration)
}

// 语义化快捷方法，贴近 naive-ui message.* 习惯，方便从桌面端逻辑迁移
export const toast = {
  info: (msg, d) => showToast(msg, 'info', d),
  success: (msg, d) => showToast(msg, 'success', d),
  error: (msg, d) => showToast(msg, 'error', d || 2600),
  warning: (msg, d) => showToast(msg, 'warning', d),
}

export function useToast() {
  return { toast, showToast }
}
