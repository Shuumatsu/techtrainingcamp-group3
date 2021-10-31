package handler

import (
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"

	"github.com/gin-gonic/gin"
)

func SnatchHandler(c *gin.Context) {
	var req models.SnatchReq
	c.Bind(&req)

	logger.Sugar.Debugw("SnatchHandler",
		"uid", req.Uid)

	envelope_id := 123
	max_count := 5
	cur_count := 3

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"envelope_id": envelope_id,
			"max_count":   max_count,
			"cur_count":   cur_count,
		},
	})
}
