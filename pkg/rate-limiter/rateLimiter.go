package ratelimiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
)

type RateLimiter struct {
	iLimiterStrategy    ILimiterStrategy
	MaxRequestsPerIP    int
	MaxRequestsPerToken int
	TimeWindowMillis    int
}

func InitRateLimiter(strategy ILimiterStrategy) *RateLimiter {
	return &RateLimiter{
		iLimiterStrategy:    strategy,
		MaxRequestsPerIP:    5,
		MaxRequestsPerToken: 10,
		TimeWindowMillis:    20000,
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
