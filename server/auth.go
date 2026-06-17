package server

import (
	"net/http"
	"os"
	"strings"

	log "go-stock/backend/logger"
)

// configuredToken 返回当前配置的访问令牌（来自环境变量 GO_STOCK_TOKEN）。
// 为空表示开发模式：不鉴权（仅限本地调试）。
func configuredToken() string {
	return strings.TrimSpace(os.Getenv("GO_STOCK_TOKEN"))
}

// authMiddleware 在配置了令牌时，对 /api/* 强制校验。
// 令牌可通过两种方式提供：
//   - 标准 HTTP 请求头 Authorization: Bearer <token>（fetch/axios）
//   - 查询参数 ?token=<token>（EventSource 无法自定义请求头，SSE 走这条）
func authMiddleware(next http.Handler) http.Handler {
	token := configuredToken()
	if token == "" {
		log.SugaredLogger.Warn("GO_STOCK_TOKEN 未设置：server 当前无鉴权(仅限本地开发)。远程部署前请务必设置令牌。")
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provided := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if provided == "" {
			provided = strings.TrimSpace(r.URL.Query().Get("token"))
		}
		if provided != token {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
