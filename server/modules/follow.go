package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/server"
)

func init() {
	// 关注列表（按分组）。groupId=0 表示默认/全部。
	server.Register("GetFollowList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().GetFollowList(server.ArgInt(args, 0)), nil
	})

	// 股票基础信息搜索（代码/名称/拼音）。
	server.Register("GetStockList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().GetStockList(server.ArgString(args, 0)), nil
	})

	server.Register("GetGroupList", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewStockGroupApi(db.Dao).GetGroupList(), nil
	})

	server.Register("GetGroupStockList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockGroupApi(db.Dao).GetGroupStockByGroupId(server.ArgInt(args, 0)), nil
	})
}
