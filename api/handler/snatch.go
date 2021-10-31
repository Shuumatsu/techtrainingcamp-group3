package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"techtrainingcamp-group3/models"
)

func SnatchHandler(c *gin.Context) {
	var req models.SnatchReq
	c.Bind(&req)
	log.Printf("snatched by %v", req.Uid)

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