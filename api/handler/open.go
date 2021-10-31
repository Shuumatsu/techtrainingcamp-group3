package handler

import (
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"

	"github.com/gin-gonic/gin"
)

func OpenHandler(c *gin.Context) {
	var req models.OpenReq
	c.Bind(&req)

	logger.Sugar.Debugw("OpenHandler",
		"envelope_id", req.EnvelopeId, "uid", req.Uid)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": 50,
		},
	})
}
