package main

import (
	"context"
	"fmt"
	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/util"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/strutil"
)

// startup 在 Web 模式下的启动逻辑
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logger.SugaredLogger.Info("Web App Bridge: startup")
	// 初始化定时任务
	a.InitCronTasks()
}

// beforeClose 在 Web 模式下不需要处理窗口关闭
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// OnSecondInstanceLaunch Web 模式下不处理单实例锁定
func OnSecondInstanceLaunch(data any) {
}

// getScreenResolution Web 模式下返回固定分辨率
func getScreenResolution() (int, int, int, int, error) {
	return 1920, 1080, 1024, 768, nil
}

// MonitorStockPrices Web 模式下不需要 GUI 相关的状态更新
func MonitorStockPrices(a *App) {
	dest := &[]data.FollowedStock{}
	db.Dao.Model(&data.FollowedStock{}).Find(dest)
	total := float64(0)

	stockInfos := GetStockInfos(*dest...)
	for _, stockInfo := range *stockInfos {
		if strutil.HasPrefixAny(stockInfo.Code, []string{"SZ", "SH", "sh", "sz"}) && (!isTradingTime(time.Now())) {
			continue
		}
		if strutil.HasPrefixAny(stockInfo.Code, []string{"hk", "HK"}) && (!IsHKTradingTime(time.Now())) {
			continue
		}
		if strutil.HasPrefixAny(stockInfo.Code, []string{"us", "US", "gb_"}) && (!IsUSTradingTime(time.Now())) {
			continue
		}

		total += stockInfo.ProfitAmountToday
		price, _ := convertor.ToFloat(stockInfo.Price)

		if stockInfo.PrePrice != price {
			go util.Emit(a.ctx, "stock_price", stockInfo)
		}
	}

	go util.Emit(a.ctx, "realtime_profit", fmt.Sprintf("  %.2f", total))
}

// onReady Web 模式下的准备就绪逻辑
func onReady(a *App) {
}

// getFrameless Web 模式下不适用
func getFrameless() bool {
	return false
}


