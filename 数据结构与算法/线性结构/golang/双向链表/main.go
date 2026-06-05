package main

import (

)

type Data interface{}

type Node struct {
	data Data
	prev *Node
	next *Node
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
	size int
}

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (list *DoublyLinkedList) IsEmpty() bool {
	return list.head == nil
}

func (list *DoublyLinkedList) InsertAt(data Data, position int) *Node {
    if list == nil {
        panic("DoublyLinkedList pointer is null")
    }
    if position < 0 || position > list.size {
        panic("invalid position")
    }

    newNode := &Node{data: data, prev: nil, next: nil}

    // 情况1：空链表
    if list.IsEmpty() {
        list.head = newNode
        list.tail = newNode
        list.size++
        return newNode
    }

    // 情况2：插入头部
    if position == 0 {
        newNode.next = list.head
        list.head.prev = newNode
        list.head = newNode
        list.size++
        return newNode
    }

    // 情况3：插入尾部
    if position == list.size-1 {
        newNode.prev = list.tail
        list.tail.next = newNode
        list.tail = newNode
        list.size++
        return newNode
    }

    // 情况4：插入中间 – 选择从头部或尾部开始遍历以缩短一半时间
    var cur *Node
    if position <= list.size/2 {
        // 从头向后找
        cur = list.head
        for i := 0; i < position; i++ {
            cur = cur.next
        }
    } else {
        // 从尾向前找
        cur = list.tail
        for i := list.size-1; i > position; i-- {
            cur = cur.prev
        }
    }
    // 此时 cur 指向原位置上的节点，新节点应插入在 cur 之前
    newNode.prev = cur.prev
    newNode.next = cur
    cur.prev.next = newNode
    cur.prev = newNode

    list.size++
    return newNode
}

func (list *DoublyLinkedList) DeleteAtPosition(position int) {
    if list == nil {
        panic("DoublyLinkedList pointer is null")
    }
    if position < 0 || position > list.size {
        panic("invalid position")
    }

	if list.IsEmpty() {
		return
	}

	if list.size == 1 {
        list.head = nil
        list.tail = nil
        list.size--
        return
    }

    if position == 0 {
        list.head = list.head.next
        list.head.prev = nil
        list.size--
        return
    }

    if position == list.size-1 {
        list.tail = list.tail.prev
        list.tail.next = nil
        list.size--
        return
    }

	var cur *Node
	if position <= list.size/2 {
		cur = list.head
		for i := 0; i < position; i++ {
			cur = cur.next
		}
	} else {
		cur = list.tail
		for i := list.size-1; i > position; i-- {
			cur = cur.prev
		}
	}

	cur.prev.next = cur.next
	cur.next.prev = cur.prev
	list.size--
}

