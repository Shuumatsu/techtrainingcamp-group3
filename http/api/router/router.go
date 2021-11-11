package router

import (
	"techtrainingcamp-group3/http/api/handler"
	"techtrainingcamp-group3/http/api/middleware"
	"techtrainingcamp-group3/http/logger"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger.Logger, time.RFC3339, true))
	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger.Logger, true))
	r.POST("/snatch", handler.SnatchHandler)
	r.POST("/open", handler.OpenHandler)
	r.POST("/get_wallet_list", handler.WalletListHandler)
	return r
}
