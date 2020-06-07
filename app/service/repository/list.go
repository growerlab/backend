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

func ListRepositories(ctx context.Context, owner string) ([]*repositoryModel.Repository, error) {
	_, currentUserNSID, err := service.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	ns, err := namespace.GetNamespaceByPath(db.DB, owner)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, errors.New(errors.NotFoundError(errors.Namespace))
	}

	repositories, err := repositoryModel.ListRepositoriesByNamespace(db.DB, ns.ID)
	if err != nil {
		return nil, err
	}

	var result []*repositoryModel.Repository
	for _, repo := range repositories {
		if err := permission.CheckViewRepository(currentUserNSID, repo.NamespaceID); err == nil {
			result = append(result, repo)
		}
	}

	return result, nil
}
