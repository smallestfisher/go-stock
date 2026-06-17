package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "go-stock/backend/logger"
)

const frontendDir = "frontend/dist"

// Start 初始化路由并阻塞运行 HTTP 服务。
// 路由分两层：
//   - /api/*  ：业务接口，统一经过 CORS + 令牌鉴权中间件
//   - /       ：前端静态资源（Vue SPA），从 frontend/dist 提供
func Start(core *Core, addr string) error {
	api := http.NewServeMux()

	api.HandleFunc("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":   true,
			"time": time.Now().Format("2006-01-02 15:04:05"),
		})
	})
	// SSE 事件总线：前端订阅它接收后端推送(行情、资讯、AI 流式等)
	api.HandleFunc("/api/events", core.Events.ServeSSE)

	// 统一 RPC 桥：前端桩 POST /api/rpc/{Method}，由各功能模块注册的 handler 处理。
	api.HandleFunc("/api/rpc/", RPCDispatcher(core))

	// Phase 2 起，流式 AI 等专用 SSE 端点在这里注册。

	root := http.NewServeMux()
	root.Handle("/api/", withCORS(authMiddleware(api)))
	root.Handle("/", withCORS(frontendHandler()))

	log.SugaredLogger.Infof("go-stock server 监听: %s", addr)
	srv := &http.Server{
		Addr:              addr,
		Handler:           root,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

// frontendHandler 从磁盘 frontend/dist 提供 Vue SPA，并做 history 模式的回退。
// dist 尚未构建时返回提示页，避免阻塞后端开发（生产构建后此目录即存在）。
func frontendHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(frontendDir); err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<!doctype html><meta charset="utf-8"><title>go-stock</title>
<body style="font-family:sans-serif;padding:2rem">
<h2>go-stock server 已启动 🚀</h2>
<p>前端尚未构建(<code>frontend/dist</code> 不存在)。</p>
<p>构建前端：<br><code>cd frontend && npm install && npm run build</code></p>
<p>构建完成后刷新本页即可。</p>
</body>`))
			return
		}

		// 路径清理 + 防目录穿越
		rel := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if rel == "." || rel == "" {
			rel = "index.html"
		}
		full := filepath.Join(frontendDir, rel)
		rootAbs, _ := filepath.Abs(frontendDir)
		fullAbs, _ := filepath.Abs(full)
		if !strings.HasPrefix(fullAbs+string(os.PathSeparator), rootAbs+string(os.PathSeparator)) {
			http.NotFound(w, r)
			return
		}

		info, err := os.Stat(full)
		if err != nil || info.IsDir() {
			// SPA history 回退：未知路由交给前端路由处理
			http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
			return
		}
		http.ServeFile(w, r, full)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// withCORS 允许跨域(开发期前端走 vite dev server 独立端口；生产同源时无副作用)。
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
