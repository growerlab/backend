package queue

import (
	"github.com/gomodule/redigo/redis"
	"github.com/growerlab/backend/app/model/db"
)

func NewList(pool *redis.Pool) *List {
	return &List{
		queuePool: pool,
	}
}

type List struct {
	queuePool *redis.Pool
}

func (l *List) Push(key string, payload []byte) (err error) {
	err = db.MemConn(l.queuePool, func(conn redis.Conn) error {
		_ = conn.Send("MULTI")
		_ = conn.Send("LPUSH", key, payload)
		_ = conn.Send("EXPIRE", key, 24*60*60) // 24h
		_, err = conn.Do("EXEC")
		return err
	})
	return
}

func (l *List) Pop(key string) (payload []byte, err error) {
	err = db.MemConn(l.queuePool, func(conn redis.Conn) error {
		payload, err = redis.Bytes(conn.Do("RPOP", key))
		if err == redis.ErrNil {
			err = nil
			return nil
		}
		return err
	})
	return
}

func (l *List) Release(key string) {
	_ = db.MemConn(l.queuePool, func(conn redis.Conn) error {
		_, err := conn.Do("DEL", key)
		return err
	})
}
