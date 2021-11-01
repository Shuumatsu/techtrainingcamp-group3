package main

import (
	"techtrainingcamp-group3/api/router"
	"techtrainingcamp-group3/tools"
)

func main() {
	r := router.Register()
	tools.PoolInit()
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
