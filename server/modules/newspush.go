package modules

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"go-stock/backend/data"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	// NewsPush：按配置(仅红消息/含自选股名)把电报经事件总线推送给前端。复刻 app.go:967。
	server.Register("NewsPush", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		news := server.ArgJSON[[]models.Telegraph](args, 0)
		follows := data.NewStockDataApi().GetFollowList(0)
		stockNames := slice.Map(*follows, func(index int, item data.FollowedStock) string { return item.Name })
		cfg := data.GetSettingConfig()
		for _, telegraph := range news {
			if cfg.EnableOnlyPushRedNews {
				if telegraph.IsRed || strutil.ContainsAny(telegraph.Content, stockNames) {
					core.Events.Emit("newsPush", telegraph)
				}
			} else {
				core.Events.Emit("newsPush", telegraph)
			}
		}
		return nil, nil
	})
}
