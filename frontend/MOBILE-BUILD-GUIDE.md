# go-stock 移动端开发 - 编译运行指南

## 📦 分支信息

**分支名称：** `feat/mobile-architecture`

**包含内容：**
- PWA 支持（可安装、离线缓存）
- 移动端全新架构（阶段 0 完成）
- 完整的设计系统和基础布局

---

## 🚀 快速开始

### 1. 拉取代码

```bash
# 克隆仓库（如果还没有）
git clone https://github.com/smallestfisher/go-stock.git
cd go-stock

# 切换到移动端分支
git checkout feat/mobile-architecture

# 查看最新提交
git log -1 --oneline
```

### 2. 安装依赖

```bash
cd frontend

# 安装前端依赖（需要 Node.js 16+）
npm install
```

**依赖说明：**
- 新增：`vite-plugin-pwa` (PWA 支持)
- 已有：Vue 3、NaiveUI、Vite 等

### 3. 启动开发服务器

```bash
# 开发模式（热更新）
npm run dev
```

**启动成功后会显示：**
```
  VITE v7.3.3  ready in 1234 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: http://192.168.x.x:5173/
```

---

## 📱 测试移动端

### 方法 1：浏览器开发者工具（推荐）

1. 打开浏览器访问 `http://localhost:5173/`
2. 按 `F12` 打开开发者工具
3. 按 `Ctrl+Shift+M`（Mac: `Cmd+Shift+M`）切换移动设备模拟
4. 选择设备：
   - iPhone 12 Pro（推荐）
   - iPhone 14 Pro Max
   - 或自定义：宽度 375px，高度 812px
5. 刷新页面

**应该看到：**
- 移动端界面（不是桌面端）
- 顶部栏：☰ go-stock 🔍 🔔
- 首页显示 "Hello Mobile!" 🎉
- 点击左上角 `☰` 图标，侧边抽屉滑出

### 方法 2：真机测试

**前提：** 手机和电脑在同一局域网

1. 启动开发服务器后，查看 Network 地址：
   ```
   ➜  Network: http://192.168.1.100:5173/
   ```

2. 手机浏览器访问该地址

3. 应该看到移动端界面

**iOS 注意事项：**
- Safari 需要允许局域网访问
- 可以点击分享 → 添加到主屏幕（测试 PWA）

**Android 注意事项：**
- Chrome 会提示"安装应用"（PWA 功能）

### 方法 3：调整浏览器窗口

1. 将浏览器窗口宽度缩小到 < 768px
2. 刷新页面
3. 应该看到移动端界面

---

## 🖥️ 测试桌面端

### 验证桌面端不受影响

1. 浏览器窗口宽度 > 768px
2. 访问 `http://localhost:5173/`
3. 应该看到原有的桌面端界面（底部导航栏）
4. 所有功能正常

---

## 🏗️ 生产构建

### 环境要求

⚠️ **重要：** 前端构建需要 **8GB+ 内存**

根据 `memory/build-constraints.md`，本机内存不足，需要在其他设备构建。

### 构建命令

```bash
cd frontend

# 方式 1：指定 Node.js 内存（推荐）
NODE_OPTIONS="--max-old-space-size=8192" npm run build

# 方式 2：如果内存充足（≥16GB）
npm run build
```

### 构建产物

构建成功后会生成：
```
frontend/dist/
├── index.html
├── assets/
│   ├── index-*.js
│   ├── index-*.css
│   └── ...
├── manifest.webmanifest      # PWA manifest
├── sw.js                      # Service Worker
├── pwa-192x192.png           # PWA 图标
├── pwa-512x512.png
└── ...
```

### 预览构建结果

```bash
npm run preview
```

访问 `http://localhost:4173/` 预览生产版本。

---

## 🧪 功能测试清单

### 移动端基础功能

- [ ] 访问 `http://localhost:5173/`，宽度 < 768px，显示移动端界面
- [ ] 访问 `http://localhost:5173/`，宽度 > 768px，显示桌面端界面
- [ ] 移动端顶部栏显示正常
- [ ] 点击左上角菜单图标，侧边抽屉滑出
- [ ] 点击遮罩层或菜单项，抽屉关闭
- [ ] 首页显示 "Hello Mobile!" 欢迎信息
- [ ] 首页显示设备信息（宽度、高度）

### 侧边抽屉功能

- [ ] 用户信息区显示正常（头像、盈亏）
- [ ] 主菜单项显示：自选、市场、K线、研究、基金、AI智能体
- [ ] 设置菜单项显示：设置、关于
- [ ] 点击菜单项时抽屉关闭
- [ ] 点击菜单项跳转路由（当前会 404，正常现象）

### PWA 功能

- [ ] Chrome 浏览器提示"安装应用"（桌面端和移动端）
- [ ] 移动端 3 秒后显示安装横幅
- [ ] 点击"立即安装"能触发系统安装对话框
- [ ] 点击"关闭"后横幅消失
- [ ] 已安装后再次访问不显示横幅

### 样式测试

- [ ] 移动端背景色为浅灰色（`#f7f8fa`）
- [ ] 卡片背景为白色，有圆角和阴影
- [ ] 涨跌色正确（绿色涨、红色跌）
- [ ] 触控反馈正常（点击按钮有视觉反馈）
- [ ] 滚动流畅（无卡顿）

### 响应式测试

- [ ] iPhone SE (375x667) 显示正常
- [ ] iPhone 12 Pro (390x844) 显示正常
- [ ] iPhone 14 Pro Max (430x932) 显示正常，安全区适配正确
- [ ] iPad (768x1024) 显示桌面端
- [ ] Android 手机显示正常

### 兼容性测试

- [ ] Chrome/Edge 显示正常
- [ ] Safari 显示正常（iOS 和 macOS）
- [ ] Firefox 显示正常
- [ ] 横屏模式显示正常

---

## 🐛 常见问题

### Q1: 移动端不显示，还是桌面端界面

**原因：** 设备检测失败

**排查：**
```bash
# 1. 检查浏览器宽度是否 < 768px
# 2. 打开开发者工具 Console，输入：
window.innerWidth
# 应该返回 < 768

# 3. 检查 useDevice 是否生效，输入：
window.matchMedia('(max-width: 768px)').matches
# 应该返回 true
```

**解决：**
- 刷新页面
- 清除浏览器缓存
- 确保开发者工具的设备模拟已启用

### Q2: 侧边抽屉菜单点击后 404

**原因：** 对应页面还未创建（这是正常的）

**说明：** 阶段 0 只完成了基础架构，菜单指向的页面（如 `/mobile/stock`）还没开发

**解决：** 等待阶段 1-9 完成后，所有页面就能正常访问

### Q3: 图标不显示

**原因：** `@vicons/ionicons5` 未安装

**检查：**
```bash
cd frontend
npm list @vicons/ionicons5
```

**解决：**
```bash
npm install @vicons/ionicons5
```

### Q4: PWA 安装提示不显示

**原因：** PWA 需要 HTTPS 或 localhost

**说明：**
- `localhost` 环境会显示
- 局域网访问（`192.168.x.x`）可能不显示
- 生产环境必须 HTTPS

**解决：**
- 使用 `localhost:5173` 访问
- 或等待 3 秒（延迟显示）

### Q5: 编译内存溢出

**错误：**
```
FATAL ERROR: Reached heap limit Allocation failed - JavaScript heap out of memory
```

**原因：** 前端构建需要 8GB 堆内存

**解决：**
```bash
NODE_OPTIONS="--max-old-space-size=8192" npm run build
```

或在内存更大的设备上构建。

### Q6: Service Worker 注册失败

**错误：** Console 显示 `Service Worker registration failed`

**原因：** 开发模式下 SW 默认禁用

**说明：** 这是正常的，`vite.config.js` 中 `devOptions: { enabled: false }`

**解决：**
- 开发模式不需要 SW
- 生产构建（`npm run build`）后 SW 会自动启用

---

## 📊 性能检查

### Lighthouse 审计（生产构建后）

```bash
# 1. 构建
NODE_OPTIONS="--max-old-space-size=8192" npm run build

# 2. 预览
npm run preview

# 3. 打开 Chrome DevTools → Lighthouse
# 4. 选择 Categories: PWA
# 5. 点击 "Analyze page load"
```

**预期评分：**
- PWA: > 90 分
- Performance: > 80 分（视图表数据量而定）
- Accessibility: > 90 分
- Best Practices: > 90 分

---

## 📁 目录结构说明

```
frontend/
├── src/
│   ├── App.vue              # 桌面端（原有）
│   ├── AppRoot.vue          # 根组件（设备切换）
│   ├── main.js              # 入口（已修改）
│   │
│   └── mobile/              # 移动端（全新）
│       ├── MobileApp.vue           # 移动端根组件
│       ├── router.js               # 移动端路由
│       ├── layouts/                # 布局组件
│       │   ├── MobileLayout.vue
│       │   ├── MobileHeader.vue
│       │   └── MobileSideDrawer.vue
│       ├── pages/                  # 页面（当前只有 HomePage）
│       │   └── HomePage.vue
│       ├── components/             # 组件（待开发）
│       │   ├── base/
│       │   ├── cards/
│       │   ├── charts/
│       │   └── widgets/
│       ├── composables/            # Hooks（待开发）
│       └── styles/                 # 样式
│           ├── variables.css
│           ├── mobile.css
│           └── animations.css
│
├── public/                  # PWA 图标
│   ├── pwa-192x192.png
│   ├── pwa-512x512.png
│   ├── apple-touch-icon-180x180.png
│   └── favicon.ico
│
├── MOBILE-REFACTOR-PLAN.md         # 完整重构计划
├── MOBILE-STAGE0-COMPLETED.md      # 阶段0完成报告
└── PWA-IMPLEMENTATION.md           # PWA改造说明
```

---

## 🔄 继续开发

### 在其他设备上继续开发

```bash
# 1. 拉取最新代码
git fetch origin
git checkout feat/mobile-architecture
git pull

# 2. 安装依赖（如果是新设备）
cd frontend
npm install

# 3. 启动开发服务器
npm run dev

# 4. 开始阶段 1：基础组件库开发
# 参考 MOBILE-REFACTOR-PLAN.md 的阶段 1 任务清单
```

### 提交代码

```bash
# 1. 查看修改
git status

# 2. 添加文件
git add .

# 3. 提交
git commit -m "feat: 完成阶段 X - XXX功能"

# 4. 推送
git push origin feat/mobile-architecture
```

---

## 📞 技术支持

### 遇到问题？

1. **查看文档：**
   - `MOBILE-REFACTOR-PLAN.md` - 完整计划
   - `MOBILE-STAGE0-COMPLETED.md` - 阶段0说明
   - `PWA-IMPLEMENTATION.md` - PWA说明

2. **检查日志：**
   - 浏览器 Console（F12）
   - 终端输出
   - Network 请求

3. **常见命令：**
   ```bash
   npm run dev          # 启动开发服务器
   npm run build        # 生产构建
   npm run preview      # 预览构建结果
   npm install          # 安装依赖
   ```

---

## ✅ 验收标准

### 移动端基础架构验收

- [x] 移动端界面能正常显示
- [x] 侧边抽屉能打开/关闭
- [x] 首页显示欢迎信息
- [x] 桌面端功能不受影响
- [x] 设计系统变量生效
- [x] PWA 安装提示能显示
- [ ] 生产构建成功（需在大内存设备上测试）
- [ ] Lighthouse PWA 评分 > 90

### 准备进入阶段 1

阶段 0 验收通过后，可以开始**阶段 1：基础组件库开发**

---

## 🎉 完成状态

**当前分支：** `feat/mobile-architecture`

**当前阶段：** 阶段 0 ✅

**下一阶段：** 阶段 1 - 基础组件库

**预计总工作量：** 53 小时（6-7 个工作日）

**已完成：** 1 小时

**剩余：** 52 小时
