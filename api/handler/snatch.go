package handler

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/sql/sqlAPI"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"
	"techtrainingcamp-group3/tools"
	"time"
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
		user, err := sqlAPI.FindOrCreateUserByUID(dbmodels.User{Uid: dbmodels.UID(req.Uid)})
		if err != nil {
			c.JSON(200, gin.H{
				"code": models.DataBaseError,
				"msg":  models.DataBaseError.Message(),
			})
			return
		}
		envelopesId, err := sqlAPI.ParseEnvelopeList(user.EnvelopeList)
		if len(envelopesId) >= max_count {
			c.JSON(200, gin.H{
				"code": models.SnatchLimit,
				"msg":  models.SnatchLimit.Message(),
			})
			return
		}
		envelope := tools.REPool.Snatch()
		err = sqlAPI.AddEnvelopeToUserByUID(dbmodels.UID(req.Uid), dbmodels.Envelope{
			EnvelopeId: dbmodels.EID(envelope.Eid),
			Uid:        user.Uid,
			Opened:     false,
			Value:      uint64(envelope.Money),
			SnatchTime: time.Now().Unix(),
		})
		if err != nil {
			c.JSON(200, gin.H{
				"code": models.DataBaseError,
				"msg":  models.DataBaseError.Message(),
			})
			return
		}
		// TODO: redis
		c.JSON(200, gin.H{
			"code": models.Success,
			"msg":  models.Success.Message(),
			"data": gin.H{
				"envelope_id": envelope.Eid,
				"max_count":   max_count,
				"cur_count":   len(envelopesId) + 1,
			},
		})
	}
}
