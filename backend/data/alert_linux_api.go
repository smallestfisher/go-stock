//go:build linux
// +build linux

package data

import (
	"go-stock/backend/logger"
	"os/exec"
)

// AlertWindowsApi @Author spark
// @Date 2025/1/8 9:40
// @Desc
// -----------------------------------------------------------------------------------
type AlertWindowsApi struct {
	AppID string
	// 窗口标题
	Title string
	// 窗口内容
	Content string
	// 窗口图标
	Icon string
}

func NewAlertWindowsApi(AppID string, Title string, Content string, Icon string) *AlertWindowsApi {
	return &AlertWindowsApi{
		AppID:   AppID,
		Title:   Title,
		Content: Content,
		Icon:    Icon,
	}
}

func (a AlertWindowsApi) SendNotification() bool {
	if GetSettingConfig().LocalPushEnable == false {
		logger.SugaredLogger.Error("本地推送未开启")
		return false
	}

	// Use notify-send for Linux
	cmd := exec.Command("notify-send", a.Title, a.Content)
	err := cmd.Run()
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to send Linux notification: %v", err)
		return false
	}
	return true
}
