package rds

import (
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/logger"
	"github.com/go-redis/redis"
	"github.com/godruoyi/go-snowflake"
	"time"
)

var DB *redis.Client

func init() {
	poolSize, err := strconv.Atoi(config.Env.RedisPoolSize)
	if err != nil {
		logger.Sugar.Fatalw("redis init error",
			"poolsize must be a number, poolsize:", config.Env.RedisPoolSize)
	}
	// init user db
	option := redis.Options{
		Addr:     config.Env.RedisHost + ":" + config.Env.RedisPort,
		Password: config.Env.RedisPasswd,
		PoolSize: poolSize,
	}
	rd := redis.NewClient(&option)
	if pingResult, err := rd.Ping().Result(); err != nil {
		logger.Sugar.Fatalw("redis init error",
			"redis config", rd,
			"ping redis result", pingResult,
			"error msg", err.Error())
	}
	DB = rd
	snowflake.SetStartTime(time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC))
	logger.Sugar.Debugw("redis init", "redis userdb config", DB)
}
