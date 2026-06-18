package modules

import (
	"context"

	"go-stock/backend/agent"
	"go-stock/server"
)

func init() {
	server.Register("CalculateNextRunTime", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		t := agent.NewCronTaskApi().CalculateNextRunTime(server.ArgString(args, 0))
		return t.Format("2006-01-02 15:04:05"), nil
	})

	server.Register("CalculateNextRunTimes", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		times := agent.NewCronTaskApi().CalculateNextRunTimes(server.ArgString(args, 0), server.ArgInt(args, 1))
		result := make([]string, 0, len(times))
		for _, t := range times {
			result = append(result, t.Format("2006-01-02 15:04:05"))
		}
		return result, nil
	})
}
