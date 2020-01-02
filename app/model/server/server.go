package server

import (
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
