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
	initScript := redis.NewScript(`
		redis.call("SET", "TotalMoney", ARGV[1])
		redis.call("SET", "MaxMoney", ARGV[2])
		redis.call("SET", "MinMoney", ARGV[3])
		redis.call("SET", "SnatchProb", ARGV[4])
		redis.call("SET", "TotalAmount", ARGV[5])
		redis.call("SET", "EnvelopeAmount", 0)
		redis.call("SET", "UsedMoney", 0)
		return 0
	`)
	_, err = initScript.Run(Ctx, DB, []string{}, config.TotalMoney, config.MaxMoney, config.MinMoney, config.SnatchProb, config.TotalAmount).Int()
	if err != nil {
		logger.Sugar.Fatalw("redis init", "Config init error", err)
	}
	// err = DB.Set(Ctx, "TotalMoney", config.TotalMoney, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "TotalMoney set error", config.TotalMoney)
	// }
	// logger.Sugar.Info("redis init", "TotalMoney set success", config.TotalMoney)

	// err = DB.Set(Ctx, "MaxMoney", config.MaxMoney, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "MaxMoney set error", config.MaxMoney)
	// }
	// logger.Sugar.Info("redis init", "MaxMoney set success", config.MaxMoney)

	// err = DB.Set(Ctx, "MinMoney", config.MinMoney, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "MinMoney set error", config.MinMoney)
	// }
	// logger.Sugar.Info("redis init", "MinMoney set success", config.MinMoney)

	// err = DB.Set(Ctx, "SnatchProb", config.SnatchProb, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "SnatchProb set error", config.SnatchProb)
	// }
	// logger.Sugar.Info("redis init", "SnatchProb set success", config.SnatchProb)

	// err = DB.Set(Ctx, "MaxSnatchAmount", config.MaxSnatchAmount, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "MaxSnatchAmount set error", config.MaxSnatchAmount)
	// }
	// logger.Sugar.Info("redis init", "MaxSnatchAmount set success", config.MaxSnatchAmount)

	// err = DB.Set(Ctx, "TotalAmount", config.TotalAmount, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "TotalAmount set error", config.TotalAmount)
	// }
	// logger.Sugar.Info("redis init", "TotalAmount set success", config.TotalAmount)

	// err = DB.Set(Ctx, "EnvelopeAmount", 0, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "EnvelopeAmount set error", 0)
	// }
	// logger.Sugar.Info("redis init", "EnvelopeAmount set success", 0)

	// err = DB.Set(Ctx, "UsedMoney", 0, 0).Err()
	// if err != nil {
	// 	logger.Sugar.Fatalw("redis init", "UsedMoney set error", 0)
	// }
	// logger.Sugar.Info("redis init", "UsedMoney set success", 0)

	snowflake.SetStartTime(time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC))
	logger.Sugar.Debugw("redis init", "redis userdb config", DB)
}
