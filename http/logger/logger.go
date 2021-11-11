package logger

import (
	"techtrainingcamp-group3/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"techtrainingcamp-group3/http/config"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func init() {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Env.LogLevel)); err != nil {
		panic(err)
	}

	Logger = logger.NewLogger(level)
	Sugar = Logger.Sugar()
}
