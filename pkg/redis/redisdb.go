package redisdb

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis(conf *viper.Viper) *redis.Client {
	redisAddr := conf.GetString("REDIS_ADDR")
	redisPassword := conf.GetString("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	return client
}
