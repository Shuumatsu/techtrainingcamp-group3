package redisAPI

import (
	"fmt"
	"github.com/go-redis/redis"
)

type redisApiError struct {
	FuncNotDefined error
	NotFound       error
}

var RedisApiError redisApiError

func init() {
	RedisApiError.FuncNotDefined = fmt.Errorf("the function is not defined")
	RedisApiError.NotFound = redis.Nil
}
