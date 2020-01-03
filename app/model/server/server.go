package server

import (
	"math/rand"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

var table = "server"

var columns = []string{
	"id",
	"summary",
	"host",
	"port",
	"status",
	"created_at",
	"deleted_at",
}

// RandNormalServer 当有多个服务器时，随机返回一个服务器
func RandNormalServer(src sqlx.Queryer) (*Server, error) {
	servers, err := ListServers(src, StatusNormal)
	if err != nil {
		return nil, err
	}
	length := len(servers)
	if length == 0 {
		return nil, nil
	}
	rand.Shuffle(length, func(i, j int) {
		servers[i], servers[j] = servers[j], servers[i]
	})
	return servers[0], nil
}

func ListServers(src sqlx.Queryer, statues ...statusType) ([]*Server, error) {
	or := sq.Or{SqlStatusNormal}
	where := sq.And{SqlNormal, &or}

	if len(statues) > 0 {
		for i := range statues {
			switch statues[i] {
			case StatusClosed:
				or = append(or, SqlStatusClosed)
			case StatusSuspend:
				or = append(or, SqlStatusSuspend)
			case StatusNormal:
				// default
			}
		}
	}

	sql, args, _ := sq.Select(columns...).
		From(table).
		Where(where).
		ToSql()

	result := make([]*Server, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}
	return result, nil
}
