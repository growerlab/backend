// KeyDB / Redis 配置

package db

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils/conf"
)

var CacheDB *redis.Client
var QueueDB *redis.Client
var PermissionDB *redis.Client

func InitMemDB() error {
	var config = conf.GetConf().Redis
	CacheDB = newPool(config, 0)
	QueueDB = newPool(config, 0)
	PermissionDB = newPool(config, 0)

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

type baseKeyBuilder struct {
	sb *strings.Builder
}

func newBaseKeyBuilder(base string) *baseKeyBuilder {
	b := &baseKeyBuilder{sb: &strings.Builder{}}
	b.sb.WriteString(base)
	return b
}

func (b *baseKeyBuilder) Append(s ...string) *baseKeyBuilder {
	for i := range s {
		b.sb.WriteString(":")
		b.sb.WriteString(s[i])
	}
	return b
}

func (b *baseKeyBuilder) String() string {
	return b.sb.String()
}

func BaseKeyBuilder(s ...string) *baseKeyBuilder {
	return newBaseKeyBuilder(conf.GetConf().Redis.Namespace).Append(s...)
}
