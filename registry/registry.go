package registry

import (
	"techtrainingcamp-group3/registry/mysql"

	"techtrainingcamp-group3/logger"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var Logger *zap.Logger = logger.NewLogger("")

type Registry struct {
	db    *mysql.Database
	redis *redis.Client
}
