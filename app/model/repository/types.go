package repository

type Repository struct {
	ID          int64  `db:"id"`
	UUID        string `db:"uuid"`
	Path        string `db:"path"`
	Name        string `db:"name"`
	NamespaceID int64  `db:"namespace_id"`
	OwnerID     int64  `db:"owner_id"`
	Description string `db:"description"`
	CreatedAt   int64  `db:"created_at"`
	ServerID    int64  `db:"server_id"`
	ServerPath  string `db:"server_path"`
}
