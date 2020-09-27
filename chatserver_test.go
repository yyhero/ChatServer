package main

import (
	"ChatServer/src/model"
	"ChatServer/src/ws"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	w sync.WaitGroup
)

func TestChat(t *testing.T) {
	fmt.Println("TestChat--------------")
	client := newClient()

	// receive msg
	receChan := make(chan ws.SocketResponseObject, 10)
	go receiveMsg(client, receChan)


	// create a room
	createRoom(client)

	// chat
	chat(client,"fcuk, boy!!!!")

	time.Sleep(5* time.Second)
}

func TestBroadChat(t *testing.T)  {
	fmt.Println("TestBroadChat--------------")

	// create 4 client
	clientSlice := make([]*ws.Client, 0)
	for i:=1; i<=4; i++{
		obj := newClient()
		clientSlice = append(clientSlice,obj)
	}

	testC := clientSlice[0]
	createRoom(testC)
	receChan := make(chan ws.SocketResponseObject, 10)
	go receiveMsg(testC, receChan)

	var roomId int64
	select {
	case resp := <-receChan:
		fmt.Printf("receive msg:%v", resp)
		if resp.CommandType == model.CreateRoom{
			data := resp.Data.(map[string]interface{})
			roomId = int64(data["RoomId"].(float64))
			break
		}
	}

	// enter room
	for i:=1; i< 4; i++ {
		joinRoom(clientSlice[i], roomId)
		go receiveMsg(clientSlice[i], nil)
	}

	//chat
	chat(testC, "fcuk, girls!!!!!")
	time.Sleep(5*time.Second)
}


// return the last 50 msgs
func TestLatestMsg(t *testing.T) {
	fmt.Println("TestLatestMsg--------------")
	testA := newClient()
	createRoom(testA)
	receChan := make(chan ws.SocketResponseObject, 10)
	go receiveMsg(testA, receChan)

	var roomId int64
	select {
	case resp := <-receChan:
		fmt.Printf("clientId:%v, receive msg:%v", testA.Id,resp)
		if resp.CommandType == model.CreateRoom{
			data := resp.Data.(map[string]interface{})
			roomId = int64(data["RoomId"].(float64))
			break
		}
	}

	// send 100 msg
	for i:= 1; i<=100; i++{
		msg	:= "fcuk,send Msg"+ strconv.Itoa(i)
		chat(testA, msg)
	}
	time.Sleep(5* time.Second)

	// new testB
	testB := newClient()
	joinRoom(testB, roomId)
	go receiveMsg(testB,nil)

	time.Sleep(5* time.Second)
}

func TestCommandPopular(t *testing.T) {
	fmt.Println("TestCommand--------------")
	testA := newClient()
	go receiveMsg(testA, nil)

	createRoom(testA)

	msg1 := "Msg aa  aa  bbb bbb bbb  home home home home"
	msg2 := "fsdaf   dfsdf  sdaf  asdf      asdf    asdf asdf asdf asdf"
	chat(testA, msg1)
	chat(testA, msg2)
	time.Sleep(2 * time.Second)

	command(testA, "/popular")
	time.Sleep(5*time.Second)
}

func TestCommandStats(t *testing.T) {
	fmt.Println("TestCommand--------------")
	testA := newClient()
	go receiveMsg(testA, nil)
	createRoom(testA)

	command(testA, "/stats", testA.Id)
	time.Sleep(5*time.Second)
}

func command(client *ws.Client, cmd string, cmdPara ...interface{})  {
	info:= make(map[string]interface{})
	info["Cmd"] = cmd
	switch  cmd {
	case "/stats":
		info["CmdPara"] = cmdPara
	case "/popular":
	}
	request := ws.NewSocketRequest(model.Command,info)
	data, err := json.Marshal(request)
	if err != nil{
		return
	}
	err = client.SendBytes(data)
	fmt.Printf("clientId:%v, send message, data:%s\n", client.Id,data)
}


func joinRoom(client *ws.Client, roomId int64)  {
	info:= make(map[string]interface{})
	info["RoomId"] = roomId
	request := ws.NewSocketRequest(model.JoinRoom,info)
	data, err := json.Marshal(request)
	if err != nil{
		return
	}
	err = client.SendBytes(data)
	fmt.Printf("clientId:%v, send message, data:%s\n", client.Id,data)
}

func createRoom(client *ws.Client) {
	request := ws.NewSocketRequest(model.CreateRoom,nil)
	data, err := json.Marshal(request)
	if err != nil{
		return
	}
	err = client.SendBytes(data)
	fmt.Printf("clientId:%v, send message, data:%s\n", client.Id,data)
}

func chat(client *ws.Client,str string) {
	info := map[string]interface{}{}
	info["Msg"] = str
	request := ws.NewSocketRequest(model.SendMessage,info)
	data, err := json.Marshal(request)
	if err != nil{
		return
	}
	fmt.Printf("clientId:%v, send message, data:%s\n", client.Id,data)
	err = client.SendBytes(data)
}

func receiveMsg(clientObj *ws.Client, msgChan chan ws.SocketResponseObject) {
	for {
		readBytes := make([]byte, 1024)

		n, err := clientObj.Conn.Read(readBytes)
		if err != nil {
			var errMsg string
			if err == io.EOF {
				errMsg = fmt.Sprintf("另一端关闭了连接：%s，读取到的字节数为：%d", err, n)
				clientObj.Conn.Close()
			} else {
				errMsg = fmt.Sprintf("读取数据错误：%s，读取到的字节数为：%d", err, n)
			}
			fmt.Printf(errMsg)
			break
		}
		clientObj.AppendReceiveData(readBytes[:n])

		for {
			// 获取有效的消息
			message, exists := clientObj.GetReceiveData()
			if !exists {
				break
			}
			if msgChan != nil {
				resp := ws.SocketResponseObject{}
				_ = json.Unmarshal(message, &resp)
				select {
					case msgChan <- resp:
				default:
				}
			}
			fmt.Printf("ClientId:%v, receive Msg:%s\n", clientObj.Id,message)
		}
	}
}

func newClient()  *ws.Client{
	var msg string
	conn, err := net.DialTimeout("tcp", ":8765", 2*time.Second)
	if err != nil {
		msg = fmt.Sprintf("Dial Error: %s", err)
	} else {
		msg = fmt.Sprintf("Connect to the server. (local address: %s)", conn.LocalAddr())
	}
	fmt.Println(msg)
	return  ws.NewClient(conn)
}

func TestMain(m *testing.M) {
	w.Add(1)
	m.Run()
	w.Wait()
}
