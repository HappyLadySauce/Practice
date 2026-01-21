package context

import (
	"context"
	"fmt"
	"time"
)

type Context interface {
	// Deadline() 用来记录到期时间，以及是否到期
	Deadline() (deadline time.Time, ok bool)
	// Done() 返回一个只读管道,且管道里不存放任何元素(struct{}),
	// 所以用这个管道仅仅是为了实现阻塞
	Done() <- chan struct{}
	// Err() 用来记录Done()管道关闭的原因,比如可能是因为超时,
	// 也可能是因为超时,也可能是因为被强行Cancel了
	Err() error
	// Value() 用来返回Key对应的value
	Value(key any) any
}

func Timeout1() {
	// ctx 为上下文, 调用cancel可以关闭ctx.Done()管道
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func ()  {
		time.Sleep(500 * time.Microsecond)
		done <- struct{}{}
	}()

	go func ()  {
		time.Sleep(200 * time.Microsecond)
		cancel()
	}()

	select {
	case <-done:
		fmt.Println("业务函数调用未超时.")
	case <-ctx.Done():
		err := ctx.Err()	// 获取错误,可能是Cancel了,也可能是超时了
		fmt.Printf("业务函数超时了: %v\n", err)
	}
}

func Timeout2() {
	// 自带超时时间的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	done := make(chan struct{})

	go func ()  {
		time.Sleep(500 * time.Microsecond)
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Println("业务函数调用未超时.")
	case <-ctx.Done():
		err := ctx.Err()	// 获取错误,可能是Cancel了,也可能是超时了
		fmt.Printf("业务函数超时了: %v\n", err)
	}
}

// 根据Go语言官方推荐,context value传递的类型最好是自定义类型
// 这样可以提供更好的可读性
type StringKey 		string
type StringValue	string
type IntValue		int

// 用 context携带value仅用于跨进程传输数据和API调用, 不用向函数传递可选参数
func RoutineID() {
	for i := 0; i < 3; i++ {
		ctx := context.WithValue(context.Background(), StringKey("gid"), IntValue(i))
		ctx = context.WithValue(ctx, StringKey("owner"), StringValue("HappyLadySauce"))
		go func (ctx context.Context)  {
			if gid, ok := ctx.Value(StringKey("gid")).(IntValue); ok {
				fmt.Printf("本协程ID为: %v\n", gid)
			}
			if owner, ok := ctx.Value(StringKey("owner")).(StringValue); ok {
				fmt.Printf("本协程owner为: %v\n", owner)
			}
		}(ctx)
	}
}




























