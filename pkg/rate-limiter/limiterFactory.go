package ratelimiter

import "github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/strategies"

func GetLimiterStrategy(limiterType string) ILimiterStrategy {
	var rateLimiterStrategy ILimiterStrategy

	if limiterType == "redis" {
		rateLimiterStrategy = &strategies.LimiterRedis{}
	} else {
		rateLimiterStrategy = &strategies.LimiterMemory{}
	}

	return rateLimiterStrategy
}
