package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"techtrainingcamp-group3/pkg/db/bloomfilter"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/rds/redisAPI"
	"techtrainingcamp-group3/pkg/db/sql/sqlAPI"
	"techtrainingcamp-group3/pkg/logger"
	"techtrainingcamp-group3/pkg/models"
	"time"
)

func WalletListHandler(c *gin.Context) {
	// Check the request parameter
	var req models.WalletListReq
	err := c.BindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw("SnatchHandler parameter bind error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar.Debugw("WalletListHandler",
		"uid", req.Uid)

	// Use bloom filter to test if user exist
	if bloomfilter.RedisTestUser(dbmodels.UID(req.Uid)) == false {
		c.JSON(200, models.WalletListResp{
			Code: models.NotFound,
			Msg:  models.NotFound.Message(),
			Data: models.WalletListData{},
		})
		logger.Sugar.Debugw("WalletListHandler: user not found in bloomfilter")
		return
	}

	// find user in redis first
	user, err := redisAPI.FindUserByUID(dbmodels.UID(req.Uid))
	if err != nil {
		logger.Sugar.Debugw("redis not found", "uid", req.Uid)
		// redis cache miss, find user in sql
		user, err = sqlAPI.FindUserByUID(dbmodels.UID(req.Uid))
	}

	// if mysql error
	if err != nil {
		switch err {
		case sqlAPI.Error.NotFound:
			// if mysql not found
			c.JSON(200, models.WalletListResp{
				Code: models.NotFound,
				Msg:  models.NotFound.Message(),
				Data: models.WalletListData{},
			})
			logger.Sugar.Debugw("WalletListHandler",
				"not found in mysql", err)
		default:
			c.JSON(200, models.WalletListResp{
				Code: models.DataBaseError,
				Msg:  models.DataBaseError.Message(),
				Data: models.WalletListData{},
			})
			logger.Sugar.Errorw("WalletListHandler",
				"mysql database error", err)
		}
		return
	}
	// find envelopes which belong to the user
	// parse envelope_list
	envelopesID, err := sqlAPI.ParseEnvelopeList(user.EnvelopeList)
	if err != nil {
		c.JSON(200, models.WalletListResp{
			Code: models.ParseError,
			Msg:  models.ParseError.Message(),
			Data: models.WalletListData{},
		})
		logger.Sugar.Errorw("WalletListHandler",
			"can't parse envelope_list", err)
		return
	}

	// Try to find envelopes using redis
	envelopes := make([]dbmodels.Envelope, len(envelopesID))
	redisNotFoundIdx := make([]int, 0)
	for i := 0; i < len(envelopesID); i++ {
		envelope, err := redisAPI.FindEnvelopeByEID(envelopesID[i])
		if err != nil {
			redisNotFoundIdx = append(redisNotFoundIdx, i)
			continue
		}
		envelopes[i] = *envelope
	}

	// If some envelopes failed to find in redis, find them in sql
	for i := 0; i < len(redisNotFoundIdx); i++ {
		idx := redisNotFoundIdx[i]
		envelope, err := sqlAPI.FindEnvelopeByUidEid(envelopesID[idx], dbmodels.UID(req.Uid))
		if err != nil {
			c.JSON(200, models.WalletListResp{
				Code: models.DataBaseError,
				Msg:  models.DataBaseError.Message(),
				Data: models.WalletListData{},
			})
			logger.Sugar.Debugw("WalletListHandler",
				"not found in mysql", err)
			return
		}
		envelopes[idx] = *envelope
	}
	//  flush results to redis
	for _, envelope := range envelopes {
		if envelope.EnvelopeId == 0 {
			continue
		}
		if err = redisAPI.SetEnvelopeByEID(&envelope, 300*time.Second); err != nil {
			logger.Sugar.Errorw("Redis set envelop opened error",
				"envelope_id", envelope.EnvelopeId, "uid", envelope.Uid)
		}
	}

	// change envelopes to resq model
	envelopesResq := make([]models.Envelope, len(envelopes))
	for i := 0; i < len(envelopes); i++ {
		envelopesResq[i] = envelopes[i].ToResqModel()
	}

	// success
	logger.Sugar.Debugw("success", "amount", user.Amount, "envelopelist", envelopesResq)
	c.JSON(200, models.WalletListResp{
		Code: models.Success,
		Msg:  models.Success.Message(),
		Data: models.WalletListData{
			Amount:       user.Amount,
			EnvelopeList: envelopesResq,
		},
	})
}
