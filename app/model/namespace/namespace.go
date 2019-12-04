package namespace

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/utils"
	"github.com/jmoiron/sqlx"
)

var tableName = "namespace"
var columns = []string{
	"id",
	"path",
	"owner_id",
}

func AddNamespace(tx sqlx.Queryer, ns *Namespace) error {
	sql, args, _ := sq.Insert(tableName).
		Columns(columns[1:]...).
		Values(
			ns.Path,
			ns.OwnerId,
		).
		Suffix(utils.Returning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&ns.ID)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}
