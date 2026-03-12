package main

import (
	"fmt"
)

func heapSort(arr []int) {
	n := len(arr)
	// 构建大根堆
	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}

	// 依次交换堆顶与未排序部分的最后一个元素，并调整堆
	for end := n - 1; end > 0; end-- {
		arr[0], arr[end] = arr[end], arr[0]
		heapify(arr, end, 0)
	}
}

func heapify(arr []int, n, root int) {
	largest := root
	left := 2*root + 1
	right := 2*root + 2

	if left < n && arr[left] > arr[largest] {
		largest = left
	}

	if right < n && arr[right] > arr[largest] {
		largest = right
	}

	if largest != root {
		arr[root], arr[largest] = arr[largest], arr[root]
		heapify(arr, n, largest)
	}
}

func main() {
	arr := []int{1, 4, 3, 5, 6, 7, 5, 832, 4, 523}
	fmt.Println("排序前", arr)
	heapSort(arr)
	fmt.Println("排序后", arr)
}