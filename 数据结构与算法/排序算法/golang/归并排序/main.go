package main

import (
	"fmt"
)

func mergeSort(arr []int, left, right int) {
	if left < right {
		mid := left + (right - left) / 2

		mergeSort(arr, left, mid)
		mergeSort(arr, mid + 1, right)

		merge(arr, left, mid, right)
	}
}


func merge(arr []int, left, mid, right int) {
	n1 := mid - left + 1
	n2 := right - mid

	L := make([]int, n1)
	R := make([]int, n2)

	for i := 0; i < n1; i++ {
		L[i] = arr[left + i]
	}
	for j := 0; j < n2; j++ {
		R[j] = arr[mid + 1 + j]
	}

	i, j := 0, 0
	k := left

	for i < n1 && j < n2 {
		if L[i] <= R[j] {
			arr[k] = L[i]
			i++
		} else {
			arr[k] = R[j]
			j++
		}
		k++
	}

	for i < n1 {
		arr[k] = L[i]
		i++
		k++
	}

	for j < n2 {
		arr[k] = R[j]
		j++
		k++
	}
}


func main() {
	arr := []int{1, 4, 3, 5, 6, 7, 5, 832, 4, 523}
	fmt.Println("排序前", arr)
	mergeSort(arr, 0, len(arr) - 1)
	fmt.Println("排序后", arr)
}