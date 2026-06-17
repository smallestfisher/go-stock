// Package modules 按功能模块拆解 server 的 RPC handler 实现。
//
// 每个文件在 init() 中调用 server.Register，把桌面端 App 的对应方法注册到 RPC 桥接。
// 这些 handler 直接复用 backend 层(data/agent/db)的现有函数，与 app.go 里的调用完全同款，
// 因此行为与桌面端一致。cmd/server 通过空白导入本包来触发注册。
package modules
