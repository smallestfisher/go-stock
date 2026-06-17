package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	// AI 分析结果历史
	server.Register("GetAIResponseResultList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[models.AIResponseResultQuery](args, 0)
		page, err := data.NewAIResponseResultService().GetAIResponseResultList(query)
		if err != nil {
			return &models.AIResponseResultPageData{}, nil
		}
		return page, nil
	})

	server.Register("DeleteAIResponseResult", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewAIResponseResultService().DeleteAIResponseResult(uint(server.ArgInt64(args, 0))); err != nil {
			return "删除失败", nil
		}
		return "删除成功", nil
	})

	server.Register("BatchDeleteAIResponseResult", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewAIResponseResultService().BatchDeleteAIResponseResult(server.ArgUintSlice(args, 0)); err != nil {
			return "删除失败", nil
		}
		return "删除成功", nil
	})

	// 提示词模板（分页管理版，与 ai.go 的 AddPrompt/DelPrompt/GetPromptTemplates 互补）
	server.Register("GetPromptTemplateList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[models.PromptTemplateQuery](args, 0)
		page, err := data.NewPromptTemplateApi().GetPromptTemplateList(&query)
		if err != nil {
			return &models.PromptTemplatePageData{}, nil
		}
		return page, nil
	})

	server.Register("AddPromptTemplate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewPromptTemplateApi().AddPrompt(server.ArgJSON[models.PromptTemplate](args, 0)), nil
	})

	server.Register("UpdatePromptTemplate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewPromptTemplateApi().AddPrompt(server.ArgJSON[models.PromptTemplate](args, 0)), nil
	})

	server.Register("DeletePromptTemplate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewPromptTemplateApi().DelPrompt(uint(server.ArgInt64(args, 0))), nil
	})
}
