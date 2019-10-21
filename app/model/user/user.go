package user

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

func AddUser(tx sqlx.Execer, user *User) error {
	user.CreatedAt = time.Now().UTC()

	sql, _, _ := squirrel.Insert("user").Columns(columns...).
		Values(
			user.Email,
			user.EncryptedPassword,
			user.Name,
			user.PublicEmail,
			user.CreatedAt,
			nil,
		).ToSql()

	_, err := tx.Exec(sql)
	return errors.Sql(err)
}

func ListUsers(src sqlx.Queryer, page, per uint64) ([]*User, error) {
	users := make([]*User, 0)

	// TODO 如果用户量很大的时候，这样分页会有性能问题
	sql, _, _ := squirrel.Select(columns...).
		From("user").
		Limit(per).
		Offset(page * per).
		ToSql()
	err := sqlx.Select(src, &users, sql)
	return users, errors.Sql(err)
}
