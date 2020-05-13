package repository

import (
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/namespace"
	"github.com/growerlab/backend/app/model/user"
)

type Repository struct {
	ID          int64  `db:"id"`
	UUID        string `db:"uuid"`         // 全站唯一ID（fork时用到）
	Path        string `db:"path"`         // 在namespace中是唯一的name
	Name        string `db:"name"`         // 目前与path字段相同
	NamespaceID int64  `db:"namespace_id"` // 仓库属于个人，还是组织
	OwnerID     int64  `db:"owner_id"`     // 仓库创建者
	Description string `db:"description"`
	CreatedAt   int64  `db:"created_at"`
	ServerID    int64  `db:"server_id"`   // 服务器id
	ServerPath  string `db:"server_path"` // 服务器中的绝对路径
	Public      int    `db:"public"`      // 共有

	ns    *namespace.Namespace
	owner *user.User
}

func (r *Repository) Namespace() *namespace.Namespace {
	if r.ns != nil {
		return r.ns
	}
	r.ns, _ = namespace.GetNamespace(db.DB, r.NamespaceID)
	return r.ns
}

func (r *Repository) Owner() *user.User {
	if r.owner != nil {
		return r.owner
	}
	r.owner, _ = user.GetUser(db.DB, r.OwnerID)
	return r.owner
}

func (r *Repository) IsPublic() bool {
	return r.Public == int(StatusPublic)
}
