package server

// Version / Commit / BuildKey 由 cmd/server 启动时注入（来自 -ldflags / 默认值）。
// BuildKey 用于机器标识派生等兼容初始化。
var (
	Version  string
	Commit   string
	BuildKey string
)
