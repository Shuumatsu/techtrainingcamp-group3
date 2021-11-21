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
	"time"
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

		//Test if user can get more envelopes
		if bloomfilter.RedisTestLimitUser(dbmodels.UID(req.Uid)) == true {
			c.JSON(200, gin.H{
				"code": models.SnatchLimit,
				"msg":  models.SnatchLimit.Message(),
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

		// if Uid is still in Processing means user may snatch too fast
		if redisAPI.SetNXUidInProcessing(dbmodels.UID(req.Uid),90*time.Second) == false{
			c.JSON(200, gin.H{
				"code": models.SnatchFast,
				"msg":  models.SnatchFast.Message(),
			})
			return
		}

		//Check if user can snatch more envelope
		envelopesId, err := sqlAPI.ParseEnvelopeList(user.EnvelopeList)
		if len(envelopesId) >= maxCount {
			//Add user id to user limiter filter, so that user cannot get more envelope
			bloomfilter.RedisAddLimitUser(dbmodels.UID(req.Uid))
			c.JSON(200, gin.H{
				"code": models.SnatchLimit,
				"msg":  models.SnatchLimit.Message(),
			})
			return
		}

		//Use redis lua script to get a random envelope
		envelope, _ := redisAPI.GetRandEnvelope(dbmodels.UID(req.Uid))
		if envelope.Value == 0 {
			c.JSON(200, gin.H{
				"code": models.NoEnvelopes,
				"msg":  models.NoEnvelopes.Message(),
			})
			return
		}

		//add envelope into user's EnvelopeList
		user.EnvelopeList += "," + envelope.EnvelopeId.String()
		user.UpdatedAt = time.Now()
		// put create the envelope in envelope table and append it to the user's envelope_list into kafka
		err = kfk.AddEnvelopeToUser(dbmodels.UID(req.Uid), envelope)
		if err != nil {
			//if fail to put user into
			logger.Sugar.Errorw("AddEnvelopeToUser kafka error","error",err)
			c.JSON(200, gin.H{
				"code": models.KafkaError,
				"msg":  models.KafkaError.Message(),
			})
			return
		}

		//Update bloom filter for envelope
		bloomfilter.RedisAddEnvelope(envelope.EnvelopeId)
		redisAPI.SetUserByUID(user,300*time.Second)

		curCount := len(envelopesId) + 1
		if curCount == maxCount{
			//Add user id to user limiter filter, so that user cannot get more envelope
			bloomfilter.RedisAddLimitUser(dbmodels.UID(req.Uid))
		}
		logger.Sugar.Debugw("snatch handler", "success", "user", user)
		c.JSON(200, gin.H{
			"code": models.Success,
			"msg":  models.Success.Message(),
			"data": gin.H{
				"envelope_id": envelope.EnvelopeId,
				"max_count":   maxCount,
				"cur_count":   curCount,
			},
		})
	}
}
