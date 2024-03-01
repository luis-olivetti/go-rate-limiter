package strategies

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
)

type rateLimiterMemory struct {
	created time.Time
	blocked time.Time
	count   int64
}

type LimiterMemory struct {
	mu     sync.Mutex
	limits map[string]*rateLimiterMemory
}

var globalLimiterMemory *LimiterMemory

func (l *LimiterMemory) Allow(ctx *gin.Context, requestParams *model.RequestParams) bool {
	if globalLimiterMemory == nil {
		globalLimiterMemory = &LimiterMemory{
			limits: make(map[string]*rateLimiterMemory),
		}
		log.Println("Memory limiter initialized")
	}

	globalLimiterMemory.mu.Lock()
	defer globalLimiterMemory.mu.Unlock()

	lim, ok := globalLimiterMemory.limits[requestParams.Key]

	if !ok {
		globalLimiterMemory.limits[requestParams.Key] = &rateLimiterMemory{
			created: time.Now(),
			count:   1,
		}
		return true
	}

	lim.count++

	if !lim.blocked.IsZero() {
		if blockExpired(lim, requestParams.BlockTime.Milliseconds()) {
			resetKey(lim)
			return true
		} else {
			log.Println("blocked")
			return false
		}
	}

	if createTimeExceededLimitInterval(lim, requestParams.Interval.Milliseconds()) {
		resetKey(lim)
	} else if lim.count > requestParams.LimitCount {
		lim.blocked = time.Now()
		return false
	}

	return true
}

func blockExpired(lim *rateLimiterMemory, blockTime int64) bool {
	return time.Since(lim.blocked).Milliseconds() > blockTime
}

func createTimeExceededLimitInterval(lim *rateLimiterMemory, interval int64) bool {
	return time.Since(lim.created).Milliseconds() > interval
}

func resetKey(lim *rateLimiterMemory) {
	lim.blocked = time.Time{}
	lim.created = time.Now()
	lim.count = 1
}
