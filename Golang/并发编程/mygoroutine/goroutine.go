package mygoroutine

import (
	"sync"
	"fmt"
	"runtime"
)

func NewGoroutine() {
	fmt.Printf("本机最多可使用的逻辑CPU数量: %v\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)		// 设置程序最多能使用的核心数
	fmt.Printf("现在的goroutine数量: %v\n", runtime.NumGoroutine())
}

func OldWaitGroup() {
	const n = 10
	w := sync.WaitGroup{}
	w.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer w.Done()
			fmt.Printf("this is %v goroutine.\n", i)
		}()

	}
	fmt.Printf("当前的协程数量为: %v\n", runtime.NumGoroutine())
	w.Wait()
	fmt.Printf("当前的协程数量为: %v\n", runtime.NumGoroutine())
}

// go 1.25 sync.WaitGroup 新写法
func NewWaitGroup() {
	const n = 10
	var w sync.WaitGroup
	for i := 0; i < n; i++ {
		w.Go(func() {
			fmt.Printf("this is %v goroutine.\n", i)
		}) 
	}
	fmt.Printf("当前的协程数量为: %v\n", runtime.NumGoroutine())
	w.Wait()
	fmt.Printf("当前的协程数量为: %v\n", runtime.NumGoroutine())
}