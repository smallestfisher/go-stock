// go-stock Web 端统一传输层：RPC 调用 + SSE 事件总线 + 访问令牌管理。
//
import { onResume } from "./lifecycle";

const BASE_URL = (import.meta && import.meta.env && import.meta.env.VITE_API_BASE) || "";
const TOKEN_KEY = "go_stock_token";
const CLIENT_ID_KEY = "go_stock_client_id";

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

export function serverURL(path) {
  const origin = typeof window !== "undefined" ? window.location.origin : "http://localhost";
  return new URL(BASE_URL + path, origin).toString();
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

export function getClientId() {
  try {
    let clientId = sessionStorage.getItem(CLIENT_ID_KEY) || "";
    if (!clientId) {
      if (typeof crypto !== "undefined" && crypto.randomUUID) {
        clientId = crypto.randomUUID();
      } else {
        clientId = Date.now().toString(36) + "-" + Math.random().toString(36).slice(2);
      }
      sessionStorage.setItem(CLIENT_ID_KEY, clientId);
    }
    return clientId;
  } catch (_) {
    return "";
  }
}

// 发起一次 RPC 调用：POST /api/rpc/{method}，body: {"args":[...]}
// 成功时 resolve 为后端返回值。
export async function rpc(method, args) {
  const headers = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) headers["Authorization"] = "Bearer " + token;
  const clientId = getClientId();
  if (clientId) headers["X-Go-Stock-Client-Id"] = clientId;
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

// --- 事件总线（浏览器本地事件 + 后端 SSE 推送）---
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

export function clearEventListeners() {
  listeners.clear();
  maybeDisconnect();
}

export function emitLocalEvent(name, ...data) {
  const set = listeners.get(name);
  if (!set) return;
  set.forEach((cb) => {
    try {
      cb(...data);
    } catch (_) {}
  });
}

function ensureConnected() {
  if (es || typeof EventSource === "undefined") return;
  const token = getToken();
  const params = new URLSearchParams();
  if (token) params.set("token", token);
  const clientId = getClientId();
  if (clientId) params.set("clientId", clientId);
  const qs = params.toString();
  const url = BASE_URL + "/api/events" + (qs ? "?" + qs : "");
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

// 页面恢复（切回前台 / 网络恢复 / bfcache 回退）时，若 SSE 未处于 OPEN 则重连。
// 这是"切走/切回或休眠后推送类事件（新闻/AI 流式）不再更新"的修复。
onResume(() => {
  if (!es || es.readyState !== 1 /* EventSource.OPEN */) {
    reconnectEvents();
  }
});

export function openURL(url) {
  if (typeof window !== "undefined" && url) {
    window.open(url, "_blank", "noopener,noreferrer");
  }
  return url;
}

export function downloadBase64File(filename, base64, mimeType) {
  if (typeof document === "undefined") {
    return "";
  }
  const clean = String(base64 || "").replace(/^data:[^;]+;base64,/, "");
  const byteChars = atob(clean);
  const chunks = [];
  for (let offset = 0; offset < byteChars.length; offset += 8192) {
    const slice = byteChars.slice(offset, offset + 8192);
    const bytes = new Uint8Array(slice.length);
    for (let i = 0; i < slice.length; i++) {
      bytes[i] = slice.charCodeAt(i);
    }
    chunks.push(bytes);
  }
  const blob = new Blob(chunks, { type: mimeType || "application/octet-stream" });
  const link = document.createElement("a");
  link.href = URL.createObjectURL(blob);
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  URL.revokeObjectURL(link.href);
  link.remove();
  return filename;
}
