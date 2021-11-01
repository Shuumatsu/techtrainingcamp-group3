package handler

import (
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

type WalletListEnvelope struct {
	EnvelopeId uint64 `gorm:"envelope_id" json:"envelope_id"`
	Opened     bool   `gorm:"opened" json:"opened"`
	Value      uint64 `gorm:"value" json:"value,omitempty"`
	SnatchTime uint64 `gorm:"snatch_time" json:"snatch_time"`
}

func (WalletListEnvelope) TableName() string {
	return "envelope"
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
			"msg":  "fail",
		})
		return
	}

	amount := 50

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": envelopes,
		},
	})
}

func GetAllEnvelopesByUID(uid uint64) ([]*WalletListEnvelope, error) {
	var user WalletListUser
	if err := db.DB.Table(WalletListUser{}.TableName()).First(
		&user, uid).Error; err != nil {
		return nil, err
	}
	envelopesID, err := db.ParseEnvelopeList(user.EnvelopeList)
	if err != nil {
		return nil, err
	}
	logger.Sugar.Debugw("debug",
		"evlpsID", envelopesID)
	var envelopes []*WalletListEnvelope
	if err = db.DB.Table(WalletListEnvelope{}.TableName()).Where(
		envelopesID).Find(&envelopes).Error; err != nil {
		return nil, err
	}
	return envelopes, nil
}
