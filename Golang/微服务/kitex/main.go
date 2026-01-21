package main

import (
	math "happyladysauce/kitex_gen/math/mathservice"
	"log"
	"net"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

func main() {
	fout, err := os.OpenFile("log/kitex.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	// set log level to debug
	klog.SetLevel(klog.LevelDebug)
	klog.SetOutput(fout)

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8083")
	svr := math.NewServer(new(MathServiceImpl), server.WithServiceAddr(addr))


	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}

	klog.Info("server is running on %s", addr.String())
}
