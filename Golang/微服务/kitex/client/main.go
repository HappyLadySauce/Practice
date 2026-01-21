package main

import (
	"context"
	"happyladysauce/kitex_gen/math"
	"happyladysauce/kitex_gen/math/mathservice"

	"github.com/cloudwego/kitex/client"
)

func main() {
	client, err := mathservice.NewClient("math", client.WithHostPorts("127.0.0.1:8083"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	req := &math.AddRequest{
		Left:  1,
		Right: 2,
	}
	resp, err := client.Add(ctx, req)
	if err != nil {
		panic(err)
	}
	println(resp.Result)

	ctx2 := context.Background()
	req2 := &math.SubRequest{
		Left:  1,
		Right: 2,
	}
	resp2, err2 := client.Sub(ctx2, req2)
	if err2 != nil {
		panic(err2)
	}
	println(resp2.Result)
}
