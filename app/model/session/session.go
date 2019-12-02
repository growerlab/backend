package session

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/utils"
	"github.com/jmoiron/sqlx"
)

const tableName = "session"

var columns = []string{
	"id",
	"user_id",
	"token",
	"created_at",
	"expired_at",
}

func AddSession(tx sqlx.Queryer, sess *Session) error {
	sql, args, _ := sq.Insert(tableName).
		Columns(columns[1:]...).
		Values(
			sess.UserID,
			sess.Token,
			sess.CreatedAt,
			sess.ExpiredAt,
		).
		Suffix(utils.Returning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&sess.ID)
	return errors.Wrap(err, errors.SQLError())
}
