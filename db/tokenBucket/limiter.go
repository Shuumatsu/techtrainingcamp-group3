package tokenBucket

import (
	"golang.org/x/time/rate"
	"strconv"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
	"time"
)

var Limiter *rate.Limiter

func init() {
	tokenInterval, err := strconv.Atoi(config.Env.TokenInterval)
	if err != nil {
		logger.Sugar.Fatalw("tokenBucket init", "tokenInterval error", err)
	}
	tokenMaxCount, err := strconv.Atoi(config.Env.TokenMaxCount)
	if err != nil {
		logger.Sugar.Fatalw("tokenBucket init", "tokenMaxCount error", err)
	}
	limit := rate.Every(time.Duration(tokenInterval) * time.Millisecond)
	Limiter = rate.NewLimiter(limit, tokenMaxCount)
}
