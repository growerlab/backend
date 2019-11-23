package activate

import "time"

const CodeMaxLen = 16

type ActivateCode struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	Code      string     `db:"code"`
	CreatedAt time.Time  `db:"created_at"`
	UsedAt    *time.Time `db:"used_at"`
	ExpiredAt time.Time  `db:"expired_at"`
}
