
package main

func shellSort(arr []int) []int {
	n = len(arr)



	for gap := n/2; gap > 0; gap /= 2 {
		for i := gap; i < n; i++ {
			key := arr[i]
			j = i

			for j >= gap && arr[j - gap] > key {
				arr[j] = arr[j - gap]
				j -= gap
			}

			arr[j] = key
		}
	}

}

func main() {
	arr := []int{1,2,5,54,32,1,46,21,512,5,127,542,145}
	fmt.Println("排序前", arr)
	shellSort(arr)
	fmt.Println("排序后", arr)
}