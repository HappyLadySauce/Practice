package main

import (
	"fmt"
)

func quickSort(arr []int, low, high int) {
	if low < high {
		pivotIndex := partition(arr, low, high)

		quickSort(arr, low, pivotIndex - 1)
		quickSort(arr, pivotIndex + 1, high)
	}
}

func partition(arr []int, low, high int) int {
	// 选择基准
	pivot := arr[low]
	i := low + 1

	for j := low + 1; j <= high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[low], arr[i - 1] = arr[i - 1], arr[low]

	return i - 1
}


func main() {
	arr := []int{1, 4, 3, 5, 6, 7, 5, 832, 4, 523}
	fmt.Println("排序前", arr)
	quickSort(arr, 0, len(arr) - 1)
	fmt.Println("排序后", arr)
}