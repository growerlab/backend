package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/utils"
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
		"public",
	}
)

func AddRepository(tx db.HookQueryer, repo *Repository) error {
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
			repo.Public,
		).
		Suffix(utils.SqlReturning("id")).
		ToSql()

	err := tx.QueryRowx(sql, args...).Scan(&repo.ID)
	if err != nil {
		return errors.Wrap(err, errors.SQLError())
	}
	return nil
}

func AreNameInNamespace(src db.HookQueryer, namespaceID int64, name string) (bool, error) {
	where := sq.And{
		sq.Eq{"namespace_id": namespaceID},
		sq.Eq{"path": name},
	}
	result, err := listRepositoriesByCond(src, []string{"id"}, where)
	if err != nil {
		return false, err
	}
	return len(result) > 0, nil
}

func ListRepositoriesByNamespace(src db.HookQueryer, namespaceID int64) ([]*Repository, error) {
	where := sq.And{sq.Eq{"namespace_id": namespaceID}}
	return listRepositoriesByCond(src, columns, where)
}

func GetRepositoryByNsWithPath(src db.HookQueryer, namespaceID int64, path string) (*Repository, error) {
	where := sq.And{sq.Eq{"namespace_id": namespaceID, "path": path}}
	repos, err := listRepositoriesByCond(src, columns, where)
	if err != nil {
		return nil, err
	}
	if len(repos) > 0 {
		return repos[0], nil
	}
	return nil, nil
}

func GetRepository(src db.HookQueryer, id int64) (*Repository, error) {
	repos, err := listRepositoriesByCond(src, columns, sq.Eq{"id": id})
	if err != nil {
		return nil, err
	}
	if len(repos) > 0 {
		return repos[0], nil
	}
	return nil, nil
}

func listRepositoriesByCond(src db.HookQueryer, tableColumns []string, cond sq.Sqlizer) ([]*Repository, error) {
	where := cond
	sql := sq.Select(tableColumns...).
		From(table).
		Where(where)

	result := make([]*Repository, 0)
	err := src.Select(&result, sql)
	if err != nil {
		return nil, errors.Wrap(err, errors.SQLError())
	}
	return result, nil
}
