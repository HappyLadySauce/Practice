package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0:n-1]
	return x
}

type MaxHeap []int

func (h MaxHeap) Len() int {
	return len(h)
}

func (h MaxHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0:n-1]
	return x
}

type PriorityQueue struct {
	items []int
	isMinHeap bool
}

func NewPriorityQueue(isMinHeap bool) *PriorityQueue {
	return &PriorityQueue{
		items: []int{},
		isMinHeap: isMinHeap,
	}
}

func parent(i int) int {
	return (i - 1) / 2
}

func leftChild(i int) int {
	return 2 * i + 1
}

func rightChild(i int) int {
	return 2 * i + 2
}

func (pq *PriorityQueue) IsEmpty() bool {
	return len(pq.items) == 0
}

func (pq *PriorityQueue) heapifyUp(i int) {
	for i > 0 {
		p := parent(i)

		if pq.isMinHeap {
			if pq.items[p] <= pq.items[i] {
				break
			}
			
		} else {
			if pq.items[p] >= pq.items[i] {
				break
			}
		}

		pq.items[p], pq.items[i] = pq.items[i], pq.items[p]
		i = p
	}
}

func (pq *PriorityQueue) heapifyDown(i int) {
	size := len(pq.items)
	for {
		smallest := i
		left := leftChild(i)
		right := rightChild(i)

		if pq.isMinHeap {
			if left < size && pq.items[left] < pq.items[smallest] {
				smallest = left
			}

			if right < size && pq.items[right] < pq.items[smallest] {
				smallest = right
			} 
		} else {
			if left < size && pq.items[left] > pq.items[smallest] {
				smallest = left
			}

			if right < size && pq.items[right] < pq.items[smallest] {
				smallest = right
			} 
		}

		if smallest == i {
			break
		}

		pq.items[i], pq.items[smallest] = pq.items[smallest], pq.items[i]
		i = smallest
	}
}

func (pq *PriorityQueue) Enqueue(item int) {
	pq.items = append(pq.items, item)
	pq.heapifyUp(len(pq.items)-1)
}

func (pq *PriorityQueue) Dequeue() int {
	if pq.IsEmpty() {
		return 0
	}

	top := pq.items[0]
	lastIndex := len(pq.items) - 1
	pq.items[0] = pq.items[lastIndex]
	pq.items = pq.items[:lastIndex]

	if len(pq.items) > 0 {
		pq.heapifyDown(0)
	}

	return top
}

func (pq *PriorityQueue) Peek() int {
	if pq.IsEmpty() {
		return 0
	}

	return pq.items[0]
}

func (pq *PriorityQueue) Size() int {
    return len(pq.items)
}


func main() {
    // 使用标准库实现
    fmt.Println("使用Go标准库heap包实现:")
    
    // 最小堆
    minHeap := &IntHeap{}
    heap.Init(minHeap)
    
    heap.Push(minHeap, 30)
    heap.Push(minHeap, 10)
    heap.Push(minHeap, 20)
    
    fmt.Println("最高优先级元素:", (*minHeap)[0]) // 输出: 10
    fmt.Println("出队的元素:", heap.Pop(minHeap)) // 输出: 10
    fmt.Println("队列大小:", minHeap.Len()) // 输出: 2
    
    // 最大堆
    maxHeap := &MaxHeap{}
    heap.Init(maxHeap)
    
    heap.Push(maxHeap, 30)
    heap.Push(maxHeap, 10)
    heap.Push(maxHeap, 20)
    
    fmt.Println("最高优先级元素:", (*maxHeap)[0]) // 输出: 30
    
    // 使用自定义实现
    fmt.Println("使用自定义优先队列实现:")
    
    // 最小堆
    minPQ := NewPriorityQueue(true)
    minPQ.Enqueue(30)
    minPQ.Enqueue(10)
    minPQ.Enqueue(20)
    
    top := minPQ.Peek()
    fmt.Println("最高优先级元素:", top) // 输出: 10
    
    popped := minPQ.Dequeue()
    fmt.Println("出队的元素:", popped) // 输出: 10
    fmt.Println("队列大小:", minPQ.Size()) // 输出: 2
    
    // 最大堆
    maxPQ := NewPriorityQueue(false)
    maxPQ.Enqueue(30)
    maxPQ.Enqueue(10)
    maxPQ.Enqueue(20)
    
    topMax := maxPQ.Peek()
    fmt.Println("最高优先级元素:", topMax) // 输出: 30
}



