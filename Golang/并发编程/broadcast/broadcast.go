package broadcast

import (
	"fmt"
	"sync"
	"time"
)

// 广播多等一
func Broadcast() {
	ch := make(chan struct{})

	const P = 3
	for i := 0; i < P; i++ {
		go func ()  {
			<-ch	// 此时 ch 为空, 触发阻塞
			fmt.Printf("%d 出发了\n", i)	// 这里可以编写具体的逻辑
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println("大伙可以出发了")
	close(ch)	// 关闭管道之后, ch 取消阻塞立马执行协程

	time.Sleep(time.Second)
}

// 一等多 WaitGroup
func CutDownLatch() {
	const P = 3
	ch := make(chan struct{}, P)
	for i := 0; i < P; i++ {
		go func ()  {
			time.Sleep(time.Duration(i) * time.Second)
			fmt.Printf("%d 完成工作了\n", i)
			ch <- struct{}{}
		}()
	}

	for i := 0; i < P; i++ {
		<-ch
	}

	fmt.Println("其他人都执行完毕，我要开始了")
}

func CondSignal() {
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func ()  {
		defer wg.Done()
		time.Sleep(time.Second)
		cond.Signal()	// 只有一个协程可以收到信号
		fmt.Printf("发送信号")
		
		time.Sleep(time.Second)
		cond.Signal()	// 只有一个协程可以收到信号
		fmt.Printf("发送信号")
	}()

	go func ()  {
		defer wg.Done()
		cond.L.Lock()	// 首先获得锁
		cond.Wait()		// Cond在调用Wait()时会释放锁,所以在调用Wait()前需要先执行Lock()
						// 而且Cond必须先Wait再接收Signal信号,不然需要重新发送Signal,否则会发生死锁
		fmt.Println("收到信号，开始执行工作")
		cond.L.Unlock()

		cond.L.Lock()	// 首先获得锁
		cond.Wait()		// Cond在调用Wait()时会释放锁,所以在调用Wait()前需要先执行Lock()
		fmt.Println("收到信号，开始执行工作")
		cond.L.Unlock()	// 释放锁
	}()

	wg.Wait()
}


