package handler

import (
	"math/rand"
	"net/http"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db/mongo"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"
	"techtrainingcamp-group3/tools"
	"time"
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
		max_count := config.MaxSnatchAmount
		user, err := mongo.SetDefaultUserByUID(req.Uid)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 3,
				"msg":  "database error",
			})
			return
		}
		if user.Wallet.Size() >= max_count {
			c.JSON(200, gin.H{
				"code": 2,
				"msg":  "too many envelopes",
			})
			return
		}
		envelope := tools.REPool.Snatch()
		err = mongo.AddEnvelopeToUserByUID(req.Uid, models.Envelope{
			models.EID(envelope.Eid), false, uint64(envelope.Money), time.Now().Unix(),
		})
		if err != nil {
			c.JSON(200, gin.H{
				"code": 3,
				"msg":  "database error",
			})
			return
		}
		// TODO: redis
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"envelope_id": envelope.Eid,
				"max_count":   max_count,
				"cur_count":   user.Wallet.Size() + 1,
			},
		})
	}
}
