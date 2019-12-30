// KeyDB / Redis 配置

package db

import (
	"net"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils/conf"
)

var CacheDB *redis.Pool
var QueueDB *redis.Pool

func InitMemDB() error {
	var config = conf.GetConf().Redis
	CacheDB = newPool(config, config.CacheDB)
	QueueDB = newPool(config, config.QueueDB)

	// Test
	err := MemConn(CacheDB, func(conn redis.Conn) error {
		reply, err := redis.String(conn.Do("PING"))
		if reply != "PONG" {
			return errors.New("memdb not ready")
		}
		return err
	})
	return err
}

func MemConn(pool *redis.Pool, callback func(redis.Conn) error) error {
	conn := pool.Get()
	defer conn.Close()

	return callback(conn)
}

func newPool(cfg *conf.Redis, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			dbOpt := redis.DialDatabase(db)
			conn, err := redis.Dial("tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)), dbOpt)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
