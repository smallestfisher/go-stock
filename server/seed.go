package server

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"go-stock/backend/data"
	"go-stock/backend/db"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
)

const (
	stockBasicPath = "build/stock_basic.json"
	stockHKPath    = "build/stock_base_info_hk.json"
	stockUSPath    = "build/stock_base_info_us.json"
)

// SeedStockBasic 在首次启动时把内置的 A股/港股/美股基础数据导入本地库（幂等：仅插入缺失项）。
// 把这些数据 //go:embed 进二进制（main.go 的 initStockData 系列）；
// 服务端改为从仓库 build/*.json 读取——部署时需随二进制带上这几个文件。
func SeedStockBasic() {
	seedAStock()
	seedHK()
	seedUS()
}

func seedAStock() {
	raw, err := os.ReadFile(stockBasicPath)
	if err != nil {
		log.SugaredLogger.Warnf("读取 %s 失败(跳过 A股 基础数据初始化)：%v", stockBasicPath, err)
		return
	}
	res := &data.TushareStockBasicResponse{}
	if err := json.Unmarshal(raw, res); err != nil {
		log.SugaredLogger.Errorf("解析 A股 基础数据失败：%v", err)
		return
	}
	log.SugaredLogger.Infof("init A股基础数据 %d 条", len(res.Data.Items))

	fields := strings.Split("ts_code,symbol,name,area,industry,cnspell,market,list_date,act_name,act_ent_type,fullname,exchange,list_status,curr_type,enname,delist_date,is_hs", ",")
	for _, item := range res.Data.Items {
		stockData := map[string]any{}
		for _, field := range fields {
			idx := slice.IndexOf(res.Data.Fields, field)
			if idx == -1 || idx >= len(item) {
				continue
			}
			stockData[field] = item[idx]
		}
		jsonData, _ := json.Marshal(stockData)
		stock := &data.StockBasic{}
		if err := json.Unmarshal(jsonData, stock); err != nil {
			continue
		}
		stock.ID = 0
		var count int64
		db.Dao.Model(&data.StockBasic{}).Where("ts_code = ?", stock.TsCode).Count(&count)
		if count > 0 {
			continue
		}
		db.Dao.Create(stock)
	}
}

func seedHK() {
	raw, err := os.ReadFile(stockHKPath)
	if err != nil {
		log.SugaredLogger.Warnf("读取 %s 失败(跳过 港股 基础数据初始化)：%v", stockHKPath, err)
		return
	}
	var v []models.StockInfoHK
	if err := json.Unmarshal(raw, &v); err != nil {
		log.SugaredLogger.Errorf("解析 港股 基础数据失败：%v", err)
		return
	}
	log.SugaredLogger.Infof("init 港股基础数据 %d 条", len(v))
	var total int64
	db.Dao.Model(&models.StockInfoHK{}).Count(&total)
	if total == int64(len(v)) {
		return
	}
	for _, item := range v {
		var count int64
		db.Dao.Model(&models.StockInfoHK{}).Where("code = ?", item.Code).Count(&count)
		if count > 0 {
			continue
		}
		db.Dao.Model(&models.StockInfoHK{}).Create(&item)
	}
}

func seedUS() {
	raw, err := os.ReadFile(stockUSPath)
	if err != nil {
		log.SugaredLogger.Warnf("读取 %s 失败(跳过 美股 基础数据初始化)：%v", stockUSPath, err)
		return
	}
	var v []models.StockInfoUS
	if err := json.Unmarshal(raw, &v); err != nil {
		log.SugaredLogger.Errorf("解析 美股 基础数据失败：%v", err)
		return
	}
	log.SugaredLogger.Infof("init 美股基础数据 %d 条", len(v))
	var total int64
	db.Dao.Model(&models.StockInfoUS{}).Count(&total)
	if total == int64(len(v)) {
		return
	}
	for _, item := range v {
		var count int64
		db.Dao.Model(&models.StockInfoUS{}).Where("code = ?", item.Code).Count(&count)
		if count > 0 {
			continue
		}
		db.Dao.Model(&models.StockInfoUS{}).Create(&item)
	}
}
