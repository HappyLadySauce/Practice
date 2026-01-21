package main

import (
	"fmt"
)

// 类型形参：T 就是一个类型形参，形参具体是什么类型取决于传进来什么类型
// 类型约束：int | float64 构成了一个类型约束，这个类型约束内规定了哪些类型是允许的，约束了类型形参的类型范围
// 类型实参：Sum[int](1,2)，手动指定了 int 类型，int 就是类型实参。
func Sum[T int | float64](a, b T) T {
	return a + b
}

func test1() {
	// 显示定义类型
	fmt.Println(Sum[int](1, 2))
	fmt.Println(Sum[float64](1.0, 2.0))
	// 编译器类型推导
	fmt.Println(Sum(1.0, 2.0))
}

func main() {
	test1()
}