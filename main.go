package main

import (
	"techtrainingcamp-group3/api/router"
)

func main() {
	r := router.Register()

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
