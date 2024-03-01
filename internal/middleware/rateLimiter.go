package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/strategies"
	"github.com/spf13/viper"
)

func RateLimiterMiddleware(conf *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		limiterRedis := &strategies.LimiterRedis{}
		rateLimiter := ratelimiter.InitRateLimiter(limiterRedis)
		if !rateLimiter.Allow() {
			msg := "Você atingiu o número máximo de requisições ou ações permitidas dentro de um determinado intervalo de tempo"
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": msg})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
