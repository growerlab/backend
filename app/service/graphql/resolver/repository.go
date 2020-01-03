package resolver

import (
	"context"

	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/service/repository"
)

func (r *mutationResolver) CreateRepository(ctx context.Context, input service.NewRepository) (*service.Result, error) {
	ok, err := repository.CreateRepository(&input)
	return &service.Result{Ok: ok}, err
}
