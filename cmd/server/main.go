package main

import (
	"fmt"

	"github.com/luis-olivetti/go-rate-limiter/internal/server"
	"github.com/luis-olivetti/go-rate-limiter/pkg/http"
	"github.com/spf13/viper"
)

var conf *viper.Viper

func init() {
	conf = viper.GetViper()

	viper.AutomaticEnv()
	// TODO: Remover esses sets
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("RATE_LIMITER_STRATEGY", "redis")
	viper.Set("RATE_LIMITER_IP_MAX_REQUESTS", 5)
	viper.Set("RATE_LIMITER_TOKEN_MAX_REQUESTS", 10)
	viper.Set("RATE_LIMITER_TIME_WINDOW_MILISECONDS", 10000)
	viper.Set("RATE_LIMITER_BLOCKING_TIME_MILLISECONDS", 2000)
	viper.Set("REDIS_ADDR", "localhost:6379")
	viper.Set("REDIS_PASSWORD", "")
}

func main() {
	engine := server.NewServerHTTP(conf)
	http.Run(engine, fmt.Sprintf(":%d", viper.GetInt("PORT")))
}
