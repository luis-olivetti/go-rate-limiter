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
	viper.SetDefault("RATE_LIMITER_STRATEGY", "memory")
}

func main() {
	engine := server.NewServerHTTP(conf)
	http.Run(engine, fmt.Sprintf(":%d", viper.GetInt("PORT")))
}
