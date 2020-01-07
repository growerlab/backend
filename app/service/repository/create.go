package repository

import (
	"context"
	"strings"
	"time"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/namespace"
	"github.com/growerlab/backend/app/model/repository"
	"github.com/growerlab/backend/app/model/server"
	"github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/utils/regex"
	"github.com/growerlab/backend/app/utils/uuid"
	"github.com/jmoiron/sqlx"
)

func CreateRepository(ctx context.Context, req *service.NewRepository) (bool, error) {
	currentUser, err := service.CurrentUser(ctx)
	if err != nil {
		return false, err
	}

	err = db.Transact(func(tx *db.DBTx) error {
		ns, err := validateAndPrepre(tx, currentUser.ID, req)
		if err != nil {
			return err
		}

		srv, err := server.RandNormalServer(tx)
		if err != nil {
			return err
		}

		repo := buildRepository(currentUser, ns, req, srv)
		err = repository.AddRepository(tx, repo)
		if err != nil {
			return err
		}

		// 真正创建仓库
		api, err := NewApi(srv, repo)
		if err != nil {
			return err
		}
		err = api.Repository().Create()
		return err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func buildRepository(
	currentUser *user.User,
	ns *namespace.Namespace,
	req *service.NewRepository,
	srv *server.Server,
) (repo *repository.Repository) {

	repo = &repository.Repository{
		NamespaceID: ns.ID,
		UUID:        uuid.UUIDv16(),
		Path:        req.Name,
		Name:        req.Name,
		OwnerID:     currentUser.ID,
		Description: "",
		CreatedAt:   time.Now().Unix(),
		ServerID:    srv.ID,
		ServerPath:  UsernameToFilePath(ns.Path, req.Name),
	}
	return repo
}

// validate
//	req.NamespacePath  TODO 这里暂时只验证namespace的owner_id 是否为用户，未来应该验证组织权限（比如是否可以选择这个组织创建仓库）
//	req.Name 名称是否合法、是否重名
func validateAndPrepre(src sqlx.Queryer, userID int64, req *service.NewRepository) (ns *namespace.Namespace, err error) {
	req.NamespacePath = strings.TrimSpace(req.NamespacePath)
	req.Name = strings.TrimSpace(req.Name)
	if len(req.NamespacePath) == 0 {
		err = errors.New(errors.InvalidParameterError(errors.Namespace, errors.Path, errors.Invalid))
		return
	}

	if len(req.Name) == 0 {
		err = errors.New(errors.InvalidParameterError(errors.Repository, errors.Name, errors.Invalid))
		return
	}

	ns, err = namespace.GetNamespaceByPath(src, req.NamespacePath)
	if err != nil {
		return nil, err
	}

	// TODO 未来应该验证权限(例如是否有权限在组织中创建权限)
	if ns.OwnerId != userID {
		return nil, errors.New(errors.AccessDenied(errors.User, errors.NotEqual))
	}

	// 验证仓库名是否合法
	if !regex.Match(req.Name, regex.RepositoryNameRegex) {
		return nil, errors.New(errors.InvalidParameterError(errors.Repository, errors.Name, errors.Invalid))
	}

	// 验证仓库名在当前namespace中是否已存在
	exists, err := repository.AreNameInNamespace(src, ns.ID, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(errors.AlreadyExistsError(errors.Repository, errors.AlreadyExists))
	}
	return ns, nil
}