package strategies

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
	"github.com/spf13/viper"
)

func TestLimiterRedis_Allow(t *testing.T) {
	requestParams := &model.RequestParams{
		Key:        "test_key",
		Interval:   time.Second * 5,
		LimitCount: 5,
		BlockTime:  time.Second * 10,
	}

	client, mock := redismock.NewClientMock()

	value := &rateLimiterRedis{
		Created: time.Now(),
		Count:   1,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	value.Count = 2
	valueByte, _ := marshalRateLimiterRedis(value)

	mock.ExpectSet("test_key", valueByte, 0).SetVal("OK")

	globalLimiterRedis = &LimiterRedis{
		redis: client,
	}

	limiter := LimiterRedis{
		redis: client,
	}

	allowed := limiter.Allow(&gin.Context{}, requestParams, viper.New())
	if !allowed {
		t.Error("Expected allow for the first request, got denied")
	}
}

func marshalRateLimiterRedis(value *rateLimiterRedis) ([]byte, error) {
	byteValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}
