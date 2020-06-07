package resolver

import (
	"context"

	repoModel "github.com/growerlab/backend/app/model/repository"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/service/repository"
)

func (r *queryResolver) Repositories(ctx context.Context, ownerPath string) ([]*repoModel.Repository, error) {
	return repository.ListRepositories(ctx, ownerPath)
}

func (r *queryResolver) Repository(ctx context.Context, ownerPath, path string) (*repoModel.Repository, error) {
	return repository.GetRepository(ctx, ownerPath, path)
}

func (r *mutationResolver) CreateRepository(ctx context.Context, input service.NewRepositoryPayload) (*service.Result, error) {
	ok, err := repository.CreateRepository(ctx, &input)
	return &service.Result{Ok: ok}, err
}
