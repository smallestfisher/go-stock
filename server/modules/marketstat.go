package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("GetTodayMarketStatistic", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewMarketStatisticApi().GetTodayData(), nil
	})

	server.Register("GetMarketStatisticByDate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketStatisticApi().GetByDate(server.ArgString(args, 0)), nil
	})

	server.Register("GetRecentDaysMarketStatistic", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketStatisticApi().GetRecentDaysData(server.ArgInt(args, 0)), nil
	})
}
