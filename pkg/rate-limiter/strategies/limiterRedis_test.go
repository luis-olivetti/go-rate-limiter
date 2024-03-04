package strategies

import (
	"encoding/json"
	"errors"
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
	requestParams := getDefaultRequestParams()

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

func TestLimiterRedis_Allow_WithSecondRequest(t *testing.T) {
	requestParams := getDefaultRequestParams()

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

func TestLimiterRedis_NotAllow_WithSixthRequest(t *testing.T) {
	requestParams := getDefaultRequestParams()

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

func TestLimiterRedis_NotAllow_WithFirstRequestButItHappensErrorOnSetKey(t *testing.T) {
	requestParams := getDefaultRequestParams()

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

	mock.ExpectSet("test_key", valueByte, 0).SetErr(errors.New("example error"))

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

func TestLimiterRedis_NotAllow_WithFirstRequestButItHappensErrorOnGetKey(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	// Utilizado para mockar o time.Now()
	timePatch := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return timePatch
	})

	mock.ExpectGet("test_key").SetErr(errors.New("example error"))

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

func TestLimiterRedis_NotAllow_WithUnmarshalFail(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	mock.ExpectGet("test_key").SetVal("invalid json")

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

func TestLimiterRedis_Allow_WithExpiredBlockedKey(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	value := &rateLimiterRedis{
		Created: time.Now(),
		Blocked: time.Now().Add(-time.Second * 11),
		Count:   5,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	timePatch := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return timePatch
	})

	value.Count = 1
	value.Created = time.Now()
	value.Blocked = time.Time{}
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

func TestLimiterRedis_NotAllow_WithExpiredBlockedKeyButItHappensErrorOnSetKey(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	value := &rateLimiterRedis{
		Created: time.Now(),
		Blocked: time.Now().Add(-time.Second * 11),
		Count:   5,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	timePatch := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return timePatch
	})

	value.Count = 1
	value.Created = time.Now()
	value.Blocked = time.Time{}
	valueByte, _ := marshalRateLimiterRedis(value)

	mock.ExpectSet("test_key", valueByte, 0).SetErr(errors.New("example error"))

	globalLimiterRedis = &LimiterRedis{
		redis: client,
	}

	limiter := LimiterRedis{
		redis: client,
	}

	allowed := limiter.Allow(&gin.Context{}, requestParams, viper.New())
	if allowed {
		t.Error("Expected now allow for the request, got allowed")
	}
}

func TestLimiterRedis_Allow_WithBlockedKey(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	value := &rateLimiterRedis{
		Created: time.Now(),
		Blocked: time.Now().Add(-time.Second * 5),
		Count:   5,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	globalLimiterRedis = &LimiterRedis{
		redis: client,
	}

	limiter := LimiterRedis{
		redis: client,
	}

	allowed := limiter.Allow(&gin.Context{}, requestParams, viper.New())
	if allowed {
		t.Error("Expected now allow for the request, got allowed")
	}
}

func TestLimiterRedis_Allow_WithCreateTimeExceedLimitInterval(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	timePatch := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return timePatch
	})

	value := &rateLimiterRedis{
		Created: time.Now().Add(-time.Second * 6),
		Count:   5,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	value.Count = 1
	value.Created = time.Now()
	value.Blocked = time.Time{}
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

func TestLimiterRedis_NotAllow_WithLimitCountExceed(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	timePatch := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return timePatch
	})

	value := &rateLimiterRedis{
		Created: time.Now().Add(-time.Second * 2),
		Count:   5,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	value.Count = 6
	value.Blocked = time.Now()
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

func TestLimiterRedis_NotAllow_WithSecondRequestButItHappensErrorOnSetKey(t *testing.T) {
	requestParams := getDefaultRequestParams()

	client, mock := redismock.NewClientMock()

	timePatch := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return timePatch
	})

	value := &rateLimiterRedis{
		Created: time.Now().Add(-time.Second * 2),
		Count:   1,
	}

	valueString, _ := marshalRateLimiterRedis(value)

	mock.ExpectGet("test_key").SetVal(string(valueString))

	mock.ExpectSet("test_key", nil, 0).SetErr(errors.New("example error"))

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

func getDefaultRequestParams() *model.RequestParams {
	return &model.RequestParams{
		Key:        "test_key",
		Interval:   time.Second * 5,
		LimitCount: 5,
		BlockTime:  time.Second * 10,
	}
}
