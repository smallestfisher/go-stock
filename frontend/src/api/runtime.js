/*
  go-stock Web 模式下的浏览器运行时工具。
  - EventsOn/Off/Emit/Once：EventsEmit 为浏览器本地事件，后端推送走 SSE（/api/events）
  - BrowserOpenURL：window.open 打开新标签
*/
import { clearEventListeners, emitLocalEvent, openURL, subscribeEvent, unsubscribeEventAll } from "./transport.js";

export function LogPrint() {}
export function LogTrace() {}
export function LogDebug() {}
export function LogInfo() {}
export function LogWarning() {}
export function LogError() {}
export function LogFatal() {}

// 监听事件，返回一个取消函数。
export function EventsOnMultiple(eventName, callback, maxCallbacks) {
  if (maxCallbacks === 1) {
    const off = subscribeEvent(eventName, (...data) => { off(); callback(...data); });
    return off;
  }
  if (maxCallbacks > 1) {
    let count = 0;
    const off = subscribeEvent(eventName, (...data) => {
      count += 1;
      if (count >= maxCallbacks) off();
      callback(...data);
    });
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
  clearEventListeners();
}

export function EventsEmit(eventName, ...data) {
  emitLocalEvent(eventName, ...data);
}

export function BrowserOpenURL(url) {
  openURL(url);
}
