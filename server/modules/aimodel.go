package modules

import (
	"context"
	"strings"

	"go-stock/backend/data"
	"go-stock/server"
)

// aiModelInfo 对应桌面端 app.go 的 AiModelInfo（字段 json 标签一致）。
type aiModelInfo struct {
	ModelName string `json:"modelName"`
	MaxTokens int    `json:"maxTokens"`
	Source    string `json:"source"`
}

func init() {
	// FetchAiModelInfo：探测模型最大 token（优先接口 /models/{name}，回退内置表）。复刻 app.go:2630。
	server.Register("FetchAiModelInfo", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		baseUrl := strings.TrimSpace(server.ArgString(args, 0))
		apiKey := strings.TrimSpace(server.ArgString(args, 1))
		modelName := strings.TrimSpace(server.ArgString(args, 2))
		if baseUrl == "" || modelName == "" {
			return nil, nil
		}
		info := &aiModelInfo{ModelName: modelName, MaxTokens: 0, Source: ""}

		if apiKey != "" {
			type modelDetail struct {
				ID             string `json:"id"`
				MaxContextLen  int    `json:"max_context_length"`
				ContextLength  int    `json:"context_length"`
				MaxOutputTok   int    `json:"max_output_tokens"`
				MaxTokensField int    `json:"max_tokens"`
			}
			var detail modelDetail
			client := data.SharedHTTPClient
			client.SetBaseURL(baseUrl)
			resp, err := client.R().
				SetHeader("Authorization", "Bearer "+apiKey).
				SetHeader("Content-Type", "application/json").
				SetResult(&detail).
				Get("/models/" + modelName)
			if err == nil && !resp.IsError() && detail.ID != "" {
				switch {
				case detail.MaxContextLen > 0:
					info.MaxTokens, info.Source = detail.MaxContextLen, "api"
				case detail.ContextLength > 0:
					info.MaxTokens, info.Source = detail.ContextLength, "api"
				case detail.MaxOutputTok > 0:
					info.MaxTokens, info.Source = detail.MaxOutputTok, "api"
				case detail.MaxTokensField > 0:
					info.MaxTokens, info.Source = detail.MaxTokensField, "api"
				}
			}
		}

		if info.MaxTokens == 0 {
			if mt := getBuiltinModelMaxTokens(modelName); mt > 0 {
				info.MaxTokens, info.Source = mt, "builtin"
			}
		}
		return info, nil
	})
}

// getBuiltinModelMaxTokens 内置模型 max_tokens 兜底表（复刻 app.go:2689）。
func getBuiltinModelMaxTokens(modelName string) int {
	modelTokenMap := map[string]int{
		"deepseek-chat": 65536, "deepseek-reasoner": 65536, "deepseek-coder": 16384,
		"deepseek-v3": 65536, "deepseek-r1": 65536,
		"gpt-4o": 16384, "gpt-4o-mini": 16384, "gpt-4o-2024-05-13": 4096,
		"gpt-4-turbo": 4096, "gpt-4-turbo-preview": 4096, "gpt-4": 8192, "gpt-4-32k": 32768,
		"gpt-3.5-turbo": 4096, "gpt-3.5-turbo-16k": 16384,
		"gpt-4.1": 32768, "gpt-4.1-mini": 32768, "gpt-4.1-nano": 32768,
		"o1": 100000, "o1-mini": 65536, "o1-preview": 32768, "o3-mini": 100000, "o4-mini": 100000,
		"claude-3-5-sonnet": 8192, "claude-3-5-haiku": 8192,
		"claude-3-opus": 4096, "claude-3-sonnet": 4096, "claude-3-haiku": 4096,
		"glm-4": 8192, "glm-4-plus": 4096, "glm-4-air": 4096, "glm-4-flash": 4096, "glm-4-long": 4096,
		"chatglm-turbo": 4096,
		"moonshot-v1-8k": 8192, "moonshot-v1-32k": 32768, "moonshot-v1-128k": 131072,
		"qwen-turbo": 8192, "qwen-plus": 131072, "qwen-max": 8192, "qwen-long": 65536,
		"qwen2.5-72b-instruct": 32768,
		"hunyuan-lite": 4096, "hunyuan-standard": 4096, "hunyuan-pro": 4096, "hunyuan-turbo": 4096,
		"spark-lite": 4096, "spark-pro": 4096, "spark-max": 4096, "spark-4.0-ultra": 4096,
		"yi-light": 16384, "yi-large": 16384, "yi-medium": 16384, "yi-spark": 16384, "yi-vision": 16384,
		"abab6.5-chat": 8192, "abab6.5s-chat": 8192, "abab5.5-chat": 4096,
		"baichuan2-turbo": 4096, "baichuan2-53b": 4096,
		"ernie-4.0": 4096, "ernie-3.5": 4096, "ernie-speed": 4096, "ernie-lite": 4096,
	}
	if mt, ok := modelTokenMap[modelName]; ok {
		return mt
	}
	for prefix, mt := range map[string]int{
		"deepseek": 65536, "gpt-4o": 16384, "gpt-4-turbo": 4096, "gpt-4-": 8192, "gpt-3.5": 4096, "gpt-4.1": 32768,
		"o1-": 65536, "o3-": 100000, "o4-": 100000, "claude-3": 8192, "glm-4": 8192, "chatglm": 4096,
		"moonshot-v1": 8192, "qwen-": 8192, "qwen2": 32768, "hunyuan-": 4096, "spark-": 4096, "yi-": 16384,
		"abab": 8192, "baichuan": 4096, "ernie-": 4096, "llama-3": 8192, "llama3": 8192,
		"mistral-": 8192, "mixtral-": 32768, "codestral-": 32768, "gemini-1.5": 8192, "gemini-2": 8192,
		"command-r": 4096, "Qwen/Qwen": 32768, "deepseek-ai/": 65536, "meta-llama/": 8192,
		"mistralai/": 32768, "Pro/deepseek-": 65536, "Pro/qwen-": 32768,
	} {
		if strings.HasPrefix(modelName, prefix) {
			return mt
		}
	}
	return 0
}
