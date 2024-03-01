package ratelimiter

import "github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/strategies"

func GetLimiterStrategy(limiterType string) ILimiterStrategy {
	var rateLimiterStrategy ILimiterStrategy

	if limiterType == "memory" {
		rateLimiterStrategy = &strategies.LimiterMemory{}
	} else {
		rateLimiterStrategy = &strategies.LimiterRedis{}
	}

	return rateLimiterStrategy
}
