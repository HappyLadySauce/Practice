package main

import (
	"fmt"
)

func bubbleSort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n - 1; i++ {
		for j := 0; j < n - 1 - 1; j++ {
			if arr[j] > arr[j + 1] {
				arr[j], arr[j + 1] = arr[j + 1], arr[j]
			}
		}
	}
	return arr
}


func main() {
	arr := []int{4,3,436,7,657,657,568}
	fmt.Println("排序前", arr)
	bubbleSort(arr)
	fmt.Println("排序后", arr)
}