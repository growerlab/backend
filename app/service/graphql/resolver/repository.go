package resolver

import (
	"context"

	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/service/repository"
)

func (r *mutationResolver) CreateRepository(ctx context.Context, input service.NewRepositoryPayload) (*service.Result, error) {
	ok, err := repository.CreateRepository(ctx, &input)
	return &service.Result{Ok: ok}, err
}
