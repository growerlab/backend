package user

import (
	"time"

	"github.com/growerlab/backend/app/model/namespace"
)

var columns = []string{
	"id",
	"email",
	"encrypted_password",
	"username",
	"name",
	"public_email",
	"created_at",
	"deleted_at",
	"verified_at",
}

type User struct {
	ID                int        `db:"id"`
	Email             string     `db:"email"`
	EncryptedPassword string     `db:"encrypted_password"`
	Username          string     `db:"username"`
	Name              string     `db:"name"`
	PublicEmail       string     `db:"public_email"`
	CreatedAt         time.Time  `db:"created_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
	VerifiedAt        *time.Time `db:"verified_at"`
}

func (u *User) Namespace() *namespace.Namespace {
	return &namespace.Namespace{}
}
