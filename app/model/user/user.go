package user

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

var columns = []string{
	"id",
	"email",
	"encrypted_password",
	"username",
	"name",
	"public_email",
	"created_at",
	"deleted_at",
	"verified_at",
}

var (
	NormalUser     = sq.And{sq.Eq{"deleted_at": nil}, sq.NotEq{"verified_at": nil}}
	InactivateUser = sq.Eq{"verified_at": nil}
	DeletedUser    = sq.NotEq{"deleted_at": nil}
)

func AddUser(tx sqlx.Execer, user *User) error {
	user.CreatedAt = time.Now().UTC()

	sql, _, _ := sq.Insert("user").
		Columns(columns[1:]...).
		Values(
			user.Email,
			user.EncryptedPassword,
			user.Username,
			user.Name,
			user.PublicEmail,
			user.CreatedAt,
			nil,
			nil,
		).ToSql()

	_, err := tx.Exec(sql)
	return errors.Sql(err)
}

func ActivateUser(tx sqlx.Execer, activateToken string) error {
	return nil
}

func ListUsers(src sqlx.Queryer, page, per uint64) ([]*User, error) {
	users := make([]*User, 0)

	// TODO 如果用户量很大的时候，这样分页会有性能问题
	sql, _, _ := sq.Select(columns...).
		From("user").
		Where(NormalUser).
		Limit(per).
		Offset(page * per).
		ToSql()
	err := sqlx.Select(src, &users, sql)
	return users, errors.Sql(err)
}
