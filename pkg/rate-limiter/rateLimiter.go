package ratelimiter

type RateLimiter struct {
	limiterStrategy LimiterStrategy
}

func InitRateLimiter(strategy LimiterStrategy) *RateLimiter {
	return &RateLimiter{
		limiterStrategy: strategy,
	}
}

func (r *RateLimiter) Allow() bool {
	return r.limiterStrategy.Allow()
}
