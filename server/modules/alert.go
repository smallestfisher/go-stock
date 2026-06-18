package modules

import (
	"context"
	"strings"
	"time"

	"go-stock/backend/data"
	log "go-stock/backend/logger"
	"go-stock/server"
)

func init() {
	// SendDingDingMessage：5 分钟内同一股票去重后发送钉钉消息（与一致）。
	server.Register("SendDingDingMessage", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		message := server.ArgString(args, 0)
		stockCode := server.ArgString(args, 1)
		if ttl, _ := core.Cache.TTL([]byte(stockCode)); ttl > 0 {
			return "", nil
		}
		if err := core.Cache.Set([]byte(stockCode), []byte("1"), 60*5); err != nil {
			log.SugaredLogger.Errorf("set cache error:%s", err.Error())
			return "", nil
		}
		return data.NewDingDingAPI().SendDingDingMessage(message), nil
	})

	// SendDingDingMessageByType：按报警类型(1 涨跌/2 股价/3 成本)去重发送，并经事件总线推送 newsPush。
	// 说明：相比，Web 版省略了本地系统通知()，newsPush 内容直接用传入 message。
	server.Register("SendDingDingMessageByType", func(_ context.Context, core *server.Core, args server.Args) (any, error) {
		message := server.ArgString(args, 0)
		stockCode := server.ArgString(args, 1)
		msgType := server.ArgInt(args, 2)

		// 交易时段闸门（基础版，无节假日日历）
		switch {
		case hasAnyPrefix(stockCode, "SZ", "SH", "sh", "sz"):
			if !isATradingSession(time.Now().In(shanghaiLoc)) {
				return "非A股交易时间", nil
			}
		case hasAnyPrefix(stockCode, "hk", "HK"):
			if !isHKTradingSession(time.Now().In(shanghaiLoc)) {
				return "非港股交易时间", nil
			}
		case hasAnyPrefix(stockCode, "us", "US", "gb_"):
			if !isUSTradingSession(time.Now()) {
				return "非美股交易时间", nil
			}
		}

		if ttl, _ := core.Cache.TTL([]byte(stockCode)); ttl > 0 {
			return "", nil
		}
		if err := core.Cache.Set([]byte(stockCode), []byte("1"), 60*5); err != nil {
			log.SugaredLogger.Errorf("set cache error:%s", err.Error())
			return "", nil
		}

		core.Events.Emit("newsPush", map[string]any{
			"time":    "📈 " + msgTypeName(msgType),
			"isRed":   true,
			"source":  "go-stock",
			"content": message,
		})

		return data.NewDingDingAPI().SendDingDingMessage(message), nil
	})
}

func hasAnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func msgTypeName(msgType int) string {
	switch msgType {
	case 1:
		return "涨跌报警"
	case 2:
		return "股价报警"
	case 3:
		return "成本价报警"
	default:
		return "消息通知"
	}
}
