package modules

import (
	"context"

	"go-stock/backend/data"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	server.Register("CreateSkill", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		s := server.ArgJSON[models.Skill](args, 0)
		if err := data.NewSkillApi().Create(&s); err != nil {
			log.SugaredLogger.Errorf("创建技能失败: %v", err)
			return "创建失败: " + err.Error(), nil
		}
		return "创建成功", nil
	})

	server.Register("UpdateSkill", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		s := server.ArgJSON[models.Skill](args, 0)
		if err := data.NewSkillApi().Update(&s); err != nil {
			return "更新失败: " + err.Error(), nil
		}
		return "更新成功", nil
	})

	server.Register("DeleteSkill", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := data.NewSkillApi().Delete(uint(server.ArgInt64(args, 0))); err != nil {
			return "删除失败: " + err.Error(), nil
		}
		return "删除成功", nil
	})

	server.Register("GetSkillByID", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		s, err := data.NewSkillApi().GetByID(uint(server.ArgInt64(args, 0)))
		if err != nil {
			return nil, nil
		}
		return s, nil
	})

	server.Register("GetSkillList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		q := server.ArgJSON[models.SkillQuery](args, 0)
		return data.NewSkillApi().List(&q), nil
	})

	server.Register("EnableSkill", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		enable := server.ArgBool(args, 1)
		if err := data.NewSkillApi().EnableSkill(uint(server.ArgInt64(args, 0)), enable); err != nil {
			return "操作失败: " + err.Error(), nil
		}
		if enable {
			return "已启用", nil
		}
		return "已禁用", nil
	})

	server.Register("GetAllSkills", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return data.NewSkillApi().GetAll(), nil
	})
}
