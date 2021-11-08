package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/rds/redisAPI"
	"techtrainingcamp-group3/db/sql/sqlAPI"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"
	"time"
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
	var req models.OpenReq
	err := c.BindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw("OpenHandler parameter bind error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar.Debugw("OpenHandler",
		"envelope_id", req.EnvelopeId, "uid", req.Uid)

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

	// To Do: add mq to update data in sql

	// Update envelope status and user amount in sql
	userP, err := sqlAPI.UpdateEnvelopeOpen(envelopeP)

	// check for envelope status again
	if errors.Is(err, dbmodels.Error.EnvelopeAlreadyOpen) {
		ConstructErrorReply(c, models.EnvelopeAlreadyOpen)
		return
	}

	// If error happened, return false
	if err != nil {
		ConstructErrorReply(c, models.DataBaseError)
		return
	}

	// If data success flush user to redis
	if err := redisAPI.SetUserByUID(userP, 300*time.Second); err != nil {
		logger.Sugar.Errorw("Redis set user error", "uid", userP.Uid)
	}

	// Update envelope status and user amount success
	logger.Sugar.Debugw("Open Handler success", "envelopeId", envelopeP.EnvelopeId)
	c.JSON(200, gin.H{
		"code": models.Success,
		"msg":  models.Success.Message(),
		"data": gin.H{
			"value": envelopeP.Value,
		},
	})
}
