package repository

import (
	"context"

	"github.com/growerlab/backend/app/model/db"
	repositoryModel "github.com/growerlab/backend/app/model/repository"
	"github.com/growerlab/backend/app/service"
)

// TODO 权限判断（公私项目区分）；分页功能；
func ListRepositories(ctx context.Context, namespaceID int64) ([]*repositoryModel.Repository, error) {
	currentUser, err := service.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	ns := currentUser.Namespace()

	// TODO 目前没有权限判断，所以目前只能取 currentUser.Namespace == namespaceID 的仓库（既自己的仓库）；以及其他人的公共仓库
	// TODO 目前不支持”组织“
	onlyPublic := ns.ID != namespaceID
	state := repositoryModel.StatePublic
	if onlyPublic {
		state = repositoryModel.StatePublic
	} else {
		state = repositoryModel.StatePrivate
	}
	repositories, err := repositoryModel.ListRepositoriesByNamespace(db.DB, state, ns.ID)
	if err != nil {
		return nil, err
	}

	return buildServiceRepositories(repositories), nil
}

func buildServiceRepositories(repos []*repositoryModel.Repository) []*repositoryModel.Repository {
	result := make([]*repositoryModel.Repository, 0)
	if len(repos) == 0 {
		return result
	}
	for _, repo := range repos {
		result = append(result, &repositoryModel.Repository{
			UUID:        repo.UUID,
			Path:        repo.Path,
			Name:        repo.Name,
			Description: repo.Description,
			CreatedAt:   repo.CreatedAt,
		})
	}
	return result
}
