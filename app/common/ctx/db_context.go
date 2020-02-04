package ctx

import (
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
)

type DBContext struct {
	Src   sqlx.Queryer
	MemDB *redis.Client
}
