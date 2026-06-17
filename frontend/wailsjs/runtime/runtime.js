/*
  go-stock Web 模式下的 wails runtime 替代实现。
  - EventsOn/Off/Emit/Eonce：走 transport 的 SSE 事件总线（/api/events）
  - BrowserOpenURL：window.open 打开新标签
  - Environment：返回固定的 web 环境
  - 窗口、通知、剪贴板等桌面能力：置为安全的 no-op（浏览器无对应概念）
  导出名与原 wails runtime 保持一致，确保前端 import 全部可解析。
*/
import { subscribeEvent, unsubscribeEvent, unsubscribeEventAll } from "../../src/api/transport.js";

export function LogPrint() {}
export function LogTrace() {}
export function LogDebug() {}
export function LogInfo() {}
export function LogWarning() {}
export function LogError() {}
export function LogFatal() {}

// 监听事件，返回一个取消函数（与 wails EventsOn 返回 cancel 的语义一致）。
export function EventsOnMultiple(eventName, callback, maxCallbacks) {
  if (maxCallbacks === 1) {
    const off = subscribeEvent(eventName, (data) => { off(); callback(data); });
    return off;
  }
  return subscribeEvent(eventName, callback);
}

export function EventsOn(eventName, callback) {
  return EventsOnMultiple(eventName, callback, -1);
}

export function EventsOnce(eventName, callback) {
  return EventsOnMultiple(eventName, callback, 1);
}

export function EventsOff(eventName, ...additionalEventNames) {
  const names = [eventName, ...additionalEventNames];
  names.forEach((n) => unsubscribeEventAll(n));
}

export function EventsOffAll() {
  // Web 端无法枚举所有事件名，安全 no-op
}

// Web 端前端发事件后端不监听，忽略。
export function EventsEmit() {}

export function BrowserOpenURL(url) {
  if (typeof window !== "undefined" && url) {
    window.open(url, "_blank");
  }
}

export function Environment() {
  return Promise.resolve({ buildType: "production", platform: "web", arch: "web" });
}

// --- 桌面专属能力，Web 下为 no-op / 兜底返回 ---
function noop() {}
function noopPromise(v) { return () => Promise.resolve(v); }

export const WindowReload = noop;
export const WindowReloadApp = noop;
export const WindowSetAlwaysOnTop = noop;
export const WindowSetSystemDefaultTheme = noop;
export const WindowSetLightTheme = noop;
export const WindowSetDarkTheme = noop;
export const WindowCenter = noop;
export const WindowSetTitle = noop;
export const WindowFullscreen = noop;
export const WindowUnfullscreen = noop;
export const WindowIsFullscreen = noopPromise(false);
export const WindowGetSize = noopPromise({ w: 0, h: 0 });
export const WindowSetSize = noop;
export const WindowSetMaxSize = noop;
export const WindowSetMinSize = noop;
export const WindowSetPosition = noop;
export const WindowGetPosition = noopPromise({ x: 0, y: 0 });
export const WindowHide = noop;
export const WindowShow = noop;
export const WindowMaximise = noop;
export const WindowToggleMaximise = noop;
export const WindowUnmaximise = noop;
export const WindowIsMaximised = noopPromise(false);
export const WindowMinimise = noop;
export const WindowUnminimise = noop;
export const WindowSetBackgroundColour = noop;
export const ScreenGetAll = noopPromise([]);
export const WindowIsMinimised = noopPromise(false);
export const WindowIsNormal = noopPromise(true);
export const Quit = noop;
export const Hide = noop;
export const Show = noop;
export const ClipboardGetText = noopPromise("");
export const ClipboardSetText = noop;
export const OnFileDrop = noop;
export const OnFileDropOff = noop;
export const CanResolveFilePaths = noopPromise(false);
export const ResolveFilePaths = noopPromise([]);
export const InitializeNotifications = noop;
export const CleanupNotifications = noop;
export const IsNotificationAvailable = noopPromise(false);
export const RequestNotificationAuthorization = noopPromise(true);
export const CheckNotificationAuthorization = noopPromise(true);
export const SendNotification = noop;
export const SendNotificationWithActions = noop;
export const RegisterNotificationCategory = noop;
export const RemoveNotificationCategory = noop;
export const RemoveAllPendingNotifications = noop;
export const RemovePendingNotification = noop;
export const RemoveAllDeliveredNotifications = noop;
export const RemoveDeliveredNotification = noop;
export const RemoveNotification = noop;
