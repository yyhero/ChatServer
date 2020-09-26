package ws

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// 包头的长度
	con_HEADER_LENGTH = 4
)

var (
	// 全局客户端的id，从1开始进行自增
	globalClientId int32 = 1000

	// 字节的大小端顺序
	byterOrder = binary.LittleEndian
)

// 获得自增的id值
func getIncrementId() int32 {
	atomic.AddInt32(&globalClientId, 1)
	return globalClientId
}

type Client struct {
	// 唯一标识
	Id int32

	RoomId int64

	LoginTs int64

	// 客户端连接对象
	Conn net.Conn

	//连接状态
	connStatus ConnStatus

	// 接收到的消息内容
	receiveData []byte

	// 待发送的数据
	sendData []*SocketResponseObject

	mutex sync.Mutex
}

// 新建客户端对象
// conn：连接对象
// 返回值：客户端对象的指针
func NewClient(_conn net.Conn) *Client {
	return &Client{
		Id:                   getIncrementId(),
		Conn:                 _conn,
		connStatus:           con_Open,
		receiveData:          make([]byte, 0, 1024),
		sendData:             make([]*SocketResponseObject, 0, 16),
	}
}

// 获取唯一标识
func (c *Client) GetId() int32 {
	return c.Id
}

func (c *Client) SetRoomId(id int64) {
	c.RoomId = 	id
}

func (c *Client) SetLoginTs(ts int64) {
	c.LoginTs = ts
}


func (c *Client) GetLoginTs() string{
	loginTime := time.Unix(c.LoginTs, 0)
	d, h, m, s := loginTime.Day(), loginTime.Hour(), loginTime.Minute(), loginTime.Second()
	return fmt.Sprintf("%02dd %02dh %02dm %02ds", d, h, m, s)
}

// 获取远程地址（IP_Port）
func (clientObj *Client) getRemoteAddr() string {
	items := strings.Split(clientObj.Conn.RemoteAddr().String(), ":")
	return fmt.Sprintf("%s_%s", items[0], items[1])
}

// 获取远程地址（IP）
func (clientObj *Client) getRemoteShortAddr() string {
	items := strings.Split(clientObj.Conn.RemoteAddr().String(), ":")
	return items[0]
}

// 追加发送的数据
func (clientObj *Client) appendSendData(responseObj *SocketResponseObject) {
	clientObj.mutex.Lock()
	defer clientObj.mutex.Unlock()

	clientObj.sendData = append(clientObj.sendData, responseObj)
}

func (clientObj *Client) SendMsg(responseObj *SocketResponseObject) {
	clientObj.appendSendData(responseObj)
}

// 追加接收到的数据
func (clientObj *Client) AppendReceiveData(receiveData []byte) {
	clientObj.receiveData = append(clientObj.receiveData, receiveData...)
	atomic.AddInt64(&totalReceiveSize, int64(len(receiveData)))
}

// 获取有效的消息
func (c *Client) GetReceiveData() ([]byte, bool) {
	// 判断是否包含头部信息
	if len(c.receiveData) < con_HEADER_LENGTH {
		return nil, false
	}

	// 获取头部信息
	header := c.receiveData[:con_HEADER_LENGTH]

	// 将头部数据转换为内部的长度
	contentLength := BytesToInt32(header, byterOrder)

	// 判断长度是否满足
	if len(c.receiveData)-con_HEADER_LENGTH < int(contentLength) {
		return nil, false
	}

	// 提取消息内容
	content := c.receiveData[con_HEADER_LENGTH : con_HEADER_LENGTH+contentLength]

	// 将对应的数据截断，以得到新的数据
	c.receiveData = c.receiveData[con_HEADER_LENGTH+contentLength:]
	return content, true
}

// 发送字节数组消息
func (clientObj *Client) SendBytes(b []byte) error {
	// 获得数组的长度
	contentLength := len(b)

	// 将长度转化为字节数组
	header := Int32ToBytes(int32(contentLength), byterOrder)

	// 将头部与内容组合在一起
	message := append(header, b...)

	// 增加发送量(包括包头的长度+内容的长度)
	atomic.AddInt64(&totalSendSize, int64(con_HEADER_LENGTH+contentLength))

	// 发送消息
	_, err := clientObj.Conn.Write(message)
	if err != nil {
		fmt.Printf("发送消息,%s,出现错误：%s", b, err)
	}

	return err
}

// 发送字节数组消息
func (clientObj *Client) sendMessage(responseObj *SocketResponseObject) error {
	b, err := json.Marshal(*responseObj)
	if err != nil {
		return errors.New("序列化response数据失败")
	}
	return  clientObj.SendBytes(b)
}

// 获取待发送的数据
func (clientObj *Client) getSendData() (responseObj *SocketResponseObject, exists bool) {
	clientObj.mutex.Lock()
	defer clientObj.mutex.Unlock()

	// 如果没有数据则直接返回
	if len(clientObj.sendData) == 0 {
		return
	}

	responseObj = clientObj.sendData[0]
	exists = true
	clientObj.sendData = clientObj.sendData[1:]
	return
}

// 获取连接状态
func (clientObj *Client) getConnStatus() ConnStatus {
	return clientObj.connStatus
}

// 设置连接状态
func (clientObj *Client) setConnStatus(status ConnStatus) {
	clientObj.connStatus = status
}


// 格式化
func (clientObj *Client) String() string {
	return fmt.Sprintf("{Id:%d, RemoteAddr:%d}", clientObj.Id, clientObj.getRemoteAddr())
}

func BytesToInt32(b []byte, order binary.ByteOrder) int32 {
	bytesBuffer := bytes.NewBuffer(b)

	var result int32
	binary.Read(bytesBuffer, order, &result)

	return result
}

func Int32ToBytes(n int32, order binary.ByteOrder) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, order, n)

	return bytesBuffer.Bytes()
}
