package strategies

import (
	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
)

type LimiterRedis struct {
}

func (l *LimiterRedis) Allow(ctx *gin.Context, requestParams *model.RequestParams) bool {
	return true
}
