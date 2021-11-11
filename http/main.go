package main

import (
	"techtrainingcamp-group3/http/api/router"
	"techtrainingcamp-group3/http/config"
	_ "techtrainingcamp-group3/http/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Env.GinMode)
	r := router.Register()
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
