package publish_test

import (
	"context"
	"happyladysauce/publish/publisher"
	"happyladysauce/publish/subscriber"
	myRedis "happyladysauce/redis"
	"testing"
	"time"
)

func TestPublish(*testing.T) {
	ctx := context.Background()
	ctx1 := context.WithValue(ctx, "publisher_name", "publisher1")
	ctx2 := context.WithValue(ctx, "publisher_name", "publisher2")
	myRedis.InitRedis()

	channel1 := "channel1"
	channel2 := "channel2"

	// 启动第一批 subscriber
	// subscribe 需要提前启动好, 在此之前频道(channel)里的消息它接收不到
	ctx3 := context.WithValue(ctx, "subscriber_name", "subscriber3")
	ctx4 := context.WithValue(ctx, "subscriber_name", "subscriber4")

	go subscriber.Subscriber(ctx3, myRedis.Client, []string{channel1})
	go subscriber.Subscriber(ctx4, myRedis.Client, []string{channel2})
	time.Sleep(time.Second)
	
	go publisher.Publish(ctx1, myRedis.Client, channel1, "你好！channel1")
	go publisher.Publish(ctx2, myRedis.Client, channel2, "你好！channel2")
	time.Sleep(time.Second)
}



