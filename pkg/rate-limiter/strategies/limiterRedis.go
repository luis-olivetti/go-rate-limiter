package strategies

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
	redisdb "github.com/luis-olivetti/go-rate-limiter/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type rateLimiterRedis struct {
	Created time.Time `json:"created"`
	Blocked time.Time `json:"blocked"`
	Count   int64     `json:"count"`
}

type LimiterRedis struct {
	mu    sync.Mutex
	redis *redis.Client
}

var globalLimiterRedis *LimiterRedis

func (l *LimiterRedis) Allow(
	ctx *gin.Context,
	requestParams *model.RequestParams,
	conf *viper.Viper,
) bool {
	if globalLimiterRedis == nil {
		globalLimiterRedis = &LimiterRedis{
			redis: redisdb.InitRedis(conf),
		}
		log.Println("Redis limiter initialized")
	}

	globalLimiterRedis.mu.Lock()
	defer globalLimiterRedis.mu.Unlock()

	val, err := globalLimiterRedis.redis.Get(ctx, requestParams.Key).Result()
	if err != nil {
		if err == redis.Nil {
			rateValue := &rateLimiterRedis{
				Created: time.Now(),
				Count:   1,
			}

			err := setKey(ctx, requestParams.Key, rateValue)
			if err != nil {
				log.Println(err)
				return false
			}

			log.Println("Key created")
			return true
		}

		log.Println(err)
		return false
	}

	var rateValue rateLimiterRedis
	err = json.Unmarshal([]byte(val), &rateValue)
	if err != nil {
		return false
	}

	rateValue.Count++

	if !rateValue.Blocked.IsZero() {
		if blockExpiredRedis(&rateValue, requestParams.BlockTime.Milliseconds()) {
			resetKeyRedis(&rateValue)
			err := setKey(ctx, requestParams.Key, &rateValue)
			if err != nil {
				log.Println(err)
				return false
			}

			return true
		} else {
			log.Println("blocked")
			return false
		}
	}

	if createTimeExceededLimitIntervalRedis(&rateValue, requestParams.Interval.Milliseconds()) {
		resetKeyRedis(&rateValue)
	} else if rateValue.Count > requestParams.LimitCount {
		rateValue.Blocked = time.Now()

		err := setKey(ctx, requestParams.Key, &rateValue)
		if err != nil {
			log.Println(err)
			return false
		}

		return false
	}

	err = setKey(ctx, requestParams.Key, &rateValue)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func blockExpiredRedis(lim *rateLimiterRedis, blockTime int64) bool {
	return time.Since(lim.Blocked).Milliseconds() > blockTime
}

func createTimeExceededLimitIntervalRedis(lim *rateLimiterRedis, interval int64) bool {
	return time.Since(lim.Created).Milliseconds() > interval
}

func resetKeyRedis(lim *rateLimiterRedis) {
	lim.Blocked = time.Time{}
	lim.Created = time.Now()
	lim.Count = 1
}

func setKey(ctx *gin.Context, key string, value *rateLimiterRedis) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = globalLimiterRedis.redis.Set(ctx, key, valueJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
