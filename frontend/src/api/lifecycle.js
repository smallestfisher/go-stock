// 页面生命周期协调：把"回到前台 / 网络恢复 / bfcache 回退"聚合成一个全局 onResume 事件。
//
// 这是"切走/切回标签页或电脑休眠后数据不再自动刷新"修复的根基：
// 恢复信号来自浏览器原生事件，不依赖网络，因此只要标签页切回前台就一定会触发。
const callbacks = new Set();

// 注册一个"页面恢复"回调，返回取消注册函数。
export function onResume(cb) {
  callbacks.add(cb);
  return () => callbacks.delete(cb);
}

function fire() {
  callbacks.forEach((cb) => {
    try {
      cb();
    } catch (_) {
      // 单个回调异常不影响其他回调
    }
  });
}

let bound = false;
function bindOnce() {
  if (bound || typeof document === "undefined") return;
  bound = true;
  // 标签页回到前台（从 hidden → visible）
  document.addEventListener("visibilitychange", () => {
    if (!document.hidden) fire();
  });
  // 网络从离线恢复在线
  window.addEventListener("online", fire);
  // bfcache 前进/后退回退（persisted=true 表示从 bfcache 恢复）
  window.addEventListener("pageshow", (e) => {
    if (e.persisted) fire();
  });
}

bindOnce();
