package ws

import (
	"sync"
)

var (
	// 收到的总数据大小，以B为单位
	totalReceiveSize int64

	// 发送的总数据大小，以B为单位
	totalSendSize int64

	// 客户端连接列表
	clientMap = make(map[int32]*Client, 1024)

	// 读写锁
	mutex sync.RWMutex
)

// 添加新的客户端
// clientObj：客户端对象
func registerClient(clientObj *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	clientMap[clientObj.GetId()] = clientObj
}

// 根据客户端Id获取对应的客户端对象
// id：客户端Id
// 返回值：客户端对象
func GetClient(id int32) (*Client, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	if clientObj, ok := clientMap[id]; ok {
		return clientObj, true
	}

	return nil, false
}



