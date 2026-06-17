package server

// Version / Commit / BuildKey 由 cmd/server 启动时注入（来自 -ldflags / 默认值，与桌面端构建一致）。
// BuildKey 用于赞助码 AES 解密等（桌面端用根 main 的 BuildKey）。
var (
	Version string
	Commit  string
	BuildKey string
)
