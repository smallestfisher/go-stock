# go-stock 服务端部署说明

本文件说明如何构建和部署 go-stock **纯服务端 + 远程 Web 访问**模式。

## 架构速览

- `server/cmd/server`：HTTP 服务入口（`net/http` + `ServeMux`），托管前端 + `/api`。
- `/api/rpc/{Method}`：统一 RPC 桥，前端 API 方法 POST 到这里。
- `/api/events`：SSE 事件总线，承载行情/资讯/AI 流式等后端推送。浏览器每个标签页会带 `clientId`，AI 流式结果只投递给发起请求的标签页。
- `/api/health`：健康检查（兼作登录令牌校验）。
- 其余 `/`：`frontend/dist` 静态资源（Vue SPA，history 回退）。
- 业务逻辑在 `backend/`，`server/modules/*` 按模块注册 RPC handler。

## 环境要求

- Go 1.26+
- Node.js（建议 18+，构建前端）
- **服务端无头 Chromium**（chromedp 依赖）：抓东方财富 cookie、雪球当前热门等数据源需要。
  - Linux：`apt install chromium-browser` 或 `google-chrome`（启动日志会打印检测到的路径）
  - 设置页里的“服务端 Chromium 路径”指部署 go-stock 的服务器/容器内浏览器可执行文件路径，不是访问 Web 页面那台电脑的浏览器。
  - 该路径留空时服务端会自动检测 Edge/Chrome/Chromium；容器部署建议把 Chromium 安装进镜像。
  - 无 Chromium 时，依赖 cookie 的“当前热门”等数据源会降级，其余功能不受影响。

## 构建

> 前端依赖较多，构建吃内存。在 ≥8GB 内存机器上构建；内存紧张时加大 Node 堆。

```bash
# 1. 前端
cd frontend
npm install --registry=https://registry.npmmirror.com   # 国内镜像，按需
NODE_OPTIONS="--max-old-space-size=8192" npm run build   # 产出 frontend/dist
cd ..

# 2. 服务端二进制
GOPROXY=https://goproxy.cn,direct go build -o bin/go-stock-server ./server/cmd/server
```

## 运行

```bash
# 必备数据文件（A股/港股/美股基础数据种子 + 使用手册），随仓库携带
#   build/stock_basic.json  build/stock_base_info_hk.json  build/stock_base_info_us.json
#   docs/go-stock使用手册.md

GO_STOCK_ADDR=:18888 \
GO_STOCK_TOKEN=请改成你自己的强随机令牌 \
  ./bin/go-stock-server
```

环境变量：

| 变量 | 默认 | 说明 |
|---|---|---|
| `GO_STOCK_ADDR` | `:18888` | 监听地址 |
| `GO_STOCK_TOKEN` | 空 | 访问令牌；**为空=不鉴权(仅本地调试)**。远程部署务必设置强令牌 |

启动后：浏览器打开 `http://<服务器IP>:18888/` → 输入令牌登录 → 进入应用。
首次启动会自动建表并把 A股/港股/美股基础数据导入 `data/stock.db`（幂等，仅插缺失）。
服务端会按顺序执行数据库初始化、迁移/种子数据导入、定时任务恢复，然后才开始监听 HTTP，避免前端请求打到未就绪的数据库或 cron 状态。

## 远程访问建议

- **HTTPS**：在前面挂 Nginx/Caddy 反向代理终结 TLS，再转发到 `:18888`。
- SSE 经反代时需关闭缓冲：Nginx 加 `proxy_buffering off;`（代码已设 `X-Accel-Buffering: no`）。
- 数据目录 `data/`（含 `stock.db`）与 `logs/` 需持久化（挂载卷）。

## 生产部署推荐流程

生产环境建议把 go-stock 放在固定目录，用 systemd 守护进程，并通过 Nginx/Caddy 提供 HTTPS。

### 1. 准备目录和依赖

示例目录使用 `/opt/go-stock`，也可以按实际运维规范调整：

```bash
sudo mkdir -p /opt/go-stock
sudo chown -R $USER:$USER /opt/go-stock

sudo apt update
sudo apt install -y git build-essential chromium-browser
```

安装 Go 1.26+ 和 Node.js 18+。如果服务器内存不足，前端可以在本地或 CI 构建好 `frontend/dist` 后上传到服务器。

### 2. 拉取代码并构建

```bash
cd /opt/go-stock
git clone <你的仓库地址> .

# 前端构建；内存紧张时建议放到 CI 或其他机器构建
cd frontend
npm install --registry=https://registry.npmmirror.com
NODE_OPTIONS="--max-old-space-size=8192" npm run build
cd ..

# 服务端构建
GOPROXY=https://goproxy.cn,direct go build -o bin/go-stock-server ./server/cmd/server
```

运行目录内需要保留仓库自带的 `build/` 和 `docs/`，首次启动会用它们导入基础数据和读取使用手册。

### 3. 配置环境文件

创建 `/etc/go-stock.env`：

```bash
GO_STOCK_ADDR=127.0.0.1:18888
GO_STOCK_TOKEN=请改成强随机访问令牌
```

生产环境建议只监听 `127.0.0.1:18888`，由反向代理对外暴露 HTTPS。不要把 `GO_STOCK_TOKEN` 留空。

### 4. 配置 systemd

创建 `/etc/systemd/system/go-stock.service`：

```ini
[Unit]
Description=go-stock web server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=/opt/go-stock
EnvironmentFile=/etc/go-stock.env
ExecStart=/opt/go-stock/bin/go-stock-server
Restart=always
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

启动并设置开机自启：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now go-stock
sudo systemctl status go-stock
```

查看日志：

```bash
journalctl -u go-stock -f
```

### 5. 配置 Nginx HTTPS 反代

Nginx 示例：

```nginx
server {
    listen 80;
    server_name your.domain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your.domain.com;

    ssl_certificate /etc/letsencrypt/live/your.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your.domain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:18888;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # SSE 流式输出需要关闭缓冲
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }
}
```

证书可以用 Certbot 或 Caddy 自动签发。使用 Caddy 时直接反代到 `127.0.0.1:18888` 即可。

### 6. 数据持久化和备份

必须持久化并定期备份：

```text
data/
logs/
frontend/dist/
```

其中 `data/stock.db` 保存自选股、分组、AI 配置、提示词、定时任务等数据。升级和重新部署时不要删除 `data/`。

备份示例：

```bash
mkdir -p /opt/backups/go-stock
tar czf /opt/backups/go-stock/go-stock-data-$(date +%F-%H%M%S).tar.gz data logs
```

### 7. 升级发布

```bash
cd /opt/go-stock
git pull

# 如前端有变化，重新构建 frontend/dist
cd frontend
npm install --registry=https://registry.npmmirror.com
NODE_OPTIONS="--max-old-space-size=8192" npm run build
cd ..

GOPROXY=https://goproxy.cn,direct go build -o bin/go-stock-server ./server/cmd/server
sudo systemctl restart go-stock
journalctl -u go-stock -n 100 --no-pager
```

升级前建议先备份 `data/`。如果使用 CI 构建前端，把新的 `frontend/dist` 和服务端二进制同步到生产目录后重启服务。

### 8. Chromium 配置

“服务端 Chromium 路径”是部署 go-stock 的服务器/容器内浏览器可执行文件路径，不是访问 Web 页面那台电脑的浏览器。留空时服务端会自动检测。

常见路径：

```text
/usr/bin/chromium-browser
/usr/bin/chromium
/usr/bin/google-chrome
```

如果“当前热门”、雪球或依赖 cookie 的数据源异常，先检查服务器是否安装 Chromium，并在设置页填入实际路径。

## 复用已有数据

把旧部署的 `data/stock.db`（及 `data/` 下其余文件）直接拷到服务端工作目录的 `data/`，
自选股、分组、AI 配置、提示词等全部沿用——无需重新配置。

## 追加 RPC 方法

新增接口时按 3 步处理：

1. 在 `server/modules/<模块>.go` 的 `init()` 里 `server.Register("方法名", func(ctx, core, args) (any, error) { ... })`。
2. handler 内调用 `backend/data` 或 `backend/agent` 中的业务函数，参数用 `server.ArgString/ArgInt/ArgFloat64/ArgBool/ArgJSON/ArgIntPtr` 按前端调用顺序取。
3. 流式方法（边收 channel 边推送）参照 `server/modules/ai.go` 的 `NewChatStream`：从 `server.ClientIDFromContext(ctx)` 取浏览器标签页 scope，goroutine 里用 `core.Events.EmitTo(clientID, 事件名, msg)` 推送；handler 立即返回，前端经 `/api/events` 收事件。全局告警/行情类事件仍使用 `core.Events.Emit` 广播。

前端 API 入口在 `frontend/src/api/app.js`，新增 RPC 后同步添加对应方法。

### 迁移覆盖情况

当前 RPC 覆盖主要用户态方法：
行情/K线、自选与分组(读写)、基金、通达信 F10、资讯/研报/龙虎榜/热门、异动与统计、板块/概念资金流、
市场统计(含采集)、策略与选股、情感分析、AI 流式分析(NewChatStream/SummaryStockNews/ChatWithAgent)、
AI 结果与推荐历史、提示词模板、MCP 服务器与技能、交易记录、定时任务(增删改查 + 启动恢复 + 单股AI分析 + 下次运行时间计算)、
钉钉告警、分享、AI 助手会话/配置、模型信息探测(FetchAiModels/FetchAiModelInfo)、免费授权兼容接口、交易时间、NewsPush、CheckStockBaseInfo、Greet 等。

### Web 运行时说明

- `frontend/src/api/runtime.js` 提供浏览器本地事件和打开外链能力。
- `EventsEmit` 是浏览器标签页内的本地事件，用于切换 tab、刷新局部列表等 UI 交互，不会发到服务端。
- `EventsOn` 同时接收本地事件和服务端 SSE 事件；后端 AI 流式事件按 `clientId` 隔离，避免多用户/多标签页串流。
- `OpenURL` 直接调用 `window.open`；`SaveImage` / `SaveWordFile` 直接触发浏览器下载。

仅余以下**内部 helper** 未注册为 RPC（它们不是前端可调用的方法形态，无需迁移）：
`AddCronTask`（返回 `func()`，被 SetStockAICron 内部使用）、`InitCronTasks`（启动时由 cmd/server 调用）、
`MonitorStockPrices`（包级函数，由 cron 调度执行）。

> 前端只应调用 `frontend/src/api/app.js` 中存在的方法；不存在的服务端能力应先注册 RPC 再添加前端 API。
