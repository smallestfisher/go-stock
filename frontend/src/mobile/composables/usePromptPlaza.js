/**
 * 提示词广场 / 问答广场 共享账号体系与 REST 封装
 *
 * 广场与问答共用同一套独立账号（promptPlazaToken，与 go-stock 主登录无关），
 * 后端通过 /api/prompt-plaza/* 反向代理到广场服务。
 * 桌面端在 promptPlaza.vue / promptQa.vue 里各自重复实现了一遍，这里抽成共享单例，
 * 移动端两页统一使用。
 */
import { ref, computed } from 'vue'
import { parsePromptPlazaResponse, promptPlazaHeaders, promptPlazaURL } from '../../api/promptPlaza'

const TOKEN_KEY = 'promptPlazaToken'
const USERNAME_KEY = 'promptPlazaUsername'
const PASSWORD_KEY = 'promptPlazaPassword'

// 单例状态：token / 当前用户在两页间共享
const token = ref(localStorage.getItem(TOKEN_KEY) || '')
const currentUser = ref(null)

const isLoggedIn = computed(() => !!token.value)

function getHeaders() {
  return promptPlazaHeaders(token.value)
}

async function apiGet(path, params = {}) {
  const resp = await fetch(promptPlazaURL(path, params), { headers: getHeaders() })
  return parsePromptPlazaResponse(resp)
}

async function apiPost(path, body = null) {
  const resp = await fetch(promptPlazaURL(path), {
    method: 'POST',
    headers: getHeaders(),
    body: body ? JSON.stringify(body) : null,
  })
  return parsePromptPlazaResponse(resp)
}

async function apiPut(path, body) {
  const resp = await fetch(promptPlazaURL(path), {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(body),
  })
  return parsePromptPlazaResponse(resp)
}

async function apiDelete(path) {
  const resp = await fetch(promptPlazaURL(path), {
    method: 'DELETE',
    headers: getHeaders(),
  })
  return parsePromptPlazaResponse(resp)
}

// 拉取当前登录用户；token 失效则清空
async function fetchCurrentUser() {
  if (!token.value) {
    currentUser.value = null
    return null
  }
  try {
    const data = await apiGet('/user/me')
    currentUser.value = data
    return data
  } catch (e) {
    token.value = ''
    localStorage.removeItem(TOKEN_KEY)
    currentUser.value = null
    return null
  }
}

async function login(username, password) {
  const data = await apiPost('/auth/login', { username, password })
  token.value = data.token
  localStorage.setItem(TOKEN_KEY, data.token)
  localStorage.setItem(USERNAME_KEY, username)
  localStorage.setItem(PASSWORD_KEY, password)
  currentUser.value = data.user
  return data
}

async function register(username, password, nickname) {
  const data = await apiPost('/auth/register', { username, password, nickname })
  token.value = data.token
  localStorage.setItem(TOKEN_KEY, data.token)
  localStorage.setItem(USERNAME_KEY, username)
  localStorage.setItem(PASSWORD_KEY, password)
  currentUser.value = data.user
  return data
}

function logout() {
  token.value = ''
  localStorage.removeItem(TOKEN_KEY)
  currentUser.value = null
}

// 显示名：昵称优先，其次用户名
function displayName(user) {
  const u = user || currentUser.value
  return (u && (u.nickname || u.username)) || '已登录'
}

// 取记住的用户名（登录表单回填）
function rememberedUsername() {
  return localStorage.getItem(USERNAME_KEY) || ''
}

function formatTime(timeStr) {
  if (!timeStr) return ''
  return String(timeStr).substring(0, 19).replace('T', ' ')
}

function timeAgo(timeStr) {
  if (!timeStr) return ''
  const now = new Date()
  const time = new Date(timeStr)
  const diff = Math.floor((now - time) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return Math.floor(diff / 60) + '分钟前'
  if (diff < 86400) return Math.floor(diff / 3600) + '小时前'
  if (diff < 2592000) return Math.floor(diff / 86400) + '天前'
  return formatTime(timeStr)
}

// 复制文本到剪贴板（带降级）
async function copyText(content) {
  if (navigator.clipboard) {
    await navigator.clipboard.writeText(content)
    return
  }
  const textarea = document.createElement('textarea')
  textarea.value = content
  document.body.appendChild(textarea)
  textarea.select()
  document.execCommand('copy')
  document.body.removeChild(textarea)
}

export function usePromptPlaza() {
  return {
    // 状态
    token,
    currentUser,
    isLoggedIn,
    // REST
    apiGet,
    apiPost,
    apiPut,
    apiDelete,
    // 账号
    fetchCurrentUser,
    login,
    register,
    logout,
    displayName,
    rememberedUsername,
    // 工具
    formatTime,
    timeAgo,
    copyText,
  }
}
