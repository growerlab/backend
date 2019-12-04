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
func Login(input *service.LoginUserPayload, clientIP string) (token string, err error) {
	err = db.Transact(func(tx *db.DBTx) error {
		user, err := userModel.GetUserByEmail(tx, input.Email)
		if err != nil {
			return err
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
		sess := buildSession(user.ID)
		err = sessionModel.AddSession(tx, sess)
		if err != nil {
			return err
		}
		token = sess.Token
		return err
	})

	if err != nil {
		return "", err
	}
	return "", nil
}

func buildSession(userID int64) *sessionModel.Session {
	return &sessionModel.Session{
		UserID:    userID,
		Token:     uuid.UUID(),
		CreatedAt: time.Now().UTC(),
		ExpiredAt: time.Now().UTC().Add(TokenExpiredTime),
	}
}