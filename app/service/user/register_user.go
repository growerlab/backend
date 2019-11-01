package user

import (
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/utils/pwd"
	"github.com/jmoiron/sqlx"
	"gopkg.in/asaskevich/govalidator.v9"
)

const (
	PasswordLenMin = 8
	PasswordLenMax = 32

	UsernameLenMin = 4
	UsernameLenMax = 40
)

func validateRegisterUser(payload *service.NewUserPayload) error {
	if !govalidator.IsEmail(payload.Email) {
		return errors.InvalidParameterError(errors.User, errors.Email, errors.InvalidParameter)
	}
	if !govalidator.IsByteLength(payload.Password, PasswordLenMin, PasswordLenMax) {
		return errors.InvalidParameterError(errors.User, errors.Password, errors.InvalidPassword)
	}
	if !govalidator.IsByteLength(payload.Username, UsernameLenMin, UsernameLenMax) {
		return errors.InvalidParameterError(errors.User, errors.Username, errors.InvalidParameter)
	}

	return nil
}

func buildUser(payload *service.NewUserPayload) (*userModel.User, error) {
	password, err := pwd.GeneratePassword(payload.Password)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &userModel.User{
		Email:             payload.Email,
		EncryptedPassword: password,
		Username:          payload.Username,
		Name:              payload.Username,
		PublicEmail:       payload.Email,
	}, nil
}

// 用户注册
// 1. 将用户信息添加到数据库中
// 2. 发送验证邮件（这里可以考虑使用KeyDB来建立邮件发送队列，避免重启进程后，发送任务丢失）
// 3. Done
//
func RegisterUser(payload *service.NewUserPayload) (*userModel.User, error) {
	var err error
	err = validateRegisterUser(payload)
	if err != nil {
		return nil, errors.Trace(err)
	}

	err = db.Transact(func(tx *sqlx.Tx) error {
		newUser, err := buildUser(payload)
		if err != nil {
			return errors.Trace(err)
		}

		err = userModel.AddUser(tx, newUser)
		// TODO 发送激活邮件
		return errors.Trace(err)
	})
	return nil, errors.Trace(err)
}
