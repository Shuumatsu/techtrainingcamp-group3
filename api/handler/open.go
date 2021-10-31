package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

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