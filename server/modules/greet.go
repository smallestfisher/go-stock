package modules

import (
	"context"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/mathutil"
	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/server"
)

func init() {
	// Greet：取某只股票的实时行情(含成本/盈亏计算)，复刻 app.go Greet + getStockInfo + addStockFollowData。
	server.Register("Greet", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		stockCode := server.ArgString(args, 0)
		follow := &data.FollowedStock{StockCode: stockCode}
		db.Dao.Model(follow).Where("stock_code = ?", stockCode).
			Preload("Groups").Preload("Groups.GroupInfo").First(follow)
		return getStockInfo(*follow), nil
	})
}

// getStockInfo 取实时行情并附加自选股的成本/盈亏信息（复刻 app.go:1633）。
func getStockInfo(follow data.FollowedStock) *data.StockInfo {
	stockDatas, err := data.NewStockDataApi().GetStockCodeRealTimeData(follow.StockCode)
	if err != nil || stockDatas == nil || len(*stockDatas) == 0 {
		return &data.StockInfo{}
	}
	stockData := (*stockDatas)[0]
	addStockFollowData(follow, &stockData)
	return &stockData
}

func addStockFollowData(follow data.FollowedStock, stockData *data.StockInfo) {
	stockData.PrePrice = follow.Price
	stockData.Sort = follow.Sort
	stockData.CostPrice = follow.CostPrice
	stockData.CostVolume = follow.Volume
	stockData.AlarmChangePercent = follow.AlarmChangePercent
	stockData.AlarmPrice = follow.AlarmPrice
	stockData.Groups = follow.Groups

	price, _ := convertor.ToFloat(stockData.Price)
	if price == 0 {
		price, _ = convertor.ToFloat(stockData.A1P)
	}
	if price == 0 {
		price, _ = convertor.ToFloat(stockData.B1P)
	}
	preClosePrice, _ := convertor.ToFloat(stockData.PreClose)
	if price == 0 {
		price = preClosePrice
	}
	highPrice, _ := convertor.ToFloat(stockData.High)
	if highPrice == 0 {
		highPrice, _ = convertor.ToFloat(stockData.Open)
	}
	lowPrice, _ := convertor.ToFloat(stockData.Low)
	if lowPrice == 0 {
		lowPrice, _ = convertor.ToFloat(stockData.Open)
	}

	if price > 0 && preClosePrice > 0 {
		stockData.ChangePrice = mathutil.RoundToFloat(price-preClosePrice, 2)
		stockData.ChangePercent = mathutil.RoundToFloat(mathutil.Div(price-preClosePrice, preClosePrice)*100, 3)
	}
	if highPrice > 0 && preClosePrice > 0 {
		stockData.HighRate = mathutil.RoundToFloat(mathutil.Div(highPrice-preClosePrice, preClosePrice)*100, 3)
	}
	if lowPrice > 0 && preClosePrice > 0 {
		stockData.LowRate = mathutil.RoundToFloat(mathutil.Div(lowPrice-preClosePrice, preClosePrice)*100, 3)
	}
	if follow.CostPrice > 0 && follow.Volume > 0 {
		if price > 0 {
			stockData.Profit = mathutil.RoundToFloat(mathutil.Div(price-follow.CostPrice, follow.CostPrice)*100, 3)
			stockData.ProfitAmount = mathutil.RoundToFloat((price-follow.CostPrice)*float64(follow.Volume), 2)
			stockData.ProfitAmountToday = mathutil.RoundToFloat((price-preClosePrice)*float64(follow.Volume), 2)
		} else {
			stockData.Profit = mathutil.RoundToFloat(mathutil.Div(preClosePrice-follow.CostPrice, follow.CostPrice)*100, 3)
			stockData.ProfitAmount = mathutil.RoundToFloat((preClosePrice-follow.CostPrice)*float64(follow.Volume), 2)
			stockData.ProfitAmountToday = 0
		}
	}

	if follow.Price != price && price > 0 {
		go db.Dao.Model(follow).Where("stock_code = ?", follow.StockCode).Updates(map[string]interface{}{
			"price": price,
		})
	}
}
