package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"techtrainingcamp-group3/db"
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/models"
)

const (
	SUCCESS = iota
	FAIL
)

type WalletListUser struct {
	Uid          uint64 `gorm:"uid"`
	EnvelopeList string `gorm:"envelope_list"`
}

func (WalletListUser) TableName() string {
	return "user"
}

func WalletListHandler(c *gin.Context) {
	var req models.WalletListReq
	c.Bind(&req)
	logger.Sugar.Debugw("WalletListHandler",
		"uid", req.Uid)

	envelopes, err := GetAllEnvelopesByUID(req.Uid)
	if err != nil {
		c.JSON(200, gin.H{
			"code": FAIL,
			"msg":  err.Error(),
		})
		logger.Sugar.Debugw("WalletListHandler",
			"get Envelopes error", err)
		return
	}
	c.JSON(200, models.WalletListResp{
		Code: SUCCESS,
		Msg:  "success",
		Data: models.WalletListData{
			Amount:       50,
			EnvelopeList: envelopes,
		},
	})
}

func GetAllEnvelopesByUID(uid uint64) ([]models.Envelope, error) {
	var user WalletListUser
	// 尝试根据uid从user表中获取envelope_list
	if err := db.DB.Table(WalletListUser{}.TableName()).First(
		&user, uid).Error; err != nil {
		return nil, fmt.Errorf("uid not found")
	}
	// 解析envelope_list
	envelopesID, err := db.ParseEnvelopeList(user.EnvelopeList)
	if err != nil {
		return nil, err
	}
	//  没有envelope
	if len(envelopesID) == 0 {
		return nil, nil
	}
	logger.Sugar.Debugw("debug",
		"evlpsID", envelopesID)
	// 尝试根据envelope_id从envelope表中获取envelop具体数据
	var envelopes []models.Envelope
	if err = db.DB.Table(models.Envelope{}.TableName()).Where(
		envelopesID).Find(&envelopes).Error; err != nil {
		return nil, err
	}
	return envelopes, nil
}
