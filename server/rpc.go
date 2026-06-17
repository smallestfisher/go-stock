package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// Args 是一次 RPC 调用的位置参数切片（按前端 JS 调用顺序）。
type Args = []json.RawMessage

// Handler 是一个 RPC 方法的服务端实现：接收位置参数，返回任意可 JSON 序列化的结果。
// ctx 取自 HTTP 请求，便于把取消信号透传给采集/流式逻辑。
type Handler func(ctx context.Context, core *Core, args Args) (any, error)

var rpcHandlers = map[string]Handler{}

// Register 把一个方法注册到 RPC 注册表（由各模块在 init() 中调用）。
func Register(name string, h Handler) {
	rpcHandlers[name] = h
}

func lookupHandler(name string) (Handler, bool) {
	h, ok := rpcHandlers[name]
	return h, ok
}

// RPCDispatcher 返回 /api/rpc/{Method} 的处理函数。
//   - 请求体: {"args": [位置参数...]}
//   - 成功: 200，响应体即结果的 JSON（与 Wails 直接返回 Go 值的语义一致）
//   - 失败: 非 2xx，响应体 {"error":"..."}，前端桩据此 reject
func RPCDispatcher(core *Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		name := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/rpc/"))
		h, ok := lookupHandler(name)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "method not found: " + name})
			return
		}
		var body struct {
			Args Args `json:"args"`
		}
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&body)
		}
		result, err := h(r.Context(), core, body.Args)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

// --- 位置参数解析辅助（i 为下标，越界/类型不符返回零值，保证健壮）---

func ArgString(args Args, i int) string {
	if i < 0 || i >= len(args) {
		return ""
	}
	var v string
	_ = json.Unmarshal(args[i], &v)
	return v
}

func ArgBool(args Args, i int) bool {
	if i < 0 || i >= len(args) {
		return false
	}
	var v bool
	_ = json.Unmarshal(args[i], &v)
	return v
}

func ArgInt(args Args, i int) int {
	if i < 0 || i >= len(args) {
		return 0
	}
	var v int
	_ = json.Unmarshal(args[i], &v)
	return v
}

func ArgInt64(args Args, i int) int64 {
	if i < 0 || i >= len(args) {
		return 0
	}
	var v int64
	_ = json.Unmarshal(args[i], &v)
	return v
}

func ArgFloat64(args Args, i int) float64 {
	if i < 0 || i >= len(args) {
		return 0
	}
	var v float64
	_ = json.Unmarshal(args[i], &v)
	return v
}

// ArgJSON 把第 i 个参数反序列化到任意类型 T（用于 Group、Prompt、SettingConfig 等结构体参数）。
func ArgJSON[T any](args Args, i int) T {
	var v T
	if i >= 0 && i < len(args) {
		_ = json.Unmarshal(args[i], &v)
	}
	return v
}

// ArgIntSlice / ArgUintSlice 解析数组参数（异动类型、批量删除 id 等）。
func ArgIntSlice(args Args, i int) []int {
	var v []int
	if i >= 0 && i < len(args) {
		_ = json.Unmarshal(args[i], &v)
	}
	return v
}

func ArgUintSlice(args Args, i int) []uint {
	var v []uint
	if i >= 0 && i < len(args) {
		_ = json.Unmarshal(args[i], &v)
	}
	return v
}

// ArgIntPtr 解析可选的 *int（前端传 number 或 null）。用于 sysPromptId 这类可空参数。
func ArgIntPtr(args Args, i int) *int {
	if i < 0 || i >= len(args) {
		return nil
	}
	raw := bytes.TrimSpace(args[i])
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil
	}
	var v int
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil
	}
	return &v
}
