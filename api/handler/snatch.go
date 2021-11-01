package handler

import (
	"log"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/tools"
	"github.com/gin-gonic/gin"
)

func SnatchHandler(c *gin.Context) {
	uid, ok := c.GetPostForm("uid")
	if !ok {
		log.Println("get form error")
	}
	log.Printf("snatched by %v", uid)
	// TODO: Query the amount of snatching for this user
	envelope := tools.REPool.Snatch()
	max_count := config.MaxSnatchAmount
	cur_count := 3

	c.JSON(200, gin.H{
		"code": 0,
		"msg":   "success",
		"data": gin.H {
			"envelope_id": envelope.Eid,
			"max_count":  max_count,
			"cur_count": cur_count,
		},
	})
}