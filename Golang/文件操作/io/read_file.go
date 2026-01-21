package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadFile(file string) []byte {
	var info []byte
	if fin, err := os.Open(file); err != nil {
		fmt.Printf("open file failed, err: %s\n", err)
	} else {
		defer fin.Close()
		bs := make([]byte, 10)

		for {
			if n, err := fin.Read(bs); err != nil {
				if err == io.EOF {
					fmt.Println("file is end")
					break
				}
				fmt.Printf("read file failed. err: %v\n", err)
			}else {
				fmt.Printf("%v\n", bs[0:n])
				info = append(info, bs[0:n]...)
			}
		}
	}
	return info
}

func ReadFileWithBuffer(file string) []byte {
	var bs []byte
	if fin, err := os.Open(file); err != nil {
		fmt.Printf("open file failed. err: %v\n", err)
	}else {
		defer fin.Close()
		reader := bufio.NewReader(fin)
		for {
			line, err := reader.ReadString('\n')
			if len(line) > 0 {
				bs = append(bs, line...)
				fmt.Println(line)
			}
			if err == io.EOF {
				break
			}
		}
	}
	return bs
}
