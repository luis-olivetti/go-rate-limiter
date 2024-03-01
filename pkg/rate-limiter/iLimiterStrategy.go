package ratelimiter

type ILimiterStrategy interface {
	Allow() bool
}
