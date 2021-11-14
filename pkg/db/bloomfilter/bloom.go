package bloomfilter

import (
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/sql"
	"techtrainingcamp-group3/pkg/logger"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/schollz/progressbar/v3"
)

var userFilter *bloom.BloomFilter
var envelopeFilter *bloom.BloomFilter

func init() {
	userFilter = bloom.NewWithEstimates(config.UserAmount, 0.001)
	envelopeFilter = bloom.NewWithEstimates(config.TotalAmount, 0.001)
	if config.Env.Bloomfilter == "true" {
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
			userFilter.AddString(user.Uid.String())
			bar.Add(1)
		}
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
			envelopeFilter.AddString(envelope.Uid.String())
			bar.Add(1)
		}
	}
	logger.Sugar.Debugw("bloom filter init success")
}

func AddUser(uid dbmodels.UID) {
	if config.Env.Bloomfilter != "true" {
		return
	}
	userFilter.AddString(uid.String())
}
func TestUser(uid dbmodels.UID) bool {
	if config.Env.Bloomfilter != "true" {
		return true
	}
	return userFilter.TestString(uid.String())
}

func AddEnvelope(eid dbmodels.EID) {
	if config.Env.Bloomfilter != "true" {
		return
	}
	envelopeFilter.AddString(eid.String())
}
func TestEnvelope(eid dbmodels.EID) bool {
	if config.Env.Bloomfilter != "true" {
		return true
	}
	return envelopeFilter.TestString(eid.String())
}
