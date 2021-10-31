package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelope_id, _ := c.GetPostForm("envelope_id")

	log.Printf("envelope %v opened by %v", envelope_id, uid)

	value := 50

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": value,
		},
	})
}