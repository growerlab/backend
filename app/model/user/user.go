package user

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

var tableName = "user"
var tableNameMark = `"user"` // user 是 pgsql中的保留关键字，所以加上引号

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

	sql, args, _ := sq.Insert(tableNameMark).
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

	ret, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, errors.SqlError)
	}
	user.ID, err = ret.LastInsertId()
	return errors.WithStack(err)
}

func AreEmailOrUsernameInUser(src sqlx.Queryer, username, email string) (bool, error) {
	if len(username) > 0 {
		user, err := getUser(src, sq.Eq{"username": username})
		if err != nil {
			return false, errors.WithStack(err)
		}
		if user != nil {
			return true, nil
		}
	}
	if len(email) > 0 {
		user, err := getUser(src, sq.Eq{"email": email})
		if err != nil {
			return false, errors.WithStack(err)
		}
		if user != nil {
			return true, nil
		}
	}
	return false, nil
}

func getUser(src sqlx.Queryer, cond sq.Sqlizer) (*User, error) {
	sql, args, _ := sq.Select(columns...).
		From(tableNameMark).
		Where(cond).
		Limit(1).
		ToSql()

	result := make([]*User, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

func ActivateUser(tx sqlx.Execer, userID int64) error {
	sql, args, _ := sq.Update(tableNameMark).
		Set("verified_at", time.Now().UTC()).
		Where(sq.Eq{"id": userID}).
		ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func ListUsers(src sqlx.Queryer, page, per uint64) ([]*User, error) {
	users := make([]*User, 0)

	// TODO 如果用户量很大的时候，这样分页会有性能问题.. 希望能碰到那一天👀
	sql, _, _ := sq.Select(columns...).
		From(tableNameMark).
		Where(NormalUser).
		Limit(per).
		Offset(page * per).
		ToSql()

	err := sqlx.Select(src, &users, sql)
	return users, errors.Wrap(err, errors.SqlError)
}
