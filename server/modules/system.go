package modules

import (
	"context"
	"os"

	"go-stock/backend/data"
	"go-stock/backend/machineid"
	"go-stock/server"
)

func init() {
	server.Register("GetConfig", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.GetSettingConfig(), nil
	})

	server.Register("GetEffectiveSponsorVip", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		level, active := data.EffectiveSponsorVipLevel()
		return map[string]any{"vipLevel": level, "active": active}, nil
	})

	server.Register("GetMachineId", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return machineid.GetMachineId(), nil
	})

	// GetVersionInfo：会附带图标/收款码 base64，Web 模式暂返回版本信息（图标等 Phase 3 补）。
	server.Register("GetVersionInfo", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return map[string]any{
			"Version":           server.Version,
			"VersionCommit":     server.Commit,
			"Content":           server.Commit,
			"Icon":              "",
			"Alipay":            "",
			"Wxpay":             "",
			"Wxgzh":             "",
			"OfficialStatement": "",
		}, nil
	})

	// GetUserManual：使用手册随仓库 docs/ 提供，运行时从磁盘读取。
	server.Register("GetUserManual", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		b, err := os.ReadFile("docs/go-stock使用手册.md")
		if err != nil {
			return "", nil
		}
		return string(b), nil
	})

	// ExportConfig：弹保存对话框写文件；Web 模式直接返回配置 JSON 文本，由前端触发下载。
	server.Register("ExportConfig", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewSettingsApi().Export(), nil
	})

}
