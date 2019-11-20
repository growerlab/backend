package queue

import (
	"github.com/garyburd/redigo/redis"
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
		_, err = conn.Do("LPUSH", key, payload)
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
