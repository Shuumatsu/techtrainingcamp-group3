package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

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
