package rds

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/logger"
	"github.com/godruoyi/go-snowflake"
	"time"
)

var DB *redis.Client
var Ctx = context.Background()

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
	if pingResult, err := rd.Ping(Ctx).Result(); err != nil {
		logger.Sugar.Fatalw("redis init error",
			"redis config", rd,
			"ping redis result", pingResult,
			"error msg", err.Error())
	}
	DB = rd
	//TODO: Handle the case of failure setting
	DB.SetNX(Ctx, "TotalMoney", config.TotalMoney, 0)
	DB.SetNX(Ctx, "MaxMoney", config.MaxMoney, 0)
	DB.SetNX(Ctx, "MinMoney", config.MinMoney, 0)
	DB.SetNX(Ctx, "SnatchProb", config.SnatchProb, 0)
	DB.SetNX(Ctx, "MaxSnatchAmount", config.MaxSnatchAmount, 0)
	DB.SetNX(Ctx, "TotalAmount", config.TotalAmount, 0)
	DB.SetNX(Ctx, "EnvelopeAmount", 0, 0)
	DB.SetNX(Ctx, "UsedMoney", 0, 0)
	snowflake.SetStartTime(time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC))
	logger.Sugar.Debugw("redis init", "redis userdb config", DB)
}
