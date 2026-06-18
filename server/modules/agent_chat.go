package modules

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/schema"
	"go-stock/backend/agent"
	log "go-stock/backend/logger"
	"go-stock/server"
)

func init() {
	// ChatWithAgent：AI 智能体流式对话。后台 goroutine 跑 agent，经事件总线推送 "agent-message"。
	// 可被 AbortChatWithAgent 中断（core.CancelAgent）。handler 立即返回。
	server.Register("ChatWithAgent", func(reqCtx context.Context, core *server.Core, args server.Args) (any, error) {
		question := server.ArgString(args, 0)
		aiConfigId := server.ArgInt(args, 1)
		sysPromptId := server.ArgIntPtr(args, 2)
		memoryMode := server.ArgBool(args, 3)
		memoryCount := server.ArgInt(args, 4)
		thinkingMode := server.ArgBool(args, 5)
		agentMode := server.ArgString(args, 6)
		clientID := server.ClientIDFromContext(reqCtx)

		ctx, cancel := context.WithCancel(context.Background())
		core.SetAgentCancel(clientID, cancel)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.SugaredLogger.Errorf("ChatWithAgent panic: %v", r)
				}
			}()
			ch := agent.NewStockAiAgentApi().ChatWithContext(ctx, question, aiConfigId, sysPromptId, memoryMode, memoryCount, thinkingMode, agentMode)
			for msg := range ch {
				emitClientEvent(core, clientID, "agent-message", agentMessageToMap(msg))
			}
			emitClientEvent(core, clientID, "agent-message", agentMessageToMap(&schema.Message{
				Role:    schema.Assistant,
				Content: "agent-DONE",
			}))
		}()
		return nil, nil
	})
}

// agentMessageToMap 把 schema.Message 序列化为 map 再推送，保证字段名与 json tag 一致。
// 复刻自 app_common.go 的 agentMessageToFrontendMap。
func agentMessageToMap(msg *schema.Message) map[string]any {
	if msg == nil {
		return map[string]any{}
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return map[string]any{
			"role":              string(msg.Role),
			"content":           msg.Content,
			"reasoning_content": msg.ReasoningContent,
		}
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{
			"role":              string(msg.Role),
			"content":           msg.Content,
			"reasoning_content": msg.ReasoningContent,
		}
	}
	return m
}
