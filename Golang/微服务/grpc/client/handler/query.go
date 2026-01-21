package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"happyladysauce/rpc/model"
	"happyladysauce/rpc/service"
)

// QueryStudent 查询单个学生
func QueryStudent(conn *grpc.ClientConn) (*model.Student, error) {
	start := time.Now()
	log.Printf("[QueryStudent] Start querying student with ID: %d, Name: %s", 1, "happyladysauce")

	// 参数校验
	if conn == nil {
		log.Printf("[QueryStudent] Error: Connection is nil")
		return nil, fmt.Errorf("connection is nil")
	}

	client := service.NewStudentClient(conn)
	resp, err := client.QueryStudent(context.Background(), &service.QueryStudentRequest{
		Id:   1,
		Name: "happyladysauce",
	})

	if err != nil {
		// 解析gRPC错误
		s, ok := status.FromError(err)
		if ok {
			log.Printf("[QueryStudent] gRPC error: Code=%s, Message=%s", s.Code(), s.Message())
			return nil, fmt.Errorf("gRPC error: code=%s, message=%s", s.Code(), s.Message())
		}
		log.Printf("[QueryStudent] Unknown error: %v", err)
		return nil, fmt.Errorf("unknown error: %v", err)
	}

	// 检查响应数据
	if resp == nil || len(resp.Students) == 0 {
		log.Printf("[QueryStudent] No student found")
		return nil, fmt.Errorf("no student found")
	}

	log.Printf("[QueryStudent] Success, found %d students, took %v ms", len(resp.Students), time.Since(start).Milliseconds())
	return resp.Students[0], nil
}

// QueryStudents 批量查询学生
func QueryStudents(conn *grpc.ClientConn) ([]*model.Student, error) {
	start := time.Now()
	log.Printf("[QueryStudents] Start batch querying students with IDs: %v, Names: %v", []int32{1, 2, 3}, []string{"student1", "student2", "student3"})

	// 参数校验
	if conn == nil {
		log.Printf("[QueryStudents] Error: Connection is nil")
		return nil, fmt.Errorf("connection is nil")
	}

	client := service.NewStudentClient(conn)
	resp, err := client.QueryStudents(context.Background(), &service.StudentIds{
		Ids: []int64{1, 2, 3},
	})

	if err != nil {
		// 解析gRPC错误
		s, ok := status.FromError(err)
		if ok {
			log.Printf("[QueryStudents] gRPC error: Code=%s, Message=%s", s.Code(), s.Message())
			return nil, fmt.Errorf("gRPC error: code=%s, message=%s", s.Code(), s.Message())
		}
		log.Printf("[QueryStudents] Unknown error: %v", err)
		return nil, fmt.Errorf("unknown error: %v", err)
	}

	// 检查响应数据
	if resp == nil {
		log.Printf("[QueryStudents] Empty response")
		return nil, fmt.Errorf("empty response")
	}

	log.Printf("[QueryStudents] Success, received %d students, took %v ms", len(resp.Students), time.Since(start).Milliseconds())
	return resp.Students, nil
}

// QueryStudentsStream 服务器流式查询学生
func QueryStudentsStream(conn *grpc.ClientConn) ([]*model.Student, error) {
	start := time.Now()
	log.Printf("[QueryStudentsStream] Start server streaming query with IDs: %v, Names: %v", []int32{1, 2, 3, 4, 5}, []string{"student1", "student2", "student3", "student4", "student5"})

	// 参数校验
	if conn == nil {
		log.Printf("[QueryStudentsStream] Error: Connection is nil")
		return nil, fmt.Errorf("connection is nil")
	}

	client := service.NewStudentClient(conn)
	stream, err := client.QueryStudentsStream(context.Background(), &service.StudentIds{
		Ids: []int64{1, 2, 3, 4, 5},
	})

	if err != nil {
		// 解析gRPC错误
		s, ok := status.FromError(err)
		if ok {
			log.Printf("[QueryStudentsStream] gRPC error: Code=%s, Message=%s", s.Code(), s.Message())
			return nil, fmt.Errorf("gRPC error: code=%s, message=%s", s.Code(), s.Message())
		}
		log.Printf("[QueryStudentsStream] Unknown error: %v", err)
		return nil, fmt.Errorf("unknown error: %v", err)
	}

	var students []*model.Student
	receiveCount := 0

	for {
		resp, err := stream.Recv()
		if err != nil {
			// 流结束时会返回io.EOF错误
			if err.Error() == "EOF" {
				log.Printf("[QueryStudentsStream] Stream ended normally")
				break
			}
			log.Printf("[QueryStudentsStream] Error receiving stream data: %v", err)
			return students, fmt.Errorf("error receiving stream data: %v", err)
		}

		receiveCount++
		student := resp
		if student != nil {
			students = append(students, student)
			log.Printf("[QueryStudentsStream] Received student #%d: ID=%d, Name=%s", receiveCount, student.Id, student.Name)
		}
	}

	log.Printf("[QueryStudentsStream] Success, total received %d students, took %v ms", len(students), time.Since(start).Milliseconds())
	return students, nil
}

// QueryStudentsStream2 客户端流式查询学生
func QueryStudentsStream2(conn *grpc.ClientConn) (*service.QueryStudentResponse, error) {
	start := time.Now()
	log.Printf("[QueryStudentsStream2] Start client streaming query")

	// 参数校验
	if conn == nil {
		log.Printf("[QueryStudentsStream2] Error: Connection is nil")
		return nil, fmt.Errorf("connection is nil")
	}

	client := service.NewStudentClient(conn)
	stream, err := client.QueryStudentsStream2(context.Background())

	if err != nil {
		// 解析gRPC错误
		s, ok := status.FromError(err)
		if ok {
			log.Printf("[QueryStudentsStream2] gRPC error: Code=%s, Message=%s", s.Code(), s.Message())
			return nil, fmt.Errorf("gRPC error: code=%s, message=%s", s.Code(), s.Message())
		}
		log.Printf("[QueryStudentsStream2] Unknown error: %v", err)
		return nil, fmt.Errorf("unknown error: %v", err)
	}

	// 发送多个请求
	requests := []*service.StudentId{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	sendCount := 0
	for _, req := range requests {
		err := stream.Send(req)
		if err != nil {
			log.Printf("[QueryStudentsStream2] Error sending request: %v, Request: %+v", err, req)
			// 尝试关闭流
			_ = stream.CloseSend()
			return nil, fmt.Errorf("error sending request: %v", err)
		}
		sendCount++
		log.Printf("[QueryStudentsStream2] Sent request #%d: ID=%d", sendCount, req.Id)
	}

	log.Printf("[QueryStudentsStream2] All %d requests sent, closing send stream", sendCount)

	// 关闭发送流并接收响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		// 解析gRPC错误
		s, ok := status.FromError(err)
		if ok {
			log.Printf("[QueryStudentsStream2] gRPC error when receiving response: Code=%s, Message=%s", s.Code(), s.Message())
			return nil, fmt.Errorf("gRPC error: code=%s, message=%s", s.Code(), s.Message())
		}
		log.Printf("[QueryStudentsStream2] Unknown error when receiving response: %v", err)
		return nil, fmt.Errorf("unknown error: %v", err)
	}

	log.Printf("[QueryStudentsStream2] Success, received response with %d students, took %v ms", len(resp.Students), time.Since(start).Milliseconds())
	return resp, nil
}

// QueryStudentsStream3 双向流式查询学生
func QueryStudentsStream3(conn *grpc.ClientConn) error {
	start := time.Now()
	log.Printf("[QueryStudentsStream3] Start bidirectional streaming query")

	// 参数校验
	if conn == nil {
		log.Printf("[QueryStudentsStream3] Error: Connection is nil")
		return fmt.Errorf("connection is nil")
	}

	client := service.NewStudentClient(conn)
	stream, err := client.QueryStudentsStream3(context.Background())

	if err != nil {
		// 解析gRPC错误
		s, ok := status.FromError(err)
		if ok {
			log.Printf("[QueryStudentsStream3] gRPC error: Code=%s, Message=%s", s.Code(), s.Message())
			return fmt.Errorf("gRPC error: code=%s, message=%s", s.Code(), s.Message())
		}
		log.Printf("[QueryStudentsStream3] Unknown error: %v", err)
		return fmt.Errorf("unknown error: %v", err)
	}

	// 发送多个请求并接收响应
	requests := []*service.StudentId{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	// 使用goroutine发送请求
	sendCount := 0
	sendErrCh := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[QueryStudentsStream3] Panic in send goroutine: %v", r)
			}
			// 关闭发送流
			if err := stream.CloseSend(); err != nil {
				log.Printf("[QueryStudentsStream3] Error closing send stream: %v", err)
			} else {
				log.Printf("[QueryStudentsStream3] Send stream closed successfully")
			}
		}()

		for _, req := range requests {
			err := stream.Send(req)
			if err != nil {
				log.Printf("[QueryStudentsStream3] Error sending request: %v, Request: %+v", err, req)
				sendErrCh <- fmt.Errorf("error sending request: %v", err)
				return
			}
			sendCount++
			log.Printf("[QueryStudentsStream3] Sent request #%d: ID=%d", sendCount, req.Id)
			// 短暂延迟，模拟实际场景
			time.Sleep(100 * time.Millisecond)
		}
		close(sendErrCh)
	}()

	// 接收响应
	receiveCount := 0
	for {
		resp, err := stream.Recv()
		if err != nil {
			// 流结束时会返回io.EOF错误
			if err.Error() == "EOF" {
				log.Printf("[QueryStudentsStream3] Stream ended normally")
				break
			}
			log.Printf("[QueryStudentsStream3] Error receiving stream data: %v", err)
			return fmt.Errorf("error receiving stream data: %v", err)
		}

		receiveCount++
		log.Printf("[QueryStudentsStream3] Received response #%d: Student=%+v", receiveCount, resp)
	}

	// 检查发送过程中是否有错误
	if sendErr, ok := <-sendErrCh; ok && sendErr != nil {
		return sendErr
	}

	log.Printf("[QueryStudentsStream3] Success, sent %d requests, received %d responses, took %v ms", sendCount, receiveCount, time.Since(start).Milliseconds())
	return nil
}
