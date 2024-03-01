package ratelimiter

type LimiterStrategy interface {
	Allow() bool
}
