package server

import (
	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/internal/middleware"
	"github.com/spf13/viper"
)

func NewServerHTTP(conf *viper.Viper) *gin.Engine {
	router := gin.Default()

	router.Use(
		middleware.RateLimiterMiddleware(conf),
	)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	return router
}
