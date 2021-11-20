package handler

import (
	"errors"
	"net/http"
	"techtrainingcamp-group3/pkg/db/bloomfilter"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/kfk"
	"techtrainingcamp-group3/pkg/db/rds/redisAPI"
	"techtrainingcamp-group3/pkg/db/sql/sqlAPI"
	"techtrainingcamp-group3/pkg/logger"
	"techtrainingcamp-group3/pkg/models"
	"time"

	"github.com/gin-gonic/gin"
)

func ConstructErrorReply(c *gin.Context, e models.ErrorCode) {
	c.JSON(200, gin.H{
		"code": e,
		"msg":  e.Message(),
		"data": gin.H{
			"value": 0,
		},
	})
}

func checkOpen(p *dbmodels.Envelope) error {
	if p.Opened == true {
		return dbmodels.Error.EnvelopeAlreadyOpen
	}
	return nil
}

func OpenHandler(c *gin.Context) {
	//Check the request parameter
	var req models.OpenReq
	err := c.BindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw("OpenHandler parameter bind error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar.Debugw("OpenHandler",
		"envelope_id", req.EnvelopeId, "uid", req.Uid)

	//Use Redis Bloom filter to check whether uid and eid are possibly available
	if bloomfilter.RedisTestUser(dbmodels.UID(req.Uid)) == false {
		ConstructErrorReply(c, models.NotFound)
		logger.Sugar.Debugw("openHandler: user not found in bloom filter")
		return
	}
	if bloomfilter.RedisTestEnvelope(dbmodels.EID(req.EnvelopeId)) == false {
		ConstructErrorReply(c, models.NotFound)
		logger.Sugar.Debugw("openHandler: envelope not found in bloom filter")
		return
	}
	var envelopeP *dbmodels.Envelope = nil

	// First find envelope by redis
	envelopeP, err = redisAPI.FindEnvelopeByEIDUID(dbmodels.EID(req.EnvelopeId), dbmodels.UID(req.Uid))
	// If cache miss search in sql
	if err != nil {
		envelopeP, err = sqlAPI.FindEnvelopeByUidEid(dbmodels.EID(req.EnvelopeId), dbmodels.UID(req.Uid))
		if errors.Is(err, sqlAPI.Error.NotFound) {
			ConstructErrorReply(c, models.NotFound)
			return
		}
	}

	// check if the owner is right
	if errors.Is(err, dbmodels.Error.ErrorEnvelopeOwner) {
		ConstructErrorReply(c, models.ErrorEnvelopeOwner)
		return
	}

	// check if there is unknown error
	if err != nil || envelopeP == nil {
		ConstructErrorReply(c, models.NotDefined)
		return
	}

	// Check if the envelope has already been opened
	err = checkOpen(envelopeP)
	if errors.Is(err, dbmodels.Error.EnvelopeAlreadyOpen) {
		ConstructErrorReply(c, models.EnvelopeAlreadyOpen)
		return
	}

	// Update envelope status in redis to prevent open twice
	envelopeP.Opened = true
	if err := redisAPI.SetEnvelopeByEID(envelopeP, 300*time.Second); err != nil {
		logger.Sugar.Errorw("Redis set envelop opened error", "envelope_id", req.EnvelopeId, "uid", req.Uid)
	}

	// Put update envelope status and user amount in sql execution to kafka
	err = kfk.OpenEnvelope(dbmodels.UID(req.Uid),envelopeP)

	// If can not put message into kafka, move envelope from redis
	if err != nil {
		redisAPI.DelEnvelopeByEID(envelopeP.EnvelopeId)
		ConstructErrorReply(c, models.KafkaError)
		return
	}

	// Open API success
	logger.Sugar.Debugw("Open Handler success", "envelopeId", envelopeP.EnvelopeId)
	c.JSON(200, gin.H{
		"code": models.Success,
		"msg":  models.Success.Message(),
		"data": gin.H{
			"value": envelopeP.Value,
		},
	})
}
