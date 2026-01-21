package main

import (
	"fmt"
)

// Exp函数的返回值是一个函数，这里将称成为grow函数，每将它调用一次，变量e就会以指数级增长一次。grow函数引用了Exp函数的两个变量：e和n，它们诞生在Exp函数的作用域内，在正常情况下随着Exp函数的调用结束，这些变量的内存会随着出栈而被回收。但是由于grow函数引用了它们，所以它们无法被回收，而是逃逸到了堆上，即使Exp函数的生命周期已经结束了，但变量e和n的生命周期并没有结束，在grow函数内还能直接修改这两个变量，grow函数就是一个闭包函数。
func Exp(n int) func() int {
	e := 1
	return func() int {
		temp := e
		e *= n
		return temp
	}
}

func test1() {
	grow := Exp(2)
	for i := range 10 {
		fmt.Printf("2^%d=%d\n", i, grow())
	}
}

// 延迟调用
func test2() {
	// defer关键字可以使得一个函数延迟一段时间调用
	defer fmt.Println("defer")
	fmt.Println("test2")

	//  当有多个 defer 描述的函数时，就会像栈一样先进后出的顺序执行。
	defer fmt.Println("defer2")
	fmt.Println("test2 end")

	// 延迟调用通常用于释放文件资源，关闭网络连接等操作，还有一个用法是捕获panic，不过这是错误处理一节中才会涉及到的东西。
}

func Fn1() int {
	fmt.Println("2")
	return 1
}

func test3() {
	// go 不会等到最后才去调用sum函数，sum函数早在延迟调用被执行以前就被调用了，并作为参数传递了fmt.Println。总结就是，对于defer直接作用的函数而言，它的参数是会被预计算的，这也就导致了第一个例子中的奇怪现象，对于这种情况，尤其是在延迟调用中将函数返回值作为参数的情况尤其需要注意。
	defer fmt.Println(Fn1())
	fmt.Println("3")
}

func main() {
	// test1()
	// test2()
	test3()
}

