package handler

import (
	"net/http"
	"techtrainingcamp-group3/db"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"

	"github.com/gin-gonic/gin"
)

type openCode int

const (
	Success        openCode = iota //success
	Opened                         // Already opened
	NoExist                        // The specified user_id and envelope_id doesn't exist
	NoUser                         // The specified user doesn't exist in user table
	UnknownFailure                 //UnknownFailure
)

func (c openCode) String() string {
	switch c {
	case Success:
		return "success"
	case Opened:
		return "opened"
	case NoExist:
		return "noExist"
	case NoUser:
		return "noUser"
	case UnknownFailure:
		return "unknownFailure"
	}
	return "N/A"
}

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

	//get envelope by envelope_id and user_id
	envelopeP := db.GetEnvelope(req.EnvelopeId, uint64(req.Uid))
	if envelopeP == nil {
		c.JSON(200, gin.H{
			"code": NoExist,
			"msg":  NoExist.String(),
			"data": gin.H{
				"value": 0,
			},
		})
		return
	}

	//get user by user_id
	userP := db.GetUser(uint64(req.Uid))
	if userP == nil {
		c.JSON(200, gin.H{
			"code": NoUser,
			"msg":  NoUser.String(),
			"data": gin.H{
				"value": 0,
			},
		})
		return
	}

	//The envelope has already been opened
	if envelopeP.Open_stat == true {
		c.JSON(200, gin.H{
			"code": Opened,
			"msg":  Opened.String(),
			"data": gin.H{
				"value": 0,
			},
		})
		return
	}

	//Update envelope status and user amount
	if err = db.UpdateEnvelopeOpen(envelopeP, userP); err != nil {
		logger.Sugar.Errorw("OpenHandler update error")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": Success,
		"msg":  Success.String(),
		"data": gin.H{
			"value": envelopeP.Value,
		},
	})

}
