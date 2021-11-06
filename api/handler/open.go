package handler

import (
	"errors"
	"net/http"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/sql/sqlAPI"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"

	"github.com/gin-gonic/gin"
)

func OpenHandler(c *gin.Context) {
	var req models.OpenReq
	err := c.Bind(&req)
	if err != nil {
		logger.Sugar.Errorw("OpenHandler parameter bind error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar.Debugw("OpenHandler",
		"envelope_id", req.EnvelopeId, "uid", req.Uid)

	// open envelope by envelope_id and user_id
	envelopeP, err := sqlAPI.OpenEnvelope(dbmodels.EID(req.EnvelopeId), dbmodels.UID(req.Uid))
	if errors.Is(err, sqlAPI.Error.NotFound) {
		c.JSON(200, gin.H{
			"code": models.NotFound,
			"msg":  models.NotFound.Message(),
			"data": gin.H{
				"value": 0,
			},
		})
		return
	}

	// check the owner
	if errors.Is(err, sqlAPI.Error.ErrorEnvelopeOwner) {
		c.JSON(200, gin.H{
			"code": models.ErrorEnvelopeOwner,
			"msg":  models.ErrorEnvelopeOwner.Message(),
			"data": gin.H{
				"value": 0,
			},
		})
		return
	}

	// The envelope has already been opened
	if errors.Is(err, sqlAPI.Error.EnvelopeAlreadyOpen) {
		c.JSON(200, gin.H{
			"code": models.EnvelopeAlreadyOpen,
			"msg":  models.EnvelopeAlreadyOpen.Message(),
			"data": gin.H{
				"value": 0,
			},
		})
		return
	}

	// Update envelope status and user amount success
	c.JSON(200, gin.H{
		"code": models.Success,
		"msg":  models.Success.Message(),
		"data": gin.H{
			"value": envelopeP.Value,
		},
	})

}
