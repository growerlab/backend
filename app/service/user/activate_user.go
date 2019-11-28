package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/activate"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/utils/conf"
	"github.com/growerlab/backend/app/utils/logger"
	"github.com/growerlab/backend/app/utils/uuid"
)

// 激活账号的前期准备
// 生成code
// 生成url
// 生成模版
// 发送邮件
//
func DoPreActivateUser(tx *db.DBTx, userID int64) error {
	code := buildActivateCode(userID)
	err := activate.AddCode(tx, code)
	if err != nil {
		return errors.Trace(err)
	}

	activateURL := buildActivateURL(code.Code)
	logger.Info("the activate url: %v", activateURL)

	// TODO 生成邮件模版(邮件模版功能应该抽出来独立，并能适配未来的其他模版)
	// TODO 发送邮件

	return nil
}

// 验证用户邮箱激活码
//
func DoActivateUser(tx *db.DBTx, code string) (bool, error) {
	acode, err := activate.GetCode(tx, code)
	if err != nil {
		return false, errors.Trace(err)
	}
	// 是否已使用过
	if acode.UsedAt != nil {
		return false, errors.New(errors.P(errors.ActivateCode, errors.Code, errors.Used))
	}
	// 是否过期
	if acode.ExpiredAt.Unix() < time.Now().UTC().Unix() {
		return false, errors.New(errors.P(errors.ActivateCode, errors.Code, errors.Expired))
	}
	// 将code改成已使用
	err = activate.UpdateCodeUsed(tx, code)
	if err != nil {
		return false, errors.Trace(err)
	}
	// 激活用户状态
	err = user.ActivateUser(tx, acode.UserID)
	if err != nil {
		return false, errors.Trace(err)
	}
	return true, nil
}

func buildActivateURL(code string) string {
	baseURL := conf.GetConf().WebsiteURL
	partURL := fmt.Sprintf("activate_user?code=%s", code)
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}
	return baseURL + partURL
}

func buildActivateCode(userID int64) *activate.ActivateCode {
	code := new(activate.ActivateCode)
	code.UserID = userID
	code.Code = uuid.UUIDv16()
	code.ExpiredAt = time.Now().UTC().Add(24 * time.Hour)
	return code
}
