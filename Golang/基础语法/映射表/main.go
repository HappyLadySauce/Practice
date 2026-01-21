package main

import (
	"fmt"
	"math/rand"
)

// 一般来说，映射表数据结构实现通常有两种，哈希表(hash table)和搜索树(search tree)，区别在于前者无序，后者有序。
// 在 Go 中，map的实现是基于哈希桶
// map 并不是一个并发安全的数据结构，Go 团队认为大多数情况下 map 的使用并不涉及高并发的场景，引入互斥锁会极大的降低性能，map 内部有读写检测机制，如果冲突会触发fatal error。

// 在 Go 中，map 的键类型必须是可比较的，比如string ，int是可比较的
func test1() {
	// map 的初始化 map[keyType]valueType{}
	// 内置函数make，对于 map 而言，接收两个参数，分别是类型与初始容量
	m := make(map[string]int, 8)
	m["a"] = 1
	m["b"] = 2
	fmt.Println(m)

	mp := map[int]string {
		1: "a",
		2: "b",
	}
	fmt.Println(mp)
}

// Set 是一种无序的，不包含重复元素的集合，Go 中并没有提供类似的数据结构实现，但是 map 的键正是无序且不能重复的，所以也可以使用 map 来替代 set
func test2() {
	set := make(map[int]struct{}, 10)
	for i := 0; i < 10; i++ {
		set[rand.Intn(100)] = struct{}{}
	}
	fmt.Println(set)
}

func main() {
	// test1()
	test2()
}



