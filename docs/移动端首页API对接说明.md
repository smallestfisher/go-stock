# 移动端首页API对接说明

## 📅 更新时间
2026/06/25

## 🎯 本次修改范围
`frontend/src/mobile/pages/HomePage.vue` - 移动端首页数据对接

---

## ✅ 已对接的功能模块

### 1. **自选股票概览卡片**
- **API**: `GetFollowList(page)`
- **功能**: 获取用户自选股票列表，首页显示前3条
- **数据字段**: 股票代码、名称、价格、涨跌幅、涨跌额
- **轮询**: 每10秒刷新一次

### 2. **市场快讯卡片**
- **API**: `GetTelegraphList('财联社电报')`
- **功能**: 获取财联社实时快讯
- **数据字段**: 标题、时间、来源
- **显示**: 最新3条快讯
- **轮询**: 每30秒刷新一次

### 3. **实时热点卡片**
- **API**: `HotTopic(10)`
- **功能**: 获取当前热门话题/概念
- **数据字段**: 话题名称、热度值、涨跌幅、相关股票数
- **显示**: 热度Top 3
- **轮询**: 每60秒刷新一次
- **特性**: 自动格式化热度值（K/M单位）

### 4. **异动监控卡片**
- **API**: `GetStockChanges([0], 1, 5)`
- **参数说明**:
  - `[0]`: 异动类型数组，0表示全部类型
  - `1`: 页码
  - `5`: 每页条数
- **功能**: 获取实时股票异动数据
- **支持的异动类型**:
  - 1: 涨停 (limit_up)
  - 2: 跌停 (limit_down)
  - 3: 急涨 (rapid_rise)
  - 4: 急跌 (rapid_fall)
  - 5: 放量 (high_volume)
  - 6: 突破 (breakthrough)
- **显示**: 最新2条异动

### 5. **行业热度卡片**
- **API**: `GetIndustryRank(market, date)`
- **参数**:
  - `market`: 0=A股, 1=港股, 2=美股
  - `date`: 日期格式 YYYYMMDD
- **功能**: 获取行业涨跌排名
- **数据字段**: 行业名称、涨跌幅、领涨个股
- **显示**: Top 5行业

### 6. **市场统计**
- **API**: `GetTodayMarketStatistic()`
- **功能**: 获取今日市场整体统计数据（传递给MarketStatusBar组件）

---

## 🔄 数据刷新机制

使用项目统一的调度器 `registerFeed` / `stopFeed` 实现轮询：

```javascript
// 注册轮询任务
registerFeed('mobile-home-stocks', {
  fetch: loadStockData,
  intervalMs: 10000  // 10秒
})

// 组件卸载时停止轮询
onBeforeUnmount(() => {
  stopFeed('mobile-home-stocks')
})
```

### 刷新频率策略
- **自选股票**: 10秒 - 价格变化快，需频繁更新
- **市场快讯**: 30秒 - 新闻更新较频繁
- **热点话题**: 60秒 - 热度变化相对缓慢
- **异动/行业**: 不自动刷新，仅在下拉刷新时更新（减少服务器压力）

---

## 🎨 数据格式映射

### 热度值格式化
```javascript
function formatHeat(heat) {
  if (typeof heat === 'string') return heat
  if (heat >= 1000000) return `${(heat / 1000000).toFixed(1)}M`
  if (heat >= 1000) return `${(heat / 1000).toFixed(0)}K`
  return String(heat)
}
```

示例：
- `1234567` → `"1.2M"`
- `45678` → `"46K"`
- `567` → `"567"`

### 异动类型映射
```javascript
function mapChangeType(type) {
  const typeMap = {
    1: 'limit_up',      // 涨停
    2: 'limit_down',    // 跌停
    3: 'rapid_rise',    // 急涨
    4: 'rapid_fall',    // 急跌
    5: 'high_volume',   // 放量
    6: 'breakthrough'   // 突破
  }
  return typeMap[type] || 'normal'
}
```

---

## ⚠️ 待完善功能

### 1. AI建议卡片
- **当前状态**: 使用静态模拟数据
- **可用API**: `GetAiRecommendStocksList(page)`
- **建议**: 后续对接AI推荐接口，根据用户自选股票生成个性化建议

### 2. 股票详情抽屉
- **当前状态**: 点击股票时仅打印日志
- **待实现**: 
  - 打开股票详情抽屉 (StockDetailSheet)
  - 显示完整行情、K线图、资金流向等

### 3. 新闻详情页
- **当前状态**: 点击新闻时仅打印日志
- **待实现**: 打开新闻详情页或外部链接

---

## 🐛 错误处理

所有API调用都包含 `try-catch` 错误处理：

```javascript
async function loadStockData() {
  try {
    const result = await GetFollowList(1)
    // ... 处理数据
  } catch (error) {
    console.error('加载自选股票失败:', error)
    // 失败时数据保持为空数组，不影响其他模块
  }
}
```

---

## 📱 下拉刷新

首页支持下拉刷新，会重新加载所有数据：

```javascript
async function handleRefresh() {
  refreshCount.value++
  await loadAllData()  // 并行加载所有6个模块的数据
}
```

---

## 🔧 后续优化建议

### 1. 缓存机制
- 实现本地缓存，减少重复请求
- 离线时显示缓存数据

### 2. 骨架屏
- 首次加载时显示骨架屏，提升用户体验
- 使用 `loading.value` 状态控制

### 3. 错误提示
- 网络失败时显示友好提示
- 提供重试按钮

### 4. 性能优化
- 使用虚拟滚动处理大量数据
- 图片懒加载

### 5. 数据推送
- 考虑使用 WebSocket 实现实时推送
- 减少轮询频率，降低服务器负载

---

## 📦 依赖的API文件

- `frontend/src/api/app.js` - 后端RPC接口定义
- `frontend/src/api/scheduler.js` - 统一轮询调度器
- `frontend/src/api/transport.js` - HTTP传输层

---

## 🧪 测试建议

### 功能测试
1. 检查首页各卡片是否正常显示数据
2. 验证下拉刷新功能
3. 测试网络断开时的错误处理
4. 验证轮询是否正常工作

### 性能测试
1. 监控API请求频率
2. 检查内存泄漏（轮询是否正确停止）
3. 测试弱网环境下的加载体验

### 兼容性测试
1. 不同屏幕尺寸
2. Android/iOS设备
3. 不同浏览器内核

---

## 📝 相关文件

- 主文件: `frontend/src/mobile/pages/HomePage.vue`
- 子组件: 
  - `frontend/src/mobile/components/cards/*.vue`
  - `frontend/src/mobile/components/widgets/*.vue`
- API定义: `frontend/src/api/app.js`

---

## 🎉 总结

本次对接完成了首页6个核心功能模块的数据接入：
- ✅ 自选股票概览
- ✅ 市场快讯
- ✅ 实时热点
- ✅ 异动监控
- ✅ 行业热度
- ✅ 市场统计

所有模块都使用真实后端API，配合合理的轮询策略，实现了数据的实时更新。
