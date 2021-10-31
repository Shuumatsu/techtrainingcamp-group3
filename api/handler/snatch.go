package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

func SnatchHandler(c *gin.Context) {
	uid, ok := c.GetPostForm("uid")
	if !ok {
		log.Println("get form error")
	}
	log.Printf("snatched by %v", uid)

	envelope_id := 123
	max_count := 5
	cur_count := 3

	c.JSON(200, gin.H{
		"code": 0,
		"msg":   "success",
		"data": gin.H {
			"envelope_id": envelope_id,
			"max_count":  max_count,
			"cur_count": cur_count,
		},
	})
}