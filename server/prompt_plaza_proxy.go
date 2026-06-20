package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-stock/backend/data"
	log "go-stock/backend/logger"
)

const (
	defaultPromptPlazaAPIBase = "http://go-stock.sparkmemory.top:1918/api"
	promptPlazaProxyPrefix    = "/api/prompt-plaza"
)

var promptPlazaHTTPClient = &http.Client{Timeout: 30 * time.Second}

func promptPlazaProxyHandler(w http.ResponseWriter, r *http.Request) {
	targetBase, err := getPromptPlazaTargetBase()
	if err != nil {
		writePromptPlazaProxyError(w, http.StatusBadGateway, err.Error())
		return
	}

	targetURL := *targetBase
	proxyPath := strings.TrimPrefix(r.URL.Path, promptPlazaProxyPrefix)
	targetURL.Path = joinURLPath(targetBase.Path, proxyPath)
	query := r.URL.Query()
	query.Del("token")
	targetURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL.String(), r.Body)
	if err != nil {
		writePromptPlazaProxyError(w, http.StatusBadGateway, "创建提示词广场请求失败: "+err.Error())
		return
	}
	copyPromptPlazaRequestHeaders(req.Header, r.Header)

	resp, err := promptPlazaHTTPClient.Do(req)
	if err != nil {
		log.SugaredLogger.Warnf("提示词广场代理请求失败: %s %s: %v", r.Method, targetURL.String(), err)
		writePromptPlazaProxyError(w, http.StatusBadGateway, "提示词广场接口不可达: "+err.Error())
		return
	}
	defer resp.Body.Close()

	copyPromptPlazaResponseHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func getPromptPlazaTargetBase() (*url.URL, error) {
	base := ""
	if cfg := data.GetSettingConfig(); cfg != nil && cfg.Settings != nil {
		base = strings.TrimSpace(cfg.PromptPlazaApiBase)
	}
	if base == "" {
		base = defaultPromptPlazaAPIBase
	}

	parsed, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("提示词广场地址配置不正确: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("提示词广场地址只支持 http/https: %s", base)
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("提示词广场地址缺少主机名: %s", base)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return parsed, nil
}

func joinURLPath(basePath, proxyPath string) string {
	basePath = strings.TrimRight(basePath, "/")
	proxyPath = strings.TrimLeft(proxyPath, "/")
	if proxyPath == "" {
		if basePath == "" {
			return "/"
		}
		return basePath
	}
	if basePath == "" {
		return "/" + proxyPath
	}
	return basePath + "/" + proxyPath
}

func copyPromptPlazaRequestHeaders(dst, src http.Header) {
	for key, values := range src {
		if isHopByHopHeader(key) || strings.EqualFold(key, "X-Go-Stock-Token") {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func copyPromptPlazaResponseHeaders(dst, src http.Header) {
	for key, values := range src {
		if isHopByHopHeader(key) || strings.HasPrefix(strings.ToLower(key), "access-control-") {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func isHopByHopHeader(key string) bool {
	switch strings.ToLower(key) {
	case "connection", "keep-alive", "proxy-authenticate", "proxy-authorization",
		"te", "trailer", "transfer-encoding", "upgrade":
		return true
	default:
		return false
	}
}

func writePromptPlazaProxyError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"code":    status,
		"message": message,
	})
}
