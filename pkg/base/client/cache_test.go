package client

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
	"testing"
	"time"
)

func initRedisTest(t *testing.T) *redis.Client {
	t.Helper()
	logger.Ignore()
	config.Init("test")
	//logger.Init("test", "info")
	r, err := NewRedisClient(0)
	if err != nil {
		t.Fatalf("NewRedisClient error: %v", err)
	}
	return r
}

func TestRedisPingPong(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}

	r := initRedisTest(t)
	now := time.Now().UnixMilli()
	times := 1000
	ctx := context.Background()

	for i := 0; i < times; i++ {
		r.Ping(ctx)
	}
	pingConsumer := time.Now().UnixMilli() - now

	now = time.Now().UnixMilli()
	for i := 0; i < times; i++ {
		r.Set(ctx, fmt.Sprintf("test:%d", i), 1, 100*time.Second)
	}
	setConsumer := time.Now().UnixMilli() - now

	fmt.Printf("ping consume: %d, \nset consume:  %d\n", pingConsumer, setConsumer)
}
