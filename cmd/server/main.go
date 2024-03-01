package main

import (
	"fmt"

	"github.com/luis-olivetti/go-rate-limiter/internal/server"
	"github.com/luis-olivetti/go-rate-limiter/pkg/http"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	// TODO: Remover esses sets
	viper.SetDefault("PORT", "8080")
	viper.Set("RATE_LIMITER_STRATEGY", "memory")
}

func main() {
	engine := server.NewServerHTTP(&viper.Viper{})
	http.Run(engine, fmt.Sprintf(":%d", viper.GetInt("PORT")))
}
