package main

import (
	"fmt"

	"github.com/luis-olivetti/go-rate-limiter/internal/server"
	"github.com/luis-olivetti/go-rate-limiter/pkg/http"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
}

func main() {
	engine := server.NewServerHTTP(&viper.Viper{})
	http.Run(engine, fmt.Sprintf(":%d", viper.GetInt("PORT")))
}
