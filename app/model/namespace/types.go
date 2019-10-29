package namespace

type Namespace struct {
	ID      int    `db:"id"`
	Path    string `db:"path"`
	OwnerId int64  `db:"owner_id"`
}
