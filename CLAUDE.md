# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

go-stock is an AI-powered stock analysis tool. This branch/deployment runs **server-only + browser access**: a Go HTTP server hosts the API and serves a Vue SPA. The frontend is mobile-first (PWA). Supports A-share, HK, and US markets, integrating many LLM providers for sentiment/technical analysis.

Note: the codebase originated as a Wails desktop app. It has been ported to the HTTP server model, but backend business logic in `backend/` still reflects the original App structure. The `server/` layer wraps it.

## Commands

All Go commands need `GOPROXY=https://goproxy.cn,direct` on this machine (see `~/.claude` memory). The domestic market-data DNS is IPv6-only and local IPv6 is broken, so `backend/data/httpclient.go` forces `tcp4` — don't remove that.

```bash
# Build the server binary
GOPROXY=https://goproxy.cn,direct go build -o bin/go-stock-server ./server/cmd/server

# Build the frontend (needs 8GB Node heap; this machine may lack memory — build in CI or another host)
cd frontend
npm install --registry=https://registry.npmmirror.com
NODE_OPTIONS="--max-old-space-size=8192" npm run build   # → frontend/dist
npm run dev   # vite dev server on :5173, proxies /api to :18888

# Run the server
GO_STOCK_ADDR=:18888 GO_STOCK_TOKEN=<strong-token> ./bin/go-stock-server
#   GO_STOCK_TOKEN empty = NO AUTH (local dev only). Always set it for remote.

# Tests (per-package; there is no repo-wide "run all" convention)
GOPROXY=https://goproxy.cn,direct go test ./server/events/
GOPROXY=https://goproxy.cn,direct go test ./backend/data/ -run TestName
GOPROXY=https://goproxy.cn,direct go vet ./server/...
```

Full deploy guide (systemd, Nginx TLS, Chromium dependency): `server/DEPLOY.md`.

## Architecture

### The RPC bridge — the central pattern

The frontend never calls arbitrary endpoints. Everything goes through one bridge:

- **Frontend**: `frontend/src/api/app.js` is a flat file of stub functions, each just `rpc("MethodName", [args])`. `transport.js` POSTs to `/api/rpc/{Method}` with body `{"args":[...]}`. When testing RPCs with curl, the body **must** be `{"args":[...]}` — a bare array returns empty results.
- **Backend**: `server/rpc.go` dispatches `/api/rpc/{Method}` to a handler in a global registry. Handlers are registered by `server.Register("Method", fn)` calls inside `init()` functions across `server/modules/*.go`. `cmd/server/main.go` blank-imports `server/modules` to trigger all registrations.
- Handlers receive positional args as `server.Args` (`[]json.RawMessage`). Use the `server.ArgString/ArgInt/ArgJSON[T]/ArgIntSlice/...` helpers in `rpc.go` to decode by index — they return zero values on out-of-range/type-mismatch, so handlers stay robust.
- Handlers are thin: they decode args and call existing functions in `backend/data`, `backend/agent`, or `backend/db`. Business logic lives in `backend/`, not in modules.

**To add an RPC method**: add the `server.Register` call in the appropriate `server/modules/*.go` file, and add the matching stub to `frontend/src/api/app.js`.

### SSE event bus

Backend pushes (quotes, news, streaming AI output) go through `server/events/hub.go`, exposed at `/api/events`.

- `hub.Emit(name, data)` broadcasts to all browsers; `hub.EmitTo(clientID, name, data)` targets one tab. Streaming AI results use `EmitTo` so they only reach the tab that made the request.
- Every browser tab generates a `clientId` (`transport.js`), sent via `X-Go-Stock-Client-Id` header on RPCs and `?clientId=` on the SSE connection. Handlers read it via `server.ClientIDFromContext(ctx)`.
- Delivery is non-blocking: a slow subscriber's events are dropped rather than blocking the broadcast (matters for high-frequency quotes).
- Frontend `runtime.js` reimplements the old Wails `EventsOn/EventsOnce/...` API on top of the SSE subscription in `transport.js`, so ported components keep working unchanged.

### Auth & bootstrap

- `server/auth.go`: if `GO_STOCK_TOKEN` is set, all `/api/*` requires it via `Authorization: Bearer`, `X-Go-Stock-Token` header, or `?token=` query (SSE uses the query form since EventSource can't set headers).
- Frontend `main.js` probes `/api/health` on load: 200 → mount `mobile/MobileApp.vue`; 401 → mount `Login.vue`, then reload into the app after a successful token entry.

### Core runtime state

`server/core.go` — `Core` holds process-level state shared across handlers: freecache, the cron scheduler, AI tool set, the events Hub, cron entry IDs, per-scope cancel funcs for streaming AI/agent tasks, and stock-alert dedup state. `main.go` builds it once and passes it into every handler.

### Startup order (main.go)

`db.Init` → `InitAnalyzeSentiment` → `server.RunMigrate` (AutoMigrate all tables + idempotent seed of A/HK/US basic data from `build/*.json`) → build `Core` → `modules.InitCronTasks(core)` (rebuild enabled cron jobs) → `server.Start` (HTTP). DB and cron are ready before HTTP listens, so early frontend requests never hit an uninitialized DB.

### Frontend layout

- Entry: `main.js` → `mobile/MobileApp.vue` + `mobile/router.js` (hash history, all pages lazy-loaded, routes under `/mobile/*`).
- `mobile/components/base/*` are the design-system primitives (MButton, MSheet, MToastHost, MobileDock, etc.). `cards/`, `charts/`, `sheets/`, `widgets/` build on them. `composables/` holds shared logic (useToast, usePromptPlaza, klineCalc, indicatorSignals, ...).
- **Mobile feedback conventions**: use `useToast`/`MSheet`/chip for feedback; the plaza Q&A shares `usePromptPlaza`; dangerous deletes use native `confirm`. Prefer toasts over native `alert`.
- Charts use echarts; markdown uses md-editor-v3. `vite.config.js` splits these into `chart-vendor`/`markdown-vendor` chunks and raises the PWA precache size limit (md-editor is ~2.3MB).

### Other components

- `ai-assistant-web/` is a **separate, standalone** Go+Vue mini-app (its own `cmd/`, embeds its built frontend via `//go:embed static`). It reuses `backend/data`, `backend/db`, `backend/models` but is not part of the main server. Don't confuse it with `frontend/`.
- `backend/agent/` — AI agent orchestration (eino-based), chat memory, model factory across providers, cron-task API, MCP tools.
- Data source details: `docs/data-sources.md`. User manual: `docs/go-stock使用手册.md`.

## Working notes

- Do the whole change, then test — don't test piecemeal. Don't write unit tests unless asked. Don't build the frontend on this machine (memory).
- Match the existing bilingual comment style (Chinese prose comments are the norm in this codebase).
- Push to a new branch, never `main`, unless told otherwise. Only commit when explicitly asked.
