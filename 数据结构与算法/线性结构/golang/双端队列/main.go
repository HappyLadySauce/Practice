package main

import (
	"fmt"
	"container/list"
)

type Data interface{}

type Deque struct {
	list *list.List
}

func NewDeque() *Deque {
	return &Deque{
		list: list.New(),
	}
}

func (deque *Deque) AddFirst(data Data) {
	deque.list.PushFront(data)
}

func (deque *Deque) AddLast(data Data) {
	deque.list.PushBack(data)
}

func (deque *Deque) IsEmpty() bool {
	return deque.list.Len() == 0
}

func (deque *Deque) RemoveFirst() Data {
	if deque.IsEmpty() {
		return nil
	}

	return deque.list.Remove(deque.list.Front())
}

func (deque *Deque) RemoveLast() Data {
	if deque.IsEmpty() {
		return nil
	}

	return deque.list.Remove(deque.list.Back())
}

// 查看前端元素
func (d *Deque) PeekFirst() interface{} {
    if d.IsEmpty() {
        return nil
    }
    return d.list.Front().Value
}

// 查看后端元素
func (d *Deque) PeekLast() interface{} {
    if d.IsEmpty() {
        return nil
    }
    return d.list.Back().Value
}

// 获取双端队列大小
func (d *Deque) Size() int {
    return d.list.Len()
}

func main() {
    deque := NewDeque()
    
    deque.AddFirst(10)
    deque.AddFirst(20)
    deque.AddLast(30)
    deque.AddLast(40)
    
    fmt.Println("前端元素:", deque.PeekFirst())  // 输出: 20
    fmt.Println("后端元素:", deque.PeekLast())   // 输出: 40
    
    fmt.Println("从前端移除:", deque.RemoveFirst())  // 输出: 20
    fmt.Println("从后端移除:", deque.RemoveLast())   // 输出: 40
    
    fmt.Println("队列大小:", deque.Size())  // 输出: 2
}
