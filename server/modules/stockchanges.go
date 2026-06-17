package modules

import (
	"context"
	"fmt"

	"go-stock/backend/data"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	// 当日异动
	server.Register("GetStockChanges", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockChangesApi().GetStockChanges(server.ArgIntSlice(args, 0), server.ArgInt(args, 1), server.ArgInt(args, 2)), nil
	})

	server.Register("GetAllStockChangesWithPaging", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		all := data.NewStockChangesApi().GetAllStockChangesWithPaging(server.ArgInt(args, 0))
		historyService := data.NewStockChangeHistoryService()
		_, _ = historyService.SaveStockChangesWithDedup(all.Data)
		return all, nil
	})

	// 异动历史与统计
	server.Register("GetStockChangeHistory", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[models.StockChangeHistoryQuery](args, 0)
		result, err := data.NewStockChangeHistoryService().GetHistoryList(query)
		if err != nil {
			return &models.StockChangeHistoryPageData{}, nil
		}
		return result, nil
	})

	server.Register("SaveStockChangesToHistory", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		changeTypes := server.ArgIntSlice(args, 0)
		result := data.NewStockChangesApi().GetStockChanges(changeTypes, 0, 500)
		if result == nil || len(result.Data) == 0 {
			return "没有获取到异动数据", nil
		}
		if err := data.NewStockChangeHistoryService().SaveStockChanges(result.Data); err != nil {
			return "保存失败: " + err.Error(), nil
		}
		return fmt.Sprintf("成功保存 %d 条异动数据", len(result.Data)), nil
	})

	server.Register("DeleteStockChangeHistory", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewStockChangeHistoryService().DeleteOldData(server.ArgInt(args, 0)); err != nil {
			return "删除失败: " + err.Error(), nil
		}
		return fmt.Sprintf("已删除 %d 天前的历史数据", server.ArgInt(args, 0)), nil
	})

	server.Register("GetDailyChangeStats", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		result, err := data.NewStockChangeHistoryService().GetDailyChangeStats(server.ArgInt(args, 0))
		if err != nil {
			return []data.DailyChangeStats{}, nil
		}
		return result, nil
	})

	server.Register("GetChangeTypeDailyStats", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		result, err := data.NewStockChangeHistoryService().GetChangeTypeDailyStats(server.ArgInt(args, 0))
		if err != nil {
			return []data.ChangeTypeDailyStats{}, nil
		}
		return result, nil
	})

	server.Register("GetChangeRank", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		result, err := data.NewStockChangeHistoryService().GetChangeRank(server.ArgInt(args, 0), server.ArgInt(args, 1))
		if err != nil {
			return &data.ChangeRankResult{}, nil
		}
		return result, nil
	})

	server.Register("GetDailyDimensionStats", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		result, err := data.NewStockChangeHistoryService().GetDailyDimensionStats(server.ArgString(args, 0), server.ArgString(args, 1), server.ArgInt(args, 2))
		if err != nil {
			return []data.DailyDimensionStats{}, nil
		}
		return result, nil
	})

	server.Register("GetTypeStatsByDate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		result, err := data.NewStockChangeHistoryService().GetTypeStatsByDate(server.ArgString(args, 0))
		if err != nil {
			return []data.TypeCountStats{}, nil
		}
		return result, nil
	})
}
