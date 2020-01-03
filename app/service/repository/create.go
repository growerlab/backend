package repository

import (
	"strings"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/namespace"
	"github.com/growerlab/backend/app/service"
	"github.com/jmoiron/sqlx"
)

func CreateRepository(reqRepo *service.NewRepository) (bool, error) {
	err := db.Transact(func(tx *db.DBTx) error {
		ns, err := prepre(tx, reqRepo)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return false, nil
}

// prepre
//	req.Path 是否是自己
//	req.Name 名称是否合法、是否重名
func prepre(src sqlx.Queryer, reqRepo *service.NewRepository) (ns *namespace.Namespace, err error) {
	ownerPath := strings.TrimSpace(reqRepo.Path)
	repoName := strings.TrimSpace(reqRepo.Name)
	if len(ownerPath) == 0 {
		err = errors.New(errors.InvalidParameterError(errors.Namespace, errors.Path, errors.Invalid))
		return
	}

	if len(repoName) == 0 {
		err = errors.New(errors.InvalidParameterError(errors.Repository, errors.Name, errors.Invalid))
		return
	}

	ns, err = namespace.GetNamespaceByPath(src, reqRepo.Path)
	if err != nil {
		return nil, err
	}

	// 验证是否重名

	return ns, nil
}
