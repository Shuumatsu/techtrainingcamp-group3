package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"techtrainingcamp-group3/db/dbmodels"
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

	// TODO: mysql
	user, err := sqlAPI.FindUserByUID(dbmodels.UID(req.Uid))
	// if mysql not found
	if errors.Is(err, sqlAPI.Error.NotFound) {
		c.JSON(200, models.WalletListResp{
			Code: models.NotFound,
			Msg:  models.NotFound.Message(),
			Data: models.WalletListData{},
		})
		logger.Sugar.Debugw("WalletListHandler",
			"not found in mysql", err)
		return
	}
	// find envelopes which belong to the user
	envelopes, err := sqlAPI.FindEnvelopesByUID(dbmodels.UID(req.Uid))
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
