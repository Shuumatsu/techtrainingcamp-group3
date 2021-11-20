package rdb

import (
	"context"
	"fmt"
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/logger"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func init() {
	poolSize, err := strconv.Atoi(config.Env.RedisPoolSize)
	if err != nil {
		logger.Sugar.Fatalw("[rdb.init]", "RedisPoolSize", config.Env.RedisPoolSize)
	}

	// init user db
	option := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Env.RedisHost, config.Env.RedisPort),
		Password: config.Env.RedisPasswd,
		PoolSize: poolSize,
	}

	client := redis.NewClient(&option)
	if result, err := client.Ping(context.TODO()).Result(); err != nil {
		logger.Sugar.Fatalw("[rdb.init] ping", "client", client, "result", result, "error", err)
	}

	Client = client
	logger.Sugar.Debugw("redis init", "redis userdb config", Client)
}
