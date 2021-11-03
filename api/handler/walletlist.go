package handler

import (
	"github.com/gin-gonic/gin"
	"techtrainingcamp-group3/db/mongo"
	"techtrainingcamp-group3/db/redis"
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

	user, err := GetUserByUID(req.Uid)
	// fail
	if err != nil {
		c.JSON(200, gin.H{
			"code": FAIL,
			"msg":  err.Error(),
		})
		logger.Sugar.Debugw("WalletListHandler",
			"get Envelopes error", err)
		return
	}
	// success
	c.JSON(200, models.WalletListResp{
		Code: SUCCESS,
		Msg:  "success",
		Data: user.Wallet,
	})
}

func GetUserByUID(uid models.UID) (*models.User, error) {
	var user *models.User
	// 查询redis缓存
	if err := redis.RD.Get(uid.String()).Scan(user); err == nil {
		logger.Sugar.Debugw("getWalletList: get in redis cache")
		return user, nil
	}
	// 查询mongodb
	return mongo.FindUserByUID(uid)
}
