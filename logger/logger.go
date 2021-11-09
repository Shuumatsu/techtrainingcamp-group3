package logger

import (
	"os"
	"path/filepath"
	"techtrainingcamp-group3/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func init() {

	if _, err := os.Stat(filepath.Dir(config.Conf.ZapLogFile)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(config.Conf.ZapLogFile), 0666)
		if err != nil {
			panic(err)
		}
	}

	var stdoutWriteSyncer zapcore.WriteSyncer
	if config.Env.GinMode == "release" {
		file, err := os.Create(config.Conf.ZapLogFile)
		if err != nil {
			panic(err)
		}
		stdoutWriteSyncer = zapcore.AddSync(file)
	} else {
		stdoutWriteSyncer = zapcore.AddSync(os.Stdout)
	}

	encoderPreset := zap.NewProductionEncoderConfig()
	encoderPreset.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderPreset)

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Env.LogLevel)); err != nil {
		panic(err)
	}
	core := zapcore.NewCore(encoder, stdoutWriteSyncer, level)

	logger := zap.New(core)

	Logger = logger
	Sugar = logger.Sugar()
}
