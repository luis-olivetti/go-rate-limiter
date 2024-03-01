package strategies

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-rate-limiter/pkg/rate-limiter/model"
)

func TestLimiterMemory_Allow(t *testing.T) {
	requestParams := &model.RequestParams{
		Key:        "test_key",
		Interval:   time.Second * 5,
		LimitCount: 5,
	}

	limiter := LimiterMemory{}

	allowed := limiter.Allow(&gin.Context{}, requestParams)
	if !allowed {
		t.Error("Expected allow for the first request, got denied")
	}

	for i := 0; i < 4; i++ {
		allowed = limiter.Allow(&gin.Context{}, requestParams)
		if !allowed {
			t.Errorf("Expected allow for request %d, got denied", i+2)
		}
	}

	allowed = limiter.Allow(&gin.Context{}, requestParams)
	if allowed {
		t.Error("Expected deny due to limit exceeded, got allowed")
	}

	time.Sleep(requestParams.Interval)

	allowed = limiter.Allow(&gin.Context{}, requestParams)
	if !allowed {
		t.Error("Expected allow after interval, got denied")
	}
}

func TestLimiterMemory_MultipleKeys(t *testing.T) {
	limiter := LimiterMemory{}

	requestParams := &model.RequestParams{
		Key:        "key1",
		Interval:   time.Second * 5,
		LimitCount: 2,
	}

	for i := 0; i < 2; i++ {
		allowed := limiter.Allow(&gin.Context{}, requestParams)
		if !allowed {
			t.Errorf("Expected allow for request %d with key1, got denied", i+1)
		}
	}

	requestParams.Key = "key2"
	allowed := limiter.Allow(&gin.Context{}, requestParams)
	if !allowed {
		t.Error("Expected allow for request with key2, got denied")
	}
}

func TestLimiterMemory_NotAllow(t *testing.T) {
	requestParams := &model.RequestParams{
		Key:        "test_key",
		Interval:   time.Second * 5,
		LimitCount: 5,
	}

	limiter := LimiterMemory{}

	for i := 0; i < 10; i++ {
		allowed := limiter.Allow(&gin.Context{}, requestParams)
		if !allowed && i < 5 {
			t.Errorf("Expected allow for request %d, got denied", i+1)
		} else if allowed && i >= 5 {
			t.Errorf("Expected deny for request %d, got allowed", i+1)
		}
	}
}
