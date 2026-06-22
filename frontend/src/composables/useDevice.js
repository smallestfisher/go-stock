import {ref} from 'vue'

// 模块级单例：全应用共享同一份设备状态，避免各组件各自 matchMedia 重复监听
const isMobile = ref(false)
const isTablet = ref(false)
const isDesktop = ref(true)

const MOBILE_QUERY = '(max-width: 768px)'
const TABLET_QUERY = '(min-width: 769px) and (max-width: 1024px)'

let mobileMql = null
let tabletMql = null
let initialized = false

function sync() {
    isMobile.value = mobileMql ? mobileMql.matches : false
    isTablet.value = tabletMql ? tabletMql.matches : false
    isDesktop.value = !isMobile.value && !isTablet.value
}

function ensureInit() {
    if (initialized || typeof window === 'undefined' || !window.matchMedia) {
        return
    }
    mobileMql = window.matchMedia(MOBILE_QUERY)
    tabletMql = window.matchMedia(TABLET_QUERY)
    sync()
    // 优先使用 addEventListener，旧 Safari 降级到 addListener
    if (mobileMql.addEventListener) {
        mobileMql.addEventListener('change', sync)
        tabletMql.addEventListener('change', sync)
    } else if (mobileMql.addListener) {
        mobileMql.addListener(sync)
        tabletMql.addListener(sync)
    }
    initialized = true
}

ensureInit()

/**
 * 全局设备状态 composable。
 * 返回响应式 isMobile / isTablet / isDesktop，监听在模块加载时建立一次。
 * 在 SSR / 预渲染环境下安全降级为桌面。
 */
export function useDevice() {
    return {isMobile, isTablet, isDesktop}
}
