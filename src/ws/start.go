package ws

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// 启动服务器
func StartServer(wg *sync.WaitGroup, address string) {
	defer func() {
		wg.Done()
	}()

	fmt.Printf("Socket服务器开始监听...")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Listen Error: %s", err)))
	} else {
		msg := fmt.Sprintf("Got listener for the server. (local address: %s)", listener.Addr())
		fmt.Println(msg)
	}

	for {
		// 阻塞直至新连接到来
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept Error: %s", err)
			continue
		}

		clientObj := NewClient(conn)
		go HandleConn(clientObj)
	}
}
