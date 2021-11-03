package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"techtrainingcamp-group3/db/mg"
	rd "techtrainingcamp-group3/db/rds"
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
	var user models.User
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
	return mg.FindUserByUID(uid)
}
