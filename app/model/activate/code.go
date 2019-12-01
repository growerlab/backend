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
		Columns(columns[1:]...).
		Values(
			code.UserID,
			code.Code,
			code.CreatedAt,
			nil,
			code.ExpiredAt,
		).ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}

func GetCode(src sqlx.Queryer, code string) (*ActivateCode, error) {
	sql, args, _ := sq.Select(columns...).
		From(tableName).
		Where(sq.Eq{"code": code}).
		Limit(1).
		ToSql()

	var data = make([]*ActivateCode, 0)
	err := sqlx.Select(src, &data, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}
	if len(data) > 0 {
		return data[0], nil
	}
	return nil, nil
}

// 修改code状态
//
func UpdateCodeUsed(tx sqlx.Execer, code string) error {
	sql, args, _ := sq.Update(tableName).
		Set("used_at", time.Now().UTC()).
		Where(sq.Eq{"code": code}).
		ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}
