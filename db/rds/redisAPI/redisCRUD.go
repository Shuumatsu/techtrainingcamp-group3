package redisAPI

import (
	"github.com/go-redis/redis"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/rds"
	"techtrainingcamp-group3/logger"
	"time"
)

// SetUserByUID
//
// key: uid, value: user, expiration: 过期时间
//
// 设置不成功返回error
func SetUserByUID(user dbmodels.User, expiration time.Duration) error {
	err := rds.DB.Set(user.Uid.Key(), user, expiration).Err()
	if err != nil {
		logger.Sugar.Errorw("redis: set user by uid", "error", err)
		return err
	}
	return nil
}

// FindUserByUID
// 根据uid在redis中查找user
//
// 如果redis中不存在该user, 返回NotFound
// 如果redis的get操作发生错误, 返回error
func FindUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	var user dbmodels.User
	err := rds.DB.Get(uid.Key()).Scan(&user)
	if err != nil {
		if err == redis.Nil {
			return nil, Error.NotFound
		}
		return nil, err
	}
	return &user, nil
}

// SetEnvelopeByEID
//
// key: envelopeId, value: envelope, expiration: 过期时间
//
// 设置不成功返回error
func SetEnvelopeByEID(envelope dbmodels.Envelope, expiration time.Duration) error {
	err := rds.DB.Set(envelope.EnvelopeId.Key(), envelope, expiration).Err()
	if err != nil {
		logger.Sugar.Errorw("redis: set envelope by eid", "error", err)
		return err
	}
	return nil
}

// FindEnvelopeByEID
// 根据envelope_id在redis中查找envelope
//
// 如果redis中不存在该envelope, 返回NotFound
// 如果redis的get操作发生错误, 返回error
func FindEnvelopeByEID(eid dbmodels.EID) (*dbmodels.Envelope, error) {
	var envelope dbmodels.Envelope
	err := rds.DB.Get(eid.Key()).Scan(&envelope)
	if err != nil {
		if err == redis.Nil {
			return nil, Error.NotFound
		}
		return nil, err
	}
	return &envelope, nil
}

func FindWalletByUID(uid dbmodels.UID) (amount uint64, envelopes []dbmodels.Envelope, err error) {
	pipe := rds.DB.TxPipeline()
	// 执行事务
	var user dbmodels.User
	err = pipe.Get(uid.Key()).Scan(&user)
	if err != nil {
		return 0, nil, err
	}
	return 0, nil, Error.FuncNotDefined
}
