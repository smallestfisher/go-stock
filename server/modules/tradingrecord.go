package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("AddTradingRecord", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		record := server.ArgJSON[data.TradingRecord](args, 0)
		id, _ := data.NewStockDataApi().AddTradingRecord(record)
		return id, nil
	})

	server.Register("GetTradingRecordList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[data.TradingRecordListQuery](args, 0)
		page, err := data.NewStockDataApi().GetTradingRecordList(query)
		if err != nil {
			return &data.TradingRecordPageData{}, nil
		}
		return page, nil
	})

	server.Register("GetTradingRecordById", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		r, err := data.NewStockDataApi().GetTradingRecordById(uint(server.ArgInt64(args, 0)))
		if err != nil {
			return nil, nil
		}
		return r, nil
	})

	server.Register("GetTradingRecordStatistics", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		stats, err := data.NewStockDataApi().GetTradingRecordStatistics()
		if err != nil {
			return &data.TradingRecordStatistics{}, nil
		}
		return stats, nil
	})

	server.Register("UpdateTradingRecord", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		_ = data.NewStockDataApi().UpdateTradingRecord(server.ArgJSON[data.TradingRecord](args, 0))
		return nil, nil
	})

	server.Register("DeleteTradingRecord", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		_ = data.NewStockDataApi().DeleteTradingRecord(uint(server.ArgInt64(args, 0)))
		return nil, nil
	})

	server.Register("CheckFrequentTrading", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		canTrade, msg := data.NewStockDataApi().CheckFrequentTrading(server.ArgString(args, 0))
		return map[string]any{"canTrade": canTrade, "msg": msg}, nil
	})
}
