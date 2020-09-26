package model

// 定义聊天命令类型对象
type CommandType int

const (
	CommandNull CommandType = 1 + iota

	// 创建房间
	CreateRoom

	// 发送消息
	SendMessage

	// 加入房间
	JoinRoom

	// cmd
	Command
)
