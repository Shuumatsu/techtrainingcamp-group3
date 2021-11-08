package handler

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/rds/redisAPI"
	"techtrainingcamp-group3/db/sql/sqlAPI"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"
	"techtrainingcamp-group3/tools"
	"time"
)

func SnatchHandler(c *gin.Context) {
	var req models.SnatchReq
	err := c.BindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw("SnatchHandler parameter bind error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if rand.Float32() > config.SnatchProb {
		c.JSON(200, gin.H{
			"code": models.SnatchFailure,
			"msg":  models.SnatchFailure.Message(),
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
		envelope := tools.GetRandEnvelope(user.Uid)
		err = sqlAPI.AddEnvelopeToUserByUID(dbmodels.UID(req.Uid), envelope)
		if err != nil {
			c.JSON(200, gin.H{
				"code": models.DataBaseError,
				"msg":  models.DataBaseError.Message(),
			})
			return
		}
		// TODO: redis
		user.EnvelopeList += "," + envelope.EnvelopeId.String()
		err = redisAPI.SetUserByUID(user, 300*time.Second)
		if err != nil {
			logger.Sugar.Debugw("snatch", "redis set error", err, "user", user)
		}
		err = redisAPI.SetEnvelopeByEID(&envelope, 300*time.Second)
		if err != nil {
			logger.Sugar.Debugw("snatch", "redis set error", err, "envelope", envelope)
		}
		c.JSON(200, gin.H{
			"code": models.Success,
			"msg":  models.Success.Message(),
			"data": gin.H{
				"envelope_id": envelope.EnvelopeId,
				"max_count":   max_count,
				"cur_count":   len(envelopesId) + 1,
			},
		})
	}
}
