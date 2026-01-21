package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"happyladysauce/rpc/service"
	"happyladysauce/server/handler"
)

// 一元拦截器 - 计时器
func timer(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Printf("[timer] 开始处理请求: %s\n", info.FullMethod)
	start := time.Now()

	// 调用实际的处理函数
	resp, err = handler(ctx, req)

	// 记录耗时
	duration := time.Since(start)
	if err != nil {
		fmt.Printf("[timer] %s 处理失败，耗时: %v, 错误: %v\n", info.FullMethod, duration, err)
	} else {
		fmt.Printf("[timer] %s 处理成功，耗时: %v\n", info.FullMethod, duration)
	}

	return
}

// 一元拦截器 - 计数器
func counter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Printf("[counter] 接收到请求: %s\n", info.FullMethod)

	// 调用实际的处理函数
	resp, err = handler(ctx, req)

	if err != nil {
		fmt.Printf("[counter] 请求 %s 处理失败\n", info.FullMethod)
	} else {
		fmt.Printf("[counter] 请求 %s 处理成功\n", info.FullMethod)
	}

	return
}

// 一元拦截器 - dev_key验证
func devkey(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 从请求元数据中获取dev_key
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("元数据获取失败")
		fmt.Printf("[devkey] %s 元数据验证失败: %v\n", info.FullMethod, err)
		return nil, status.Errorf(status.Code(err), "%v", err)
	}

	devKey := md.Get("dev_key")
	if len(devKey) == 0 {
		err = errors.New("dev_key缺失")
		fmt.Printf("[devkey] %s 元数据验证失败: %v\n", info.FullMethod, err)
		return nil, status.Errorf(status.Code(err), "%v", err)
	}

	if devKey[0] != "dev123" {
		err = errors.New("dev_key值无效")
		fmt.Printf("[devkey] %s 元数据验证失败: %v\n", info.FullMethod, err)
		return nil, status.Errorf(status.Code(err), "%v", err)
	}

	fmt.Printf("[devkey] %s 元数据验证成功\n", info.FullMethod)
	return handler(ctx, req)
}

// 流式拦截器 - 计时器
func streamTimer(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	fmt.Printf("[streamTimer] 开始处理流式请求: %s\n", info.FullMethod)
	start := time.Now()

	// 调用实际的处理函数
	err := handler(srv, ss)

	// 记录耗时
	duration := time.Since(start)
	if err != nil {
		fmt.Printf("[streamTimer] %s 处理失败，耗时: %v, 错误: %v\n", info.FullMethod, duration, err)
	} else {
		fmt.Printf("[streamTimer] %s 处理成功，耗时: %v\n", info.FullMethod, duration)
	}

	return err
}

// 流式拦截器 - 计数器
func streamCounter(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	fmt.Printf("[streamCounter] 接收到流式请求: %s\n", info.FullMethod)

	// 调用实际的处理函数
	err := handler(srv, ss)

	if err != nil {
		fmt.Printf("[streamCounter] 流式请求 %s 处理失败\n", info.FullMethod)
	} else {
		fmt.Printf("[streamCounter] 流式请求 %s 处理成功\n", info.FullMethod)
	}

	return err
}

// 流式拦截器 - dev_key验证
func streamDevkey(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 从请求元数据中获取dev_key
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		err := errors.New("元数据获取失败")
		fmt.Printf("[streamDevkey] %s 元数据验证失败: %v\n", info.FullMethod, err)
		return status.Errorf(status.Code(err), "%v", err)
	}

	devKey := md.Get("dev_key")
	if len(devKey) == 0 {
		err := errors.New("dev_key缺失")
		fmt.Printf("[streamDevkey] %s 元数据验证失败: %v\n", info.FullMethod, err)
		return status.Errorf(status.Code(err), "%v", err)
	}

	if devKey[0] != "dev123" {
		err := errors.New("dev_key值无效")
		fmt.Printf("[streamDevkey] %s 元数据验证失败: %v\n", info.FullMethod, err)
		return status.Errorf(status.Code(err), "%v", err)
	}

	fmt.Printf("[streamDevkey] %s 元数据验证成功\n", info.FullMethod)
	return handler(srv, ss)
}

func RunServer() {
	fmt.Println("[服务器] 开始初始化gRPC服务器...")

	// 计算证书文件的绝对路径
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("[服务器] 获取当前工作目录失败: %v", err))
	}
	fmt.Printf("[服务器] 当前工作目录: %s\n", currentDir)

	// 构建证书和密钥的绝对路径
	certPath := filepath.Join(currentDir, "cert", "server.crt")
	keyPath := filepath.Join(currentDir, "cert", "server.key")

	// 检查证书文件是否存在
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("[服务器] 证书文件不存在: %s", certPath))
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("[服务器] 密钥文件不存在: %s", keyPath))
	}
	fmt.Printf("[服务器] 证书和密钥文件检查通过\n")

	// 从文件加载服务器证书和私钥
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		panic(fmt.Sprintf("[服务器] 加载证书失败: %v", err))
	}
	fmt.Println("[服务器] 证书加载成功")

	// 创建gRPC服务器实例
	server := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(timer, counter, devkey),
		grpc.ChainStreamInterceptor(streamTimer, streamCounter, streamDevkey),
	)
	fmt.Println("[服务器] gRPC服务器创建成功，已配置拦截器")

	// 注册服务
	service.RegisterStudentServer(server, &handler.Student{})
	fmt.Println("[服务器] 学生服务注册成功")

	// 监听TCP端口
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(fmt.Sprintf("[服务器] 监听端口失败: %v", err))
	}
	fmt.Printf("[服务器] 开始监听端口: %s\n", listen.Addr())

	// 启动服务器
	fmt.Println("[服务器] 启动gRPC服务器...")
	if err := server.Serve(listen); err != nil {
		panic(fmt.Sprintf("[服务器] 服务器启动失败: %v", err))
	}
}
