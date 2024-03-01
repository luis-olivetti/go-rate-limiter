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

		rateLimiter := ratelimiter.InitRateLimiter(iLimiterStrategy)
		if !rateLimiter.Allow(ctx) {
			msg := "Você atingiu o número máximo de requisições ou ações permitidas dentro de um determinado intervalo de tempo"
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": msg})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
