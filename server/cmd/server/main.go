// go-stock 纯服务端入口：放弃 Wails 桌面形态，以 HTTP 服务承载全部功能，供浏览器远程访问。
package main

import (
	"os"

	"go-stock/backend/data"
	"go-stock/backend/db"
	log "go-stock/backend/logger"
	"go-stock/backend/machineid"
	"go-stock/server"
	"go-stock/server/modules" // 触发各功能模块的 init() 注册 RPC handler，并提供 InitCronTasks
)

// 由 -ldflags 注入（与桌面端构建一致）。
var (
	Version       string
	VersionCommit string
	BuildKey      string
)

const defaultBuildKey = "cc1e0d684e32f176c56ff1fcf384dcd9"

func main() {
	ensureDir("data")
	ensureDir("logs")

	if BuildKey == "" {
		BuildKey = defaultBuildKey
	}
	machineid.Init(BuildKey)
	data.SponsorDecryptKeyHex = BuildKey

	db.Init("")
	data.InitAnalyzeSentiment()
	go server.RunMigrate()

	log.SugaredLogger.Infof("go-stock server 启动: version=%s commit=%s", Version, VersionCommit)

	server.Version = Version
	server.Commit = VersionCommit
	server.BuildKey = BuildKey

	core := server.NewCore()
	// 启动时重建所有已启用的定时任务（异动保存、自定义 cron 等）。
	go modules.InitCronTasks(core)

	addr := os.Getenv("GO_STOCK_ADDR")
	if addr == "" {
		addr = ":18888"
	}
	if err := server.Start(core, addr); err != nil {
		log.SugaredLogger.Fatalf("server 运行失败: %v", err)
	}
}

func ensureDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, 0o755)
	}
}
