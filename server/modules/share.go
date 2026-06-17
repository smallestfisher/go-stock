package modules

import (
	"context"
	"strings"
	"time"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	// ShareText：把文本分享到社区。
	server.Register("ShareText", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		text := strings.TrimSpace(server.ArgString(args, 0))
		title := strings.TrimSpace(server.ArgString(args, 1))
		if text == "" {
			return "内容为空", nil
		}
		if title == "" {
			title = "AI助手"
		}
		analysisTime := time.Now().Format("2006/01/02")
		resp, err := data.SharedHTTPClient.R().SetHeader("ua-x", "go-stock").SetFormData(map[string]string{
			"text":         text,
			"stockCode":    title,
			"stockName":    title,
			"analysisTime": analysisTime,
		}).Post("http://go-stock.sparkmemory.top:16688/upload")
		if err != nil {
			return err.Error(), nil
		}
		return resp.String(), nil
	})

	// ShareAnalysis：分享该股票最近一次 AI 分析结果。
	server.Register("ShareAnalysis", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		stockCode := server.ArgString(args, 0)
		stockName := server.ArgString(args, 1)
		res := data.NewDeepSeekOpenAi(context.Background(), 0).GetAIResponseResult(stockCode)
		if res == nil || len(res.Content) <= 100 {
			return "分析结果异常", nil
		}
		analysisTime := res.CreatedAt.Format("2006/01/02")
		resp, err := data.SharedHTTPClient.R().SetHeader("ua-x", "go-stock").SetFormData(map[string]string{
			"text":         res.Content,
			"stockCode":    stockCode,
			"stockName":    stockName,
			"analysisTime": analysisTime,
		}).Post("http://go-stock.sparkmemory.top:16688/upload")
		if err != nil {
			return err.Error(), nil
		}
		return resp.String(), nil
	})

	// SaveAsMarkdown：桌面端弹保存框写文件；Web 版直接返回 markdown 内容，
	// 由前端桩/组件触发浏览器下载（前端需把返回字符串作为 Blob 下载）。
	server.Register("SaveAsMarkdown", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		stockCode := server.ArgString(args, 0)
		res := data.NewDeepSeekOpenAi(context.Background(), 0).GetAIResponseResult(stockCode)
		if res == nil || len(res.Content) <= 100 {
			return "分析结果异常,无法保存。", nil
		}
		return res.Content, nil
	})
}
