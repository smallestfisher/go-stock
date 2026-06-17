package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	// 板块资金流
	server.Register("GetBKFundFlowList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewBKFundFlowApi().GetBKFundFlowList(server.ArgString(args, 0), server.ArgInt(args, 1)), nil
	})
	server.Register("GetBKFundFlowListByDate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewBKFundFlowApi().GetBKFundFlowListByDate(server.ArgString(args, 0), server.ArgString(args, 1)), nil
	})
	server.Register("GetBKFundFlowTopList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewBKFundFlowApi().GetBKFundFlowTopList(server.ArgInt(args, 0)), nil
	})
	server.Register("GetBKFundFlowTopListByDate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewBKFundFlowApi().GetBKFundFlowTopListByDate(server.ArgString(args, 0), server.ArgInt(args, 1)), nil
	})
	server.Register("GetAllBKCodes", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewBKFundFlowApi().GetAllBKCodes(), nil
	})

	// 概念资金流
	server.Register("GetConceptFundFlowList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewConceptFundFlowApi().GetConceptFundFlowList(server.ArgString(args, 0), server.ArgInt(args, 1)), nil
	})
	server.Register("GetConceptFundFlowListByDate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewConceptFundFlowApi().GetConceptFundFlowListByDate(server.ArgString(args, 0), server.ArgString(args, 1)), nil
	})
	server.Register("GetConceptFundFlowTopList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewConceptFundFlowApi().GetConceptFundFlowTopList(server.ArgInt(args, 0)), nil
	})
	server.Register("GetConceptFundFlowTopListByDate", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewConceptFundFlowApi().GetConceptFundFlowTopListByDate(server.ArgString(args, 0), server.ArgInt(args, 1)), nil
	})
	server.Register("GetAllConceptCodes", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewConceptFundFlowApi().GetAllConceptCodes(), nil
	})
}
