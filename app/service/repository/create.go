package repository

import (
	"strings"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/namespace"
	"github.com/growerlab/backend/app/model/repository"
	"github.com/growerlab/backend/app/service"
	"github.com/jmoiron/sqlx"
)

func CreateRepository(reqRepo *service.NewRepository) (bool, error) {
	err := db.Transact(func(tx *db.DBTx) error {
		ns, err := validateAndPrepre(tx, reqRepo)
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
func validateAndPrepre(src sqlx.Queryer, reqRepo *service.NewRepository) (ns *namespace.Namespace, err error) {
	reqRepo.Path = strings.TrimSpace(reqRepo.Path)
	reqRepo.Name = strings.TrimSpace(reqRepo.Name)
	if len(reqRepo.Path) == 0 {
		err = errors.New(errors.InvalidParameterError(errors.Namespace, errors.Path, errors.Invalid))
		return
	}

	if len(reqRepo.Name) == 0 {
		err = errors.New(errors.InvalidParameterError(errors.Repository, errors.Name, errors.Invalid))
		return
	}

	ns, err = namespace.GetNamespaceByPath(src, reqRepo.Path)
	if err != nil {
		return nil, err
	}

	// 验证仓库名在当前namespace中是否已存在
	exists, err := repository.AreNameInNamespace(src, ns.ID, reqRepo.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(errors.AlreadyExistsError(errors.Repository, errors.AlreadyExists))
	}
	return ns, nil
}
