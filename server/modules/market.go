package modules

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("GetIndustryRank", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		res := data.NewMarketNewsApi().GetIndustryRank(server.ArgString(args, 0), server.ArgInt(args, 1))
		return res["data"].([]any), nil
	})

	server.Register("GetIndustryMoneyRankSina", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().GetIndustryMoneyRankSina(server.ArgString(args, 0), server.ArgString(args, 1)), nil
	})

	server.Register("GetMoneyRankSina", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().GetMoneyRankSina(server.ArgString(args, 0)), nil
	})

	server.Register("GetStockMoneyTrendByDay", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		res := data.NewMarketNewsApi().GetStockMoneyTrendByDay(server.ArgString(args, 0), server.ArgInt(args, 1))
		slice.Reverse(res)
		return res, nil
	})
}
