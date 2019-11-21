package user

import (
	"time"

	"github.com/growerlab/backend/app/model/activate"
	"github.com/growerlab/backend/app/utils/uuid"
	"github.com/jmoiron/sqlx"
)

// 生成code
// 生成url
// 生成模版
// 发送邮件
//
func DoPreActivateUser(tx *sqlx.Tx, userID int64) error {
	// code := buildActivateCode(userID)

	return nil
}

func buildActivateCode(userID int64) *activate.ActivateCode {
	code := new(activate.ActivateCode)
	code.UserID = userID
	code.Code = uuid.UUIDv16()
	code.ExpiredAt = time.Now().UTC().Add(24 * time.Hour)
	return code
}
