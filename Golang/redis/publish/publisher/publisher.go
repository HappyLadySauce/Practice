package publisher

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func Publish(ctx context.Context, client *redis.Client, channel string, message any) {
	cmd := client.Publish(ctx, channel, message)
	if cmd.Err() == nil {
		n := cmd.Val()	// 订阅数量
		fmt.Printf("%v向频道%s发送了消息, 此时该频道有%d个订阅者\n", ctx.Value("publisher_name"), channel, n)
	} else {
		fmt.Printf("%v向频道%s发送消息失败%v\n", ctx.Value("publisher_name"), channel, cmd.Err())
	}
}

