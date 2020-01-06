package repository

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
}
