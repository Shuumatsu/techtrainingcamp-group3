package redisAPI

import (
	"github.com/go-redis/redis/v8"
	"github.com/godruoyi/go-snowflake"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/rds"
	"techtrainingcamp-group3/pkg/logger"
	"time"
	"math/rand"
)

// SetUserByUID
//
// key: uid, value: user, expiration: 过期时间
//
// 设置不成功返回error
func SetUserByUID(user *dbmodels.User, expiration time.Duration) error {
	err := rds.DB.Set(rds.Ctx,user.Uid.Key(), user, expiration).Err()
	if err != nil {
		logger.Sugar.Errorw("redis: set user by uid", "error", err)
		return err
	}
	return nil
}

// FindUserByUID
// Find user in redis according to UID
//
// 如果redis中不存在该user, 返回NotFound
// 如果redis的get操作发生错误, 返回error
func FindUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	var user dbmodels.User
	err := rds.DB.Get(rds.Ctx,uid.Key()).Scan(&user)
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
	err := rds.DB.Set(rds.Ctx,envelope.EnvelopeId.Key(), envelope, expiration).Err()
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
	err := rds.DB.Get(rds.Ctx,eid.Key()).Scan(&envelope)
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
	err := rds.DB.Get(rds.Ctx,eid.Key()).Scan(&envelope)
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
		local TotalAmount = tonumber(redis.call('GET', 'TotalAmount'))
		local EnvelopeAmount = tonumber(redis.call('GET','EnvelopeAmount'))
		
		if (EnvelopeAmount >= TotalAmount)
		then
			return 0
		end
		
		local LeftMoney = tonumber(redis.call('GET', 'TotalMoney')) - tonumber(redis.call('GET', 'UsedMoney'))
		
		if (LeftMoney <= 0)
		then
			return 0
		end
		
		local MinMoney = tonumber(redis.call('GET', 'MinMoney'))
		
		if (MinMoney > LeftMoney)
		then
			return 0
		end
		
		local MaxMoney = math.min(tonumber(redis.call('GET', 'MaxMoney')), LeftMoney)
		
		local mean = LeftMoney / (TotalAmount - EnvelopeAmount)
		local std = (MaxMoney - MinMoney) / 4
		
		local Money = math.floor(tonumber(ARGV[1]) * mean + std)
		Money = math.min(Money, MaxMoney)
		Money = math.max(Money, MinMoney)
		
		redis.call('INCRBY', 'UsedMoney', Money)
		redis.call('INCR', 'EnvelopeAmount')
		
		return Money
	`)
	rand.Seed(time.Now().Unix())
	value, err := script.Run(rds.Ctx, rds.DB, []string{}, rand.NormFloat64()).Uint64()
	if err != nil {
		logger.Sugar.Debugw("lua script error", "err msg", err)
	}
	return dbmodels.Envelope{
		EnvelopeId: dbmodels.EID(snowflake.ID()),
		Uid:        uid,
		Value:      value,
		SnatchTime: time.Now().Unix(),
	}, nil
}
// DeleteEnvelopeByEID
func DelEnvelopeByEID(eid dbmodels.EID) error {
	err := rds.DB.Del(rds.Ctx,eid.Key()).Err()
	return err
}
