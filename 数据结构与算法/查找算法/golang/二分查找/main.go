package main

import "fmt"

func BinarySearch(arr []int, target int) int {
    left := 0
    right := len(arr) - 1

    for left <= right {
        mid := left + (right-left)/2   // 防止整数溢出

        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {  // 目标在右半区
            left = mid + 1
        } else {                       // 目标在左半区
            right = mid - 1
        }
    }
    return -1
}

func main() {
	arr := []int{0,1,2,3,5,6,7,8,9}
	fmt.Println(arr)
	if index := BinarySearch(arr, 6); index >= 0 {
		fmt.Println("find index:", index) 
	}
}