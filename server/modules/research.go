package modules

import (
	"context"

	"github.com/duke-git/lancet/v2/convertor"
	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("GetTimezone", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return map[string]any{"offset": 8 * 60 * 60, "location": "Asia/Shanghai"}, nil
	})

	server.Register("LongTigerRank", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().LongTiger(server.ArgString(args, 0)), nil
	})

	server.Register("StockResearchReport", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().StockResearchReport(server.ArgString(args, 0), 7), nil
	})

	server.Register("StockNotice", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().StockNotice(server.ArgString(args, 0)), nil
	})

	server.Register("IndustryResearchReport", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().IndustryResearchReport(server.ArgString(args, 0), 7), nil
	})

	// EMDictCode 用 core.Cache 做板块字典缓存（桌面端用 a.cache）。
	server.Register("EMDictCode", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().EMDictCode(server.ArgString(args, 0), core.Cache), nil
	})

	server.Register("HotStock", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().XUEQIUHotStock(100, server.ArgString(args, 0)), nil
	})

	server.Register("HotEvent", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		size := server.ArgInt(args, 0)
		if size <= 0 {
			size = 10
		}
		return data.NewMarketNewsApi().HotEvent(size), nil
	})

	server.Register("HotTopic", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		size := server.ArgInt(args, 0)
		if size <= 0 {
			size = 10
		}
		return data.NewMarketNewsApi().HotTopic(size), nil
	})

	server.Register("InvestCalendarTimeLine", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().InvestCalendar(server.ArgString(args, 0)), nil
	})

	server.Register("ClsCalendar", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewMarketNewsApi().ClsCalendar(), nil
	})

	server.Register("GetUplimitHot", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMarketNewsApi().GetUplimitHot(server.ArgString(args, 0), server.ArgInt(args, 1)), nil
	})

	server.Register("SearchStock", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewSearchStockApi(server.ArgString(args, 0)).SearchStock(5000), nil
	})

	server.Register("GetStockRealTimePrice", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		code := server.ArgString(args, 0)
		stockDatas, err := data.NewStockDataApi().GetStockCodeRealTimeData(code)
		if err != nil || stockDatas == nil || len(*stockDatas) == 0 {
			return map[string]any{"code": -1, "message": "获取股票价格失败", "price": 0}, nil
		}
		stock := (*stockDatas)[0]
		price, _ := convertor.ToFloat(stock.Price)
		if price == 0 {
			price, _ = convertor.ToFloat(stock.A1P)
		}
		if price == 0 {
			price, _ = convertor.ToFloat(stock.B1P)
		}
		if price == 0 {
			price, _ = convertor.ToFloat(stock.PreClose)
		}
		return map[string]any{"code": 0, "message": "success", "price": price, "name": stock.Name}, nil
	})
}
