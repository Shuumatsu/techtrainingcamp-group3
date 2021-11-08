package bloomfilter

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/schollz/progressbar/v3"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/sql"
	"techtrainingcamp-group3/logger"
)

var User *bloom.BloomFilter
var Envelope *bloom.BloomFilter

func init() {
	User = bloom.NewWithEstimates(config.UserAmount, 0.001)
	userRows, err := sql.DB.Model(&dbmodels.User{}).Rows()
	defer userRows.Close()
	if err != nil {
		logger.Sugar.Errorw("bloom filter init", "user filter error", err)
		panic(err)
	}
	bar := progressbar.Default(config.UserAmount, "init user bloom filter")
	for userRows.Next() {
		var user dbmodels.User
		// ScanRows 方法用于将一行记录扫描至结构体
		sql.DB.ScanRows(userRows, &user)

		// 业务逻辑...
		User.AddString(user.Uid.String())
		bar.Add(1)
	}
	Envelope = bloom.NewWithEstimates(config.MaxSnatchAmount*config.UserAmount, 0.001)
	envelopeRows, err := sql.DB.Model(&dbmodels.Envelope{}).Rows()
	defer envelopeRows.Close()
	if err != nil {
		logger.Sugar.Errorw("bloom filter init", "envelope filter error", err)
		panic(err)
	}
	var envelopeCount int64
	err = sql.DB.Table(dbmodels.Envelope{}.TableName()).Count(&envelopeCount).Error
	if err != nil {
		logger.Sugar.Errorw("bloom filter init", "get envelope count error", err)
		panic(err)
	}
	bar = progressbar.Default(envelopeCount, "init envelope bloom filter")
	for envelopeRows.Next() {
		var envelope dbmodels.Envelope
		// ScanRows 方法用于将一行记录扫描至结构体
		sql.DB.ScanRows(envelopeRows, &envelope)

		// 业务逻辑...
		User.AddString(envelope.Uid.String())
		bar.Add(1)
	}
	logger.Sugar.Debugw("bloom filter init success")
}
