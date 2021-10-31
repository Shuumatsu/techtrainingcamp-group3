package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"techtrainingcamp-group3/models"
)

func OpenHandler(c *gin.Context) {
	var req models.OpenReq
	c.Bind(&req)

	log.Printf("envelope %v opened by %v", req.EnvelopeId, req.Uid)

	value := 50

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": value,
		},
	})
}