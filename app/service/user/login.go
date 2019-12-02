package user

import (
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/utils/pwd"
)

// 用户登录
//  用户邮箱是否已验证
//	更新用户最后的登录时间/IP
//	生成用户登录token
//
func Login(input *service.LoginUserPayload, clientIP string) (token string, err error) {
	err = db.Transact(func(tx *db.DBTx) error {
		user, err := userModel.GetUserByEmail(tx, input.Email)
		if err != nil {
			return errors.Trace(err)
		}

		ok := pwd.ComparePassword(user.EncryptedPassword, input.Password)
		if !ok {
			return errors.New(errors.InvalidParameterError(errors.User, errors.Password, errors.NotEqual))
		}

		err = userModel.UpdateLogin(tx, user.ID, clientIP)
		if err != nil {
			return errors.Trace(err)
		}

		// TODO 生成TOKEN返回给客户端

		return errors.Trace(err)
	})

	if err != nil {
		return "", errors.Trace(err)
	}
	return "", nil
}
