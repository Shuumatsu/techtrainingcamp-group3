package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/snatch", SnatchHandler)
	r.POST("/open", OpenHandler)
	r.POST("/get_wallet_list", WalletListHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("snatched by %v", uid)

	envelop_id := 123
	max_count := 5
	cur_count := 3

	c.JSON(200, gin.H{
		"code": 0,
		"msg":   "success",
		"data": gin.H{
			"envelop_id": envelop_id,
			"max_count":  max_count,
			"curr_count": cur_count,
		},
	})
}

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelop_id, _ := c.GetPostForm("envelop_id")

	log.Printf("envelop %v opened by %v", envelop_id, uid)

	value := 50

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": value,
		},
	})
}

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("query %v's wallet", uid)

	envelopes := []gin.H{
		{
			"envelop_id":  123,
			"value":       50,
			"opened":      true,
			"snatch_time": 1634551711,
		},
		{
			"envelop_id":  123,
			"value":       50,
			"opened":      false,
			"snatch_time": 1634551711,
		},
	}

	amount := 50

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"mount":        amount,
			"envelop_list": envelopes,
		},
	})
}
