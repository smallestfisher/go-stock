package modules

import (
	"context"
	"strings"

	"go-stock/backend/data"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	server.Register("GetAiConfigs", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.GetSettingConfig().AiConfigs, nil
	})

	server.Register("GetAiAssistantSession", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		res, err := data.GetAiAssistantSession(server.ArgString(args, 0))
		if err != nil {
			return nil, nil
		}
		return res, nil
	})

	server.Register("SaveAiAssistantSession", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		sessionId := server.ArgString(args, 0)
		messages := server.ArgJSON[[]models.AiAssistantMessage](args, 1)
		if err := data.SaveAiAssistantSession(sessionId, messages); err != nil {
			return err.Error(), nil
		}
		return nil, nil
	})

	// FetchAiModels：根据接口地址+apiKey 拉 /models 列表（OpenAI/DeepSeek 兼容）。
	server.Register("FetchAiModels", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		baseUrl := strings.TrimSpace(server.ArgString(args, 0))
		apiKey := strings.TrimSpace(server.ArgString(args, 1))
		clientProfile := strings.TrimSpace(server.ArgString(args, 2))
		clientVersion := strings.TrimSpace(server.ArgString(args, 3))
		if baseUrl == "" || apiKey == "" {
			return []string{}, nil
		}
		type modelItem struct {
			ID string `json:"id"`
		}
		var respData struct {
			Data []modelItem `json:"data"`
		}
		client := data.SharedHTTPClient
		client.SetBaseURL(baseUrl)
		req := client.R().
			SetHeader("Authorization", "Bearer "+apiKey).
			SetHeader("Content-Type", "application/json").
			SetResult(&respData)
		req = data.ApplyAIClientProfileHeaders(req, clientProfile, clientVersion)
		resp, err := req.Get("/models")
		if err != nil {
			log.SugaredLogger.Errorf("FetchAiModels error: %v", err)
			return []string{}, nil
		}
		if resp.IsError() {
			return []string{}, nil
		}
		out := make([]string, 0, len(respData.Data))
		for _, m := range respData.Data {
			if strings.TrimSpace(m.ID) != "" {
				out = append(out, m.ID)
			}
		}
		return out, nil
	})
}
