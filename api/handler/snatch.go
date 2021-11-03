package handler

import (
	"math/rand"
	"net/http"
	"strings"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"
	"github.com/gin-gonic/gin"
)

func SnatchHandler(c *gin.Context) {
	var req models.SnatchReq
	err := c.Bind(&req)
	if err != nil {
		logger.Sugar.Errorw("SnatchHandler parameter bind error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if rand.Float32() > config.SnatchProb {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "fail",
		})
	} else {
		var user db.User
		max_count := config.MaxSnatchAmount
		cur_count := 0
		if err := db.DB.Table(db.User{}.TableName()).First(
			&user, req.Uid).Error; err == nil {
			cur_count = strings.Count(user.EnvelopeList, ",") + 1
		}
		if cur_count >= max_count {
			c.JSON(200, gin.H{
				"code": 2,
				"msg":  "too many envelopes",
			})
		} else {
			if eid, err := db.UpdateUsersEnvelope(user, cur_count); err != nil {
				c.JSON(200, gin.H{
					"code": 3,
					"msg":  "database error",
				})
			} else {
				c.JSON(200, gin.H{
					"code": 0,
					"msg":  "success",
					"data": gin.H{
						"envelope_id": eid,
						"max_count":   max_count,
						"cur_count":   cur_count + 1,
					},
				})
			}
		}
	}
}

