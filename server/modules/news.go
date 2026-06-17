package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("GetTelegraphList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().GetTelegraphList(server.ArgString(args, 0)), nil
	})

	// ReFleshTelegraphList：后台刷新各资讯源，并立即返回当前电报列表。
	server.Register("ReFleshTelegraphList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		go data.NewMarketNewsApi().TelegraphList(30)
		go data.NewMarketNewsApi().GetSinaNews(30)
		go data.NewMarketNewsApi().TradingViewNews()
		return data.NewMarketNewsApi().GetTelegraphList(server.ArgString(args, 0)), nil
	})

	server.Register("GlobalStockIndexes", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewMarketNewsApi().GlobalStockIndexes(30), nil
	})

	server.Register("GlobalStockIndexesReadable", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewMarketNewsApi().GlobalStockIndexesReadable(30), nil
	})
}
