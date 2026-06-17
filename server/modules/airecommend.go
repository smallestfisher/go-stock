package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	server.Register("GetAiRecommendStocksList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[models.AiRecommendStocksQuery](args, 0)
		page, err := data.NewAiRecommendStocksService().GetAiRecommendStocksList(&query)
		if err != nil {
			return &models.AiRecommendStocksPageData{}, nil
		}
		return page, nil
	})

	server.Register("DeleteAiRecommendStocks", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewAiRecommendStocksService().DeleteAiRecommendStocks(uint(server.ArgInt64(args, 0))); err != nil {
			return "删除失败", nil
		}
		return "删除成功", nil
	})

	server.Register("UpdateAiRecommendStocksAlert", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewAiRecommendStocksService().UpdateAiRecommendStocksAlert(uint(server.ArgInt64(args, 0)), server.ArgBool(args, 1)); err != nil {
			return "更新预警状态失败", nil
		}
		return "更新预警状态成功", nil
	})
}
