package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"ChatServer/src/model"
)

// 处理需要客户端发送的数据
func handleSendData(clientObj *Client) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	for {
		//连接是否断开
		if clientObj.getConnStatus() == con_Close {
			break
		}

		handled := false
		for {
			if responseObject, exists := clientObj.getSendData(); exists {
				handled = true
				if err := clientObj.sendMessage(responseObject); err != nil {
					return
				}
			} else {
				break
			}
		}

		if !handled {
			time.Sleep(5 * time.Millisecond)
		}
	}
}

// 处理客户端连接
func HandleConn(clientObj  *Client){
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()
	registerClient(clientObj)
	go handleSendData(clientObj)

	for {
		readBytes := make([]byte, 1024)
		n, err := clientObj.Conn.Read(readBytes)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("conn closed：%s，本次读取的字节数为：%d", err, n)
			} else {
				fmt.Printf("conn closed:%s，本次读取的字节数为：%d\n", err, n)
			}
			break
		}
		clientObj.AppendReceiveData(readBytes[:n])
		HandleReceiveData(clientObj)
	}
}

// 处理客户端收到的数据
func HandleReceiveData(clientObj *Client) {
	for {
		message, exists := clientObj.GetReceiveData()
		if !exists {
			break
		}
		fmt.Printf("receive Msg:%s\n", message)

		if len(message) == 0 {
			continue
		} else {
			handleRequest(clientObj, message)
		}
	}
}

// 处理客户端请求
func handleRequest(clientObj *Client, request []byte) {
	responseObj := NewSocketResponseObject(model.CommandNull)

	defer func() {
		clientObj.appendSendData(responseObj)
	}()

	// 定义变量
	var requestMap map[string]interface{}
	var ok bool
	var err error

	// 解析请求字符串
	if err = json.Unmarshal(request, &requestMap); err != nil {
		fmt.Printf("反序列化出错，错误信息为：%s", err)
		responseObj.SetClientDataError()
		return
	}

	// 解析CommandType
	if commandType_float, ok := requestMap["CommandType"].(float64); !ok {
		fmt.Printf("CommandType不是int类型")
		responseObj.SetClientDataError()
		return
	} else {
		responseObj.SetCommandType(model.CommandType(int(commandType_float)))
	}

	if responseObj.CommandType != model.CommandNull{
		if requestMap["Command"] != nil {
			if _, ok = requestMap["Command"].(map[string]interface{}); !ok {
				fmt.Println("commandMap不是map类型")
				responseObj.SetClientDataError()
				return
			}
		}
	}

	// 调用方法
	requestFuncObj, exists := funcMap[responseObj.CommandType]
	if !exists {
		responseObj.SetResultStatus(model.Con_CommandTypeNotDefined)
		return
	}
	responseObj = requestFuncObj.funcDefinition(clientObj, requestMap)
}
