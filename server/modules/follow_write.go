package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("Follow", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().Follow(server.ArgString(args, 0)), nil
	})

	server.Register("UnFollow", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().UnFollow(server.ArgString(args, 0)), nil
	})

	server.Register("SetCostPriceAndVolume", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().SetCostPriceAndVolume(server.ArgFloat64(args, 1), server.ArgInt64(args, 2), server.ArgString(args, 0)), nil
	})

	server.Register("SetTradingPrice", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().SetTradingPrice(
			server.ArgFloat64(args, 1), // entryPrice
			server.ArgFloat64(args, 2), // takeProfitPrice
			server.ArgFloat64(args, 3), // stopLossPrice
			server.ArgFloat64(args, 4), // costPrice
			server.ArgString(args, 0),  // stockCode
		), nil
	})

	server.Register("SetAlarmChangePercent", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().SetAlarmChangePercent(server.ArgFloat64(args, 0), server.ArgFloat64(args, 1), server.ArgString(args, 2)), nil
	})

	server.Register("SetStockSort", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		data.NewStockDataApi().SetStockSort(server.ArgInt64(args, 0), server.ArgString(args, 1))
		return nil, nil
	})
}
