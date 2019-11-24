package resolver

import (
	"context"

	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/service/user"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input service.NewUserPayload) (*service.Result, error) {
	ok, err := user.RegisterUser(&input)
	return &service.Result{Ok: ok}, err
}

func (r *mutationResolver) ActivateUser(ctx context.Context, input service.AcitvateCodePayload) (*service.Result, error) {
	ok, err := user.ActivateUser(&input)
	return &service.Result{Ok: ok}, err
}

func (r *queryResolver) Users(ctx context.Context) ([]*userModel.User, error) {

	// fmt.Println("---- ", currentID)

	panic("not implemented")
}
