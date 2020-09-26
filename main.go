package main

import (
	"fmt"
	"sync"
)

import (
	"ChatServer/src/ws"
	_ "ChatServer/src/bll/handle"
	"os"
	"os/signal"
	"syscall"
)

var (
	wg sync.WaitGroup
)

func init() {
	wg.Add(1)
}

// 处理系统信号
func signalProc() {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r)
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	for {
		// 准备接收信息
		<-sigs
		os.Exit(0)
	}
}

func main() {
	// 处理系统信号
	go signalProc()

	// 启动socket服务器
	go ws.StartServer(&wg, ":8765")

	// 阻塞等待，以免main线程退出
	wg.Wait()
}
