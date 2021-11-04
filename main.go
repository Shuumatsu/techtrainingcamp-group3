package main

import (
	"io"
	"os"
	"techtrainingcamp-group3/api/router"
	"techtrainingcamp-group3/config"
	_ "techtrainingcamp-group3/db/mysql"
	_ "techtrainingcamp-group3/db/redis"
	"techtrainingcamp-group3/logger"

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

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
