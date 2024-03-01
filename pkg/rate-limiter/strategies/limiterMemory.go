package strategies

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
)

type rateLimiterMemory struct {
	lastAccess time.Time
	count      int64
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
			lastAccess: time.Now(),
			count:      1,
		}
		return true
	}

	if time.Since(lim.lastAccess).Milliseconds() > requestParams.Interval.Milliseconds() {
		lim.count = 0
		lim.lastAccess = time.Now()
	}

	lim.count++

	return lim.count <= requestParams.LimitCount
}
