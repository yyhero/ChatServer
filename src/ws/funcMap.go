package ws

import (
	"ChatServer/src/model"
)

var (
	// 所有对外提供的方法列表
	funcMap = make(map[model.CommandType]*requestFunc)
)

// 请求方法对象
type requestFunc struct {
	funcName string
	funcDefinition func(*Client, map[string]interface{}) *SocketResponseObject
}

// 创建新的请求方法对象
func newRequestFunc(_funcName string, _funcDefinition func(*Client, map[string]interface{}) *SocketResponseObject) *requestFunc {
	return &requestFunc{
		funcName:       _funcName,
		funcDefinition: _funcDefinition,
	}
}

//注册回调函数
func RegisterFuncMap(commandType model.CommandType, funcName string, funcDefinition func(*Client, map[string]interface{}) *SocketResponseObject) {
	funcMap[commandType] = newRequestFunc(funcName, funcDefinition)
}
