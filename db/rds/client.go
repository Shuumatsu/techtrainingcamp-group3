package rds

import (
	"github.com/go-redis/redis"
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
)

var UserRds *redis.Client
var EnvelopeRds *redis.Client

func init() {
	// init user db
	redisUserDB, err := strconv.Atoi(config.Env.RedisUserDB)
	if err != nil {
		logger.Sugar.Fatal("the redis user db must be a number")
	}
	option := redis.Options{
		Addr:     config.Env.RedisHost + ":" + config.Env.RedisPort,
		Password: config.Env.RedisPasswd,
		DB:       redisUserDB,
	}
	rd := redis.NewClient(&option)
	if pingResult, err := rd.Ping().Result(); err != nil {
		logger.Sugar.Fatalw("redis init error",
			"redis config", rd,
			"ping redis result", pingResult,
			"error msg", err.Error())
	}
	UserRds = rd
	logger.Sugar.Debugw("redis init", "redis userdb config", UserRds)
	// init envelope db
	redisEnvelopeDB, err := strconv.Atoi(config.Env.RedisEnvelopeDB)
	if err != nil {
		logger.Sugar.Fatal("the redis envelope db must be a number")
	}
	option = redis.Options{
		Addr:     config.Env.RedisHost + ":" + config.Env.RedisPort,
		Password: config.Env.RedisPasswd,
		DB:       redisEnvelopeDB,
	}
	rd = redis.NewClient(&option)
	if pingResult, err := rd.Ping().Result(); err != nil {
		logger.Sugar.Fatalw("redis init error",
			"redis config", rd,
			"ping redis result", pingResult,
			"error msg", err.Error())
	}
	EnvelopeRds = rd
	logger.Sugar.Debugw("redis init", "redis envelopedb config", EnvelopeRds)
}
