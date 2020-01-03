package namespace

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/utils"
	"github.com/jmoiron/sqlx"
)

var table = "namespace"
var columns = []string{
	"id",
	"path",
	"owner_id",
	"type",
}

func AddNamespace(tx sqlx.Queryer, ns *Namespace) error {
	sql, args, _ := sq.Insert(table).
		Columns(columns[1:]...).
		Values(
			ns.Path,
			ns.OwnerId,
			ns.Type,
		).
		Suffix(utils.Returning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&ns.ID)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}

func GetNamespaceByPath(src sqlx.Queryer, path string) (*Namespace, error) {
	sql, args, _ := sq.Select(columns...).From(table).Where(sq.Eq{"path": path}).ToSql()

	result := make([]*Namespace, 0, 1)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}
