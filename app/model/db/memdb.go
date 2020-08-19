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

var MemDB *MemDBClient
var PermissionDB *MemDBClient

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

func newPool(cfg *conf.Redis, db int) *MemDBClient {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	idleTimeout := time.Duration(cfg.IdleTimeout) * time.Second

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           db,
		PoolSize:     cfg.MaxActive,
		MinIdleConns: cfg.MaxIdle,
		IdleTimeout:  idleTimeout,
	})

	memDB := &MemDBClient{
		client,
		NewKeyBuilder(conf.GetConf().Redis.Namespace),
	}
	return memDB
}

type MemDBClient struct {
	*redis.Client
	*KeyBuilder
}

type KeyBuilder struct {
	namespaceKey string
}

func NewKeyBuilder(namespaceKey string) *KeyBuilder {
	return &KeyBuilder{
		namespaceKey: namespaceKey,
	}
}

func (b *KeyBuilder) PartMaker() *KeyPart {
	var sb strings.Builder
	sb.WriteString(b.namespaceKey)
	sb.WriteString(sep)

	return &KeyPart{
		sb: sb,
	}
}

const sep = ":"

type KeyPart struct {
	sb strings.Builder
}

func NewKeyPart() *KeyPart {
	return &KeyPart{}
}

func (b *KeyPart) Append(s ...string) *KeyPart {
	if len(s) == 0 {
		return b
	}

	b.sb.WriteString(s[0])
	for i := range s[1:] {
		b.sb.WriteString(sep)
		b.sb.WriteString(s[i])
	}
	return b
}

func (b *KeyPart) String() string {
	return b.sb.String()
}
