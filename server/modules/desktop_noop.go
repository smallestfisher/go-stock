package modules

import (
	"context"

	"go-stock/server"
)

// 桌面专属能力：Web 下无意义（系统托盘、本地文件保存），注册为 no-op，避免前端调用时 404。
// SaveImage/SaveWordFile 在 Web 下应改为前端触发浏览器下载，此处先置 no-op。
func init() {
	noOp := func(name string) {
		server.Register(name, func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
			return nil, nil
		})
	}
	noOp("RestartAsAdmin")
	noOp("QuitApp")
	noOp("HideToTray")
	noOp("ShowFromTray")
	noOp("SaveImage")
	noOp("SaveWordFile")
}
