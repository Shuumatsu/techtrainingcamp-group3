package redisAPI

import (
	"fmt"
	"github.com/go-redis/redis"
)

type redisApiError struct {
	FuncNotDefined error
	NotFound       error
}

var Error redisApiError

func init() {
	Error.FuncNotDefined = fmt.Errorf("the function is not defined")
	Error.NotFound = redis.Nil
}
