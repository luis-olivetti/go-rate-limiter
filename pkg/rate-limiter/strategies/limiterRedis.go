package strategies

type LimiterRedis struct {
}

func (l *LimiterRedis) Allow() bool {
	return true
}
