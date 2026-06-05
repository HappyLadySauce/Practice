
package main

import (
	"fmt"
)

type Data interface{}

type Stack struct {
	items []Data
}

func NewStack() *Stack {
	return &Stack{}
}

func (stack *Stack) IsEmpty() bool {
	return len(stack.items) == 0
}

func (stack *Stack) Push(item Data) {
	stack.items = append(stack.items, item)
}

func (stack *Stack) Pop() Data {
	if stack.IsEmpty() {
		return nil
	}

	n := len(stack.items)
	item := stack.items[n - 1]
	stack.items = stack.items[:n - 1]
	return item
}

func (stack *Stack) Peek() Data {
	if stack.IsEmpty() {
		return nil
	}
	n := len(stack.items)
	item := stack.items[n - 1]
	return item
}

func (stack *Stack) Size() int {
    return len(stack.items)
}

func main() {
    stack := NewStack()
    
    stack.Push(10)
    stack.Push(20)
    stack.Push(30)
    
    top := stack.Peek()
    fmt.Println("栈顶元素:", top) // 输出: 30
    
    popped := stack.Pop()
    fmt.Println("弹出的元素:", popped) // 输出: 30
    
    fmt.Println("栈大小:", stack.Size()) // 输出: 2
    fmt.Println("栈是否为空:", stack.IsEmpty()) // 输出: false
}