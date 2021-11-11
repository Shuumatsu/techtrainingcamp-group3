package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level) *zap.Logger {
	stdoutWriteSyncer := zapcore.AddSync(os.Stdout)

	encoderPreset := zap.NewProductionEncoderConfig()
	encoderPreset.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderPreset)

	core := zapcore.NewCore(encoder, stdoutWriteSyncer, level)

	return zap.New(core)
}
