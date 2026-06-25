## 🎨 设计系统

### 颜色系统

```css
/* 主色 - 绿色（涨） */
--m-color-rise: #18a058;
--m-color-rise-light: rgba(24, 160, 88, 0.12);
--m-color-rise-hover: #16935e;

/* 红色（跌） */
--m-color-fall: #d03050;
--m-color-fall-light: rgba(208, 48, 80, 0.12);

/* 中性色 */
--m-color-gray: #666;
--m-color-flat: #999;

/* 背景色 */
--m-bg-primary: #f7f8fa;
--m-bg-card: #ffffff;
--m-bg-elevated: #ffffff;

/* 文字色 */
--m-text-primary: #1e2227;
--m-text-secondary: #6b7077;
--m-text-tertiary: #9ca3af;

/* 边框色 */
--m-border-color: #e7e3da;
--m-divider-color: #f0f0f0;
```

### 尺寸系统

```css
/* 触控区域（最小） */
--m-touch-min: 44px;

/* 间距 */
--m-space-xs: 4px;
--m-space-sm: 8px;
--m-space-md: 12px;
--m-space-lg: 16px;
--m-space-xl: 24px;

/* 圆角 */
--m-radius-sm: 6px;
--m-radius-md: 12px;
--m-radius-lg: 16px;

/* 卡片 */
--m-card-padding: 16px;
--m-card-gap: 12px;
```

### 字体系统

```css
/* 字号 */
--m-font-xs: 11px;
--m-font-sm: 12px;
--m-font-md: 14px;
--m-font-lg: 16px;
--m-font-xl: 18px;
--m-font-2xl: 24px;

/* 字重 */
--m-font-weight-normal: 400;
--m-font-weight-medium: 500;
--m-font-weight-bold: 600;
```

### 布局系统

```css
/* 顶部固定栏高度 */
--m-header-height: 56px;

/* 安全区 */
--m-safe-top: env(safe-area-inset-top);
--m-safe-bottom: env(safe-area-inset-bottom);

/* 内容区域高度 */
--m-content-height: calc(100dvh - var(--m-header-height) - var(--m-safe-top));
```

### 阴影系统

```css
--m-shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.06);
--m-shadow-md: 0 2px 8px rgba(0, 0, 0, 0.08);
--m-shadow-lg: 0 4px 16px rgba(0, 0, 0, 0.12);
```

### 动画系统

```css
--m-duration-fast: 150ms;
--m-duration-normal: 300ms;
--m-duration-slow: 500ms;
--m-ease-out: cubic-bezier(0.16, 1, 0.3, 1);
```

---

## 📊 页面功能清单

### 优先级分级

**P0 - 核心功能（首批开发）**
1. ✅ 基础架构（MobileApp、MobileLayout、Header、Drawer）
2. HomePage - 首页信息流
3. StockListPage - 自选股票列表
4. StockDetailSheet - 股票详情抽屉
5. NewsListPage - 市场快讯
6. KlineAnalysisPage - K线分析

**P1 - 高频功能（第二批）**
7. MarketPage - 市场行情主页（Tab切换框架）
8. GlobalIndexPage - 全球股指
9. IndustryRankPage - 行业排名
10. HotStockPage - 当前热门
11. DragonTigerPage - 龙虎榜
12. StockMoneyFlowPage - 个股资金流向

**P2 - 研究功能（第三批）**
13. ResearchPage - 研究中心主页（Tab切换框架）
14. AiReportPage - AI分析报告
15. StockRecommendPage - 股票推荐记录
16. AbnormalMonitorPage - 异动监控
17. LimitUpLadderPage - 涨停梯队
18. StockReportPage - 个股研报
19. CompanyNoticePage - 公司公告

**P3 - 辅助功能（第四批）**
20. SettingsPage - 设置
21. FundFollowPage - 基金自选
22. FundRankingPage - 基金排行
23. AgentChatPage - AI智能体
24. SearchPage - 全局搜索

**P4 - 低频功能（最后批次）**
25. PromptTemplatePage - 提示词模板
26. PromptPlazaPage - 提示词广场
27. QaPlazaPage - 问答广场
28. PatternSelectPage - 形态选股
29. IndicatorSelectPage - 指标选股
30. CronTaskPage - 定时任务
31. TradingLogPage - 交易日志
32. McpServerPage - MCP服务
33. SkillManagerPage - 技能管理
34. AboutPage - 关于

---

## 🔧 通用组件开发顺序

### 第一批：基础组件（必需）
1. MCard - 卡片容器
2. MButton - 按钮
3. MSheet - 底部抽屉
4. MTabs - 横向滚动标签页
5. MLoading - 加载状态
6. MEmpty - 空状态

### 第二批：列表组件
7. MList - 列表容器
8. MListItem - 列表项
9. MPullRefresh - 下拉刷新
10. VirtualList - 虚拟滚动列表

### 第三批：业务卡片
11. StockCard - 股票卡片
12. NewsCard - 新闻卡片
13. RankCard - 排行榜卡片
14. AlertCard - 预警卡片

### 第四批：图表组件
15. MiniKlineChart - 迷你K线
16. MiniSparkline - 迷你走势线
17. FullKlineChart - 全屏K线
18. FenshiChart - 分时图

### 第五批：功能组件
19. StockSearchBar - 搜索框
20. MarketStatusBar - 市场状态条
21. PriceTag - 价格标签
22. PercentTag - 百分比标签

---

## 🚀 实施步骤（分阶段）

### 阶段 0：准备工作（当前）✅
- [x] 创建目录结构
- [x] 编写重构计划文档
- [ ] 创建设计系统样式文件
- [ ] 修改 App.vue 支持设备切换

### 阶段 1：基础架构（2-3小时）
**目标：搭建移动端框架，能跑起来**

1. 创建 `mobile/styles/variables.css`（设计变量）
2. 创建 `mobile/styles/mobile.css`（全局样式）
3. 创建 `mobile/MobileApp.vue`（移动端根组件）
4. 创建 `mobile/router.js`（移动端路由）
5. 创建 `mobile/layouts/MobileLayout.vue`（主布局）
6. 创建 `mobile/layouts/MobileHeader.vue`（顶部栏）
7. 创建 `mobile/layouts/MobileSideDrawer.vue`（侧边抽屉）
8. 修改 `src/App.vue`（根据 isMobile 切换组件）
9. 创建 `mobile/pages/HomePage.vue`（临时首页，显示"Hello Mobile"）

**验证标准：**
- 移动端设备访问，显示移动端界面
- 侧边抽屉能打开/关闭
- 路由能正常跳转

---

### 阶段 2：基础组件库（3-4小时）
**目标：建立可复用的移动端组件**

1. `mobile/components/base/MCard.vue`
2. `mobile/components/base/MButton.vue`
3. `mobile/components/base/MSheet.vue`
4. `mobile/components/base/MTabs.vue`
5. `mobile/components/base/MLoading.vue`
6. `mobile/components/base/MEmpty.vue`
7. `mobile/components/base/MPullRefresh.vue`
8. `mobile/composables/usePullRefresh.js`

**验证标准：**
- 每个组件独立可用
- 触控交互流畅（无延迟、无误触）
- 适配安全区（刘海屏、横屏）

---

### 阶段 3：首页信息流（4-5小时）
**目标：完成首页，展示核心数据**

1. 重构 `mobile/pages/HomePage.vue`（信息流布局）
2. 创建 `mobile/components/cards/StockSummaryCard.vue`（自选概览）
3. 创建 `mobile/components/cards/NewsCard.vue`（快讯卡片）
4. 创建 `mobile/components/cards/HotTopicCard.vue`（热点卡片）
5. 创建 `mobile/components/cards/AlertCard.vue`（异动预警）
6. 创建 `mobile/components/widgets/MarketStatusBar.vue`（市场状态）
7. 集成下拉刷新
8. 集成数据加载（调用现有 API）

**验证标准：**
- 首页显示真实数据（自选、快讯、热点）
- 下拉刷新能更新数据
- 点击卡片能跳转到对应页面（先跳空页面）

---

### 阶段 4：自选股票模块（5-6小时）
**目标：完成核心交易功能**

1. `mobile/pages/stock/StockListPage.vue`（自选列表）
2. `mobile/components/cards/StockCard.vue`（股票卡片，带迷你K线）
3. `mobile/components/charts/MiniSparkline.vue`（迷你走势线）
4. `mobile/components/sheets/StockDetailSheet.vue`（股票详情抽屉）
5. `mobile/components/charts/FenshiChart.vue`（分时图）
6. `mobile/components/charts/FullKlineChart.vue`（全屏K线）
7. `mobile/components/widgets/PriceTag.vue`（价格标签）
8. `mobile/components/widgets/PercentTag.vue`（百分比标签）
9. 集成虚拟滚动（处理大量股票）
10. 集成搜索、分组功能

**验证标准：**
- 自选列表显示所有股票
- 点击股票打开详情抽屉
- 详情抽屉显示完整信息（价格、K线、盘口等）
- 滚动流畅（300+股票不卡顿）

---

### 阶段 5：市场行情模块（6-8小时）
**目标：完成市场信息查看功能**

1. `mobile/pages/market/MarketPage.vue`（市场主页，横向Tab切换）
2. `mobile/pages/market/NewsListPage.vue`（市场快讯列表）
3. `mobile/pages/market/GlobalIndexPage.vue`（全球股指）
4. `mobile/pages/market/IndustryRankPage.vue`（行业排名）
5. `mobile/pages/market/HotStockPage.vue`（当前热门）
6. `mobile/pages/market/DragonTigerPage.vue`（龙虎榜）
7. `mobile/components/cards/NewsTimelineCard.vue`（新闻时间轴卡片）
8. `mobile/components/cards/IndustryCard.vue`（行业卡片）
9. `mobile/components/cards/RankCard.vue`（排行榜卡片）
10. `mobile/components/sheets/NewsDetailSheet.vue`（新闻详情）

**验证标准：**
- 市场主页能左右滑动切换Tab
- 每个Tab显示对应数据
- 新闻列表按时间轴展示
- 行业排名、龙虎榜等用卡片展示

---

### 阶段 6：K线分析模块（4-5小时）
**目标：完成专业图表分析功能**

1. `mobile/pages/kline/KlineAnalysisPage.vue`（K线分析主页）
2. 优化 `mobile/components/charts/FullKlineChart.vue`（全屏K线）
3. `mobile/components/sheets/ChartControlSheet.vue`（图表参数控制）
4. `mobile/composables/useSwipe.js`（左右滑动切换周期）
5. 集成技术指标（MA、MACD、KDJ等）
6. 集成多周期切换（日K、周K、月K）

**验证标准：**
- K线图全屏显示，无多余UI
- 左右滑动切换周期
- 底部抽屉控制指标参数
- 支持缩放、拖拽

---

### 阶段 7：研究中心模块（6-8小时）
**目标：完成数据分析功能**

1. `mobile/pages/research/ResearchPage.vue`（研究主页，Tab切换）
2. `mobile/pages/research/AiReportPage.vue`（AI分析报告）
3. `mobile/pages/research/StockRecommendPage.vue`（推荐记录）
4. `mobile/pages/research/AbnormalMonitorPage.vue`（异动监控）
5. `mobile/pages/research/LimitUpLadderPage.vue`（涨停梯队）
6. `mobile/pages/market/StockReportPage.vue`（个股研报）
7. `mobile/pages/market/CompanyNoticePage.vue`（公司公告）
8. `mobile/components/cards/ReportCard.vue`（研报卡片）
9. `mobile/components/sheets/ReportDetailSheet.vue`（研报详情）

**验证标准：**
- 研究主页Tab切换流畅
- AI报告、推荐记录正常显示
- 异动监控、涨停梯队实时更新
- 研报、公告卡片展示完整

---

### 阶段 8：辅助功能（3-4小时）
**目标：完成设置、基金、AI等辅助功能**

1. `mobile/pages/settings/SettingsPage.vue`（设置）
2. `mobile/pages/settings/AboutPage.vue`（关于）
3. `mobile/pages/fund/FundFollowPage.vue`（基金自选）
4. `mobile/pages/fund/FundRankingPage.vue`（基金排行）
5. `mobile/pages/agent/AgentChatPage.vue`（AI智能体）
6. `mobile/pages/search/SearchPage.vue`（全局搜索）

**验证标准：**
- 设置页单列表单，分组清晰
- 基金模块功能完整
- AI聊天全屏显示，输入框固定底部
- 全局搜索快速响应

---

### 阶段 9：低频功能（4-5小时）
**目标：补全所有剩余功能**

1. `mobile/pages/research/PromptTemplatePage.vue`
2. `mobile/pages/research/PromptPlazaPage.vue`
3. `mobile/pages/research/QaPlazaPage.vue`
4. `mobile/pages/research/PatternSelectPage.vue`
5. `mobile/pages/research/IndicatorSelectPage.vue`
6. `mobile/pages/research/CronTaskPage.vue`
7. `mobile/pages/research/TradingLogPage.vue`
8. `mobile/pages/research/McpServerPage.vue`
9. `mobile/pages/research/SkillManagerPage.vue`
10. `mobile/pages/market/SectorMoneyFlowPage.vue`
11. `mobile/pages/market/ConceptMoneyFlowPage.vue`
12. `mobile/pages/market/IndustryResearchPage.vue`
13. `mobile/pages/market/FeaturedSitesPage.vue`

**验证标准：**
- 所有页面可访问
- 功能逻辑完整
- 无明显Bug

---

### 阶段 10：优化与测试（3-4小时）
**目标：性能优化和兼容性测试**

1. 性能优化
   - 虚拟滚动优化（大列表）
   - 图片懒加载
   - 路由懒加载
   - 按需引入组件

2. 兼容性测试
   - iOS Safari（刘海屏、横屏）
   - Android Chrome
   - 各种屏幕尺寸

3. 用户体验优化
   - 加载状态优化
   - 错误提示优化
   - 动画流畅度优化
   - 触控反馈优化

4. PWA集成
   - 确保PWA功能在新架构下正常
   - 测试离线访问
   - 测试安装流程

**验证标准：**
- Lighthouse PWA 评分 > 90
- 首屏加载 < 2s
- 列表滚动 60fps
- 无明显兼容性问题

---

## 📝 开发规范

### 命名规范

**文件命名：**
- 页面组件：`XxxPage.vue`（PascalCase + Page后缀）
- 布局组件：`MobileXxx.vue`（Mobile前缀）
- 基础组件：`MXxx.vue`（M前缀，代表Mobile）
- 业务组件：`XxxCard.vue`、`XxxSheet.vue`
- Composables：`useXxx.js`（use前缀）

**变量命名：**
- CSS变量：`--m-xxx-yyy`（m-前缀，代表mobile）
- JS变量：小驼峰 `camelCase`
- 组件Props：小驼峰 `camelCase`

### 代码规范

**组件结构：**
```vue
<script setup>
// 1. 导入
import { ref, computed } from 'vue'

// 2. Props/Emits
const props = defineProps({...})
const emit = defineEmits([...])

// 3. 响应式数据
const data = ref(null)

// 4. 计算属性
const computed = computed(() => ...)

// 5. 方法
function method() {...}

// 6. 生命周期
onMounted(() => {...})
</script>

<template>
  <!-- 模板 -->
</template>

<style scoped>
/* 样式 */
</style>
```

**样式规范：**
- 优先使用设计系统变量（`var(--m-xxx)`）
- 使用 `scoped` 避免样式污染
- 移动端专用样式写在组件内，不要写在全局

**性能规范：**
- 长列表必须用虚拟滚动
- 图片必须懒加载
- 路由必须懒加载
- 避免不必要的响应式数据

---

## 🔗 与桌面端的关系

### 代码隔离

**完全独立：**
- 移动端代码在 `src/mobile/` 目录
- 桌面端代码在 `src/views/`、`src/components/` 目录
- 两者不共享组件

**可共享：**
- API调用（`src/api/`）
- 工具函数（`src/utils/`）
- Composables（`src/composables/useDevice.js` 等）
- 样式变量（颜色、字体等常量）

### 路由切换

在 `src/App.vue` 根据设备类型切换：

```vue
<script setup>
import { computed } from 'vue'
import { useDevice } from './composables/useDevice'
import DesktopApp from './views/DesktopApp.vue' // 桌面端（现有的）
import MobileApp from './mobile/MobileApp.vue'   // 移动端（新建的）

const { isMobile } = useDevice()
</script>

<template>
  <MobileApp v-if="isMobile" />
  <DesktopApp v-else />
</template>
```

---

## ⏱️ 预估工作量

| 阶段 | 内容 | 预估时间 | 累计时间 |
|------|------|---------|---------|
| 0 | 准备工作 | 1h | 1h |
| 1 | 基础架构 | 3h | 4h |
| 2 | 基础组件库 | 4h | 8h |
| 3 | 首页信息流 | 5h | 13h |
| 4 | 自选股票模块 | 6h | 19h |
| 5 | 市场行情模块 | 8h | 27h |
| 6 | K线分析模块 | 5h | 32h |
| 7 | 研究中心模块 | 8h | 40h |
| 8 | 辅助功能 | 4h | 44h |
| 9 | 低频功能 | 5h | 49h |
| 10 | 优化与测试 | 4h | 53h |

**总计：约 53 小时（6-7 个工作日）**

---

## ✅ 完成标准

### 功能完整性
- [ ] 所有 34 个页面均可访问
- [ ] 所有核心功能（自选、行情、K线、研究）正常运行
- [ ] 所有API调用正常
- [ ] 无明显功能缺失

### 用户体验
- [ ] 触控交互流畅（无延迟、无误触）
- [ ] 滚动流畅（60fps）
- [ ] 动画自然（符合移动端习惯）
- [ ] 加载状态清晰
- [ ] 错误提示友好

### 兼容性
- [ ] iOS Safari 完美支持
- [ ] Android Chrome 完美支持
- [ ] 刘海屏适配正常
- [ ] 横屏适配正常
- [ ] PWA 功能正常

### 性能
- [ ] Lighthouse PWA 评分 > 90
- [ ] 首屏加载 < 2s
- [ ] 列表滚动 60fps
- [ ] 内存占用合理

---

## 📌 注意事项

1. **不要复用桌面端组件**：桌面端组件是为大屏设计的，强行复用会有很多显示问题
2. **触控优先设计**：所有按钮、链接最小 44px，避免误触
3. **安全区适配**：所有固定定位元素都要考虑刘海屏和底部手势条
4. **性能优先**：长列表必须虚拟滚动，图片必须懒加载
5. **渐进增强**：先保证基础功能可用，再添加动画、手势等增强体验
6. **设备测试**：开发过程中在真机上测试，不要只在浏览器开发者工具测试

---

## 🎯 下一步

当前已完成：
- [x] 创建目录结构
- [x] 编写重构计划文档

接下来：
1. 创建设计系统文件（`mobile/styles/variables.css` 等）
2. 搭建基础架构（`MobileApp.vue`、`MobileLayout.vue` 等）
3. 开发基础组件库
4. 按优先级实施各功能模块

**建议从阶段 1 开始，先把框架搭起来，能在移动端看到"Hello Mobile"再继续。**
