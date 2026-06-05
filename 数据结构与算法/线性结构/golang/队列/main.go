package main

import (
	"fmt"
)

type Data interface{}

type Queue struct {
	items []Data
}

func NewQueue() *Queue {
	return &Queue{}
}

func (queue *Queue) IsEmpty() bool {
	return len(queue.items) == 0
}

func (queue *Queue) Enqueue(item Data) {
	queue.items = append(queue.items, item)
}

func (queue *Queue) Dequeue() Data {
	if queue.IsEmpty() {
		return nil
	}

	item := queue.items[0]
	queue.items = queue.items[1:]
	return item
}

func (queue *Queue) Front() Data {
	if queue.IsEmpty() {
		return nil
	}

	return queue.items[0]
}

func (queue *Queue) Size() int {
	return len(queue.items)
}


func main() {
    queue := NewQueue()
    
    queue.Enqueue(10)
    queue.Enqueue(20)
    queue.Enqueue(30)
    
    front := queue.Front()
    fmt.Println("队头元素:", front) // 输出: 10
    
    dequeued := queue.Dequeue()
    fmt.Println("出队的元素:", dequeued) // 输出: 10
    
    fmt.Println("队列大小:", queue.Size()) // 输出: 2
    fmt.Println("队列是否为空:", queue.IsEmpty()) // 输出: false
}