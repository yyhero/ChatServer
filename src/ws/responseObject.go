package ws

import (
	"ChatServer/src/model"
)

// Socket服务器的响应对象
type SocketResponseObject struct {
	// 响应结果的状态值
	Code model.ResultStatus

	// 响应结果的状态值所对应的描述信息
	Message string

	// 响应结果的数据
	Data interface{}

	// 响应结果对应的请求命令类型
	CommandType model.CommandType
}


func NewSocketResponseObject(ct model.CommandType) *SocketResponseObject {
	return &SocketResponseObject{
		Code:        model.Con_Success,
		Message:     model.Con_Success.String(),
		Data:        nil,
		CommandType: ct,
	}
}

func (this *SocketResponseObject) SetDataError() *SocketResponseObject {
	return this.SetResultStatus(model.Con_DataError)
}

func (this *SocketResponseObject) SetAPIDataError() *SocketResponseObject {
	return this.SetResultStatus(model.Con_APIDataError)
}

func (this *SocketResponseObject) SetClientDataError() *SocketResponseObject {
	return this.SetResultStatus(model.Con_APIDataError)
}

func (this *SocketResponseObject) SetResultStatus(rs model.ResultStatus) *SocketResponseObject {
	this.Code = rs
	this.Message = rs.String()
	return this
}

func (this *SocketResponseObject) SetCommandType(ct model.CommandType) *SocketResponseObject {
	this.CommandType = ct
	return this
}

func (this *SocketResponseObject) SetData(data interface{}) *SocketResponseObject {
	this.Data = data
	return  this
}

type SocketRequest struct {
	CommandType model.CommandType
	Command map[string]interface{}
}

func NewSocketRequest(commandType model.CommandType, data map[string]interface{})  *SocketRequest{
	return &SocketRequest{CommandType:commandType, Command:data}
}
