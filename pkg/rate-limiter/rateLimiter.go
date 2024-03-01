package ratelimiter

type RateLimiter struct {
	iLimiterStrategy ILimiterStrategy
}

func InitRateLimiter(strategy ILimiterStrategy) *RateLimiter {
	return &RateLimiter{
		iLimiterStrategy: strategy,
	}
}

func (r *RateLimiter) Allow() bool {
	return r.iLimiterStrategy.Allow()
}
