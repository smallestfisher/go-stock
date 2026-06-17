// go-stock Web 端统一传输层：RPC 调用 + SSE 事件总线 + 访问令牌管理。
// 取代 Wails 的 window.go.* 绑定与 runtime.Events*。
const BASE_URL = (import.meta && import.meta.env && import.meta.env.VITE_API_BASE) || "";
const TOKEN_KEY = "go_stock_token";

export class AuthError extends Error {
  constructor(message) {
    super(message);
    this.name = "AuthError";
  }
}

export function getToken() {
  try {
    return localStorage.getItem(TOKEN_KEY) || "";
  } catch (_) {
    return "";
  }
}

export function setToken(token) {
  try {
    if (token) localStorage.setItem(TOKEN_KEY, token);
    else localStorage.removeItem(TOKEN_KEY);
  } catch (_) {}
  reconnectEvents();
}

export function clearToken() {
  setToken("");
}

export function isAuthed() {
  return !!getToken();
}

// 发起一次 RPC 调用：POST /api/rpc/{method}，body: {"args":[...]}
// 成功时 resolve 为后端返回值（与 Wails 直接返回 Go 值的语义一致）。
export async function rpc(method, args) {
  const headers = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) headers["Authorization"] = "Bearer " + token;
  let resp;
  try {
    resp = await fetch(BASE_URL + "/api/rpc/" + encodeURIComponent(method), {
      method: "POST",
      headers,
      body: JSON.stringify({ args: args || [] }),
    });
  } catch (e) {
    throw new Error("网络请求失败：" + method + " (" + (e && e.message ? e.message : e) + ")");
  }
  if (resp.status === 401) throw new AuthError("未授权：访问令牌无效或缺失");
  if (!resp.ok) {
    let msg = "RPC " + method + " 失败 (" + resp.status + ")";
    try {
      const e = await resp.json();
      if (e && e.error) msg = e.error;
    } catch (_) {}
    throw new Error(msg);
  }
  return resp.json();
}

// --- SSE 事件总线（取代 wails runtime.EventsOn/Off）---
// 后端 hub 把每条事件包成 {"event":<name>,"data":<payload>} 作为默认 message 发送，
// 这里用一个全局 EventSource 订阅 /api/events，并按 frame.event 分发给对应回调。
const listeners = new Map(); // eventName -> Set<callback>
let es = null;

export function subscribeEvent(name, cb) {
  if (!listeners.has(name)) listeners.set(name, new Set());
  listeners.get(name).add(cb);
  ensureConnected();
  return () => unsubscribeEvent(name, cb);
}

export function unsubscribeEvent(name, cb) {
  const set = listeners.get(name);
  if (set) {
    set.delete(cb);
    if (set.size === 0) listeners.delete(name);
  }
  maybeDisconnect();
}

export function unsubscribeEventAll(name) {
  listeners.delete(name);
  maybeDisconnect();
}

function ensureConnected() {
  if (es || typeof EventSource === "undefined") return;
  const token = getToken();
  const url = BASE_URL + "/api/events" + (token ? "?token=" + encodeURIComponent(token) : "");
  es = new EventSource(url);
  es.onmessage = (ev) => {
    try {
      const frame = JSON.parse(ev.data);
      const set = listeners.get(frame.event);
      if (set) set.forEach((cb) => { try { cb(frame.data); } catch (_) {} });
    } catch (_) {}
  };
  es.onerror = () => {
    // EventSource 会自动重连，这里不主动关闭
  };
}

function maybeDisconnect() {
  if (es && listeners.size === 0) {
    try { es.close(); } catch (_) {}
    es = null;
  }
}

function reconnectEvents() {
  if (es) {
    try { es.close(); } catch (_) {}
    es = null;
  }
  if (listeners.size > 0) ensureConnected();
}
