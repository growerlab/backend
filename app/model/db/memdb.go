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

var MemDB *redis.Client
var PermissionDB *redis.Client

func InitMemDB() error {
	var config = conf.GetConf().Redis
	MemDB = newPool(config, 0)
	PermissionDB = newPool(config, 0)

	// Test
	reply, err := MemDB.Ping().Result()
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

type KeyBuilder struct {
	sb *strings.Builder
}

func NewKeyBuilder(base string) *KeyBuilder {
	b := &KeyBuilder{sb: &strings.Builder{}}
	b.sb.WriteString(base)
	return b
}

func (b *KeyBuilder) Append(s ...string) *KeyBuilder {
	for i := range s {
		b.sb.WriteString(":")
		b.sb.WriteString(s[i])
	}
	return b
}

func (b *KeyBuilder) String() string {
	return b.sb.String()
}

func BaseKeyBuilder(s ...string) *KeyBuilder {
	return NewKeyBuilder(conf.GetConf().Redis.Namespace).Append(s...)
}
