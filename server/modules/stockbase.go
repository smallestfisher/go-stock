package modules

import (
	"context"
	"time"

	"go-stock/backend/data"
	"go-stock/backend/db"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	// CheckStockBaseInfo：从远端拉取并全量替换 A股/港股/美股 基础信息，完成后推送 loadingMsg=done。
	server.Register("CheckStockBaseInfo", func(_ context.Context, core *server.Core, _ server.Args) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				log.SugaredLogger.Errorf("CheckStockBaseInfo panic: %v", r)
			}
			go core.Events.Emit("loadingMsg", "done")
		}()

		stockBasics := &[]data.StockBasic{}
		data.SharedHTTPClient.R().SetHeader("user", "go-stock").SetResult(stockBasics).
			Get("http://8.134.249.145:18080/go-stock/stock_basic.json")
		db.Dao.Unscoped().Model(&data.StockBasic{}).Where("1=1").Delete(&data.StockBasic{})
		if err := db.Dao.CreateInBatches(stockBasics, 400).Error; err != nil {
			log.SugaredLogger.Errorf("保存StockBasic股票基础信息失败:%s", err.Error())
		}

		stockHKBasics := &[]models.StockInfoHK{}
		data.SharedHTTPClient.R().SetHeader("user", "go-stock").SetResult(stockHKBasics).
			Get("http://8.134.249.145:18080/go-stock/stock_base_info_hk.json")
		db.Dao.Unscoped().Model(&models.StockInfoHK{}).Where("1=1").Delete(&models.StockInfoHK{})
		if err := db.Dao.CreateInBatches(stockHKBasics, 400).Error; err != nil {
			log.SugaredLogger.Errorf("保存StockInfoHK股票基础信息失败:%s", err.Error())
		}

		stockUSBasics := &[]models.StockInfoUS{}
		data.SharedHTTPClient.R().SetHeader("user", "go-stock").SetResult(stockUSBasics).
			Get("http://8.134.249.145:18080/go-stock/stock_base_info_us.json")
		db.Dao.Unscoped().Model(&models.StockInfoUS{}).Where("1=1").Delete(&models.StockInfoUS{})
		if err := db.Dao.CreateInBatches(stockUSBasics, 400).Error; err != nil {
			log.SugaredLogger.Errorf("保存StockInfoUS股票基础信息失败:%s", err.Error())
		}
		return nil, nil
	})

	// FetchAndSaveMarketStatistic：交易时间内采集并保存当日市场统计。
	server.Register("FetchAndSaveMarketStatistic", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		if !isATradingSession(time.Now().In(shanghaiLoc)) {
			return nil, nil
		}
		if err := data.NewMarketStatisticApi().FetchAndSave(); err != nil {
			log.SugaredLogger.Errorf("获取市场统计数据失败: %v", err)
		}
		return nil, nil
	})
}
