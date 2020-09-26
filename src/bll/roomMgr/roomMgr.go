package roomMgr

import (
	"ChatServer/src/ws"
	"sync"
)

var (
	RoomManger = &RoomMgr{rooms:make(map[int64]*Room), globalRoomId:1000}
)

type RoomMgr struct {
	globalRoomId int64
	rooms  map[int64]*Room
	mutex  sync.RWMutex
}

func (this *RoomMgr)  GetRoomByPlayer(player *ws.Client) (*Room, bool){
	this.mutex.RLock()
	defer  this.mutex.RUnlock()

	room ,exist := this.rooms[player.RoomId]
	return room,exist
}

func (this *RoomMgr)  CreateRoom(player *ws.Client) *Room{
	this.mutex.RLock()
	defer  this.mutex.RUnlock()

	this.globalRoomId++
	player.SetRoomId(this.globalRoomId)
	room := newRoom(this.globalRoomId)
	room.AddPlayer(player)
	this.rooms[room.RoomId] = room

	return room
}

func (this *RoomMgr)  JoinByRoomId(roomId int64, player *ws.Client) (*Room, bool){
	this.mutex.RLock()
	defer  this.mutex.RUnlock()

	room, exist := this.rooms[roomId]
	if !exist{
		return nil, false
	}else{
		room.AddPlayer(player)
	}
	return room, true
}
