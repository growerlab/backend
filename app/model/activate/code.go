package activate

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

var tableName = "activate_code"
var columns = []string{
	"id",
	"user_id",
	"code",
	"created_at",
	"used_at",
	"expired_at",
}

func AddCode(tx sqlx.Execer, code *ActivateCode) error {
	code.CreatedAt = time.Now().UTC()

	sql, args, _ := sq.Insert(tableName).
		Columns(columns...).
		Values(
			nil,
			code.UserID,
			code.Code,
			code.CreatedAt,
			nil,
			code.ExpiredAt,
		).ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.Sql(err)
	}
	return nil
}
