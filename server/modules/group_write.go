package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/server"
)

func init() {
	server.Register("AddGroup", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		group := server.ArgJSON[data.Group](args, 0)
		if data.NewStockGroupApi(db.Dao).AddGroup(group) {
			return "添加成功", nil
		}
		return "添加失败", nil
	})

	server.Register("UpdateGroupSort", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockGroupApi(db.Dao).UpdateGroupSort(server.ArgInt(args, 0), server.ArgInt(args, 1)), nil
	})

	server.Register("InitializeGroupSort", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewStockGroupApi(db.Dao).InitializeGroupSort(), nil
	})

	server.Register("AddStockGroup", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if data.NewStockGroupApi(db.Dao).AddStockGroup(server.ArgInt(args, 0), server.ArgString(args, 1)) {
			return "添加成功", nil
		}
		return "添加失败", nil
	})

	server.Register("RemoveStockGroup", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if data.NewStockGroupApi(db.Dao).RemoveStockGroup(server.ArgString(args, 0), server.ArgString(args, 1), server.ArgInt(args, 2)) {
			return "移除成功", nil
		}
		return "移除失败", nil
	})

	server.Register("RemoveGroup", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if data.NewStockGroupApi(db.Dao).RemoveGroup(server.ArgInt(args, 0)) {
			return "移除成功", nil
		}
		return "移除失败", nil
	})
}
