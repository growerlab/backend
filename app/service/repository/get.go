package repository

import (
	"context"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/permission"
	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/backend/app/model/repository"
	"github.com/growerlab/backend/app/service"
)

func GetRepository(ctx context.Context, ownerPath, path string) (*repositoryModel.Repository, error) {
	_, currentUserNSID, err := service.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	if ownerPath == "" {
		return nil, errors.New(errors.InvalidParameterError(errors.Namespace, errors.Path, errors.Empty))
	}
	if path == "" {
		return nil, errors.New(errors.InvalidParameterError(errors.Repository, errors.Path, errors.Empty))
	}

	ns, err := namespace.GetNamespaceByPath(db.DB, ownerPath)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, errors.New(errors.NotFoundError(errors.Namespace))
	}

	repo, err := repositoryModel.GetRepositoryByNsWithPath(db.DB, ns.ID, path)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, errors.New(errors.NotFoundError(errors.Repository))
	}

	err = permission.CheckViewRepository(currentUserNSID, repo.ID)
	if err != nil {
		return nil, err
	}
	return repo, err
}
