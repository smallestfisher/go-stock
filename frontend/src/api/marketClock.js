// 全局唯一的"是否交易时段"判定源。
//
// 收口原本散落在 App.vue / market.vue / HotStockList.vue 各自的
// IsTradingTime().catch(()=>false) 调用，并修复其核心 bug：
// 网络失败时沿用上次成功结果，绝不把"网络错"误判成"收盘 → 停轮询"。
//
// 导入即自启动：每 60s 刷新一次交易时段状态；页面恢复时由 scheduler 立即重判。
import { ref } from "vue";
import { IsTradingTime, IsHKTradingTime, IsUSTradingTime } from "./app";
import { registerFeed } from "./scheduler";

export const cnOpen = ref(false); // A股是否开市
export const hkOpen = ref(false); // 港股是否开市
export const usOpen = ref(false); // 美股是否开市
export const anyOpen = ref(false); // 任意一个市场开市

// 上次成功查询的结果（网络抖动时保留，避免状态跳变）
let lastCn = false;
let lastHk = false;
let lastUs = false;

export async function tickMarketClock() {
  try {
    const [cn, hk, us] = await Promise.all([
      IsTradingTime(),
      IsHKTradingTime(),
      IsUSTradingTime(),
    ]);
    lastCn = !!cn;
    lastHk = !!hk;
    lastUs = !!us;
  } catch (_) {
    // 网络失败：保留 last* 不变，绝不停任何依赖此状态的轮询
  }
  cnOpen.value = lastCn;
  hkOpen.value = lastHk;
  usOpen.value = lastUs;
  anyOpen.value = lastCn || lastHk || lastUs;
}

// 自启动（ES 模块只会求值一次，故只注册一个 feed）
registerFeed("marketClock", { fetch: tickMarketClock, intervalMs: 60000 });
tickMarketClock(); // 首次立即拉取，让 anyOpen 尽快可用
