package model

// 服务端响应结果的状态
type ResultStatus int

const (
	// 成功
	Con_Success ResultStatus = -1 * iota

	// 数据错误
	Con_DataError

	// API数据错误
	Con_APIDataError

	// 命令类型未定义
	Con_CommandTypeNotDefined

	// 玩家不存在
	Con_PlayerNotExist

)

var status = []string{
	"Success",
	"DataError",
	"APIDataError",
	"CommandTypeNotDefined",
	"PlayerNotExist",
}

func (rs ResultStatus) String() string {
	return status[rs*-1]
}
