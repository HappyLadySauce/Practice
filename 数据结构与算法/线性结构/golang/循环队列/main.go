package main

import (
	"fmt"
)

type MyCircularQueue struct {
    items	[]int
	front	int
	rear	int
	size	int
	capacity int
}


func Constructor(k int) MyCircularQueue {
    return MyCircularQueue{
		items: make([]int, k),
		front:  0,
		rear:  0,
		size:  0,
		capacity: k,
	}
}


func (this *MyCircularQueue) EnQueue(value int) bool {
	if this.IsFull() {
		return false
	}

	this.items[this.rear] = value
	this.rear = (this.rear + 1) % (this.capacity)
	this.size++
	return true
}


func (this *MyCircularQueue) DeQueue() bool {
    if this.IsEmpty() {
		return false
	}

	this.front = (this.front + 1) % this.capacity
	this.size--
	return true
}


func (this *MyCircularQueue) Front() int {
    if this.IsEmpty() {
		return -1
	}

	return this.items[this.front]
}


func (this *MyCircularQueue) Rear() int {
    if this.IsEmpty() {
		return -1
	}

	return this.items[(this.rear - 1 + this.capacity) % this.capacity]
}


func (this *MyCircularQueue) IsEmpty() bool {
	return this.size == 0
}


func (this *MyCircularQueue) IsFull() bool {
	return this.size == this.capacity
}

/**
 * Your MyCircularQueue object will be instantiated and called as such:
 * obj := Constructor(k);
 * param_1 := obj.EnQueue(value);
 * param_2 := obj.DeQueue();
 * param_3 := obj.Front();
 * param_4 := obj.Rear();
 * param_5 := obj.IsEmpty();
 * param_6 := obj.IsFull();
 */

func main() {
    // 测试示例
    fmt.Println("=== 测试官方示例 ===")
    circularQueue := Constructor(3)
    
    fmt.Println(circularQueue.EnQueue(1)) // 返回 true
    fmt.Println(circularQueue.EnQueue(2)) // 返回 true
    fmt.Println(circularQueue.EnQueue(3)) // 返回 true
    fmt.Println(circularQueue.EnQueue(4)) // 返回 false，队列已满
    fmt.Println(circularQueue.Rear())     // 返回 3
    fmt.Println(circularQueue.IsFull())   // 返回 true
    fmt.Println(circularQueue.DeQueue())  // 返回 true
    fmt.Println(circularQueue.EnQueue(4)) // 返回 true
    fmt.Println(circularQueue.Rear())     // 返回 4
    
    fmt.Println("\n=== 测试边界情况 ===")
    // 测试空队列
    emptyQueue := Constructor(2)
    fmt.Println("空队列 Front:", emptyQueue.Front())  // -1
    fmt.Println("空队列 Rear:", emptyQueue.Rear())    // -1
    fmt.Println("空队列 DeQueue:", emptyQueue.DeQueue()) // false
    
    fmt.Println("\n=== 测试循环利用空间 ===")
    cycleQueue := Constructor(3)
    cycleQueue.EnQueue(10)
    cycleQueue.EnQueue(20)
    cycleQueue.EnQueue(30)
    fmt.Println("队列已满:", cycleQueue.IsFull()) // true
    fmt.Println("出队:", cycleQueue.DeQueue())    // true
    fmt.Println("入队40:", cycleQueue.EnQueue(40)) // true
    fmt.Println("Front:", cycleQueue.Front())    // 20
    fmt.Println("Rear:", cycleQueue.Rear())      // 40
    
    fmt.Println("\n=== 测试多次出入队 ===")
    multiQueue := Constructor(2)
    multiQueue.EnQueue(1)
    multiQueue.EnQueue(2)
    fmt.Println("Front:", multiQueue.Front()) // 1
    fmt.Println("Rear:", multiQueue.Rear())   // 2
    multiQueue.DeQueue()
    multiQueue.EnQueue(3)
    fmt.Println("Front:", multiQueue.Front()) // 2
    fmt.Println("Rear:", multiQueue.Rear())   // 3
    multiQueue.DeQueue()
    multiQueue.EnQueue(4)
    fmt.Println("Front:", multiQueue.Front()) // 3
    fmt.Println("Rear:", multiQueue.Rear())   // 4
}