package main

import (
	"fmt"
	"net"
	"os"
	"test/client"
	"test/server"
)

// 网络地址解析
func test1() {
	// MAC地址解析
	mac, err := net.ParseMAC("01:23:45:67:89:ab")
	if err != nil {
		fmt.Println("MAC地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("MAC地址解析成功:", mac)

	// IP地址解析
	ip := net.ParseIP("192.168.1.1")
	if ip == nil {
		fmt.Println("IP地址解析失败")
		os.Exit(1)
	}
	fmt.Println("IP地址解析成功:", ip)

	// CIDR地址解析
	ip, ipnet, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		fmt.Println("CIDR地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("IP地址解析成功:", ip)
	fmt.Println("CIDR地址网段解析成功:", ipnet)
}

// ResolveIPAddr IP地址解析函数 对比 ParseCIDR 函数功能简单
func test2() {
	// IPv4地址解析
	ipv4Addr, err := net.ResolveIPAddr("ip4", "192.168.1.1")
	if err != nil {
		fmt.Println("ipv4地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("ipv4地址解析成功:", ipv4Addr)

	// IPv6地址解析
	ipv6Addr, err := net.ResolveIPAddr("ip6", "2001:db8::1")
	if err != nil {
		fmt.Println("ipv6地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("ipv6地址解析成功:", ipv6Addr)
}

// TCP地址解析
func test3() {
	tcp4Addr, err := net.ResolveTCPAddr("tcp4", "192.168.1.1:8080")
	if err != nil {
		fmt.Println("tcp4地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("tcp4地址解析成功:", tcp4Addr)

	tcp6Addr, err := net.ResolveTCPAddr("tcp6", "[2001:db8::1]:8080")
	if err != nil {
		fmt.Println("tcp6地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("tcp6地址解析成功:", tcp6Addr)
}

// UDP地址解析
func test4() {
	udp4Addr, err := net.ResolveUDPAddr("udp4", "192.168.1.1:8080")
	if err != nil {
		fmt.Println("udp4地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("udp4地址解析成功:", udp4Addr)

	udp6Addr, err := net.ResolveUDPAddr("udp6", "[2001:db8::1]:8080")
	if err != nil {
		fmt.Println("udp6地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("udp6地址解析成功:", udp6Addr)
}

// unix地址解析
func test5() {
	// Unix 地址支持 unix，unixgram 和 unixpacket 三种网络类型
	uniAddr, err := net.ResolveUnixAddr("unix", "/tmp/socket")
	if err != nil {
		fmt.Println("unix地址解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("unix地址解析成功:", uniAddr)
}

// DNS解析
func test6() {
	addrs, err := net.LookupHost("www.example.com")
	if err != nil {
		fmt.Println("DNS解析失败:", err)
		os.Exit(1)
	}
	fmt.Println("DNS解析成功:", addrs)
}

//////////////////////////////////////////////////////////////////////////////////

func test7() {
	fmt.Println("启动TCP服务器...")
	// 在单独的 goroutine 启动服务器，这样主 goroutine 可以继续启动客户端用于测试
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.NewTcpServer()
	}()

	// 给服务器一点时间启动并监听
	// 注意：生产中应使用更可靠的同步方法（如信号或重试连接）。这里为演示使用简单延迟
	// sleep 使用 time 包
	// ...existing code...
	// 等待短暂时间后启动客户端
	// 如果服务器提前返回错误，打印并退出
	select {
	case err := <-srvErr:
		if err != nil {
			fmt.Println("启动TCP服务器失败:", err)
			os.Exit(1)
		}
	default:
	}

	// 使用短延迟让服务器完成 Listen
	// 引入 time 包
	// ...existing code...
	// 现在启动客户端
	// 直接调用 client.NewTcpClient
	// 如果需要，可以在 client.NewTcpClient 内部实现重试
	// 为了避免导入冲突，下面会在文件顶部添加 time 的导入

	// Sleep 200ms
	// ...existing code...

	// 启动客户端
	if err := client.NewTcpClient(); err != nil {
		fmt.Println("启动TCP客户端失败:", err)
		os.Exit(1)
	}
	fmt.Println("启动TCP客户端成功")

	// 等待服务器 goroutine 结束（如果它返回错误或退出）
	if err := <-srvErr; err != nil {
		fmt.Println("服务器退出，错误:", err)
		os.Exit(1)
	}
}

func main() {
	// 网络地址解析
	// test1()
	// test2()
	// test3()
	// test4()
	// test5()
	// test6()

	// 网络编程
	test7()
}
