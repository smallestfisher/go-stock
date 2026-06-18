package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-stock/backend/data"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	// --- 非流式 AI 结果存取 ---
	server.Register("SaveAIResponseResult", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		data.NewDeepSeekOpenAi(context.Background(), server.ArgInt(args, 5)).SaveAIResponseResult(
			server.ArgString(args, 0), server.ArgString(args, 1), server.ArgString(args, 2),
			server.ArgString(args, 3), server.ArgString(args, 4),
		)
		return nil, nil
	})

	server.Register("GetAIResponseResult", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewDeepSeekOpenAi(context.Background(), 0).GetAIResponseResult(server.ArgString(args, 0)), nil
	})

	// --- 提示词模板 ---
	server.Register("GetPromptTemplates", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewPromptTemplateApi().GetPromptTemplates(server.ArgString(args, 0), server.ArgString(args, 1)), nil
	})

	server.Register("AddPrompt", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		p := server.ArgJSON[models.Prompt](args, 0)
		tpl := models.PromptTemplate{ID: p.ID, Content: p.Content, Name: p.Name, Type: p.Type}
		return data.NewPromptTemplateApi().AddPrompt(tpl), nil
	})

	server.Register("DelPrompt", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewPromptTemplateApi().DelPrompt(uint(server.ArgInt64(args, 0))), nil
	})

	// --- 流式 AI 分析：NewChatStream ---
	// 该方法为 fire-and-forget，边收 channel 边 EventsEmit("newChatStream")。
	// Web 版在后台 goroutine 中跑流，并通过事件总线 hub 推送；handler 立即返回，
	// 前端调用后经 /api/events 监听 "newChatStream" 收数据。
	server.Register("NewChatStream", func(ctx context.Context, core *server.Core, args server.Args) (any, error) {
		stock := server.ArgString(args, 0)
		stockCode := server.ArgString(args, 1)
		question := server.ArgString(args, 2)
		aiConfigId := server.ArgInt(args, 3)
		sysPromptId := server.ArgIntPtr(args, 4)
		enableTools := server.ArgBool(args, 5)
		think := server.ArgBool(args, 6)
		clientID := server.ClientIDFromContext(ctx)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.SugaredLogger.Errorf("NewChatStream panic: %v", r)
					emitClientEvent(core, clientID, "newChatStream", map[string]any{
						"code":    0,
						"content": fmt.Sprintf("AI分析异常: %v", r),
					})
					emitClientEvent(core, clientID, "newChatStream", "DONE")
				}
			}()
			var tools []data.Tool
			if enableTools {
				tools = core.AiTools
			}
			msgs := data.NewDeepSeekOpenAi(context.Background(), aiConfigId).NewChatStream(stock, stockCode, question, sysPromptId, tools, think)
			for msg := range msgs {
				emitClientEvent(core, clientID, "newChatStream", msg)
			}
			emitClientEvent(core, clientID, "newChatStream", "DONE")
		}()
		return nil, nil
	})

	// --- 流式市场资讯总结：SummaryStockNews ---
	// 支持自定义事件名(默认 summaryStockNews)与对话历史；可被前端 AbortSummaryStockNews 中断。
	server.Register("SummaryStockNews", func(reqCtx context.Context, core *server.Core, args server.Args) (any, error) {
		question := server.ArgString(args, 0)
		aiConfigId := server.ArgInt(args, 1)
		sysPromptId := server.ArgIntPtr(args, 2)
		enableTools := server.ArgBool(args, 3)
		think := server.ArgBool(args, 4)
		eventName := strings.TrimSpace(server.ArgString(args, 5))
		if eventName == "" {
			eventName = "summaryStockNews"
		}
		history := parseHistory(server.ArgString(args, 6))
		clientID := server.ClientIDFromContext(reqCtx)

		ctx, cancel := context.WithCancel(context.Background())
		core.SetSummaryCancel(clientID, cancel)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.SugaredLogger.Errorf("SummaryStockNews panic: %v", r)
					emitClientEvent(core, clientID, eventName, map[string]any{
						"code":    0,
						"content": fmt.Sprintf("AI分析异常: %v", r),
					})
				}
			}()
			var msgs <-chan map[string]any
			if enableTools {
				msgs = data.NewDeepSeekOpenAi(ctx, aiConfigId).NewSummaryStockNewsStreamWithTools(question, sysPromptId, core.AiTools, think, history)
			} else {
				msgs = data.NewDeepSeekOpenAi(ctx, aiConfigId).NewSummaryStockNewsStream(question, sysPromptId, think, history)
			}
			for msg := range msgs {
				emitClientEvent(core, clientID, eventName, msg)
			}
			emitClientEvent(core, clientID, eventName, "DONE")
		}()
		return nil, nil
	})

	server.Register("AbortSummaryStockNews", func(ctx context.Context, core *server.Core, _ server.Args) (any, error) {
		core.CancelSummary(server.ClientIDFromContext(ctx))
		return nil, nil
	})

	server.Register("AbortChatWithAgent", func(ctx context.Context, core *server.Core, _ server.Args) (any, error) {
		core.CancelAgent(server.ClientIDFromContext(ctx))
		return nil, nil
	})
}

func emitClientEvent(core *server.Core, clientID string, eventName string, data any) {
	if clientID == "" {
		core.Events.Emit(eventName, data)
		return
	}
	core.Events.EmitTo(clientID, eventName, data)
}

// parseHistory 解析 AI 助手对话历史 JSON。
func parseHistory(historyJSON string) []map[string]interface{} {
	historyJSON = strings.TrimSpace(historyJSON)
	if historyJSON == "" {
		return nil
	}
	var list []models.AiAssistantMessage
	if err := json.Unmarshal([]byte(historyJSON), &list); err != nil || len(list) == 0 {
		return nil
	}
	history := make([]map[string]interface{}, 0, len(list))
	for _, m := range list {
		item := map[string]interface{}{"role": m.Role, "content": m.Content}
		if m.Role == "assistant" && m.Reasoning != "" {
			item["reasoning_content"] = m.Reasoning
		}
		history = append(history, item)
	}
	return history
}
