package handler

import (
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
	if err != nil {
		switch err {
		case sqlAPI.Error.NotFound:
			c.JSON(200, gin.H{
				"code": models.NotFound,
				"msg":  models.NotFound.Message(),
				"data": gin.H{
					"value": 0,
				},
			})
		case sqlAPI.Error.ErrorEnvelopeOwner:
			// check the owner
			c.JSON(200, gin.H{
				"code": models.ErrorEnvelopeOwner,
				"msg":  models.ErrorEnvelopeOwner.Message(),
				"data": gin.H{
					"value": 0,
				},
			})
		case sqlAPI.Error.EnvelopeAlreadyOpen:
			// The envelope has already been opened
			c.JSON(200, gin.H{
				"code": models.EnvelopeAlreadyOpen,
				"msg":  models.EnvelopeAlreadyOpen.Message(),
				"data": gin.H{
					"value": 0,
				},
			})
		default:
			c.JSON(200, gin.H{
				"code": models.DataBaseError,
				"msg":  models.DataBaseError.Message(),
				"data": gin.H{
					"value": 0,
				},
			})
		}
		logger.Sugar.Debugw("openHandler", "error", err)
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
