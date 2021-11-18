package logger

import (
	"os"

	"techtrainingcamp-group3/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level) *zap.Logger {
	stdoutWriteSyncer := zapcore.AddSync(os.Stdout)

	encoderPreset := zap.NewProductionEncoderConfig()
	encoderPreset.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderPreset)
	// encoder := zapcore.NewJSONEncoder(encoderPreset)

	core := zapcore.NewCore(encoder, stdoutWriteSyncer, level)

	return zap.New(core)
}

func SetLevel(level zapcore.Level) {
	Logger.Sync()
	Logger = NewLogger(level)
	Sugar = Logger.Sugar()
}

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func init() {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Env.LogLevel)); err != nil {
		panic(err)
	}

	Logger = NewLogger(level)
	Sugar = Logger.Sugar()
}
