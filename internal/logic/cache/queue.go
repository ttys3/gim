package cache

import (
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const (
	RoomQueue = "room_queue" // 房间消息队列
	AllQueue  = "all_queue"  // 全服消息队列
)

type queue struct{}

var Queue = new(queue)

func (queue) Publish(topic string, bytes []byte) error {
	_, err := db.RedisCli.Publish(topic, bytes).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}
