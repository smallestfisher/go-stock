package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	// 通达信(TDX)协议封装：集合竞价、F10 公司资料、财务、除权除息、板块归属等。
	server.Register("GetTdxCallAuction", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		start := uint32(server.ArgInt(args, 1))
		count := uint32(server.ArgInt(args, 2))
		if count == 0 {
			count = 500
		}
		return data.NewTdxKLineApi().GetCallAuction(server.ArgString(args, 0), start, count), nil
	})

	server.Register("GetTdxCompanyInfo", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewTdxKLineApi().GetF10Data(server.ArgString(args, 0)), nil
	})

	server.Register("GetTdxFinanceInfo", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewTdxKLineApi().GetFinanceInfo(server.ArgString(args, 0)), nil
	})

	server.Register("GetTdxXDXRInfo", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewTdxKLineApi().GetXDXRInfo(server.ArgString(args, 0)), nil
	})

	server.Register("GetTdxCompanyCategoryList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewTdxKLineApi().GetF10CategoryList(server.ArgString(args, 0)), nil
	})

	server.Register("GetTdxCompanyCategoryContent", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewTdxKLineApi().GetF10CategoryContent(server.ArgString(args, 0), server.ArgString(args, 1)), nil
	})

	server.Register("GetTdxSymbolBelongBoard", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewTdxKLineApi().GetMACSymbolBelongBoard(server.ArgString(args, 0)), nil
	})
}
