package ratelimiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
	"github.com/spf13/viper"
)

type RateLimiter struct {
	ILimiterStrategy    ILimiterStrategy
	MaxRequestsPerIP    int64
	MaxRequestsPerToken int64
	TimeWindowMillis    int64
	BlockTimeMillis     int64
}

func InitRateLimiter(strategy ILimiterStrategy, conf *viper.Viper) *RateLimiter {
	return &RateLimiter{
		ILimiterStrategy:    strategy,
		MaxRequestsPerIP:    conf.GetInt64("RATE_LIMITER_IP_MAX_REQUESTS"),
		MaxRequestsPerToken: conf.GetInt64("RATE_LIMITER_TOKEN_MAX_REQUESTS"),
		TimeWindowMillis:    conf.GetInt64("RATE_LIMITER_TIME_WINDOW_MILISECONDS"),
		BlockTimeMillis:     conf.GetInt64("RATE_LIMITER_BLOCKING_TIME_MILLISECONDS"),
	}
}

func (r *RateLimiter) Allow(ctx *gin.Context, conf *viper.Viper) bool {
	var key string
	var limitCount int64

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
		Interval:   time.Duration(r.TimeWindowMillis) * time.Millisecond,
		BlockTime:  time.Duration(r.BlockTimeMillis) * time.Millisecond,
	}

	return r.ILimiterStrategy.Allow(ctx, requestParams, conf)
}
