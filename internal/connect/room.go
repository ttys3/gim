package connect

import (
	"container/list"
	"gim/pkg/pb"
	"sync"
)

var RoomsManager sync.Map

// subscribedRoom 订阅房间
func SubscribedRoom(conn *Conn, roomId string) {
	if roomId == conn.RoomId {
		return
	}

	oldRoomId := conn.RoomId
	// 取消订阅
	if oldRoomId != "" {
		value, ok := RoomsManager.Load(oldRoomId)
		if !ok {
			return
		}
		room := value.(*Room)
		room.Unsubscribe(conn)

		if room.Conns.Front() == nil {
			RoomsManager.Delete(oldRoomId)
		}
		return
	}

	// 订阅
	if roomId != "" {
		value, ok := RoomsManager.Load(roomId)
		var room *Room
		if !ok {
			room = NewRoom(roomId)
			RoomsManager.Store(roomId, room)
		} else {
			room = value.(*Room)
		}
		room.Subscribe(conn)
		return
	}
}

func PushToRoom(roomId string, message pb.MessageSend) {
	value, ok := RoomsManager.Load(roomId)
	if !ok {
		return
	}

	value.(*Room).Push(message)
}

type Room struct {
	RoomId string     // 房间ID
	Conns  *list.List // 订阅房间消息的连接
	lock   sync.RWMutex
}

func NewRoom(roomId string) *Room {
	return &Room{
		RoomId: roomId,
		Conns:  list.New(),
	}
}

// subscribe 订阅房间
func (r *Room) Subscribe(conn *Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()

	conn.Element = r.Conns.PushBack(conn)
	conn.RoomId = r.RoomId
}

// unsubscribe 取消订阅
func (r *Room) Unsubscribe(conn *Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Conns.Remove(conn.Element)
	conn.Element = nil
	conn.RoomId = ""
}

// Push 推送消息到房间
func (r *Room) Push(message pb.MessageSend) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	element := r.Conns.Front()
	for {
		conn := element.Value.(*Conn)
		conn.Send(pb.PackageType_PT_MESSAGE, 0, nil, &message)

		element = element.Next()
		if element == nil {
			break
		}
	}
}
