package session

import "time"

type Session struct {
	ID        int64     `db:"id,omitempty"`
	UserID    int64     `db:"user_id,omitempty"`
	Token     string    `db:"token,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	ExpiredAt time.Time `db:"expired_at,omitempty"`
}
