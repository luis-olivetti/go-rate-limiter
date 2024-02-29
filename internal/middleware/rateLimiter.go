package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RateLimiterMiddleware(conf *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
