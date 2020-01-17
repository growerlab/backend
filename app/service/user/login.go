package user

import (
	"time"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	sessionModel "github.com/growerlab/backend/app/model/session"
	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/utils/pwd"
	"github.com/growerlab/backend/app/utils/uuid"
)

const TokenExpiredTime = 24 * time.Hour * 30 // 30天过期

// 用户登录
//  用户邮箱是否已验证
//	更新用户最后的登录时间/IP
//	生成用户登录token
//
func Login(input *service.LoginUserPayload, clientIP string) (
	result *service.UserLoginResult,
	err error,
) {
	err = db.Transact(func(tx *db.DBTx) error {
		user, err := userModel.GetUserByEmail(tx, input.Email)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New(errors.NotFoundError(errors.User))
		}
		if !user.Verified() {
			return errors.New(errors.AccessDenied(errors.User, errors.NotActivated))
		}

		ok := pwd.ComparePassword(user.EncryptedPassword, input.Password)
		if !ok {
			return errors.New(errors.InvalidParameterError(errors.User, errors.Password, errors.NotEqual))
		}

		err = userModel.UpdateLogin(tx, user.ID, clientIP)
		if err != nil {
			return err
		}

		// 生成TOKEN返回给客户端
		sess := buildSession(user.ID, clientIP)
		err = sessionModel.AddSession(tx, sess)
		if err != nil {
			return err
		}

		// namespace
		ns := user.Namespace()
		result = &service.UserLoginResult{
			Token:         sess.Token,
			NamespacePath: ns.Path,
			Name:          user.Name,
			Email:         user.Email,
			PublicEmail:   user.PublicEmail,
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return
}

func buildSession(userID int64, clientIP string) *sessionModel.Session {
	return &sessionModel.Session{
		UserID:    userID,
		Token:     uuid.UUID(),
		ClientIP:  clientIP,
		CreatedAt: time.Now().Unix(),
		ExpiredAt: time.Now().Add(TokenExpiredTime).Unix(),
	}
}
