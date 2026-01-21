package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisLock(ctx context.Context, client *redis.Client, expire time.Duration, key string) bool {
	cmd := client.SetNX(ctx, key, "ok", expire)
	if cmd.Err() != nil {
		return false
	} else {
		return cmd.Val()
	}
}

func RedisUnlock(ctx context.Context, client *redis.Client, key string) {
	client.Del(ctx, key)
}