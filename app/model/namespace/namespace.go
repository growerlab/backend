package namespace

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/jmoiron/sqlx"
)

var tableName = "namespace"
var columns = []string{
	"id",
	"path",
	"owner_id",
}

func AddNamespace(tx sqlx.Execer, ns *Namespace) error {
	sql, args, _ := sq.Insert(tableName).
		Columns(columns[1:]...).
		Values(
			ns.Path,
			ns.OwnerId,
		).
		ToSql()

	ret, err := tx.Exec(sql, args)
	if err != nil {
		return errors.Sql(err)
	}
	ns.ID, err = ret.LastInsertId()
	return errors.Trace(err)
}
