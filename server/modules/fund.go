package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("GetfundList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewFundApi().GetFundList(server.ArgString(args, 0)), nil
	})

	server.Register("GetFollowedFund", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewFundApi().GetFollowedFund(), nil
	})

	server.Register("FollowFund", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewFundApi().FollowFund(server.ArgString(args, 0)), nil
	})

	server.Register("UnFollowFund", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewFundApi().UnFollowFund(server.ArgString(args, 0)), nil
	})

	server.Register("GetFundKLine", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewFundKLineApi().GetFundKLineWithFallback(server.ArgString(args, 0), server.ArgString(args, 1), server.ArgInt(args, 2)), nil
	})

	server.Register("GetFundHistoryNetValue", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		res, _ := data.NewFundApi().GetFundHistoryNetValue(server.ArgString(args, 0), 1, server.ArgInt(args, 1), server.ArgString(args, 2), server.ArgString(args, 3))
		if res == nil {
			return []data.FundHistoryNetValue{}, nil
		}
		return res, nil
	})

	server.Register("GetFundTop10Holdings", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		res, err := data.NewFundApi().GetFundTop10Holdings(server.ArgString(args, 0))
		if err != nil || res == nil {
			return []data.FundHoldingStock{}, nil
		}
		return res, nil
	})

	server.Register("GetFundRanking", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		res, err := data.NewFundApi().GetFundRanking(
			server.ArgString(args, 0), server.ArgString(args, 1),
			server.ArgString(args, 2), server.ArgString(args, 3),
			server.ArgInt(args, 4), server.ArgInt(args, 5),
		)
		if err != nil || res == nil {
			return &data.FundRankingResult{}, nil
		}
		return res, nil
	})

	server.Register("SearchFundCodes", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewFundApi().SearchFundCodes(server.ArgString(args, 0)), nil
	})

	server.Register("GetFollowedFundPaged", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewFundApi().GetFollowedFundPaged(server.ArgInt(args, 0), server.ArgInt(args, 1), server.ArgString(args, 2)), nil
	})
}
