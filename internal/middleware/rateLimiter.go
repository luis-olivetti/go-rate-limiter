package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter"
	"github.com/spf13/viper"
)

func RateLimiterMiddleware(conf *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		iLimiterStrategy := ratelimiter.GetLimiterStrategy(conf.GetString("RATE_LIMITER_STRATEGY"))

		rateLimiter := ratelimiter.InitRateLimiter(iLimiterStrategy, conf)
		if !rateLimiter.Allow(ctx, conf) {
			msg := "You have reached the maximum number of requests or actions allowed within a certain time frame"
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": msg})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
