package user

import (
	"github.com/growerlab/backend/app/common/errors"
	activateModel "github.com/growerlab/backend/app/model/activate"
	"github.com/growerlab/backend/app/model/db"
	nsModel "github.com/growerlab/backend/app/model/namespace"
	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/utils/pwd"
	"github.com/growerlab/backend/app/utils/regex"
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
		return errors.New(errors.P(errors.User, errors.Email, errors.InvalidParameter))
	}
	if !govalidator.IsByteLength(payload.Password, PasswordLenMin, PasswordLenMax) {
		return errors.New(errors.P(errors.User, errors.Password, errors.InvalidLength))
	}
	if !govalidator.IsByteLength(payload.Username, UsernameLenMin, UsernameLenMax) {
		return errors.New(errors.P(errors.User, errors.Username, errors.InvalidLength))
	}
	if !regex.Match(payload.Username, regex.UsernameRegex) {
		return errors.New(errors.P(errors.User, errors.Username, errors.InvalidParameter))
	}
	if !regex.Match(payload.Password, regex.PasswordRegex) {
		return errors.New(errors.P(errors.User, errors.Password, errors.InvalidParameter))
	}

	// 不允许使用的关键字
	if _, invalidUsername := userModel.InvalidUsernameSet[payload.Username]; invalidUsername {
		return errors.New(errors.AlreadyExistsError(errors.User, ""))
	}

	// email, username是否已经存在
	exists, err := userModel.AreEmailOrUsernameInUser(db.DB, payload.Username, payload.Email)
	if err != nil {
		return errors.Trace(err)
	}
	if exists {
		return errors.New(errors.AlreadyExistsError(errors.User, ""))
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

func buildNamespace(user *userModel.User) *nsModel.Namespace {
	return &nsModel.Namespace{
		Path:    user.Username,
		OwnerId: user.ID,
	}
}

// 用户注册
// 1. 将用户信息添加到数据库中
// 2. 发送验证邮件（这里可以考虑使用KeyDB来建立邮件发送队列，避免重启进程后，发送任务丢失）
// 3. Done
//
func RegisterUser(payload *service.NewUserPayload) (bool, error) {
	var err error
	err = validateRegisterUser(payload)
	if err != nil {
		return false, errors.Trace(err)
	}

	err = db.Transact(func(tx *db.DBTx) error {
		user, err := buildUser(payload)
		if err != nil {
			return errors.Trace(err)
		}

		err = userModel.AddUser(tx, user)
		if err != nil {
			return errors.Trace(err)
		}

		// 创建namespace
		ns := buildNamespace(user)
		err = nsModel.AddNamespace(tx, ns)
		if err != nil {
			return errors.Trace(err)
		}

		// 激活用户
		err = DoPreActivateUser(tx, user.ID)
		if err != nil {
			return errors.Trace(err)
		}
		return nil
	})
	if err != nil {
		return false, errors.Trace(err)
	}
	return true, nil
}

// 激活用户
func ActivateUser(payload *service.AcitvateCodePayload) (result bool, err error) {
	if !govalidator.IsByteLength(payload.Code, activateModel.CodeMaxLen, activateModel.CodeMaxLen) {
		return false, errors.New(errors.P(errors.ActivateCode, errors.Code, errors.InvalidParameter))
	}

	err = db.Transact(func(tx *db.DBTx) error {
		result, err = DoActivateUser(tx, payload.Code)
		return err
	})
	return
}
