package handler

import (
	"github.com/gin-gonic/gin"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/rds/redisAPI"
	"techtrainingcamp-group3/db/sql/sqlAPI"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"
)

func WalletListHandler(c *gin.Context) {
	var req models.WalletListReq
	c.Bind(&req)
	logger.Sugar.Debugw("WalletListHandler",
		"uid", req.Uid)
	// TODO: redis
	user, err := redisAPI.FindUserByUID(dbmodels.UID(req.Uid))
	if err != nil {
		// TODO: mysql
		// redis缓存未命中
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
			logger.Sugar.Debugw("WalletListHandler",
				"mysql database error", err)
		}
		return
	}
	// find envelopes which belong to the user
	// TODO: redis
	// parse envelope_list
	envelopesID, err := sqlAPI.ParseEnvelopeList(user.EnvelopeList)
	if err != nil {
		c.JSON(200, models.WalletListResp{
			Code: models.ParseError,
			Msg:  models.ParseError.Message(),
			Data: models.WalletListData{},
		})
		logger.Sugar.Debugw("WalletListHandler",
			"can't parse envelope_list", err)
		return
	}
	envelopes := make([]dbmodels.Envelope, 0)
	for i := 0; i < len(envelopesID); i++ {
		envelope, err := redisAPI.FindEnvelopeByEID(envelopesID[i])
		if err != nil {
			break
		}
		envelopes = append(envelopes, *envelope)
	}
	if len(envelopes) != len(envelopesID) {
		// TODO:mysql
		envelopes, err = sqlAPI.FindEnvelopesByUID(dbmodels.UID(req.Uid))
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
	}
	// change envelopes to resq model
	envelopesResq := make([]models.Envelope, len(envelopes))
	for i := 0; i < len(envelopes); i++ {
		envelopesResq[i] = envelopes[i].ToResqModel()
	}
	// success
	c.JSON(200, models.WalletListResp{
		Code: models.Success,
		Msg:  models.Success.Message(),
		Data: models.WalletListData{
			Amount:       user.Amount,
			EnvelopeList: envelopesResq,
		},
	})
}
