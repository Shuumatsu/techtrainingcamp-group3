package router

import (
	"github.com/gin-gonic/gin"
	"techtrainingcamp-group3/api/handler"
	"techtrainingcamp-group3/api/middlewares"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.POST("/snatch", handler.SnatchHandler)
	r.POST("/open", handler.OpenHandler)
	r.POST("/get_wallet_list", handler.WalletListHandler)
	return r
}
