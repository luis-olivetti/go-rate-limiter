package ratelimiter

import (
	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
)

type ILimiterStrategy interface {
	Allow(ctx *gin.Context, requestParams *model.RequestParams) bool
}
