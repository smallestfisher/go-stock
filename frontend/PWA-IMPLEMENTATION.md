# PWA 改造完成说明

## ✅ 已完成的工作

### 1. PWA 核心配置
- ✅ 安装 `vite-plugin-pwa` 插件
- ✅ 配置 `vite.config.js` 的 PWA 选项
- ✅ 生成 Web App Manifest（应用名称、图标、主题色等）
- ✅ 配置 Workbox Service Worker（静态资源缓存策略）

### 2. PWA 图标
- ✅ 生成多尺寸 PWA 图标（192x192、512x512）
- ✅ 生成 iOS 专用图标（180x180）
- ✅ 生成 Favicon
- ⚠️ **注意**：当前图标是临时生成的（直接复制源图标），浏览器会自动缩放
- 📌 **生产环境建议**：使用在线工具 https://realfavicongenerator.net/ 生成正确尺寸的图标

### 3. HTML Meta 标签
- ✅ 添加 PWA 主题色（`theme-color`）
- ✅ 添加 iOS Safari 专用配置（全屏模式、状态栏样式）
- ✅ 添加应用描述和 Favicon 链接
- ✅ 修正 `lang` 属性为 `zh-CN`

### 4. 安装提示组件
- ✅ 创建 `PwaInstallPrompt.vue` 组件
- ✅ 监听 `beforeinstallprompt` 事件
- ✅ 3 秒延迟显示（避免首次访问立即打扰）
- ✅ 记录用户选择（防止重复骚扰）
- ✅ 已集成到 `App.vue`（仅移动端显示）

### 5. 移动端体验优化
- ✅ PWA 独立窗口模式的安全区适配
- ✅ iOS 刘海屏顶部安全区支持
- ✅ 横屏模式布局优化（底部导航栏自动缩小）
- ✅ 系统手势遮挡防护（底部留出额外间距）

## 📦 修改的文件清单

**新增文件：**
- `frontend/public/pwa-192x192.png`
- `frontend/public/pwa-512x512.png`
- `frontend/public/apple-touch-icon-180x180.png`
- `frontend/public/favicon.ico`
- `frontend/src/components/PwaInstallPrompt.vue`
- `frontend/scripts/generate-pwa-icons.js`

**修改文件：**
- `frontend/package.json` — 添加 vite-plugin-pwa 依赖
- `frontend/vite.config.js` — 配置 VitePWA 插件
- `frontend/index.html` — 添加 PWA meta 标签
- `frontend/src/App.vue` — 引入安装提示组件
- `frontend/src/style.css` — PWA 独立模式安全区适配

## 🚀 下一步操作（需在有足够内存的设备上执行）

### 1. 安装依赖
```bash
cd frontend
npm install
```

### 2. 开发环境测试（可选）
```bash
npm run dev
```
打开浏览器开发者工具 → Application：
- Manifest 标签页 — 检查 manifest 配置
- Service Workers 标签页 — 检查 SW 注册状态

**注意**：开发模式下 Service Worker 默认禁用（`devOptions: { enabled: false }`），如需测试可临时改为 `true`

### 3. 生产构建
```bash
# 需要 8GB+ 内存
NODE_OPTIONS="--max-old-space-size=8192" npm run build
```

### 4. 预览生产版本
```bash
npm run preview
```

在浏览器访问 `http://localhost:4173`，测试 PWA 功能。

## ✅ PWA 功能测试清单

### 基础功能
- [ ] Chrome DevTools → Lighthouse → PWA 审计（评分 > 90）
- [ ] 访问应用时浏览器显示"安装应用"提示
- [ ] 点击"立即安装"按钮能触发系统安装对话框
- [ ] 安装后应用图标显示在主屏幕
- [ ] 从主屏幕启动后显示独立窗口（无浏览器地址栏）

### 移动端测试
- [ ] 手机浏览器访问，3秒后出现安装横幅
- [ ] 安装后图标正常显示
- [ ] 刘海屏设备顶部安全区正常（内容不被遮挡）
- [ ] 底部导航不被系统手势遮挡
- [ ] 横屏模式下底部导航自动缩小
- [ ] 启动画面（Splash Screen）显示应用图标和名称

### iOS Safari 专项测试
- [ ] 分享菜单 → "添加到主屏幕"
- [ ] 主屏幕图标清晰（180x180 尺寸）
- [ ] 启动后全屏显示（无 Safari 工具栏）
- [ ] 状态栏样式正确（`black-translucent`）

### 离线功能测试
- [ ] 断开网络后刷新页面，静态资源（HTML/CSS/JS）正常加载
- [ ] 断开网络后 API 请求失败（符合预期：股票数据不缓存）
- [ ] 恢复网络后自动恢复正常

## 🔧 可选的后续优化

### 1. 优化图标质量（推荐）
当前图标是临时生成的（200x200 源图标直接复制为多个尺寸），建议：

**方式 A - 在线工具（最简单）：**
1. 访问 https://realfavicongenerator.net/
2. 上传 `build/appicon.png`
3. 下载生成的图标包
4. 解压并替换 `frontend/public/` 中的图标

**方式 B - ImageMagick（如有可用）：**
```bash
cd frontend/public
convert ../../build/appicon.png -resize 192x192 pwa-192x192.png
convert ../../build/appicon.png -resize 512x512 pwa-512x512.png
convert ../../build/appicon.png -resize 180x180 apple-touch-icon-180x180.png
convert ../../build/appicon.png -resize 32x32 favicon.ico
```

### 2. 安装提示时机调整
如需修改安装提示的时机（当前为首次访问 3 秒后），编辑 `PwaInstallPrompt.vue`：

```javascript
// 方案 1：多次访问后提示
const visitCount = parseInt(localStorage.getItem('visit-count') || '0') + 1
localStorage.setItem('visit-count', visitCount)
if (visitCount >= 3) {
  showPrompt.value = true
}

// 方案 2：仅手动触发（删除 setTimeout）
// 在设置页提供"安装应用"按钮
```

### 3. API 缓存策略（可选）
如需缓存 API 数据（离线显示历史数据），在 `vite.config.js` 的 `workbox.runtimeCaching` 添加：

```javascript
{
  urlPattern: /^\/api\//,
  handler: 'NetworkFirst', // 优先网络，失败时用缓存
  options: {
    cacheName: 'api-cache',
    expiration: { maxAgeSeconds: 5 * 60 }, // 5 分钟过期
    cacheableResponse: { statuses: [0, 200] }
  }
}
```

### 4. 推送通知（需要后端支持）
如需实现股价预警推送，需要：
1. 后端生成 VAPID 密钥
2. Service Worker 监听 `push` 事件
3. 用户订阅推送（`navigator.serviceWorker.register().then(reg => reg.pushManager.subscribe())`）
4. 后端调用 Web Push API 发送通知

**建议后续单独规划**，当前先完成基础 PWA 功能。

## 📝 注意事项

1. **内存限制**：本机无法执行 `npm run build`（需要 8GB 堆），代码改完即可，实际构建在其他设备进行
2. **图标优化**：生产环境建议用正确尺寸的图标（当前临时图标可用于开发测试）
3. **HTTPS 要求**：PWA 需要 HTTPS（或 localhost），部署时确保使用 HTTPS
4. **浏览器支持**：Chrome/Edge/Safari 完整支持，Firefox 部分支持（无安装提示）
5. **Service Worker 更新**：修改代码后 SW 会自动更新（`registerType: 'autoUpdate'`）

## 🎉 完成状态

✅ **PWA 核心功能已全部实现**

可在有内存的设备上执行 `npm install` 和 `npm run build` 进行测试。测试通过后即可部署到生产环境，用户即可在移动端将 go-stock 安装为类原生 App！
