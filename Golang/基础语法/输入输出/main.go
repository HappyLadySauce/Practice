package main

import (
	"bufio"
	"fmt"
	"os"
)

// test1 函数演示了正确的缓冲写入器使用方式
func test1() {
	// 直接输出到标准输出，立即显示
	println("test1")
	
	// 使用 fmt.Println 输出，会自动添加换行符，立即显示
	fmt.Println("this is error")
	
	// 创建一个带缓冲的写入器，写入的内容会先存储在缓冲区
	writer := bufio.NewWriter(os.Stdout)
	
	// defer 语句：函数返回前执行，确保缓冲区内容被刷新到标准输出
	defer writer.Flush()
	
	// 向缓冲区写入字符串，此时不会立即显示
	writer.WriteString("hello world\n")  // 添加换行符
}

// test2 函数演示了缓冲写入器的另一种使用方式
func test2() {
	// 直接输出到标准输出，立即显示
	println("test2")
	
	// 使用 fmt.Println 输出，会自动添加换行符，立即显示
	fmt.Println("this is error")
	
	// 创建一个带缓冲的写入器
	writer := bufio.NewWriter(os.Stdout)
	
	// 向缓冲区写入字符串，此时不会立即显示
	writer.WriteString("hello world\n")  // 添加换行符
	
	// defer 语句：函数返回前执行，将缓冲区内容刷新到标准输出
	defer writer.Flush()
}

// bufio 与 fmt 包结合
func test3() {
	// 直接输出到标准输出，立即显示
	println("test3")
	
	// 使用 fmt.Println 输出，会自动添加换行符，立即显示
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	// 使用 fmt.Fprintln 输出，会自动添加换行符，立即显示
	fmt.Fprintln(writer, "Hello World")
}


func main() {
	test1()  // 先执行 test1
	test2()  // 再执行 test2
	test3()  // 再执行 test3
}
