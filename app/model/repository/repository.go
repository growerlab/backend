package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/utils"
	"github.com/jmoiron/sqlx"
)

var (
	table   = "repository"
	columns = []string{
		"id",
		"uuid",
		"path",
		"name",
		"namespace_id",
		"owner_id",
		"description",
		"created_at",
		"server_id",
		"server_path",
	}
)

func AddRepository(tx sqlx.Queryer, repo *Repository) error {
	sql, args, _ := sq.Insert(table).
		Columns(columns[1:]...).
		Values(
			repo.UUID,
			repo.Path,
			repo.Name,
			repo.NamespaceID,
			repo.OwnerID,
			repo.Description,
			repo.CreatedAt,
			repo.ServerID,
			repo.ServerPath,
		).
		Suffix(utils.SqlReturning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&repo.ID)
	if err != nil {
		return errors.New(errors.SQLError())
	}
	return nil
}

func AreNameInNamespace(src sqlx.Queryer, namespaceID int64, name string) (bool, error) {
	where := sq.And{
		sq.Eq{"namespace_id": namespaceID},
		sq.Eq{"path": name},
	}
	sql, args, _ := sq.Select(columns[0]).
		From(table).
		Where(where).
		ToSql()

	result := make([]int, 0)
	err := sqlx.Select(src, &result, sql, args...)
	if err != nil {
		return false, errors.Wrap(err, errors.SQLError())
	}
	if len(result) > 0 {
		return true, nil
	}
	return false, nil
}
