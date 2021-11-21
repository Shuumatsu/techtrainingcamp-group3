package rds

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/godruoyi/go-snowflake"
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/logger"
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
	// TODO: Handle the case of failure setting
	err = DB.SetNX(Ctx, "TotalMoney", config.TotalMoney, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "TotalMoney set error", config.TotalMoney)
	}
	logger.Sugar.Infow("redis init TotalMoney set success", "config", config.TotalMoney, "redis", DB.Get(Ctx, "TotalMoney").String())

	err = DB.SetNX(Ctx, "MaxMoney", config.MaxMoney, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "MaxMoney set error", config.MaxMoney)
	}
	logger.Sugar.Infow("redis init MaxMoney set success", "config", config.MaxMoney, "redis", DB.Get(Ctx, "MaxMoney").String())

	err = DB.SetNX(Ctx, "MinMoney", config.MinMoney, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "MinMoney set error", config.MinMoney)
	}
	logger.Sugar.Infow("redis init MinMoney set success", "config", config.MinMoney, "redis", DB.Get(Ctx, "MinMoney").String())

	err = DB.SetNX(Ctx, "SnatchProb", config.SnatchProb, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "SnatchProb set error", config.SnatchProb)
	}
	logger.Sugar.Infow("redis init SnatchProb set success", "config", config.SnatchProb, "redis", DB.Get(Ctx, "SnatchProb").String())

	err = DB.SetNX(Ctx, "MaxSnatchAmount", config.MaxSnatchAmount, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "MaxSnatchAmount set error", config.MaxSnatchAmount)
	}
	logger.Sugar.Infow("redis init MaxSnatchAmount set success", "config", config.MaxSnatchAmount, "redis", DB.Get(Ctx, "MaxSnatchAmount").String())

	err = DB.SetNX(Ctx, "TotalAmount", config.TotalAmount, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "TotalAmount set error", config.TotalAmount)
	}
	logger.Sugar.Infow("redis init TotalAmount set success", "config", config.TotalAmount, "redis", DB.Get(Ctx, "TotalAmount").String())

	err = DB.SetNX(Ctx, "EnvelopeAmount", 0, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "EnvelopeAmount set error", 0)
	}
	logger.Sugar.Infow("redis init EnvelopeAmount set success", "config", 0, "redis", DB.Get(Ctx, "EnvelopeAmount").String())

	err = DB.SetNX(Ctx, "UsedMoney", 0, 0).Err()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "UsedMoney set error", 0)
	}
	logger.Sugar.Infow("redis init UsedMoney set success", "config", 0, "redis", DB.Get(Ctx, "UsedMoney").String())

	snowflake.SetStartTime(time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC))
	logger.Sugar.Debugw("redis init", "redis userdb config", DB)
}
