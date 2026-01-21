package subscriber

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func Subscriber(ctx context.Context, client *redis.Client, channel []string) {
	ps := client.Subscribe(ctx, channel...)	// 返回一个订阅者
	defer ps.Close()	// 函数结束关闭订阅者

	for {
		if msg, err := ps.ReceiveMessage(ctx); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("%s从频道%s里接收到消息:%s\n", ctx.Value("subscriber_name"), msg.Channel, msg.Payload)
		}
	}
}