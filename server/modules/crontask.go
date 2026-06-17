package modules

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-stock/backend/agent"
	"go-stock/backend/data"
	log "go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/server"
)

func init() {
	server.Register("CreateCronTask", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		task := server.ArgJSON[models.CronTask](args, 0)
		if err := agent.NewCronTaskApi().Create(&task); err != nil {
			return fmt.Sprintf("创建失败：%v", err), nil
		}
		if err := scheduleCronTask(core, &task); err != nil {
			return "任务创建成功,但定时失败", nil
		}
		return "创建成功", nil
	})

	server.Register("UpdateCronTask", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		task := server.ArgJSON[models.CronTask](args, 0)
		if err := agent.NewCronTaskApi().Update(&task); err != nil {
			return fmt.Sprintf("更新失败：%v", err), nil
		}
		if eid, ok := core.GetCronEntry(cronKey(task.ID, task.Name)); ok {
			core.Cron.Remove(eid)
		}
		if err := scheduleCronTask(core, &task); err != nil {
			return fmt.Sprintf("更新失败：%v", err), nil
		}
		return "更新成功", nil
	})

	server.Register("DeleteCronTask", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		id := uint(server.ArgInt64(args, 0))
		_ = agent.NewCronTaskApi().Delete(id)
		if t, err := agent.NewCronTaskApi().GetByID(id); err == nil {
			if eid, ok := core.GetCronEntry(cronKey(id, t.Name)); ok {
				core.Cron.Remove(eid)
			}
		}
		return "删除成功", nil
	})

	server.Register("GetCronTaskByID", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		task, err := agent.NewCronTaskApi().GetByID(uint(server.ArgInt64(args, 0)))
		if err != nil {
			return nil, nil
		}
		return task, nil
	})

	server.Register("GetCronTaskList", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		query := server.ArgJSON[models.CronTaskQuery](args, 0)
		return agent.NewCronTaskApi().List(&query), nil
	})

	server.Register("EnableCronTask", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		id := uint(server.ArgInt64(args, 0))
		enable := server.ArgBool(args, 1)
		_ = agent.NewCronTaskApi().EnableTask(id, enable)
		t, err := agent.NewCronTaskApi().GetByID(id)
		if err == nil {
			if eid, ok := core.GetCronEntry(cronKey(id, t.Name)); ok {
				core.Cron.Remove(eid)
			}
			if enable {
				if err := scheduleCronTask(core, t); err != nil {
					return "操作成功,但定时失败", nil
				}
			}
		}
		return "操作成功", nil
	})

	server.Register("ExecuteCronTaskNow", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		id := uint(server.ArgInt64(args, 0))
		task, err := agent.NewCronTaskApi().GetByID(id)
		if err != nil {
			return fmt.Sprintf("任务不存在：%v", err), nil
		}
		go func() {
			if err := agent.NewCronTaskApi().ExecuteTask(context.Background(), task); err != nil {
				log.SugaredLogger.Errorf("执行任务失败：%v %s", err, task.Name)
			}
		}()
		return "任务执行中", nil
	})

	server.Register("GetCronTaskTypes", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return agent.NewCronTaskApi().GetTaskTypes(), nil
	})

	server.Register("ValidateCronExpr", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		if err := agent.NewCronTaskApi().ValidateCronExpr(server.ArgString(args, 0)); err != nil {
			return fmt.Sprintf("无效表达式：%v", err), nil
		}
		return "有效表达式", nil
	})

	server.Register("SearchCronTasks", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return agent.NewCronTaskApi().SearchTasks(server.ArgString(args, 0)), nil
	})

	// SetStockAICron：为单只股票设置定时 AI 分析。
	server.Register("SetStockAICron", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		cronText := server.ArgString(args, 0)
		stockCode := server.ArgString(args, 1)
		data.NewStockDataApi().SetStockAICron(cronText, stockCode)

		sc := stockCode
		if strings.HasPrefix(strings.ToLower(sc), "gb_") {
			sc = strings.ToUpper(sc)
			sc = strings.Replace(sc, "GB_", "us", 1)
		}
		if eid, ok := core.GetCronEntry(sc); ok {
			core.Cron.Remove(eid)
		}
		follow := data.NewStockDataApi().GetFollowedStockByStockCode(sc)
		entryID, err := core.Cron.AddFunc(cronText, aiAnalysisJob(core, follow))
		if err != nil {
			return nil, nil
		}
		core.SetCronEntry(sc, entryID)
		return nil, nil
	})
}

// cronKey 定时任务在 core 中的注册键（与桌面端 convertor.ToString(id)+"_"+name 一致）。
func cronKey(id uint, name string) string {
	return strconv.FormatUint(uint64(id), 10) + "_" + name
}

// scheduleCronTask 把一个通用 CronTask 注册到 core.Cron，执行体为 agent.ExecuteTask。
func scheduleCronTask(core *server.Core, task *models.CronTask) error {
	taskCopy := *task
	entryID, err := core.Cron.AddFunc(taskCopy.CronExpr, func() {
		if err := agent.NewCronTaskApi().ExecuteTask(context.Background(), &taskCopy); err != nil {
			log.SugaredLogger.Errorf("执行任务失败：%v %s", err, taskCopy.Name)
		}
	})
	if err != nil {
		return err
	}
	core.SetCronEntry(cronKey(task.ID, task.Name), entryID)
	return nil
}

// aiAnalysisJob 构造单只股票定时 AI 分析的执行体（复刻 app.go AddCronTask）。
// 与桌面端的区别：用事件总线 hub 推送 warnMsg（桌面端用 runtime.EventsEmit）。
func aiAnalysisJob(core *server.Core, follow data.FollowedStock) func() {
	return func() {
		core.Events.Emit("warnMsg", "开始自动分析"+follow.Name+"_"+follow.StockCode)
		ai := data.NewDeepSeekOpenAi(context.Background(), follow.AiConfigId)
		thinking := data.GetSettingConfig().GetAIConfigThinking(follow.AiConfigId)
		msgs := ai.NewChatStream(follow.Name, follow.StockCode, "", nil, core.AiTools, thinking)
		var res strings.Builder
		chatId := ""
		question := ""
		for msg := range msgs {
			if v, ok := msg["extraContent"].(string); ok && v != "" {
				res.WriteString(v + "\n")
			}
			if v, ok := msg["content"].(string); ok && v != "" {
				res.WriteString(v)
			}
			if v, ok := msg["chatId"].(string); ok {
				chatId = v
			}
			if v, ok := msg["question"].(string); ok {
				question = v
			}
		}
		data.NewDeepSeekOpenAi(context.Background(), follow.AiConfigId).SaveAIResponseResult(follow.StockCode, follow.Name, res.String(), chatId, question)
		core.Events.Emit("warnMsg", "AI分析完成："+follow.Name+"_"+follow.StockCode)
	}
}

// InitCronTasks 在服务启动时重建所有已启用的通用定时任务（含自动创建"异动数据保存"任务）。
// 由 cmd/server 在 NewCore 之后调用。对应桌面端 app.go 的 InitCronTasks。
func InitCronTasks(core *server.Core) {
	cronApi := agent.NewCronTaskApi()
	if !cronApi.ExistsByTaskType("stock_change_save") {
		task := &models.CronTask{
			Name:        "异动数据保存",
			CronExpr:    "0 */1 * * * *",
			TaskType:    "stock_change_save",
			Enable:      true,
			Status:      "active",
			Description: "每分钟自动保存A股异动数据，交易时间外自动跳过",
		}
		if err := cronApi.Create(task); err != nil {
			log.SugaredLogger.Errorf("自动创建异动数据保存任务失败：%v", err)
		} else {
			log.SugaredLogger.Info("已自动创建异动数据保存定时任务")
		}
	}
	tasks := cronApi.GetAll()
	for _, t := range tasks {
		taskCopy := t
		entryID, err := core.Cron.AddFunc(taskCopy.CronExpr, func() {
			if err := agent.NewCronTaskApi().ExecuteTask(context.Background(), &taskCopy); err != nil {
				log.SugaredLogger.Errorf("启动任务失败：%v %s", err, taskCopy.Name)
			}
		})
		if err != nil {
			log.SugaredLogger.Errorf("自动创建定时任务失败：%v %s", err, taskCopy.Name)
			continue
		}
		core.SetCronEntry(cronKey(taskCopy.ID, taskCopy.Name), entryID)
	}
}
