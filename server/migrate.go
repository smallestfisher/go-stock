package server

import (
	"go-stock/backend/data"
	"go-stock/backend/db"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
)

// RunMigrate 建表。内容与 main.go 的 AutoMigrate 保持一致（幂等）。
// 由 cmd/server 在启动 HTTP 前同步调用，确保数据库结构和内置数据先就绪。
func RunMigrate() {
	db.Dao.AutoMigrate(&data.StockInfo{})
	db.Dao.AutoMigrate(&data.StockBasic{})
	db.Dao.AutoMigrate(&data.FollowedStock{})
	db.Dao.AutoMigrate(&data.IndexBasic{})
	db.Dao.AutoMigrate(&data.Settings{})
	db.Dao.AutoMigrate(&models.AIResponseResult{})
	db.Dao.AutoMigrate(&models.StockInfoHK{})
	db.Dao.AutoMigrate(&models.StockInfoUS{})
	db.Dao.AutoMigrate(&data.FollowedFund{})
	db.Dao.AutoMigrate(&data.FollowedStock{})
	db.Dao.AutoMigrate(&data.FundBasic{})
	db.Dao.AutoMigrate(&models.PromptTemplate{})
	db.Dao.AutoMigrate(&data.Group{})
	db.Dao.AutoMigrate(&data.GroupStock{})
	db.Dao.AutoMigrate(&models.Tags{})
	db.Dao.AutoMigrate(&models.Telegraph{})
	db.Dao.AutoMigrate(&models.TelegraphTags{})
	db.Dao.AutoMigrate(&models.LongTigerRankData{})
	db.Dao.AutoMigrate(&data.AIConfig{})
	db.Dao.AutoMigrate(&models.BKDict{})
	db.Dao.AutoMigrate(&models.WordAnalyze{})
	db.Dao.AutoMigrate(&models.SentimentResultAnalyze{})
	db.Dao.AutoMigrate(&models.AiRecommendStocks{})
	db.Dao.AutoMigrate(&models.AllStockInfo{})
	db.Dao.AutoMigrate(&models.CronTask{})
	db.Dao.AutoMigrate(&models.AiAssistantSession{})
	db.Dao.AutoMigrate(&models.GlobalStockIndex{})
	db.Dao.AutoMigrate(&data.TradingRecord{})
	db.Dao.AutoMigrate(&models.MCPServer{})
	db.Dao.AutoMigrate(&models.MCPServerTool{})
	db.Dao.AutoMigrate(&models.Skill{})
	db.Dao.AutoMigrate(&models.CustomStrategy{})
	db.Dao.AutoMigrate(&models.BKFundFlow{})
	db.Dao.AutoMigrate(&models.ConceptFundFlow{})

	initGlobalStockIndexCacheTask()

	// 建表完成后，导入内置 A股/港股/美股基础数据（幂等）。
	SeedStockBasic()
}

// initGlobalStockIndexCacheTask 确保存在“全球指数缓存”定时任务记录（与一致）。
func initGlobalStockIndexCacheTask() {
	var count int64
	db.Dao.Model(&models.CronTask{}).Where("task_type = ?", "global_stock_index_cache").Count(&count)
	if count > 0 {
		return
	}
	task := &models.CronTask{
		Name:        "全球指数缓存",
		CronExpr:    "0 0/5 * * * *",
		TaskType:    "global_stock_index_cache",
		Target:      "",
		Params:      `{"crawlTimeOut": 30}`,
		Enable:      true,
		Status:      "active",
		Description: "自动缓存全球股票指数数据",
	}
	if err := db.Dao.Create(task).Error; err != nil {
		log.SugaredLogger.Errorf("创建 global_stock_index_cache 定时任务失败：%v", err)
	} else {
		log.SugaredLogger.Info("创建 global_stock_index_cache 定时任务成功")
	}
}
