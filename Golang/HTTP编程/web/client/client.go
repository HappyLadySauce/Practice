package client

import (
	"bufio"
	"fmt"
	"happyladysauce/web/utils/encode"
	"io"
	"net/http"
)

func HttpObservation() {
	resp, err := http.Get("http://127.0.0.1:8081/obs?" + encode.EncodeUrlParams(map[string]string{
		"name": "happyladysauce",
		"age": "18 还是 17 呢？",
	}))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("response proto: %s\n", resp.Proto)
	if major, minor, ok := http.ParseHTTPVersion(resp.Proto); ok {
		fmt.Printf("http major version: %d, minor version: %d\n", major, minor)
	}
	fmt.Printf("response status: %s\n", resp.Status)
	fmt.Printf("response status code: %d\n", resp.StatusCode)

	for key, values := range resp.Header {
		fmt.Printf("%s: %v\n", key, values)
	}

	fmt.Println("respone body:")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(body))
}

func HttpStream() {
	resp, err := http.Get("http://127.0.0.1:8081/stream")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // 将defer移到错误检查后，确保始终关闭响应体
	
	totalSize := resp.ContentLength
	fmt.Printf("total size: %d\n", totalSize)

	reader := bufio.NewReader(resp.Body)
	var totalReadSize int64
	for {
		bs, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(bs) > 0 {
					totalReadSize += int64(len(bs))
					fmt.Printf("read %d%% bytes ", totalReadSize*100/totalSize)
					fmt.Print(string(bs))
				}
				break
			} else {
				fmt.Printf("err: %v\n", err)
				break
			}
		} else {
			totalReadSize += int64(len(bs))
			fmt.Printf("read %d%% bytes ", totalReadSize*100/totalSize)
			fmt.Print(string(bs))
		}
	}
}

func NewHttpClient() {
	// HttpObservation()
	HttpStream()
}

