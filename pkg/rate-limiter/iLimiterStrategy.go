package ratelimiter

import (
	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
	"github.com/spf13/viper"
)

type ILimiterStrategy interface {
	Allow(ctx *gin.Context, requestParams *model.RequestParams, conf *viper.Viper) bool
}
