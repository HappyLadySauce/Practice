package main

import (
	"fmt"
)

type Data interface{}

type Node struct {
	data Data
	next *Node
}

type SinglyLinkedList struct {
	Head *Node
	num  int
}

func NewSinglyLinkedList() *SinglyLinkedList {
	return &SinglyLinkedList{nil, 0}
}

func (l *SinglyLinkedList) IsEmpty() bool {
	if l == nil {
		panic("SinglyLinkedList pointer is null")
	}
	return l.Head == nil
}

func (l *SinglyLinkedList) InsertAtHead(data Data) {
	if l == nil {
		panic("SinglyLinkedList pointer is null")
	}

	newNode := &Node{data, nil}

	if l.IsEmpty() {
		l.Head = newNode
	} else {
		newNode.next = l.Head
		l.Head = newNode
	}

	l.num++
}

func (l *SinglyLinkedList) InsertAtTail(data Data) {
	if l == nil {
		panic("SinglyLinkedList pointer is null")
	}

	newNode := &Node{data, nil}

	if l.IsEmpty() {
		l.Head = newNode
	} else {
		cur := l.Head
		for cur.next != nil {
			cur = cur.next
		}
		cur.next = newNode
	}

	l.num++
}

func (l *SinglyLinkedList) InsertAtPosition(data Data, position int) {
	if l == nil {
		panic("SinglyLinkedList pointer is null")
	}

	if position < 0 || position > l.num {
		panic("position invalid")
	}

	newNode := &Node{data, nil}

	if position == 0 {
		newNode.next = l.Head
		l.Head = newNode
	} else {
		prev := l.Head
		for position > 1 {
			prev = prev.next
			position--
		}

		newNode.next = prev.next 
		prev.next = newNode
	}

	l.num++
}

func (l *SinglyLinkedList) Delete(data Data) {
	if l == nil {
		panic("SinglyLinkedList pointer is null")
	}

	if l.IsEmpty() {
		return
	}

	if l.Head.data == data {
		l.Head = l.Head.next
		l.num--
	} else {
		prev := l.Head
		for prev.next != nil {
			if prev.next.data == data {
				prev.next = prev.next.next
				l.num--
				return
			}
			prev = prev.next
		}
	}
}

func (l* SinglyLinkedList) Search(data Data) bool {
	if l == nil {
		panic("SinglyLinkedList pointer is null")
	}

	if l.IsEmpty() {
		return false
	}

	if l.Head.data == data {
		return true
	} else {
		prev := l.Head
		for prev.next != nil {
			if prev.next.data == data {
				return true
			}
			prev = prev.next
		}
	}
	return false
}

func (l *SinglyLinkedList) Size() int {
	return l.num
}


func printList(l *SinglyLinkedList) {
    cur := l.Head
    for cur != nil {
        fmt.Printf("%v ", cur.data)
        cur = cur.next
    }
    fmt.Println()
}

func main() {
    list := NewSinglyLinkedList()

    fmt.Println("=== 测试插入 ===")
    list.InsertAtHead(10)
    list.InsertAtHead(20)
    list.InsertAtTail(30)
    list.InsertAtPosition(15, 2) // 在位置2插入15（0-based）
    fmt.Print("链表: ")
    printList(list) // 预期: 20 10 15 30
    fmt.Println("大小:", list.Size()) // 4

    fmt.Println("\n=== 测试搜索 ===")
    fmt.Println("搜索15:", list.Search(15)) // true
    fmt.Println("搜索99:", list.Search(99)) // false

    fmt.Println("\n=== 测试删除 ===")
    list.Delete(10)
    fmt.Print("删除10后: ")
    printList(list) // 20 15 30

    list.Delete(20) // 删除头节点
    fmt.Print("删除20后: ")
    printList(list) // 15 30

    list.Delete(30) // 删除尾节点
    fmt.Print("删除30后: ")
    printList(list) // 15

    list.Delete(15)
    fmt.Print("删除15后: ")
    printList(list) // (空)
    fmt.Println("是否为空:", list.IsEmpty()) // true
    fmt.Println("大小:", list.Size()) // 0

    fmt.Println("\n=== 边界测试 ===")
    list.InsertAtHead(100)
    list.InsertAtTail(200)
    list.InsertAtPosition(150, 1)
    fmt.Print("链表: ")
    printList(list) // 100 150 200

    // 删除不存在的元素（无影响）
    list.Delete(999)
    fmt.Print("删除不存在的999后: ")
    printList(list) // 100 150 200

    // 空链表删除
    emptyList := NewSinglyLinkedList()
    emptyList.Delete(1) // 不会panic
    fmt.Println("空链表删除后大小:", emptyList.Size())
}
