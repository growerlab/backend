// KeyDB / Redis 配置

package db

import (
	"net"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/growerlab/backend/app/utils/conf"
)

var CacheDB *redis.Pool
var QueueDB *redis.Pool

func InitMemDB() error {
	var config = conf.GetConf().Redis
	CacheDB = newPool(config, config.CacheDB)
	QueueDB = newPool(config, config.QueueDB)
	return nil
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
