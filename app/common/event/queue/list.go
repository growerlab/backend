package queue

import (
	"github.com/go-redis/redis/v7"
	"time"
)

func NewList(pool *redis.Client) *List {
	return &List{
		client: pool,
	}
}

type List struct {
	client *redis.Client
}

func (l *List) Push(key string, payload []byte) (err error) {
	pipe := l.client.TxPipeline()
	_ = pipe.LPush(key, payload)
	_ = pipe.Expire(key, time.Duration(24*60*60)*time.Second)
	_, err = pipe.Exec()
	return err
}

func (l *List) Pop(key string) (payload []byte, err error) {
	result := l.client.RPop(key)
	if result.Err() == redis.Nil {
		return
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result.Bytes()
}

func (l *List) Release(key string) {
	l.client.Del(key)
}
