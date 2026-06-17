package modules

import (
	"context"
	"strings"
	"time"

	"go-stock/server"
)

// shanghaiLoc 东八区。交易时段判断统一基于上海时间。
var shanghaiLoc = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("CST", 8*3600)
	}
	return loc
}()

func init() {
	// 注：这里的交易时间判断为基础版（工作日 + 时段），不含法定节假日日历。
	// 节假日精确版本（复用 data 的节假日数据）在 Phase 3 补齐。
	server.Register("IsTradingTime", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return isATradingSession(time.Now().In(shanghaiLoc)), nil
	})
	server.Register("IsHKTradingTime", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return isHKTradingSession(time.Now().In(shanghaiLoc)), nil
	})
	server.Register("IsUSTradingTime", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return isUSTradingSession(time.Now()), nil
	})
	server.Register("IsTradingDay", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return isTradingDayBasic(server.ArgString(args, 0)), nil
	})
	server.Register("GetLatestTradingDay", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return latestTradingDay(), nil
	})
}

func isWeekday(d time.Time) bool {
	wd := d.Weekday()
	return wd != time.Saturday && wd != time.Sunday
}

func isATradingSession(t time.Time) bool {
	if !isWeekday(t) {
		return false
	}
	h, m, _ := t.Clock()
	mins := h*60 + m
	// 9:30-11:30，13:00-15:00
	return (mins >= 570 && mins <= 690) || (mins >= 780 && mins <= 900)
}

func isHKTradingSession(t time.Time) bool {
	if !isWeekday(t) {
		return false
	}
	h, m, _ := t.Clock()
	mins := h*60 + m
	// 9:30-12:00，13:00-16:00
	return (mins >= 570 && mins <= 720) || (mins >= 780 && mins <= 960)
}

func isUSTradingSession(t time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.FixedZone("EST", -5*3600)
	}
	t = t.In(loc)
	if !isWeekday(t) {
		return false
	}
	h, m, _ := t.Clock()
	mins := h*60 + m
	// 9:30-16:00 ET
	return mins >= 570 && mins <= 960
}

func isTradingDayBasic(date string) bool {
	date = strings.TrimSpace(date)
	if date == "" {
		return false
	}
	t, err := time.ParseInLocation("2006-01-02", date, shanghaiLoc)
	if err != nil {
		return false
	}
	return isWeekday(t)
}

func latestTradingDay() string {
	now := time.Now().In(shanghaiLoc)
	if isWeekday(now) {
		h, m, _ := now.Clock()
		if h < 15 || (h == 15 && m == 0) {
			return now.AddDate(0, 0, -1).Format("2006-01-02")
		}
		return now.Format("2006-01-02")
	}
	for i := 1; i <= 7; i++ {
		d := now.AddDate(0, 0, -i)
		if isWeekday(d) {
			return d.Format("2006-01-02")
		}
	}
	return now.Format("2006-01-02")
}
