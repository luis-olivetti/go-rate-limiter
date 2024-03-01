package strategies

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
)

type LimiterMemory struct {
}

func (l *LimiterMemory) Allow(ctx *gin.Context, requestParams *model.RequestParams) bool {

	println(fmt.Sprintf("LimiterMemory.Allow: key=%s, limitCount=%d, interval=%d", requestParams.Key, requestParams.LimitCount, requestParams.Interval))

	return false
}
