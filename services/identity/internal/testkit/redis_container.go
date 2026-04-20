package testkit

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

type redisState struct {
	container *tcredis.RedisContainer
	client    *redis.Client
	err       error
}

var (
	redisOnce   sync.Once
	sharedRedis redisState
)

// SharedRedis returns the shared Redis state for integration tests.
// SharedRedis 返回集成测试共享的 Redis 状态。
func SharedRedis(t *testing.T) *redisState {
	t.Helper()

	redisOnce.Do(func() {
		ctx := t.Context()
		container, err := tcredis.Run(ctx, "redis:7-alpine")
		if err != nil {
			sharedRedis.err = fallbackRedis(ctx, err)
			return
		}

		connString, err := container.ConnectionString(ctx)
		if err != nil {
			sharedRedis.err = err
			return
		}

		opts, err := redis.ParseURL(connString)
		if err != nil {
			sharedRedis.err = err
			return
		}
		client := redis.NewClient(opts)
		if err := client.Ping(ctx).Err(); err != nil {
			sharedRedis.err = err
			return
		}

		sharedRedis.container = container
		sharedRedis.client = client
	})

	if sharedRedis.err != nil {
		t.Skipf("skip Redis integration test: %v", sharedRedis.err)
	}

	return &sharedRedis
}

// ResetRedis flushes the shared Redis database before a test.
// ResetRedis 在测试前清空共享 Redis 数据库。
func ResetRedis(t *testing.T) {
	t.Helper()

	state := SharedRedis(t)
	if err := state.client.FlushDB(context.Background()).Err(); err != nil {
		t.Fatalf("failed to flush Redis test DB: %v", err)
	}
}

// RedisClient returns the shared Redis client for integration tests.
// RedisClient 返回集成测试共享的 Redis 客户端。
func RedisClient(t *testing.T) *redis.Client {
	t.Helper()
	return SharedRedis(t).client
}

func fallbackRedis(ctx context.Context, containerErr error) error {
	cfg, ok, err := loadRedisEnv()
	if err != nil {
		return fmt.Errorf("load Redis fallback environment failed after container error %v: %w", containerErr, err)
	}
	if !ok {
		return containerErr
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Database,
	})
	if pingErr := client.Ping(ctx).Err(); pingErr != nil {
		return fmt.Errorf("start Redis container failed: %v; fallback Redis ping failed: %w", containerErr, pingErr)
	}

	sharedRedis.client = client

	return nil
}
