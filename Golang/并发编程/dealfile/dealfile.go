package dealfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	READ_FILE_ROUTINE_NUM = 10
	PROCESS_ROUTINE_NUM = 15
)

var (
	// 文件处理的总和
	sum int64 = 0

	// 定义 channel 缓冲通道
	filelist = make(chan string, 100)
	lineBuffer = make(chan string, 1000)

	// 定义工作协程数量
	walkWg = sync.WaitGroup{}
	readWg = sync.WaitGroup{}
	processWg = sync.WaitGroup{}
)

// 递归目录
func walkDir(dir string) {
	defer walkWg.Done()
	filepath.WalkDir(dir, func(subPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			filelist <- subPath
		}
		return nil
	})
}

func readFile() {
	defer readWg.Done()
	for infile := range filelist {
		fin, err := os.Open(infile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// bufio 缓存读取每一行文件进入 lineBuffer
		reader := bufio.NewReader(fin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					if len(line) > 0 {
						lineBuffer <- strings.TrimSpace(line)
					}
					break
				} else {
					fmt.Println(err)
					break
				}
			} else {
				lineBuffer <- strings.TrimSpace(line)
			}
		}
		fin.Close()
	}
}

func processLine() {
	defer processWg.Done()
	for line := range lineBuffer {
		if i, err := strconv.Atoi(line); err != nil {
			fmt.Printf("%s not unmber\n", line)
		} else {
			atomic.AddInt64(&sum, int64(i))
		}
	}
}


func DealMassFile(dir string) {

	go func() {
		tk := time.NewTicker(time.Second)
		defer tk.Stop()
		for {
			<- tk.C
			fmt.Printf("堆积了%d个文件未处理，堆积了%d行内容未处理\n",
		len(filelist), len(lineBuffer))
		}
	}()

	walkWg.Add(1)
	readWg.Add(READ_FILE_ROUTINE_NUM)
	processWg.Add(PROCESS_ROUTINE_NUM)


	// 开启协程读取目录文件路径到 filelist
	go walkDir(dir)
	time.Sleep(time.Second)
	fmt.Printf("当前文件缓冲数量：%v\n", len(filelist))

	// 开启协程读取文件内容到 lineBuffer
	for i := 0; i < READ_FILE_ROUTINE_NUM; i++ {
		go readFile()
	}

	// 开启协程处理每一行文件
	for i := 0; i < PROCESS_ROUTINE_NUM; i++ {
		go processLine()
	}

	walkWg.Wait()
	close(filelist)
	readWg.Wait()
	close(lineBuffer)
	processWg.Wait()

	fmt.Printf("sum=%v\n", sum)
}






