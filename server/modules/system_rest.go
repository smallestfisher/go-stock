package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("GetSponsorInfo", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		level, active := data.EffectiveSponsorVipLevel()
		return map[string]any{
			"vipLevel":        level,
			"active":          active,
			"sponsorCode":     "",
			"vipStartTime":    "",
			"vipEndTime":      "",
			"vipAuthTime":     "",
			"freeEntitlement": true,
		}, nil
	})

	server.Register("CheckSponsorCode", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return map[string]any{
			"code": 1,
			"msg":  "全部本地功能已免费开放，无需赞助码。",
		}, nil
	})

	server.Register("CheckDeviceBinding", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return map[string]any{
			"bound":       true,
			"deviceCount": 0,
			"maxDevices":  0,
		}, nil
	})

}
