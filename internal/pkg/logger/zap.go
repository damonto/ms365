package logger

import (
	"go.uber.org/zap"
)

// Sugar is the zap sugared logger
var Sugar *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()

	Sugar = logger.Sugar()
}
