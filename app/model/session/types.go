package session

type Session struct {
	ID        int64  `db:"id"`
	UserID    int64  `db:"user_id"`
	Token     string `db:"token"`
	ClientIP  string `db:"client_ip"` // 未来可能用来检验token是否被劫持
	CreatedAt int64  `db:"created_at"`
	ExpiredAt int64  `db:"expired_at"`
}
