package redis

import (
	"github.com/go-redis/redis"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
)

var userRds *redis.Client
var envelopeRds *redis.Client

func init() {
	option := redis.Options{
		Addr:     config.Env.RedisHost + ":" + config.Env.RedisPort,
		Password: config.Env.RedisPasswd,
	}
	rd := redis.NewClient(&option)
	if pingResult, err := rd.Ping().Result(); err != nil {
		logger.Sugar.Errorw("redis init error",
			"redis config", rd,
			"ping redis result", pingResult,
			"error msg", err.Error())
		panic(err)
	}
	userRds = rd
	logger.Sugar.Debugw("redis init", "redis config", userRds)
}
