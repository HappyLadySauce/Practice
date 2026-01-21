package main

import (
	"fmt"
	"slices"
)

// 数组
func test1() {
	// 数组在声明时长度只能是一个常量，不能是变量，你不能在声明一个变量让后用变量作为数组的长度值
	// 在数组初始化时，需要注意的是，长度必须为一个常量表达式，否则将无法通过编译，常量表达式即表达式的最终结果是一个常量
	var a [5]int

	// 数组初始化
	b := [5]int{1,2,3}
	c := [...]int{1,2,3,4}

	// 获得一个数组指针
	d := new([5]int)
	
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)

	// 获取数组的长度
	fmt.Println(len(a))
	// 获取数组的容量
	fmt.Println(cap(a))
}

// 切片
func test2() {
	// 切割数组的格式为arr[startIndex:endIndex]，切割的区间为左闭右开
	nums := [5]int{1,2,3,4,5}
	// 切割数组时，startIndex可以省略，默认从0开始
	fmt.Println(nums[1:])
	// 切割数组时，endIndex可以省略，默认到数组的最后一个元素
	fmt.Println(nums[:3])

	// 数组在切割后，就会变为切片类型
	fmt.Printf("%T\n",nums)
	fmt.Printf("%T\n", nums[1:3])

	// 若要将数组转换为切片类型，不带参数进行切片即可，转换后的切片与原数组指向的是同一片内存，修改切片会导致原数组内容的变化
	s := nums[:]
	fmt.Printf("%T\n",s)
	// 如果要对转换后的切片进行修改，建议使用下面这种方式进行转换
	slice := slices.Clone(nums[:])
	slice[0] = 0
	fmt.Printf("%v\n",slice)
	// 可以看到，修改切片并不会影响到原数组
	fmt.Printf("%v\n",nums)
}

// 切片的初始化
func test3() {
	// 切片在 Go 中的应用范围要比数组广泛的多，它用于存放不知道长度的数据，且后续使用过程中可能会频繁的插入和删除元素。
	var nums1 []int // 这种方式声明的切片，默认值为nil，所以不会为其分配内存
	nums2 := []int{1,2,3,4,5}
	fmt.Printf("%v\n",nums1)
	fmt.Printf("%v\n",nums2)
	// 通常情况下，推荐使用make来创建一个空切片，只是对于切片而言，make函数接收三个参数：类型，长度，容量。
	nums3 := make([]int, 0 ,0) // 使用make进行初始化时，建议预分配一个足够的容量，可以有效减少后续扩容的内存消耗
	nums4 := new([]int)
	fmt.Printf("%v\n",nums3)
	fmt.Printf("%v\n",nums4)
}

// append
func test4() {
	// 切片可以通过append函数实现许多操作，函数签名如下，slice是要添加元素的目标切片，elems是待添加的元素，返回值是添加后的切片
	// func append(slice []Type, elems ...Type) []Type
	nums := make([]int, 0, 0)
	nums = append(nums, 1,2,3,4,5)
	fmt.Println(len(nums), cap(nums))

	// 切片元素的插入
	// 头插
	nums = append([]int{-1,0}, nums...)
	fmt.Println(nums)

	// 从中间下标 i 插入元素
	nums = append(nums[:2+1], append([]int{999,999}, nums[2+1:]...)...)
	fmt.Println(nums)

	// 尾插
	nums = append(nums, 6,7,8,9,10)
	fmt.Println(nums)

	// 切片的删除
	// 删除头元素
	nums = nums[1:]
	fmt.Println(nums)

	// 删除中间元素
	nums = append(nums[:2], nums[3:]...)
	fmt.Println(nums)

	// 删除尾元素
	nums = nums[:len(nums)-1]
	fmt.Println(nums)

	for index, val := range nums {
		fmt.Printf("%d:%d ", index, val)
	}
}

// 多维切片
func test5() {
	var nums [5][5]int
	for _, num := range nums {
		fmt.Println(num)
	}
	slices := make([][]int, 5)
	for _, slice := range slices {
		fmt.Println(slice)
	}

	// 二维切片的初始化
	for i := range slices {
		slices[i] = make([]int, 5)
	}
	for _, slice := range slices {
		fmt.Println(slice)
	}
}

func main() {
	// test1()
	// test2()
	// test3()
	// test4()
	test5()
}