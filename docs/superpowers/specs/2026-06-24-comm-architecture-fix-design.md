# Web 通信架构刷新失效修复设计

- **日期**:2026-06-24
- **状态**:已批准,待写实现计划
- **分支**:`feat/mobile-responsive`
- **范围**:前端连接生命周期 + 定时轮询调度;后端 SSE 心跳。**不含**架构重选。

## 1. 背景

go-stock 原为 Wails 桌面 App,现转为"Go 后端 + Vue/NaiveUI 前端"的纯 web 模式。
当前 web 端出现一个故障:**行情数据看一段时间后不再自动刷新,必须手动刷新页面**。

### 1.1 现状架构(经代码确认)

- **RPC**:`POST /api/rpc/{method}`,body `{"args":[...]}`,见 `frontend/src/api/transport.js:62`。
- **SSE 推送**:`EventSource → /api/events`,后端 `server/events/hub.go` 广播事件;前端按 `frame.event` 分发(`transport.js:142`)。
  - 当前 SSE **只**承载:新闻推送(`newsPush`)、AI 流式输出(`EmitTo`)、loading/warn 消息。
- **实时行情数据不走 SSE**,而是前端用 **`setInterval` 轮询** RPC 拿,散落在 ~12 个组件里。

### 1.2 根因(经确认,对应故障场景"切走/切回标签页或电脑休眠后")

1. **完全没有页面生命周期处理**:`grep` 确认前端零处监听 `visibilitychange` / `online` / `pageshow`。浏览器对后台标签页 timer 节流/挂起,休眠或网络变化后 SSE 静默断开,切回前台不会自动重连重刷。
2. **交易时段门控把网络错误误判为"非交易时段"**:`App.vue:113` 与 `market.vue:165` 都是 `IsTradingTime().catch(() => false)`。一次网络抖动 → 判定非交易时段 → `stopTradingTimers()` 停掉轮询,且 60s 才重判。
3. **SSE 无心跳**:`hub.go` 仅在连接建立时写一次 `: connected`,空载时反代(nginx 默认 60s、Cloudflare 100s)会掐断连接。

### 1.3 决策:不重选架构

RPC + SSE 是业界标准 web 架构,选型本身正确。故障是运行时生命周期 bug,非架构不适用。
`server/DEPLOY.md:252` 本就声明"行情类事件应使用 `core.Events.Emit` 广播"——即设计者原意是行情走推送,Wails→Web 迁移时遗留为前端轮询。
**使用规模为自用/内网、1–几个标签页**,上游数据源压力不是问题,因此不做"服务端推送实时行情"(方案 C)这种较大改动(YAGNI)。

## 2. 方案

采用**方案 A + B**:既修复运行时生命周期 bug(A),又把散落的 ~12 个 `setInterval` 收敛到统一调度器(B)。

## 3. 架构

### 3.1 模块划分

```
后端                          前端
─────                         ─────────────────────────────
events/hub.go                 transport.js (改)
  + SSE 心跳 15s keepalive       + 监听 visibilitychange/online/pageshow
                                 + SSE 未连上则 reconnectEvents()
                               ┌─────────────────────────┐
                               │ scheduler.js (新增, 单例) │ ← 通用刷新调度
                               │  registerFeed(key,{      │
                               │    fetch, intervalMs,    │
                               │    activeWhen? })        │
                               └────────────┬────────────┘
                                            │ 订阅 activeWhen
                               ┌────────────┴────────────┐
                               │ marketClock.js (新增)    │ ← 唯一"是否交易时段"源
                               │  60s 查 IsTradingTime    │   容错:失败≠停轮询
                               │  → ref: anyOpen          │
                               └─────────────────────────┘
                               lifecycle.js (新增) ← 把"恢复"聚合成全局事件
```

### 3.2 职责切分(单一职责)

| 模块 | 干什么 | 不干什么 |
|------|--------|----------|
| `lifecycle.js` | 把 visibilitychange/online/pageshow 聚合成 `onResume(cb)` | 不知有"行情" |
| `transport.js` | SSE 连接 + 自身的生命周期重连 | 不知有"行情" |
| `scheduler.js` | 通用 feed 注册表;管 timer;恢复时立即刷新 | 不知"交易时段"业务 |
| `marketClock.js` | 唯一判断"是否交易时段";容错 | 不管 timer、不管 UI |
| `hub.go` | SSE keepalive | 不管业务 |

### 3.3 恢复数据流(修复核心)

```
浏览器触发 visibilitychange(visible) / online / pageshow
        │
        ├─→ transport.js: 若 es.readyState≠OPEN → reconnectEvents()
        │      └→ SSE 恢复(新闻推送、AI 流式)
        │
        └─→ scheduler.js: 遍历每个活跃 feed
              ├→ 重新 arm interval(避开浏览器对陈旧 timer 的节流)
              └→ 立即 fetch()  ← 行情数据马上刷新
        (marketClock 一直在跑;它失败不会停任何 feed)
```

## 4. 详细设计

### 4.1 `lifecycle.js`(新增,与 transport.js 同目录)

```js
// 页面恢复:可见 或 网络在线 或 bfcache 回退
const callbacks = new Set()
export function onResume(cb) { callbacks.add(cb); return () => callbacks.delete(cb) }

function fire() { callbacks.forEach(cb => { try { cb() } catch (_) {} }) }

if (typeof document !== 'undefined') {
  document.addEventListener('visibilitychange', () => { if (!document.hidden) fire() })
  window.addEventListener('online', fire)
  window.addEventListener('pageshow', e => { if (e.persisted) fire() })  // bfcache
}
```

### 4.2 `scheduler.js`(新增,单例)

```js
import { onResume } from './lifecycle.js'

const feeds = new Map()   // key -> { fetch, intervalMs, timer, activeWhen? }

export function registerFeed(key, { fetch, intervalMs, activeWhen }) { ... }
export function stopFeed(key) { ... }   // onBeforeUnmount 调用

function armTimer(f) {
  clearInterval(f.timer)
  f.timer = setInterval(() => { if (!f.activeWhen || f.activeWhen()) { try { f.fetch() } catch (_) {} } }, f.intervalMs)
}

// 恢复时:重新 arm + 立即 fetch(绕开节流的陈旧 timer)
onResume(() => {
  for (const f of feeds.values()) {
    if (f.activeWhen && !f.activeWhen()) continue
    armTimer(f)
    try { f.fetch() } catch (_) {}
  }
})
```

设计点:
- `activeWhen` 为函数引用而非布尔;`marketClock` 的 `anyOpen` 变化后,下次 tick 自然读到新值,组件无需监听交易时段。
- 每次恢复重新 `armTimer`,彻底绕开浏览器对陈旧 timer 的节流。

### 4.3 `marketClock.js`(新增,收口交易时段判定 + 修 bug)

```js
import { ref } from 'vue'
import { IsTradingTime, IsHKTradingTime, IsUSTradingTime } from '@/api/app.js'

export const anyOpen = ref(false)
let lastGood = false                 // 上次成功结果

export async function tickMarketClock() {
  try {
    const [cn, hk, us] = await Promise.all([IsTradingTime(), IsHKTradingTime(), IsUSTradingTime()])
    lastGood = !!(cn || hk || us)
  } catch (_) {
    // 网络失败:沿用 lastGood,绝不把"网络错"误判为"收盘→停轮询"
  }
  anyOpen.value = lastGood
}
```
由 scheduler 注册一个 feed 驱动:`registerFeed('marketClock', { fetch: tickMarketClock, intervalMs: 60000 })`(**不带 `activeWhen`**,始终运行——交易时段判定本身不能依赖交易时段)。页面恢复时也会立即重判一次。`market.vue`/`App.vue` 原来的 tradingCheck 逻辑改为直接读 `anyOpen.value`,删除各自的 `IsTradingTime().catch(()=>false)` 调用。

### 4.4 `transport.js`(改)

在模块顶层 import `onResume`,并在合适位置:`onResume(() => { if (es && es.readyState !== 1 /*OPEN*/) reconnectEvents() })`。
其余 SSE 逻辑不变。

### 4.5 `hub.go`(改,SSE 心跳)

`ServeSSE` 的 `for-select` 内增加 ticker:

```go
ticker := time.NewTicker(15 * time.Second)
defer ticker.Stop()
for {
  select {
  case <-ctx.Done():
    return
  case p, ok := <-s.ch:
    ... // 原逻辑
  case <-ticker.C:
    _, _ = w.Write([]byte(": ping\n\n"))
    flusher.Flush()
  }
}
```
心跳为注释行(`:` 开头),EventSource 忽略内容、仅重置"最后收到时间"。

### 4.6 组件迁移(机械替换)

旧(`market.vue` 现状):
```js
tradingCheckInterval.value = setInterval(async () => {
  const anyTrading = (await IsTradingTime().catch(() => false)) || ...  // bug
  if (anyTrading && !indexInterval.value) startTradingTimers()
  else if (!anyTrading && indexInterval.value) stopTradingTimers()
}, 60000)
indexInterval.value = setInterval(getIndex, 3000)
```

新:
```js
import { registerFeed, stopFeed } from '@/api/scheduler'
import { anyOpen } from '@/api/marketClock'

onMounted(() => {
  registerFeed('market.index',    { fetch: getIndex,  intervalMs: 3000,  activeWhen: () => anyOpen.value })
  registerFeed('market.industry', { fetch: () => { industryRank(); ReFlesh('财联社电报'); ReFlesh('新浪财经'); ReFlesh('外媒') }, intervalMs: 10000, activeWhen: () => anyOpen.value })
})
onBeforeUnmount(() => { stopFeed('market.index'); stopFeed('market.industry') })
```

需迁移的取数据轮询点(~12 处):
`stock.vue`(feishiInterval)、`market.vue`(tradingCheck/index/indexIndustryRank)、`HotStockList.vue`(fetchHotStock、check)、`App.vue`(marketStatusTimer)、`rankTable.vue`、`FundFollow.vue`(ticker/countdown 若取数据)、`stockChangesMonitor.vue`、`conceptFundFlowChart.vue`、`HotEvents.vue`。
**纯 UI 计时器不迁移**(`newsList.vue` 倒计时显示、`FundFollow` popover 延迟等),只迁移真正取数据的轮询。

## 5. 错误处理与降级

| 故障 | 行为 | 用户感知 |
|------|------|----------|
| 单次 `feed.fetch()` 抛错 | scheduler try/catch 吞掉,timer 照常跑,下 tick 重试 | 本次未更新,下次自动恢复;绝不永久停 |
| SSE 断开 | `onResume` + EventSource 自带重连;心跳加速发现 | 推送类暂缺,恢复后补 |
| 交易时段查询失败 | `marketClock` 沿用 `lastGood`,不停任何 feed | 行情照刷(可能多刷几次非交易时段) |
| scheduler 整个挂 | feed 不调度 → 数据静止 | 与现状 bug 症状相同,不会更差 |

**核心保证**:恢复路径依赖浏览器原生事件,不依赖网络。标签页切回前台必触发重连 + 立即刷新。

### 5.1 交易时段"刷"与"不刷"

保留"收盘→停轮询"意图(`activeWhen: () => anyOpen.value`)。收盘时标签页在前台不会困扰用户(本就该停);问题仅在"开市中 + 切走切回",该门控既省请求又不影响修复目标。

## 6. 测试策略

遵循项目规则:**不写单元测试**(dev-workflow-rules)。验证靠手动 + 现有集成测试:

- 现有 `server/core_test.go`、`server/events/hub_test.go`、`server/modules/*_test.go`。hub.go 心跳改动需确认不破坏 `hub_test.go`。
- 手动验证(改完一起跑):
  1. 开页盯 2 分钟,行情动 → 切走 5 分钟/锁屏 → 切回,**30 秒内数据恢复跳动**。
  2. 拔网 10s 再恢复 → SSE/行情恢复。
  3. 收盘时段验证不误触发"停轮询";开市再验恢复。

## 7. 渐进迁移顺序(每步可独立验证,可随时停)

| 步 | 改动 | 解决 | 完成标志 |
|----|------|------|----------|
| 1 | `lifecycle.js` + `transport.js` 恢复重连 + SSE 心跳 | SSE 推送类切回不恢复 | ✅ 已见效 |
| 2 | `marketClock.js`(收口 + 修 `catch(()=>false)`) | 网络抖动误停轮询 | ✅ |
| 3 | `scheduler.js` + 迁移 `market.vue`、`stock.vue` 两个高频组件 | **行情切回不刷新(主诉)** | ✅ 核心修复完成 |
| 4 | 迁移剩余 ~8 个组件 | 收口所有 timer | ✅ 收敛清理 |

第 3 步完成后,用户最初报告的 bug 即修复;第 4 步为收敛清理,纯收益。

## 8. 不做的事(YAGNI)

- 服务端推送实时行情(方案 C):自用/内网规模不需要。
- 引入 WebSocket:SSE 已满足推送需求且已实现。
- 重写 RPC 层:现状 RPC 工作正常。
- 迁移纯 UI 计时器。
