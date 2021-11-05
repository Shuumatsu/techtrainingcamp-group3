package rds

import (
	"github.com/go-redis/redis"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
)

var Rds *redis.Client

func init() {
	// init user db
	option := redis.Options{
		Addr:     config.Env.RedisHost + ":" + config.Env.RedisPort,
		Password: config.Env.RedisPasswd,
	}
	rd := redis.NewClient(&option)
	if pingResult, err := rd.Ping().Result(); err != nil {
		logger.Sugar.Fatalw("redis init error",
			"redis config", rd,
			"ping redis result", pingResult,
			"error msg", err.Error())
	}
	Rds = rd
	logger.Sugar.Debugw("redis init", "redis userdb config", Rds)
}
