package server

import (
	"fmt"
	"happyladysauce/web/utils/encode"
	"io"
	"net/http"
	"strings"
)

func HttpObservation(w http.ResponseWriter,r *http.Request) {
	// 通过 r *http.Request 获取客户的请求头信息
	fmt.Printf("request method: %s\n", r.Method)
	fmt.Printf("request host: %s\n", r.Host)
	fmt.Printf("request proto: %s\n", r.Proto)
	fmt.Printf("request url: %s\n", r.URL)
	params := encode.ParseUrlParams(r.URL.RawQuery)
	fmt.Println("request header.")

	// 请求头是一个 map
	for key, values := range r.Header {
		fmt.Printf("%s: %v\n", key, values)
	}
	fmt.Println()

	// 请求体则是和文件一样的"流"
	fmt.Println("request body:")
	// io.Copy(os.Stdout, r.Body)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(body))

	
	// 响应头,状态码,响应体都必须按照顺序进行赋值填写,如果混乱了则无效
	// 通过 w *http.ResponseWrite 写入http回复, 返回给客户端
	w.Header().Add("happyladysauce-name", "666")	// 在http头中写入数据
	w.WriteHeader(http.StatusAccepted)	// 设置 http 状态码
	// 如果在调用 Write 时之前没有调用 WriteHeader 则默认赋值状态码为 200
	w.Write([]byte("Hello World!\n"))	// 写入数据体
	w.Write([]byte("Hello World!\n"))
	fmt.Fprintf(w, "your name is %s, age is %s\n", params["name"], params["age"])
	w.Header().Add("happyladysauce-name", "666")	// 无效
	fmt.Println(strings.Repeat("*", 60))
}

func HttpStream(w http.ResponseWriter, r *http.Request) {
	line := []byte("这个是一个大数据的一部分\n")
	totalSize := len(line)
	const P = 10
	totalSize = totalSize * P

	w.Header().Add("Content-Length", fmt.Sprintf("%d", totalSize))

	for i := 0; i < P; i++ {
		if _, err := w.Write(line); err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
		// time.Sleep(time.Second)
	}
	fmt.Println("stream done.")
}


func NewWebServer() {
	http.HandleFunc("/obs", HttpObservation)
	http.HandleFunc("/stream", HttpStream)
	
	// 启动http服务端
	err := http.ListenAndServe("127.0.0.1:8081", nil)
	if err != nil {
		panic(err)
	}
}

