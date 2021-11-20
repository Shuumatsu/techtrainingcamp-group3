package handler

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/db/bloomfilter"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/kfk"
	"techtrainingcamp-group3/pkg/db/rds/redisAPI"
	"techtrainingcamp-group3/pkg/db/sql/sqlAPI"
	"techtrainingcamp-group3/pkg/logger"
	"techtrainingcamp-group3/pkg/models"
	"techtrainingcamp-group3/pkg/tools"
)

func SnatchHandler(c *gin.Context) {
	//Check the request parameter
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
		//maxCount is max envelope number that user can snatch
		maxCount := config.MaxSnatchAmount

		//Test if there is a user in bloom filter
		if bloomfilter.RedisTestUser(dbmodels.UID(req.Uid)) == false {
			c.JSON(200, gin.H{
				"code": models.NotFound,
				"msg":  models.NotFound.Message(),
			})
			return
		}

		//Find user information in redis First
		user,err := redisAPI.FindUserByUID(dbmodels.UID(req.Uid))
		if err != nil {
			user, err = sqlAPI.FindUserByUID(dbmodels.UID(req.Uid))
			if err != nil {
				c.JSON(200, gin.H{
					"code": models.DataBaseError,
					"msg":  models.DataBaseError.Message(),
				})
				return
			}
		}

		//Todo add logic for check that user cannot snatch envelope twice in x seconds
		//Todo or add logic for check one user can not snatch envelopes at same time

		//Check if user can snatch more envelope
		envelopesId, err := sqlAPI.ParseEnvelopeList(user.EnvelopeList)
		if len(envelopesId) >= maxCount {
			c.JSON(200, gin.H{
				"code": models.SnatchLimit,
				"msg":  models.SnatchLimit.Message(),
			})
			return
		}
		//Generate envelope
		envelope := tools.GetRandEnvelope(user.Uid)

		// put create the envelope in envelope table and append it to the user's envelope_list into kafka
		err = kfk.AddEnvelopeToUser(dbmodels.UID(req.Uid), envelope)
		if err != nil {
			logger.Sugar.Errorw("AddEnvelopeToUser kafka error","error",err)
			c.JSON(200, gin.H{
				"code": models.KafkaError,
				"msg":  models.KafkaError.Message(),
			})
			return
		}

		//Update bloom filter for envelope
		bloomfilter.RedisAddEnvelope(envelope.EnvelopeId)

		user.EnvelopeList += "," + envelope.EnvelopeId.String()
		logger.Sugar.Debugw("snatch handler", "success", "user", user)
		c.JSON(200, gin.H{
			"code": models.Success,
			"msg":  models.Success.Message(),
			"data": gin.H{
				"envelope_id": envelope.EnvelopeId,
				"max_count":   maxCount,
				"cur_count":   len(envelopesId) + 1,
			},
		})
	}
}
