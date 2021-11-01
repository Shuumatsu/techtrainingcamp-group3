package main

import (
	"io"
	"os"
	"techtrainingcamp-group3/api/router"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
	"time"

	_ "techtrainingcamp-group3/db"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func main() {
	defer logger.Sugar.Sync()
	defer logger.Logger.Sync()

	gin.SetMode(config.Env.GinMode)
	if config.Env.GinMode == gin.ReleaseMode {
		ginLogFile, err := os.Create(config.Conf.GinLogFile)
		if err != nil {
			panic(err)
		}
		gin.DefaultWriter = io.MultiWriter(ginLogFile)
	}

	r := router.Register()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger.Logger, time.RFC3339, true))
	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger.Logger, true))

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
