package redis

import (
	"errors"

	"github.com/go-redis/redis"
)

var ErrFuncNotDefined = errors.New("the function is not defined")
var ErrNotFound = errors.New("not found")

func NewRedis(host, port string, password string, poolSize int) *redis.Client {
	option := redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		PoolSize: poolSize,
	}

	client := redis.NewClient(&option)
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	return client
}
