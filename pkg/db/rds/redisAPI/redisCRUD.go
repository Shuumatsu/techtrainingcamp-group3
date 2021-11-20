package redisAPI

import (
	"github.com/go-redis/redis"
	"github.com/godruoyi/go-snowflake"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/rds"
	"techtrainingcamp-group3/pkg/logger"
	"time"
)

// SetUserByUID
//
// key: uid, value: user, expiration: 过期时间
//
// 设置不成功返回error
func SetUserByUID(user *dbmodels.User, expiration time.Duration) error {
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
func SetEnvelopeByEID(envelope *dbmodels.Envelope, expiration time.Duration) error {
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

// FindEnvelopeByEIDUID
//	find the envelope and check the envelope owner
func FindEnvelopeByEIDUID(eid dbmodels.EID, uid dbmodels.UID) (*dbmodels.Envelope, error) {
	var envelope dbmodels.Envelope
	err := rds.DB.Get(eid.Key()).Scan(&envelope)
	if err != nil {
		if err == redis.Nil {
			return nil, Error.NotFound
		}
		return nil, err
	}
	if envelope.Uid != uid {
		return nil, dbmodels.Error.ErrorEnvelopeOwner
	}
	return &envelope, nil
}

// GetRandEnvelope
// Return an envelope with random value
func GetRandEnvelope(uid dbmodels.UID) (dbmodels.Envelope, error) {
	script := redis.NewScript(`
		local TotalAmount = redis.call('GET', 'TotalAmount')
		local EnvelopeAmount = redis.call('INCR','EnvelopeAmount')
		
		if (EnvelopeAmount > TotalAmount)
		then
			return 0
		end
		local LeftMoney = redis.call('GET', 'TotalMoney') - redis.call('GET', 'UsedMoney')
		
		if (LeftMoney <= 0)
		then
			return 0
		end
		
		local MinMoney = redis.call('GET', 'MinMoney')
		
		if (MinMoney > LeftMoney)
		then
			return 0
		end
		
		local MaxMoney = math.min(redis.call('GET', 'MaxMoney'), LeftMoney)
		
		math.randomseed(os.time())
		local Money =  math.random(MinMoney, MaxMoney)
		
		redis.call('INCRBY', 'UsedMoney', Money)
		redis.call('INCR', 'EnvelopeAmount')
		
		return Money
	`)

	value, _ := script.Run(rds.DB, []string{}).Uint64()

	return dbmodels.Envelope{
		EnvelopeId: dbmodels.EID(snowflake.ID()),
		Uid:        uid,
		Value:      value,
		SnatchTime: time.Now().Unix(),
	}, nil
}
