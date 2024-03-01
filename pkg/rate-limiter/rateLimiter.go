package ratelimiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
	"github.com/spf13/viper"
)

type RateLimiter struct {
	iLimiterStrategy    ILimiterStrategy
	MaxRequestsPerIP    int64
	MaxRequestsPerToken int64
	TimeWindowMillis    int64
}

func InitRateLimiter(strategy ILimiterStrategy, conf *viper.Viper) *RateLimiter {
	return &RateLimiter{
		iLimiterStrategy:    strategy,
		MaxRequestsPerIP:    conf.GetInt64("RATE_LIMITER_IP_MAX_REQUESTS"),
		MaxRequestsPerToken: conf.GetInt64("RATE_LIMITER_TOKEN_MAX_REQUESTS"),
		TimeWindowMillis:    conf.GetInt64("RATE_LIMITER_TIME_WINDOW_MILISECONDS"),
	}
}

func (r *RateLimiter) Allow(ctx *gin.Context) bool {
	var key string
	var limitCount int64
	interval := time.Duration(r.TimeWindowMillis) * time.Millisecond

	apiKey := ctx.GetHeader("API_KEY")

	if apiKey != "" {
		key = apiKey
		limitCount = int64(r.MaxRequestsPerToken)
	} else {
		key = ctx.ClientIP()
		limitCount = int64(r.MaxRequestsPerIP)
	}

	requestParams := &model.RequestParams{
		Key:        key,
		LimitCount: limitCount,
		Interval:   int64(interval),
	}

	return r.iLimiterStrategy.Allow(ctx, requestParams)
}
