package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"techtrainingcamp-group3/models"
)

func WalletListHandler(c *gin.Context) {
	var req models.WalletListReq
	c.Bind(&req)
	log.Printf("query %v's wallet", req.Uid)

	envelopes := []gin.H{
		{
			"envelope_id":  123,
			"value":       50,
			"opened":      true,
			"snatch_time": 1634551711,
		},
		{
			"envelope_id":  123,
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
			"amount":        amount,
			"envelope_list": envelopes,
		},
	})
}
