package rds

import (
	"github.com/go-redis/redis"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
)

var RD *redis.Client

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
	RD = rd
	logger.Sugar.Debugw("redis init", "redis config", RD)
}
