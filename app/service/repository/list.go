package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/service/common/session"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/permission"
	"github.com/growerlab/backend/app/model/db"
	namespaceModel "github.com/growerlab/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/backend/app/model/repository"
)

func ListRepositories(c *gin.Context, namespace string) ([]*repositoryModel.Repository, error) {
	currentUserNSID := session.New(c).UserNamespace()

	ns, err := namespaceModel.GetNamespaceByPath(db.DB, namespace)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, errors.NotFoundError(errors.Namespace)
	}

	repositories, err := repositoryModel.ListRepositoriesByNamespace(db.DB, ns.ID)
	if err != nil {
		return nil, err
	}

	var result []*repositoryModel.Repository
	for _, repo := range repositories {
		err := permission.CheckViewRepository(currentUserNSID, repo.ID)
		if err == nil {
			result = append(result, repo)
		}
	}

	return result, nil
}
