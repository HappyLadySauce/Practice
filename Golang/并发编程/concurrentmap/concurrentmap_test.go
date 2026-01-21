package concurrentmap_test

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"happyladysauce/concurrentmap"
)

func TestConcurrentMap(*testing.T) {
	fmt.Println("this is TestConcurrentMap.")
	concurrentMap := concurrentmap.NewConcurrentMap[string, int](10)
	concurrentMap.Store("age", 18)
	age, exists := concurrentMap.Load("age")
	if exists != true {
		fmt.Println("not found age for concurrentMap.")
		os.Exit(1)
	}
	fmt.Printf("%v=%v\n", "age", age)
}

func TestForConcurrentMap(*testing.T) {
	fmt.Println("this is TestForConcurrentMap.")
	concurrentMap := concurrentmap.NewConcurrentMap[string, int](1000)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Go(func ()  {
			concurrentMap.Store(fmt.Sprintf("ID%d", i), i)
		})
	}

	wg.Wait()
	for i := 0; i < 100; i++ {
		wg.Go(func ()  {
			for j := 0; j < 1000; j++ {
				value, exists := concurrentMap.Load(fmt.Sprintf("ID%d", j))
				if exists != true {
					fmt.Printf("not found ID%d for concurrentMap.\n", j)
					return
				}
				fmt.Printf("concurrentMap[%v] ID%v=%v\n", i, j, value)
			}
		})
		time.Sleep(10 * time.Microsecond)
	}
	fmt.Printf("now goroutine number is: %d\n.", runtime.NumGoroutine())
	wg.Wait()
}