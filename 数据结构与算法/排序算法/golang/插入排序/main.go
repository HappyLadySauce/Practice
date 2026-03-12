// 插入排序示例
package main

import "fmt"

func insertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i -1

		for j >= 0 && arr[j] > key {
			arr[j + 1] = arr[j]
			j--
		}

		arr[j+1] = key
	}
}

// 二分插入排序
func binaryInsertSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		key = arr[i]

		left, right := 0, i - 1
		for left <= right {
			mid := left + (right - left) / 2
			if arr[mid] > key {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		for k := i; j > left; j-- {
			arr[j] = arr[j - 1]
		}

		arr[left] = key
	}


} 

func main() {
	arr := []int{1, 4, 3, 5, 6, 7, 5, 832, 4, 523}
	fmt.Println("排序前", arr)
	insertionSort(arr)
	fmt.Println("排序后", arr)
	arr := []int{1, 4, 3, 5, 6, 7, 5, 832, 4, 523}
	fmt.Println("排序前", arr)
	binaryInsertSort(arr)
	fmt.Println("排序后", arr)
}
