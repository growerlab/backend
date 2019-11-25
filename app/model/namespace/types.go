package namespace

type Namespace struct {
	ID      int64  `db:"id"`
	Path    string `db:"path"`
	OwnerId int64  `db:"owner_id"`
}
