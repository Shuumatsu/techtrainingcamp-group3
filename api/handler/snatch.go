package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("snatched by %v", uid)

	envelop_id := 123
	max_count := 5
	cur_count := 3

	c.JSON(200, gin.H{
		"coide": 0,
		"msg":   "success",
		"data": gin.H{
			"envelop_id": envelop_id,
			"max_count":  max_count,
			"curr_count": cur_count,
		},
	})
}