package handler

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"happyladysauce/rpc/model"
	"happyladysauce/rpc/service"
)

// 实现了StudentServer接口的空结构体
type Student struct {
	service.UnimplementedStudentServer // 提供所有方法的默认实现
}

func (s *Student) QueryStudent(ctx context.Context, req *service.QueryStudentRequest) (resp *service.QueryStudentResponse, err error) {
	// 日志记录请求信息
	fmt.Printf("[QueryStudent] 接收到请求: %+v\n", req)

	// 参数校验
	if req.Id <= 0 && req.Name == "" {
		err = fmt.Errorf("无效的请求参数: id和name不能同时为空")
		fmt.Printf("[QueryStudent] 错误: %v\n", err)
		return nil, err
	}

	// 构造响应
	resp = &service.QueryStudentResponse{
		Students: []*model.Student{
			{Id: 123, Name: "Happyladysauce", Age: 18},
		},
	}

	fmt.Printf("[QueryStudent] 请求处理成功\n")
	return
}

// 批量查询学生
func (s *Student) QueryStudents(ctx context.Context, req *service.StudentIds) (resp *service.QueryStudentResponse, err error) {
	// 日志记录请求信息
	fmt.Printf("[QueryStudents] 接收到请求: %+v\n", req)

	// 参数校验
	if len(req.Ids) == 0 {
		err = fmt.Errorf("无效的请求参数: ids不能为空")
		fmt.Printf("[QueryStudents] 错误: %v\n", err)
		return nil, err
	}

	// 构造响应
	resp = &service.QueryStudentResponse{
		Students: []*model.Student{
			{Id: 123, Name: "Happyladysauce", Age: 18},
		},
	}

	fmt.Printf("[QueryStudents] 请求处理成功, 返回 %d 个学生\n", len(resp.Students))
	return
}

// 流式批量查询学生
func (s *Student) QueryStudentsStream(req *service.StudentIds, stream grpc.ServerStreamingServer[model.Student]) error {
	// 日志记录请求信息
	fmt.Printf("[QueryStudentsStream] 接收到请求: %+v\n", req)

	// 参数校验
	if len(req.Ids) == 0 {
		err := fmt.Errorf("无效的请求参数: ids不能为空")
		fmt.Printf("[QueryStudentsStream] 错误: %v\n", err)
		return err
	}

	// 逐个发送学生信息
	for i, id := range req.Ids {
		student := &model.Student{Id: id, Name: fmt.Sprintf("Happyladysauce-%d", id), Age: 18}
		if err := stream.Send(student); err != nil {
			err = fmt.Errorf("发送第 %d 个学生信息失败: %w", i+1, err)
			fmt.Printf("[QueryStudentsStream] 错误: %v\n", err)
			return err
		}
		fmt.Printf("[QueryStudentsStream] 成功发送学生信息: id=%d\n", id)
	}

	fmt.Printf("[QueryStudentsStream] 请求处理完成, 共发送 %d 个学生信息\n", len(req.Ids))
	return nil
}

// 客户端流式查询学生
func (s *Student) QueryStudentsStream2(stream grpc.ClientStreamingServer[service.StudentId, service.QueryStudentResponse]) error {
	fmt.Printf("[QueryStudentsStream2] 开始接收客户端流请求\n")
	var students []*model.Student

	// 接收所有请求
	for i := 0; ; i++ {
		// 接收单个请求
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// 客户端已完成发送
				fmt.Printf("[QueryStudentsStream2] 客户端流结束，共接收 %d 个请求\n", i)
				break
			}
			err = fmt.Errorf("接收第 %d 个请求失败: %w", i+1, err)
			fmt.Printf("[QueryStudentsStream2] 错误: %v\n", err)
			return err
		}

		// 参数校验
		if req.Id <= 0 {
			err = fmt.Errorf("第 %d 个请求参数无效: id必须大于0，当前值为 %d", i+1, req.Id)
			fmt.Printf("[QueryStudentsStream2] 错误: %v\n", err)
			return err
		}

		// 处理请求
		students = append(students, &model.Student{
			Id:   req.Id,
			Name: fmt.Sprintf("Happyladysauce-%d", req.Id),
			Age:  18,
		})
		fmt.Printf("[QueryStudentsStream2] 成功接收请求: id=%d\n", req.Id)
	}

	// 构建响应
	resp := &service.QueryStudentResponse{
		Students: students,
	}

	// 发送响应
	if err := stream.SendAndClose(resp); err != nil {
		err = fmt.Errorf("发送响应失败: %w", err)
		fmt.Printf("[QueryStudentsStream2] 错误: %v\n", err)
		return err
	}

	fmt.Printf("[QueryStudentsStream2] 请求处理完成, 返回 %d 个学生信息\n", len(students))
	return nil
}

// 双向流式查询学生
func (s *Student) QueryStudentsStream3(stream grpc.BidiStreamingServer[service.StudentId, model.Student]) error {
	fmt.Printf("[QueryStudentsStream3] 开始处理双向流请求\n")
	counter := 0

	for {
		// 接收请求
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// 客户端已完成发送
				fmt.Printf("[QueryStudentsStream3] 客户端流结束，共处理 %d 个请求\n", counter)
				return nil
			}
			err = fmt.Errorf("接收请求失败: %w", err)
			fmt.Printf("[QueryStudentsStream3] 错误: %v\n", err)
			return err
		}

		// 参数校验
		if req.Id <= 0 {
			err = fmt.Errorf("无效的请求参数: id必须大于0，当前值为 %d", req.Id)
			fmt.Printf("[QueryStudentsStream3] 错误: %v\n", err)
			return err
		}

		// 构建响应
		student := &model.Student{
			Id:   req.Id,
			Name: fmt.Sprintf("Happyladysauce-%d", req.Id),
			Age:  18,
		}

		// 发送响应
		if err := stream.Send(student); err != nil {
			err = fmt.Errorf("发送响应失败: %w", err)
			fmt.Printf("[QueryStudentsStream3] 错误: %v\n", err)
			return err
		}

		counter++
		fmt.Printf("[QueryStudentsStream3] 成功处理请求: id=%d\n", req.Id)
	}
}
