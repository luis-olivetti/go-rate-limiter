package strategies

import (
	"encoding/json"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func TestLimiterRedis_Allow_WithFirstRequest(t *testing.T) {
	requestParams := &model.RequestParams{
		Key:        "test_key",
		Interval:   time.Second * 5,
		LimitCount: 5,
		BlockTime:  time.Second * 10,
	}

	client, mock := redismock.NewClientMock()

	// Utilizado para mockar o time.Now()
	timePatch := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return timePatch
	})

	mock.ExpectGet("test_key").SetErr(redis.Nil)

	value := &rateLimiterRedis{
		Created: time.Now(),
		Count:   1,
	}
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
		t.Error("Expected allow for the request, got denied")
	}
}

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
		t.Error("Expected allow for the request, got denied")
	}
}

func TestLimiterRedis_NotAllow(t *testing.T) {
	requestParams := &model.RequestParams{
		Key:        "test_key",
		Interval:   time.Second * 5,
		LimitCount: 5,
		BlockTime:  time.Second * 10,
	}

	client, mock := redismock.NewClientMock()

	value := &rateLimiterRedis{
		Created: time.Now(),
		Count:   5,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	value.Count = 6
	valueByte, _ := marshalRateLimiterRedis(value)

	mock.ExpectSet("test_key", valueByte, 0).SetVal("OK")

	globalLimiterRedis = &LimiterRedis{
		redis: client,
	}

	limiter := LimiterRedis{
		redis: client,
	}

	allowed := limiter.Allow(&gin.Context{}, requestParams, viper.New())
	if allowed {
		t.Error("Expected not allow for the request, got allowed")
	}
}

func marshalRateLimiterRedis(value *rateLimiterRedis) ([]byte, error) {
	byteValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

// func unmarshalRateLimiterRedis(value []byte) (*rateLimiterRedis, error) {
// 	var rateValue rateLimiterRedis
// 	err := json.Unmarshal(value, &rateValue)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &rateValue, nil
// }
