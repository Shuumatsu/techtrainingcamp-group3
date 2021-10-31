package logger

import (
	"os"
	"techtrainingcamp-group3/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func init() {
	file, err := os.Create(config.Conf.ZapLogFile)
	if err != nil {
		panic(err)
	}
	stdoutWriteSyncer := zapcore.AddSync(file)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Env.LogLevel)); err != nil {
		panic(err)
	}
	core := zapcore.NewCore(encoder, stdoutWriteSyncer, level)

	logger := zap.New(core)
	if err != nil {
		panic(err)
	}

	Logger = logger
	Sugar = logger.Sugar()
}
