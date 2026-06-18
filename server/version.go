package server

// Version / Commit / BuildKey 由 cmd/server 启动时注入（来自 -ldflags / 默认值）。
// BuildKey 用于赞助码 AES 解密等。
var (
	Version string
	Commit  string
	BuildKey string
)
