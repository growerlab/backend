// KeyDB / Redis 配置

package db

import (
	"net"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils/conf"
)

var CacheDB *redis.Client
var QueueDB *redis.Client

func InitMemDB() error {
	var config = conf.GetConf().Redis
	CacheDB = newPool(config, config.CacheDB)
	QueueDB = newPool(config, config.QueueDB)

	// Test
	reply, err := CacheDB.Ping().Result()
	if err != nil || reply != "PONG" {
		return errors.New("memdb not ready")
	}
	return err
}

func newPool(cfg *conf.Redis, db int) *redis.Client {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	idleTimeout := time.Duration(cfg.IdleTimeout) * time.Second

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           db,
		PoolSize:     cfg.MaxActive,
		MinIdleConns: cfg.MaxIdle,
		IdleTimeout:  idleTimeout,
	})
	return client
}
