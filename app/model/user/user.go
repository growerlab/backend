package user

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/session"
	"github.com/growerlab/backend/app/model/utils"
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
	"last_login_at",
	"last_login_ip",
	"register_ip",
}

func AddUser(tx sqlx.Queryer, user *User) error {
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
			nil,
			nil,
			user.RegisterIP,
		).
		Suffix(utils.Returning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&user.ID)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}

func AreEmailOrUsernameInUser(src sqlx.Queryer, username, email string) (bool, error) {
	if len(username) > 0 {
		user, err := getUser(src, sq.Eq{"username": username})
		if err != nil {
			return false, err
		}
		if user != nil {
			return true, nil
		}
	}
	if len(email) > 0 {
		user, err := getUser(src, sq.Eq{"email": email})
		if err != nil {
			return false, err
		}
		if user != nil {
			return true, nil
		}
	}
	return false, nil
}

func GetUserByEmail(src sqlx.Queryer, email string) (*User, error) {
	user, err := getUser(src, sq.Eq{"email": email})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getUser(src sqlx.Queryer, cond sq.Sqlizer) (*User, error) {
	sql, args, _ := sq.Select(columns...).
		From(tableNameMark).
		Where(sq.And{cond, NormalUser}).
		Limit(1).
		ToSql()

	result := make([]*User, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

func ActivateUser(tx sqlx.Execer, userID int64) error {
	sql, args, _ := sq.Update(tableNameMark).
		Set("verified_at", time.Now().Unix()).
		Where(sq.And{sq.Eq{"id": userID}, InactivateUser}).
		ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}

func ListAllUsers(src sqlx.Queryer, page, per uint64) ([]*User, error) {
	users := make([]*User, 0)

	// TODO 如果用户量很大的时候，这样分页会有性能问题.. 希望能碰到那一天👀
	sql, _, _ := sq.Select(columns...).
		From(tableNameMark).
		Where(NormalUser).
		Limit(per).
		Offset(page * per).
		ToSql()

	err := sqlx.Select(src, &users, sql)
	return users, errors.Wrap(err, errors.SQLError())
}

func UpdateLogin(tx sqlx.Execer, userID int64, clientIP string) error {
	sql, args, _ := sq.Update(tableNameMark).
		SetMap(map[string]interface{}{
			"last_login_at": time.Now().Unix(),
			"last_login_ip": clientIP,
		}).
		Where(sq.Eq{"id": userID}).
		ToSql()

	_, err := tx.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}

func GetUserByUserToken(src sqlx.Queryer, userToken string) (*User, error) {
	sql, args, _ := sq.Select(columns...).
		From(tableNameMark).
		Join(fmt.Sprintf("%s ON %s.token = ? AND %s.expired_at <= ?", session.TableName(), session.TableName(), session.TableName()), userToken, time.Now().Unix()).
		Where(fmt.Sprintf("%s.id = %s.user_id", tableNameMark, session.TableName())).
		ToSql()

	users := make([]*User, 0, 1)

	err := sqlx.Select(src, &users, sql, args...)
	if err != nil {
		return nil, errors.New(errors.SQLError())
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, nil
}
