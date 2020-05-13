package user

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	sessionModel "github.com/growerlab/backend/app/model/session"
	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/utils/pwd"
	"github.com/growerlab/backend/app/utils/uuid"
	"github.com/jmoiron/sqlx"
)

const TokenExpiredTime = 24 * time.Hour * 30 // 30天过期

// 用户登录
//  用户邮箱是否已验证
//	更新用户最后的登录时间/IP
//	生成用户登录token
//
func Login(input *service.LoginUserPayload, ctx *gin.Context) (
	result *service.UserLoginResult,
	err error,
) {
	clientIP := ctx.ClientIP()

	err = db.Transact(func(tx *db.DBTx) error {
		user, err := Validate(tx, input.Email, input.Password)
		if err != nil {
			return err
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
		ctx.SetCookie("token", sess.Token, 0, "/", ctx.Request.Host, false, false)
		return err
	})
	if err != nil {
		return nil, err
	}
	return
}

func Validate(src sqlx.Queryer, usernameOrEmail, password string) (user *userModel.User, err error) {
	if strings.Contains(usernameOrEmail, "@") {
		user, err = userModel.GetUserByEmail(src, usernameOrEmail)
		if err != nil {
			return nil, err
		}
	} else {
		user, err = userModel.GetUserByUsername(src, usernameOrEmail)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, errors.New(errors.NotFoundError(errors.User))
	}
	if !user.Verified() {
		return nil, errors.New(errors.AccessDenied(errors.User, errors.NotActivated))
	}

	ok := pwd.ComparePassword(user.EncryptedPassword, password)
	if !ok {
		return nil, errors.New(errors.InvalidParameterError(errors.User, errors.Password, errors.NotEqual))
	}
	return user, nil
}

func buildSession(userID int64, clientIP string) *sessionModel.Session {
	return &sessionModel.Session{
		OwnerID:   userID,
		Token:     uuid.UUID(),
		ClientIP:  clientIP,
		CreatedAt: time.Now().Unix(),
		ExpiredAt: time.Now().Add(TokenExpiredTime).Unix(),
	}
}
