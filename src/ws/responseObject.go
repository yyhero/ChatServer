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

func (responseObject *SocketResponseObject) SetDataError() *SocketResponseObject {
	return responseObject.SetResultStatus(model.Con_DataError)
}

func (responseObject *SocketResponseObject) SetAPIDataError() *SocketResponseObject {
	return responseObject.SetResultStatus(model.Con_APIDataError)
}

func (responseObject *SocketResponseObject) SetClientDataError() *SocketResponseObject {
	return responseObject.SetResultStatus(model.Con_APIDataError)
}

func (responseObject *SocketResponseObject) SetResultStatus(rs model.ResultStatus) *SocketResponseObject {
	responseObject.Code = rs
	responseObject.Message = rs.String()

	return responseObject
}

func (responseObject *SocketResponseObject) SetCommandType(ct model.CommandType) *SocketResponseObject {
	responseObject.CommandType = ct

	return responseObject
}

func (responseObject *SocketResponseObject) SetData(data interface{}) *SocketResponseObject {
	responseObject.Data = data

	return responseObject
}

type SocketRequest struct {
	CommandType model.CommandType
	Command map[string]interface{}
}

func NewSocketRequest(commandType model.CommandType, data map[string]interface{})  *SocketRequest{
	//if data == nil{
	//	data = make(map[string]interface{},0)
	//}
	return &SocketRequest{CommandType:commandType, Command:data}
}
