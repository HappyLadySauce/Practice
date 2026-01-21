package lock

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type TestLock struct {
	n atomic.Int32
	mu sync.Mutex	// 设置一个写锁
	
	// 读写锁中的读锁可以被多个协程同时持有; 而写锁只能被一个协程持有,且此时会对其他锁进行阻塞,直到写锁释放. 
	// ru sync.RWMutex	// 设置一个读写锁
}

func NewTestLock() *TestLock {
	return &TestLock{}
}

func (l *TestLock) Add(n int32) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.n.Add(n)
	time.Sleep(1 * time.Millisecond)
	fmt.Printf("n = %v\n", l.n.Load())
}

func RunLock() {
	l := NewTestLock()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Go(func ()  {
			l.Add(1)
		})
	}
	fmt.Printf("当前携程数量为: %v\n", runtime.NumGoroutine())
	wg.Wait()
}