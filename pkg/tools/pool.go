package tools

import (
	"math"
	"math/rand"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"time"
	"techtrainingcamp-group3/pkg/db/rds"
	"github.com/godruoyi/go-snowflake"
)

type UnopenedRedEnvelope struct {
	Money int
	Eid   uint64
}

func GetRandEnvelope(uid dbmodels.UID) dbmodels.Envelope {
	mean := float64(config.TotalMoney) / float64(config.TotalAmount)
	stdDev := math.Min(float64(config.MaxMoney)-mean, mean-float64(config.MinMoney)) / 3
	value := uint64(rand.NormFloat64()*stdDev + mean)
	if value > config.MaxMoney {
		value = config.MaxMoney
	}
	if value < config.MinMoney {
		value = config.MinMoney
	}
	v, err := rds.DB.Incr("TotalAmount").Result()
	if err != nil || v > config.TotalAmount {
		return dbmodels.Envelope{}
	}
	m, err := rds.DB.IncrBy("TotalMoney", int64(value)).Result()
	if err != nil || m > config.TotalMoney {
		return dbmodels.Envelope{}
	}
	return dbmodels.Envelope{
		EnvelopeId: dbmodels.EID(snowflake.ID()),
		Uid:        uid,
		Value:      value,
		SnatchTime: time.Now().Unix(),
	}
}

// func init() {
// 	snowflake.SetStartTime(time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC))
// }
