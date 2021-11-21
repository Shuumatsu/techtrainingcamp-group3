package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"techtrainingcamp-group3/cmd/http/api/router"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/logger"
)

func main() {
	defer logger.Sugar.Sync()
	defer logger.Logger.Sync()

	gin.SetMode(config.Env.GinMode)
	r := router.Register()

	addr := fmt.Sprintf("%s:%s", config.Env.HttpHost, config.Env.HttpPort)
	logger.Sugar.Infow("start server", "addr", addr)
	r.Run(addr)
}
