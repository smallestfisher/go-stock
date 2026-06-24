// 通用刷新调度器（单例）。接管所有定时轮询：
//  - 统一管理 timer，避免散落在十几个组件里各自 setInterval；
//  - 页面恢复（lifecycle.onResume）时重新 arm timer 并立即 fetch，
//    绕开浏览器对后台标签页/休眠后陈旧 timer 的节流。
//
// 单次 fetch 失败只吞掉、不打断调度（下个 tick 再试），绝不"永久停"。
import { onResume } from "./lifecycle";

// key -> feed
// feed: { fetch, intervalMs, timer, activeWhen? }
const feeds = new Map();

function armTimer(f) {
  if (f.timer) clearInterval(f.timer);
  f.timer = setInterval(() => {
    // activeWhen 返回 false（如已收盘）时跳过本次拉取
    if (f.activeWhen && !f.activeWhen()) return;
    try {
      f.fetch();
    } catch (_) {
      // 单次拉取失败不打断调度
    }
  }, f.intervalMs);
}

// 注册一个数据源。fetch 为同步或异步函数；intervalMs 为轮询间隔(ms)；
// activeWhen 可选，返回 false 时暂停拉取（通常由 marketClock 的 anyOpen 驱动）。
// 返回停止函数。注册不自动拉取，需要初始数据请在组件里显式调一次。
export function registerFeed(key, opts) {
  stopFeed(key);
  const f = {
    fetch: opts.fetch,
    intervalMs: opts.intervalMs,
    timer: null,
    activeWhen: opts.activeWhen,
  };
  feeds.set(key, f);
  armTimer(f);
  return () => stopFeed(key);
}

// 停止并移除一个数据源。
export function stopFeed(key) {
  const f = feeds.get(key);
  if (!f) return;
  if (f.timer) clearInterval(f.timer);
  f.timer = null;
  feeds.delete(key);
}

// 立即触发一次指定 feed 的拉取（满足 activeWhen 时）。
export function refreshNow(key) {
  const f = feeds.get(key);
  if (!f) return;
  if (f.activeWhen && !f.activeWhen()) return;
  try {
    f.fetch();
  } catch (_) {}
}

// 页面恢复：重新 arm 所有活跃 feed 的 timer（避开被节流的陈旧 timer），并各拉取一次。
onResume(() => {
  for (const f of feeds.values()) {
    armTimer(f);
    if (f.activeWhen && !f.activeWhen()) continue;
    try {
      f.fetch();
    } catch (_) {}
  }
});
