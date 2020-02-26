package repository

import (
	"context"

	"github.com/growerlab/backend/app/model/db"
	"github.com/growerlab/backend/app/model/namespace"
	repositoryModel "github.com/growerlab/backend/app/model/repository"
	"github.com/growerlab/backend/app/service"
)

// TODO 权限判断（公私项目区分）；分页功能；
func ListRepositories(ctx context.Context, ownerPath string) ([]*repositoryModel.Repository, error) {
	currentUser, err := service.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	currentNamespace := currentUser.Namespace()

	// TODO 目前没有权限判断，所以目前只能取 currentUser.Namespace == namespaceID 的仓库（既自己的仓库）；以及其他人的公共仓库
	// TODO 目前不支持”组织“
	onlyPublic := currentNamespace.Path != ownerPath
	state := repositoryModel.StatusAll
	if onlyPublic {
		state = repositoryModel.StatusPublic
	} else {
		state = repositoryModel.StatusAll
	}

	ns, err := namespace.GetNamespaceByPath(db.DB, ownerPath)
	if err != nil {
		return nil, err
	}

	repositories, err := repositoryModel.ListRepositoriesByNamespace(db.DB, state, ns.ID)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}
