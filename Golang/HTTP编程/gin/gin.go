package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HomeHandler(ctx *gin.Context) {
	// gin.Context 中封装了标准库中的 request 和 response, 这意味着其完全支持标准库的所有功能, 并且可导出 request 和 response
	fmt.Println("请求头")
	for k, v := range ctx.Request.Header {
		fmt.Printf("%s=%s\n", k, v)
	}

	fmt.Println("请求体")
	io.Copy(os.Stdout, ctx.Request.Body)

	// 生成响应头和响应体
	ctx.Writer.Header().Add("language1", "go")
	ctx.Header("language2", "golang")

	// 设置状态码
	ctx.Writer.WriteHeader(http.StatusOK)

	// 最后设置响应体
	ctx.Writer.WriteString("welcome")

	// gin 的封装函数
	ctx.String(200, "to BeiJing")
	
	ctx.JSON(200, map[string]any{"welcome": "world!"})	// 支持结构体和map参数
	ctx.JSON(200, gin.H{"welcome": "world!"})
}

func G1(ctx *gin.Context) {
	ctx.String(200, "你好")
	ctx.Next()
	ctx.String(200 ,"g1中间件结束")
}

func main() {
	// engine := gin.Default()	// 默认使用了 Logger 和 Recovery 中间件
	engine := gin.New()	// 不使用默认的中间件
	{
		// 路由分组
		g1 := engine.Group("/v1")
		g1.Use(G1)	// 路由组公共中间件
		g1.GET("/a", func (ctx *gin.Context) {	// 路径为 /v1/a
			ctx.String(200, "/a")
		})
	}

	engine.Use(gin.Logger())		// 使用日志中间件
	engine.Use(gin.Recovery())	// 使用 recovery 中间件
	// 定义路由
	engine.GET("/home", HomeHandler)

	if err := engine.Run("127.0.0.1:8082"); err != nil {
		panic(err)
	}
}

// Gin 获取参数
func Url(engine *gin.Engine) {
	// http://127.0.0.1:8082/student?name=zcy&add=bj
	engine.GET("/student", func (ctx *gin.Context)  {
		a := ctx.Query("name")
		b := ctx.DefaultQuery("addr", "China")	// 如果没传addr参数, 则默认为China
		ctx.String(200, a + "live in" + b)
	})
}

// 从 Restful 风格的 url 中获取参数
func Restful(engine *gin.Engine) {
	// : 只能对应一级路径, * 则可以匹配多级路径
	engine.GET("/student/:name/*addr", func(ctx *gin.Context) {
		name := ctx.Param("name")
		addr := ctx.Param("addr")
		ctx.String(200, name + "live in" + addr)
	})
}

// 从 Post 请求中获取参数
func PostForm(engine *gin.Engine) {
	engine.POST("/student/form", func(ctx *gin.Context) {
		name := ctx.PostForm("username")
		addr := ctx.DefaultPostForm("addr", "China")
		ctx.String(200, name + "live in" + addr)
	})
}

type Student struct {
	name string
	addr string
}

// 从 JSON 请求体中获取参数
func PostJson(engine *gin.Engine) {
	engine.POST("/student/json", func(ctx *gin.Context) {
		var stu Student
		bs, _ := io.ReadAll(ctx.Request.Body)
		if err := json.Unmarshal(bs, &stu); err != nil {
			name := stu.name
			addr := stu.addr
			ctx.String(200, name + "live in" + addr)
		}
	})
}

// 上传文件
func Upload_file(engine * gin.Engine) {
	engine.MaxMultipartMemory = 8 << 20	// 设置表单上传大小为 8M, 默认上限是 32M
	engine.POST("/upload", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")	// 绑定请求体键值
		if err != nil {
			fmt.Printf("get file error %v\n", err)
			ctx.String(http.StatusInternalServerError, "upload file failed")
		} else {
			if err = ctx.SaveUploadedFile(file, "./data" + file.Filename); err != nil {
				fmt.Printf("save file to %s failed: %v\n", "./data/" + file.Filename, err)
			} else {
				ctx.String(http.StatusOK, file.Filename)
			}
		}
	})
}

