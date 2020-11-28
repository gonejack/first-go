package main

import (
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap/zapcore"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction(
		zap.AddStacktrace(zapcore.WarnLevel),
		zap.AddCaller(),
	)
	spew.Dump(err)
	defer logger.Sync()

	url := "http://example.org/api"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
