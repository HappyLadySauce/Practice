package lock_test

import (
    "context"
    "strconv"
    "testing"
    "time"

    lock "happyladysauce/lock"
    myRedis "happyladysauce/redis"
)

func requireRedis(t *testing.T) {
    t.Helper()
    myRedis.InitRedis()
    if myRedis.Client == nil {
        t.Skip("Redis 客户端未初始化，跳过测试")
    }
    if err := myRedis.Client.Ping(context.Background()).Err(); err != nil {
        t.Skipf("Redis 无法连接，跳过测试：%v", err)
    }
}

func TestRedisLockUnlock(t *testing.T) {
    requireRedis(t)
    ctx := context.Background()
    key := "test_lock_" + strconv.FormatInt(time.Now().UnixNano(), 10)
    expire := 3 * time.Second

    if got := lock.RedisLock(ctx, myRedis.Client, expire, key); !got {
        t.Fatalf("首次加锁应该成功")
    }
    if got := lock.RedisLock(ctx, myRedis.Client, expire, key); got {
        t.Fatalf("已加锁的 key 不应再次成功加锁")
    }

    lock.RedisUnlock(ctx, myRedis.Client, key)

    if got := lock.RedisLock(ctx, myRedis.Client, expire, key); !got {
        t.Fatalf("解锁后再次加锁应该成功")
    }
}

func TestRedisLockExpire(t *testing.T) {
    requireRedis(t)
    ctx := context.Background()
    key := "test_lock_expire_" + strconv.FormatInt(time.Now().UnixNano(), 10)
    expire := 1 * time.Second

    if got := lock.RedisLock(ctx, myRedis.Client, expire, key); !got {
        t.Fatalf("首次加锁应该成功")
    }

    time.Sleep(expire + 500*time.Millisecond)

    if got := lock.RedisLock(ctx, myRedis.Client, expire, key); !got {
        t.Fatalf("过期后应可再次加锁")
    }

    lock.RedisUnlock(ctx, myRedis.Client, key)
}

