package ws

// 客户端连接状态
type ConnStatus int

const (
	con_Open ConnStatus = 1 + iota

	con_Close
)
