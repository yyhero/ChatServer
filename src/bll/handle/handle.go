package handle

import (
	"ChatServer/src/bll/roomMgr"
	"ChatServer/src/bll/sensitive"
	"ChatServer/src/model"
	"ChatServer/src/ws"
	"fmt"
	"strconv"
	"time"
)

func init() {
	ws.RegisterFuncMap(model.JoinRoom, "JoinRoom", handlerJoinRoomByRoomId)

	ws.RegisterFuncMap(model.CreateRoom, "CreateRoom", handlerCreateRoom)

	ws.RegisterFuncMap(model.SendMessage, "SendMessage", handlerSendMsg)

	ws.RegisterFuncMap(model.Command, "Command", handlerCommand)
}


func handlerCreateRoom(clientObj *ws.Client, parameters map[string]interface{}) *ws.SocketResponseObject {
	ct, _ := parameters["CommandType"].(float64)

	room := roomMgr.RoomManger.CreateRoom(clientObj)
	data := make(map[string]interface{})
	data["RoomId"] = room.RoomId
	clientObj.SetLoginTs(time.Now().Unix())

	responseObj := ws.NewSocketResponseObject(model.CommandType(ct))
	responseObj.SetData(data)
	return responseObj
}

func handlerJoinRoomByRoomId(clientObj *ws.Client, parameters map[string]interface{}) *ws.SocketResponseObject {
	ct, _ := parameters["CommandType"].(float64)
	commandMap, _ := parameters["Command"].(map[string]interface{})
	responseObj := ws.NewSocketResponseObject(model.CommandType(ct))

	roomId := int64(commandMap["RoomId"].(float64))
	room,exist := roomMgr.RoomManger.JoinByRoomId(roomId, clientObj)
	if !exist{
		responseObj.SetDataError()
		return responseObj
	}
	clientObj.SetLoginTs(time.Now().Unix())
	data := room.GetHistory()
	if len(data) >0{
		responseObj.SetData(data)
	}

	return responseObj
}

func handlerSendMsg(clientObj *ws.Client, parameters map[string]interface{}) *ws.SocketResponseObject {
	ct, _ := parameters["CommandType"].(float64)
	commandMap, _ := parameters["Command"].(map[string]interface{})
	msg := commandMap["Msg"].(string)
	responseObj := ws.NewSocketResponseObject(model.CommandType(ct))

	room,exist := roomMgr.RoomManger.GetRoomByPlayer(clientObj)
	if !exist{
		fmt.Printf("i am here")
		responseObj.SetResultStatus(model.Con_PlayerNotExist)
		return responseObj
	}

	str := sensitive.Replace(msg, '*')
	data := room.AddToMsgs(clientObj, str)

	players := room.GetAllPlayerExceptMe(clientObj)
	for _, p := range players{
		obj:= ws.NewSocketResponseObject(model.CommandType(ct))
		obj.SetData(data)
		p.SendMsg(obj)
	}

	return responseObj
}

func handlerCommand(clientObj *ws.Client, parameters map[string]interface{}) *ws.SocketResponseObject {
	ct, _ := parameters["CommandType"].(model.CommandType)
	commandMap, _ := parameters["Command"].(map[string]interface{})
	cmd := commandMap["Cmd"].(string)
	responseObj := ws.NewSocketResponseObject(ct)

	switch  cmd{
	case "/stats":
		paras := commandMap["CmdPara"].([]interface{})
		para := int32(paras[0].(float64))
		client,exsit := ws.GetClient(para)
		if !exsit{
			responseObj.SetResultStatus(model.Con_PlayerNotExist)
			return responseObj
		}
		paraStr := strconv.Itoa(int(para))
		str := cmd + " " + paraStr + "    " +  client.GetLoginTs()
		responseObj.SetData(str)
		return responseObj

	case "/popular":
		room,exist := roomMgr.RoomManger.GetRoomByPlayer(clientObj)
		if !exist{
			responseObj.SetResultStatus(model.Con_PlayerNotExist)
			return responseObj
		}else {
			data := room.GetPopularWorld()
			responseObj.SetData(data)
		}

	default:
		responseObj.SetDataError()
	}

	return responseObj
}

