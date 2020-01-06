package user

import (
	"github.com/growerlab/backend/app/model/namespace"
)

type User struct {
	ID                int64   `db:"id"`
	Email             string  `db:"email"`
	EncryptedPassword string  `db:"encrypted_password"`
	Username          string  `db:"username"`
	Name              string  `db:"name"`
	PublicEmail       string  `db:"public_email"`
	CreatedAt         int64   `db:"created_at"`
	DeletedAt         *int64  `db:"deleted_at"`
	VerifiedAt        *int64  `db:"verified_at"`
	LastLoginAt       *int64  `db:"last_login_at"`
	LastLoginIP       *string `db:"last_login_ip"`
	RegisterIP        string  `db:"register_ip"`
}

func (u *User) Namespace() *namespace.Namespace {
	return &namespace.Namespace{}
}
