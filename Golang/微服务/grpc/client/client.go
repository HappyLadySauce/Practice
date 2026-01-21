package client

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"happyladysauce/client/handler"
)

// timer是gRPC客户端的一元拦截器，用于记录调用耗时
func timer(ctx context.Context, method string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	begin := time.Now()
	err := invoker(ctx, method, req, resp, cc, opts...)
	log.Printf("[TimerInterceptor] Method: %s, Cost: %v ms, Error: %v", method, time.Since(begin).Milliseconds(), err)
	return err
}

// counter是gRPC客户端的一元拦截器，用于计数和日志记录
func counter(ctx context.Context, method string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("[CounterInterceptor] Calling method: %s", method)
	err := invoker(ctx, method, req, resp, cc, opts...)
	if err != nil {
		log.Printf("[CounterInterceptor] Method %s failed with error: %v", method, err)
	} else {
		log.Printf("[CounterInterceptor] Method %s completed successfully", method)
	}
	return err
}

// devKey是gRPC客户端的一元拦截器，用于为上下文添加一个键值对，用于传递给服务器 模拟实现 rpc认证
func devKey(ctx context.Context, method string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("[DevKeyInterceptor] Adding authentication header for method: %s", method)
	ctx = metadata.AppendToOutgoingContext(ctx, "dev_key", "dev123") // 为上下文添加一个键值对，用于传递给服务器
	err := invoker(ctx, method, req, resp, cc, opts...)
	if err != nil {
		// 解析认证相关的错误
		s, ok := status.FromError(err)
		if ok && s.Code() == grpc.Code(errors.New("unauthenticated")) {
			log.Printf("[DevKeyInterceptor] Authentication failed for method: %s", method)
		}
	}
	return err
}

// streamTimer是gRPC客户端的流式拦截器，用于记录流式方法调用耗时
func streamTimer(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	begin := time.Now()
	stream, err := streamer(ctx, desc, cc, method, opts...)
	log.Printf("[StreamTimerInterceptor] Method: %s, Cost: %v ms, Error: %v", method, time.Since(begin).Milliseconds(), err)
	return stream, err
}

// streamCounter是gRPC客户端的流式拦截器，用于计数和日志记录
func streamCounter(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("[StreamCounterInterceptor] Calling stream method: %s", method)
	stream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		log.Printf("[StreamCounterInterceptor] Stream method %s failed with error: %v", method, err)
	} else {
		log.Printf("[StreamCounterInterceptor] Stream method %s connection established", method)
	}
	return stream, err
}

// streamDevKey是gRPC客户端的流式拦截器，用于为上下文添加dev_key
func streamDevKey(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("[StreamDevKeyInterceptor] Adding authentication header for stream method: %s", method)
	ctx = metadata.AppendToOutgoingContext(ctx, "dev_key", "dev123")
	stream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		// 解析认证相关的错误
		s, ok := status.FromError(err)
		if ok && s.Code() == grpc.Code(errors.New("unauthenticated")) {
			log.Printf("[StreamDevKeyInterceptor] Authentication failed for stream method: %s", method)
		}
	}
	return stream, err
}

func RunClient() {
	log.Printf("[Client] Starting gRPC client initialization...")
	// 计算证书文件的绝对路径
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("[Client] Failed to get current working directory: %v", err)
	}
	log.Printf("[Client] Current working directory: %s", currentDir)

	// 构建证书的绝对路径
	certPath := filepath.Join(currentDir, "cert", "server.crt")
	log.Printf("[Client] Certificate path: %s", certPath)

	// 检查证书文件是否存在
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		log.Fatalf("[Client] Certificate file does not exist: %s", certPath)
	}
	log.Printf("[Client] Certificate file exists")

	// 从文件加载服务器证书进行身份验证
	creds, err := credentials.NewClientTLSFromFile(certPath, "localhost") // 验证服务器证书
	if err != nil {
		log.Fatalf("[Client] Failed to load certificate: %v", err)
	}
	log.Printf("[Client] Certificate loaded successfully")

	log.Printf("[Client] Establishing gRPC connection to 127.0.0.1:9090...")
	conn, err := grpc.NewClient(
		"127.0.0.1:9090",                     // 端口与服务器保持一致
		grpc.WithTransportCredentials(creds), // 使用TLS加密
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024), // 最大接收消息大小为1MB
			grpc.MaxCallSendMsgSize(1024*1024), // 最大发送消息大小为1MB
		),
		// grpc.WithUnaryInterceptor(timer), // UnaryInterceptor是gRPC客户端的一元拦截器
		grpc.WithChainUnaryInterceptor(timer, counter, devKey),                    // 链式拦截器，先调用timer，再调用counter
		grpc.WithChainStreamInterceptor(streamTimer, streamCounter, streamDevKey), // 链式流式拦截器
	)
	if err != nil {
		log.Fatalf("[Client] Failed to establish gRPC connection: %v", err)
	}
	log.Printf("[Client] gRPC connection established successfully")
	defer func() {
		log.Printf("[Client] Closing gRPC connection")
		conn.Close()
	}()

	// 测试所有方法
	testMethods(conn)
}

// testMethods 测试所有可用的gRPC方法
func testMethods(conn *grpc.ClientConn) {
	// 测试单个查询
	log.Printf("[Client] Testing QueryStudent method...")
	student, err := handler.QueryStudent(conn)
	if err != nil {
		log.Printf("[Client] ERROR: QueryStudent failed: %v", err)
	} else {
		log.Printf("[Client] QueryStudent succeeded, student: %+v", student)
	}

	// 测试批量查询
	log.Printf("[Client] Testing QueryStudents method...")
	students, err := handler.QueryStudents(conn)
	if err != nil {
		log.Printf("[Client] ERROR: QueryStudents failed: %v", err)
	} else {
		log.Printf("[Client] QueryStudents succeeded, count: %d", len(students))
	}

	// 测试服务器流式查询
	log.Printf("[Client] Testing QueryStudentsStream method...")
	streamStudents, err := handler.QueryStudentsStream(conn)
	if err != nil {
		log.Printf("[Client] ERROR: QueryStudentsStream failed: %v", err)
	} else {
		log.Printf("[Client] QueryStudentsStream succeeded, count: %d", len(streamStudents))
	}

	// 测试客户端流式查询
	log.Printf("[Client] Testing QueryStudentsStream2 method...")
	stream2Resp, err := handler.QueryStudentsStream2(conn)
	if err != nil {
		log.Printf("[Client] ERROR: QueryStudentsStream2 failed: %v", err)
	} else {
		log.Printf("[Client] QueryStudentsStream2 succeeded, count: %d", len(stream2Resp.Students))
	}

	// 测试双向流式查询
	log.Printf("[Client] Testing QueryStudentsStream3 method...")
	err = handler.QueryStudentsStream3(conn)
	if err != nil {
		log.Printf("[Client] ERROR: QueryStudentsStream3 failed: %v", err)
	} else {
		log.Printf("[Client] QueryStudentsStream3 succeeded")
	}

	log.Printf("[Client] All methods tested")
}
