package main

import (
	"fmt"
)

// 在 Go 中，有仅有一种循环语句：for，Go 抛弃了while语句，for语句可以被当作while来使用。
func test1() {
	// for 循环控制
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	// while 循环
	var n int = 1
	for n < 3 {
		fmt.Println(n)
		n++
	}
}

// for range 可以更加方便的遍历一些可迭代的数据结构，如数组，切片，字符串，映射表，通道。
func test2() {
	var iterable = []int{1, 2, 3, 4}
	// index为可迭代数据结构的索引，value则是对应索引下的值，例如使用for range遍历一个字符串。
	for index, value := range iterable {
		fmt.Println(index, value)
	}

	n := 10
	for i := range n {
		fmt.Println(i)
	}
}

func main() {
	// test1()
	test2()
}
