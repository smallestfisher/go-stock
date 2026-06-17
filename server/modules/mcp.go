package modules

import (
	"context"

	"go-stock/backend/data"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	server.Register("CreateMCPServer", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		s := server.ArgJSON[models.MCPServer](args, 0)
		if err := data.NewMCPServerApi().Create(&s); err != nil {
			log.SugaredLogger.Errorf("创建MCP服务器失败: %v", err)
			return "创建失败: " + err.Error(), nil
		}
		return "创建成功", nil
	})

	server.Register("UpdateMCPServer", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		s := server.ArgJSON[models.MCPServer](args, 0)
		if err := data.NewMCPServerApi().Update(&s); err != nil {
			return "更新失败: " + err.Error(), nil
		}
		return "更新成功", nil
	})

	server.Register("DeleteMCPServer", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewMCPServerApi().Delete(uint(server.ArgInt64(args, 0))); err != nil {
			return "删除失败: " + err.Error(), nil
		}
		return "删除成功", nil
	})

	server.Register("GetMCPServerByID", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		s, err := data.NewMCPServerApi().GetByID(uint(server.ArgInt64(args, 0)))
		if err != nil {
			return nil, nil
		}
		return s, nil
	})

	server.Register("GetMCPServerList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		q := server.ArgJSON[models.MCPServerQuery](args, 0)
		return data.NewMCPServerApi().List(&q), nil
	})

	server.Register("EnableMCPServer", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		enable := server.ArgBool(args, 1)
		if err := data.NewMCPServerApi().EnableServer(uint(server.ArgInt64(args, 0)), enable); err != nil {
			return "操作失败: " + err.Error(), nil
		}
		if enable {
			return "已启用", nil
		}
		return "已禁用", nil
	})

	server.Register("TestMCPServer", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		result, err := data.NewMCPServerApi().TestConnection(uint(server.ArgInt64(args, 0)))
		if err != nil {
			return "测试失败: " + err.Error(), nil
		}
		return result, nil
	})

	server.Register("GetMCPToolsByServerID", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.NewMCPServerApi().GetToolsByServerID(uint(server.ArgInt64(args, 0))), nil
	})

	server.Register("GetAllMCPTools", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewMCPServerApi().GetAllTools(), nil
	})
}
