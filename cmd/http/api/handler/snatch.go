package handler

import (
	"math/rand"
	"net/http"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/db/bloomfilter"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/rds/redisAPI"
	"techtrainingcamp-group3/pkg/db/sql/sqlAPI"
	"techtrainingcamp-group3/pkg/logger"
	"techtrainingcamp-group3/pkg/models"
	"techtrainingcamp-group3/pkg/tools"
	"time"

	"github.com/gin-gonic/gin"
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
		if bloomfilter.TestUser(dbmodels.UID(req.Uid)) == false {
			c.JSON(200, gin.H{
				"code": models.NotFound,
				"msg":  models.NotFound.Message(),
			})
			return
		}
		max_count := config.MaxSnatchAmount
		user, err := sqlAPI.FindUserByUID(dbmodels.UID(req.Uid))
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
		// TODO: bloom filter
		bloomfilter.AddEnvelope(envelope.EnvelopeId)
		logger.Sugar.Debugw("snarch handler", "success", "user", user)
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
