package limit_test

import (
	"fmt"
	"happyladysauce/limit"
	"sync"
	"testing"
	"time"
)

func TestLimit(*testing.T) {
	go func () {
		for {
			fmt.Printf("当前的接口访问数量: %v\n", limit.Count)
			time.Sleep(time.Second)
		}
	}()

	const P = 1000
	wg := sync.WaitGroup{}
	wg.Add(P)

	for i := 0; i < P; i++{
		go func () {
			defer wg.Done()
			limit.Handler()
		}()
	}
	wg.Wait()
}

func testlimit() {
	time.Sleep(time.Second)
}

func TestGoroutineLimit(*testing.T) {
	go func () {
		for {
			fmt.Printf("当前的接口访问数量: %v\n", limit.Count)
			time.Sleep(time.Second)
		}
	}()

	const P = 1000
	wg := sync.WaitGroup{}
	wg.Add(P)

	limit := limit.NewGoroutineLimit(100)

	for i := 0; i < P; i++{
		go func () {
			defer wg.Done()
			limit.Run(testlimit)
		}()
	}
	wg.Wait()
}