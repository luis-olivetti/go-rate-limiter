package strategies

type LimiterMemory struct {
}

func (l *LimiterMemory) Allow() bool {
	return true
}
