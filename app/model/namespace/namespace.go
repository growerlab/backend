package namespace

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/utils"
)

var table = "namespace"
var columns = []string{
	"id",
	"path",
	"owner_id",
	"type",
}

func AddNamespace(tx db.HookQueryer, ns *Namespace) error {
	sql, args, _ := sq.Insert(table).
		Columns(columns[1:]...).
		Values(
			ns.Path,
			ns.OwnerID,
			ns.Type,
		).
		Suffix(utils.SqlReturning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&ns.ID)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}

func GetNamespaceByPath(src db.HookQueryer, path string) (*Namespace, error) {
	return getNamespaceByCond(src, sq.Eq{"path": path})
}

func GetNamespaceByOwnerID(src db.HookQueryer, ownerID int64) (*Namespace, error) {
	return getNamespaceByCond(src, sq.Eq{"owner_id": ownerID})
}

func GetNamespace(src db.HookQueryer, id int64) (*Namespace, error) {
	return getNamespaceByCond(src, sq.Eq{"id": id})
}

func getNamespaceByCond(src db.HookQueryer, cond sq.Sqlizer) (*Namespace, error) {
	ns, err := listNamespaceByCond(src, cond)
	if err != nil {
		return nil, err
	}
	if len(ns) > 0 {
		return ns[0], nil
	}
	return nil, nil
}

func ListNamespacesByOwner(src db.HookQueryer, userType NamespaceType, ownerIDs ...int64) ([]*Namespace, error) {
	where := sq.And{
		sq.Eq{"owner_id": ownerIDs},
		sq.Eq{"type": userType},
	}
	return listNamespaceByCond(src, where)
}

func listNamespaceByCond(src db.HookQueryer, cond sq.Sqlizer) ([]*Namespace, error) {
	sql := sq.Select(columns...).From(table).Where(cond)
	result := make([]*Namespace, 0)

	err := src.Select(&result, sql)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}
	return result, nil
}
