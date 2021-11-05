package redisAPI

import (
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
	err := rds.UserRds.Set(user.Uid.String(), user, expiration).Err()
	if err != nil {
		logger.Sugar.Errorw("redis: set user by uid", "error", err)
		return err
	}
	return nil
}

func FindUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, RedisApiError.FuncNotDefined
}

// SetEnvelopeByEID
//
// key: envelopeId, value: envelope, expiration: 过期时间
//
// 设置不成功返回error
func SetEnvelopeByEID(envelope dbmodels.Envelope, expiration time.Duration) error {
	err := rds.UserRds.Set(envelope.EnvelopeId.String(), envelope, expiration).Err()
	if err != nil {
		logger.Sugar.Errorw("redis: set envelope by eid", "error", err)
		return err
	}
	return nil
}
func FindEnvelopeByEID(eid dbmodels.EID) (dbmodels.Envelope, error) {
	return dbmodels.Envelope{}, RedisApiError.FuncNotDefined
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return RedisApiError.FuncNotDefined
}
