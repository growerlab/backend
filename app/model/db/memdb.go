// KeyDB / Redis 配置

package db

import (
	"net"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/growerlab/backend/app/utils/conf"
	"github.com/growerlab/backend/app/utils/logger"
)

var CacheDB *redis.Pool
var QueueDB *redis.Pool

func InitMemDB() error {
	var config = conf.GetConf().Redis
	CacheDB = newPool(config, config.CacheDB)
	QueueDB = newPool(config, config.QueueDB)

	// Test
	MemConn(CacheDB, func(conn redis.Conn) error {
		reply, err := redis.String(conn.Do("PING"))
		if reply != "PONG" {
			panic(err)
		}
		return nil
	})
	return nil
}

func MemConn(pool *redis.Pool, callback func(redis.Conn) error) {
	conn := pool.Get()
	defer conn.Close()

	err := callback(conn)
	if err != nil {
		logger.Error("mem db callback() has err: %v", err)
	}
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
