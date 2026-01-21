package limit

import (
	"time"
)

// 定义最多可以同时访问多少次接口
var qps = make(chan struct{}, 100)
var Count = 0


func Handler() {
	qps <- struct{}{}
	Count++
	defer func ()  {
		<- qps
		Count--	
	}()
	// TODO: 业务逻辑
	time.Sleep(time.Second)
}

// 协程限制器
type GoroutineLimit struct {
	limit	int
	ch 		chan struct{}
}

func NewGoroutineLimit(n int) *GoroutineLimit {
	return &GoroutineLimit{
		limit:	n,
		ch:		make(chan struct{}, n),
	}
}

func (g *GoroutineLimit) Run(f func()) {
	g.ch <- struct{}{}
	go func ()  {
		f()
		<- g.ch 
	}()
}