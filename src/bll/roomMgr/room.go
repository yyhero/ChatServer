package roomMgr

import (
	"ChatServer/src/model"
	"ChatServer/src/ws"
	"strings"
	"sync"
	"time"
)

var(
	validTime = int64(5)
)

type Room struct {
	RoomId int64
	players   map[int32]*ws.Client
	msgs []*model.Message
	msgLen int
	mutex sync.RWMutex
}

func newRoom(id int64)  *Room{
	return &Room{
		players:make(map[int32]*ws.Client),
		msgs:make([]*model.Message, 0),
		RoomId:id,
		msgLen:1000,
	}
}

func (this *Room)AddPlayer(player *ws.Client) {
	this.mutex.RLock()
	defer  this.mutex.RUnlock()
	if _, exist := this.players[player.Id]; exist{
		return
	}
	this.players[player.Id] = 	player
}

func (this *Room)GetHistory() []*model.Message{
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	length := len(this.msgs)
	if length<= 50{
		return this.msgs
	}
	return this.msgs[length-50:]
}

func (this *Room) AddToMsgs(player *ws.Client, info string) *model.Message{
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	if len(this.msgs) >= this.msgLen{
		this.msgs = this.msgs[1:]
	}

	msg := &model.Message{Ts:time.Now().Unix(), Msg:info}
	this.msgs = append(this.msgs, msg)
	return msg
}

func (this *Room) GetAllPlayerExceptMe(player *ws.Client) []*ws.Client{
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	result := make([]*ws.Client, 0)
	for _, obj := range this.players{
		if obj.Id != player.Id{
			result = append(result, obj)
		}
	}
	return result
}


func (this *Room) GetPopularWorld() *model.Word{
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	maxWord := &model.Word{Msg:"", Num:0}
	msgs := make(map[string]int64)
	curTs := time.Now().Unix()
	for _, obj := range this.msgs{
		if curTs - obj.Ts <= validTime{
			words := strings.Fields(obj.Msg)
			for _, word := range  words{
				if _, exist := msgs[word]; !exist{
					msgs[word] = 1
				}else{
					msgs[word]++
				}

				if msgs[word] > maxWord.Num{
					maxWord.Num = msgs[word]
					maxWord.Msg = word
				}
			}
		}
	}
	return maxWord
}
