package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/mysql/sqlAPI"
	rd "techtrainingcamp-group3/db/redis"
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

	user, err := getUserByUID(req.Uid)
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
	// hide value if the envelope is not open
	err = hideValueByOpened(user)
	if err != nil {
		c.JSON(200, gin.H{
			"code": FAIL,
			"msg":  err.Error(),
		})
		logger.Sugar.Debugw("WalletListHandler",
			"hide value error", err)
		return
	}
	// success
	c.JSON(200, models.WalletListResp{
		Code: SUCCESS,
		Msg:  "success",
		Data: user.Wallet,
	})
}

func hideValueByOpened(wallet *models.WalletListData) error {
	if wallet == nil {
		return fmt.Errorf("the wallet is nil")
	}
	for i := 0; i < len(wallet.EnvelopeList); i++ {
		if wallet.EnvelopeList[i].Opened == false {
			wallet.EnvelopeList[i].Value = 0
		}
	}
	return nil
}

func getUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	var user dbmodels.User
	// 查询redis缓存
	err := rd.RD.Get(uid.String()).Scan(&user)
	if err != nil && err != redis.Nil {
		// redis error
		logger.Sugar.Debugw("redis", "error", err)
	}
	if err != redis.Nil {
		// 命中缓存
		return &user, nil
	}
	// 查询mongodb
	return sqlAPI.FindUserByUID(uid)
}
