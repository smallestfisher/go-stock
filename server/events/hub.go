// Package events 提供 server 端的 SSE 事件总线。
// 后端任意 goroutine 调用 hub.Emit(name, data),所有连接到 /api/events 的浏览器
// 客户端会收到一条对应的 SSE 消息(event: name)，前端按事件名分发即可。
package events

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type eventPayload struct {
	name string
	data any
}

type subscriber struct {
	clientID string
	ch       chan eventPayload
}

// Hub 维护一组 SSE 订阅者，并把 Emit 的事件广播给它们。
type Hub struct {
	mu   sync.RWMutex
	subs map[*subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{subs: make(map[*subscriber]struct{})}
}

// Emit 向所有订阅者广播一个事件。采用非阻塞投递：若某订阅者缓冲已满则丢弃该条，
// 避免一个慢客户端拖垮整个广播（高频行情推送下尤其重要）。
func (h *Hub) Emit(name string, data any) {
	h.deliver("", name, data)
}

// EmitTo 向指定浏览器会话投递事件。用于 AI 流式输出这类请求发起者专属数据，
// 避免多个浏览器/标签页之间串流。
func (h *Hub) EmitTo(clientID string, name string, data any) {
	if clientID == "" {
		h.Emit(name, data)
		return
	}
	h.deliver(clientID, name, data)
}

func (h *Hub) deliver(clientID string, name string, data any) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for s := range h.subs {
		if clientID != "" && s.clientID != clientID {
			continue
		}
		select {
		case s.ch <- eventPayload{name: name, data: data}:
		default:
		}
	}
}

func (h *Hub) subscribe(clientID string) *subscriber {
	s := &subscriber{clientID: clientID, ch: make(chan eventPayload, 256)}
	h.mu.Lock()
	h.subs[s] = struct{}{}
	h.mu.Unlock()
	return s
}

func (h *Hub) unsubscribe(s *subscriber) {
	h.mu.Lock()
	delete(h.subs, s)
	h.mu.Unlock()
	close(s.ch)
}

// ServeSSE 是 /api/events 的处理函数：建立长连接，持续把事件以 SSE 格式写回客户端。
func (h *Hub) ServeSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// 关闭代理缓冲，确保事件能即时推送（nginx 等反代需额外配置 X-Accel-Buffering: no）
	w.Header().Set("X-Accel-Buffering", "no")

	clientID := r.URL.Query().Get("clientId")
	if clientID == "" {
		clientID = r.Header.Get("X-Go-Stock-Client-Id")
	}
	s := h.subscribe(clientID)
	defer h.unsubscribe(s)

	ctx := r.Context()
	// 先写一条注释行，让客户端立刻确认连接已建立
	_, _ = w.Write([]byte(": connected\n\n"))
	flusher.Flush()

	// 心跳：每 15s 写一行注释。EventSource 会忽略 ":" 开头的内容，仅重置"最后收到时间"。
	// 既防止反代（nginx 默认 60s、Cloudflare 100s）空闲超时掐断连接，
	// 也让客户端在连接静默死亡后更快被发现并重连。
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = w.Write([]byte(": ping\n\n"))
			flusher.Flush()
		case p, ok := <-s.ch:
			if !ok {
				return
			}
			// 把事件名与载荷打包成一帧，作为默认 message 下发；前端 EventSource.onmessage
			// 解析出 frame.event 后按名分发（这样前端无需为每个事件名预先注册监听器）。
			frame := struct {
				Event string `json:"event"`
				Data  any    `json:"data"`
			}{Event: p.name, Data: p.data}
			raw, err := json.Marshal(frame)
			if err != nil {
				continue
			}
			_, _ = w.Write([]byte("data: " + string(raw) + "\n\n"))
			flusher.Flush()
		}
	}
}
