package modules

import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)

func init() {
	server.Register("AnalyzeSentiment", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		return data.AnalyzeSentiment(server.ArgString(args, 0)), nil
	})

	server.Register("AnalyzeSentimentWithFreqWeight", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		result, cleanFrequencies := data.NewsAnalyze(server.ArgString(args, 0), false)
		return map[string]any{"result": result, "frequencies": cleanFrequencies}, nil
	})
}
