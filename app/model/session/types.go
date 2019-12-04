package session

type Session struct {
	ID        int64  `db:"id,omitempty"`
	UserID    int64  `db:"user_id,omitempty"`
	Token     string `db:"token,omitempty"`
	CreatedAt int64  `db:"created_at,omitempty"`
	ExpiredAt int64  `db:"expired_at,omitempty"`
}
