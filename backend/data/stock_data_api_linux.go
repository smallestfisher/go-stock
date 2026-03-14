//go:build linux
// +build linux

package data

import (
	"os/exec"
)

// CheckBrowser 在 Linux 系统上检查浏览器是否存在
func CheckBrowser() (string, bool) {
	// 依次检查 common 浏览器
	browsers := []string{"google-chrome", "google-chrome-stable", "chromium-browser", "chromium", "microsoft-edge-stable", "microsoft-edge", "firefox"}
	for _, b := range browsers {
		path, err := exec.LookPath(b)
		if err == nil {
			return path, true
		}
	}
	return "", false
}

func CheckChrome() (string, bool) {
	path, err := exec.LookPath("google-chrome")
	if err == nil {
		return path, true
	}
	return "", false
}
