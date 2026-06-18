package modules

import (
	"context"
	"fmt"
	"strings"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("GetStockKLine", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().GetHK_KLineData(server.ArgString(args, 0), "day", server.ArgInt64(args, 2)), nil
	})

	server.Register("GetStockMinutePriceLineData", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		code := server.ArgString(args, 0)
		name := server.ArgString(args, 1)
		priceData, date := data.NewStockDataApi().GetStockMinutePriceData(code)
		return map[string]any{
			"priceData": priceData,
			"date":      date,
			"stockName": name,
			"stockCode": code,
		}, nil
	})

	server.Register("GetStockCommonKLine", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewStockDataApi().GetCommonKLineData(server.ArgString(args, 0), "day", server.ArgInt64(args, 2)), nil
	})

	server.Register("GetStockEastMoneyKLine", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return eastMoneyKLinePage(server.ArgString(args, 0), server.ArgString(args, 2), server.ArgInt(args, 3), ""), nil
	})

	server.Register("GetStockEastMoneyKLinePage", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return eastMoneyKLinePage(server.ArgString(args, 0), server.ArgString(args, 2), server.ArgInt(args, 3), server.ArgString(args, 4)), nil
	})

	server.Register("GetStockKLineWithFallback", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return kLineWithFallback(server.ArgString(args, 0), server.ArgString(args, 1), server.ArgString(args, 2), server.ArgInt(args, 3), ""), nil
	})

	server.Register("GetStockKLinePageWithFallback", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return kLineWithFallback(server.ArgString(args, 0), server.ArgString(args, 1), server.ArgString(args, 2), server.ArgInt(args, 3), server.ArgString(args, 4)), nil
	})

	server.Register("GetChipDistribution", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return chipDistribution(server.ArgString(args, 0), server.ArgInt(args, 1), server.ArgInt(args, 2), server.ArgString(args, 3))
	})
}

// eastMoneyKLinePage 归一化参数并调用东方财富分页 K 线接口。
func eastMoneyKLinePage(stockCode, klt string, limit int, end string) *[]data.KLineData {
	if limit <= 0 {
		limit = 500
	}
	if limit > 5000 {
		limit = 5000
	}
	klt = strings.TrimSpace(klt)
	if klt == "" {
		klt = "1"
	}
	api := data.NewEastMoneyKLineApi(data.GetSettingConfig())
	end = strings.TrimSpace(end)
	return api.GetKLineDataBefore(stockCode, klt, "", limit, end)
}

// kLineWithFallback 归一化参数并在东方财富失败时自动切换到新浪。
func kLineWithFallback(stockCode, stockName, klt string, limit int, end string) *data.KLineSourceResult {
	if limit <= 0 {
		limit = 500
	}
	if limit > 5000 {
		limit = 5000
	}
	klt = strings.TrimSpace(klt)
	if klt == "" {
		klt = "101"
	}
	end = strings.TrimSpace(end)
	return data.FetchKLineWithFallback(stockCode, stockName, klt, limit, end)
}

// chipDistribution 
func chipDistribution(stockCode string, days int, bins int, adjustFlag string) (*data.ChipDistributionResult, error) {
	stockCode = strings.TrimSpace(stockCode)
	if stockCode == "" {
		return nil, fmt.Errorf("stockCode 不能为空")
	}
	if days <= 0 {
		days = 120
	}
	if bins <= 0 {
		bins = 80
	}
	adjustFlag = strings.TrimSpace(strings.ToLower(adjustFlag))
	if adjustFlag != "" && adjustFlag != "qfq" && adjustFlag != "hfq" {
		adjustFlag = "qfq"
	}

	api := data.NewEastMoneyKLineApi(data.GetSettingConfig())
	if !api.ValidateStockCode(stockCode) {
		return nil, fmt.Errorf("股票代码无效：%s", stockCode)
	}

	var kLines *[]data.KLineData
	if adjustFlag != "" {
		kLines = api.GetKLineData(stockCode, "101", adjustFlag, days)
	} else {
		result := data.FetchKLineWithFallback(stockCode, "", "101", days, "")
		if result != nil && result.Data != nil {
			kLines = result.Data
		}
	}
	if kLines == nil || len(*kLines) == 0 {
		return nil, fmt.Errorf("未获取到K线数据")
	}
	calculator := data.NewChipDistributionCalculator()
	return calculator.Calculate(stockCode, *kLines, bins)
}
