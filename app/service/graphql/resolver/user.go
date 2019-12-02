package resolver

import (
	"context"

	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/service/user"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input service.NewUserPayload) (*service.Result, error) {
	ok, err := user.Register(&input)
	return &service.Result{Ok: ok}, err
}

func (r *mutationResolver) ActivateUser(ctx context.Context, input service.ActivateCodePayload) (*service.Result, error) {
	ok, err := user.Activate(&input)
	return &service.Result{Ok: ok}, err
}

func (r *mutationResolver) LoginUser(ctx context.Context, input service.LoginUserPayload) (*service.UserToken, error) {
	var session = GetSession(ctx)
	var clientIP string
	if session != nil {
		clientIP = session.GetContext().ClientIP()
	}
	token, err := user.Login(&input, clientIP)
	return &service.UserToken{Token: token}, err
}

func (r *queryResolver) Users(ctx context.Context) ([]*userModel.User, error) {

	// fmt.Println("---- ", currentID)

	panic("not implemented")
}
