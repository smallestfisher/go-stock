package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	// UpdateConfig：落库即可。还会重排 MonitorStockPrices 定时任务，
	// 该 cron 调度在 Phase 2 接入 core.Cron 后补齐。
	server.Register("UpdateConfig", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		cfg := server.ArgJSON[data.SettingConfig](args, 0)
		return data.UpdateConfig(&cfg), nil
	})
}
