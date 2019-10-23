package user

import "time"

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
	ID                int64      `db:"id,omitempty"`
	Email             string     `db:"email,omitempty"`
	EncryptedPassword string     `db:"encrypted_password,omitempty"`
	Username          string     `db:"username,omitempty"`
	Name              string     `db:"name,omitempty"`
	PublicEmail       *string    `db:"public_email,omitempty"`
	CreatedAt         time.Time  `db:"created_at,omitempty"`
	DeletedAt         *time.Time `db:"deleted_at,omitempty"`
	VerifiedAt        *time.Time `db:"verified_at,omitempty"`
}
