package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	// 热门选股策略
	server.Register("GetHotStrategy", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewSearchStockApi("").HotStrategy(), nil
	})

	// 自定义策略
	server.Register("GetCustomStrategyList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[models.CustomStrategyQuery](args, 0)
		page, err := data.NewCustomStrategyApi().GetCustomStrategyList(&query)
		if err != nil {
			return &models.CustomStrategyPageData{}, nil
		}
		return page, nil
	})

	server.Register("GetAllCustomStrategies", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewCustomStrategyApi().GetAllCustomStrategies(), nil
	})

	server.Register("SaveCustomStrategy", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewCustomStrategyApi().SaveCustomStrategy(server.ArgJSON[models.CustomStrategy](args, 0)), nil
	})

	server.Register("DeleteCustomStrategy", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewCustomStrategyApi().DeleteCustomStrategy(uint(server.ArgInt64(args, 0))), nil
	})

	// 指标选股
	server.Register("GetAllStocks", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().GetAllStocks(
			server.ArgInt(args, 0), server.ArgInt(args, 1), server.ArgString(args, 2),
			server.ArgJSON[models.TechnicalIndicators](args, 3),
		), nil
	})

	// 全市场股票库（选股池）
	server.Register("GetAllStockInfoList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[data.AllStockInfoQuery](args, 0)
		page, err := data.NewStockDataApi().GetAllStockInfoList(&query)
		if err != nil {
			return &data.AllStockInfoPageData{}, nil
		}
		return page, nil
	})

	server.Register("GetAllStockInfoById", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		stock, err := data.NewStockDataApi().GetAllStockInfoById(uint(server.ArgInt64(args, 0)))
		if err != nil {
			return &models.AllStockInfo{}, nil
		}
		return stock, nil
	})

	server.Register("AddAllStockInfo", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewStockDataApi().AddAllStockInfo(server.ArgJSON[models.AllStockInfo](args, 0)); err != nil {
			return "操作失败: " + err.Error(), nil
		}
		return "操作成功", nil
	})

	server.Register("DeleteAllStockInfo", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewStockDataApi().DeleteAllStockInfo(uint(server.ArgInt64(args, 0))); err != nil {
			return "删除失败: " + err.Error(), nil
		}
		return "删除成功", nil
	})

	server.Register("BatchDeleteAllStockInfo", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewStockDataApi().BatchDeleteAllStockInfo(server.ArgUintSlice(args, 0)); err != nil {
			return "批量删除失败: " + err.Error(), nil
		}
		return "批量删除成功", nil
	})

	server.Register("GetAllMarkets", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		markets, err := data.NewStockDataApi().GetAllMarkets()
		if err != nil {
			return []string{}, nil
		}
		return markets, nil
	})

	server.Register("GetAllIndustries", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		industries, err := data.NewStockDataApi().GetAllIndustries()
		if err != nil {
			return []string{}, nil
		}
		return industries, nil
	})

	server.Register("GetAllConcepts", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		concepts, err := data.NewStockDataApi().GetAllConcepts()
		if err != nil {
			return []string{}, nil
		}
		return concepts, nil
	})
}
