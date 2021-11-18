package rdb

import (
	"context"
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/logger"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

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

	client := redis.NewClient(&option)
	if result, err := client.Ping(context.TODO()).Result(); err != nil {
		logger.Sugar.Fatalw("redis init error",
			"redis config", client,
			"ping redis result", result,
			"error msg", err.Error())
	}

	Client = client
	logger.Sugar.Debugw("redis init", "redis userdb config", Client)
}
