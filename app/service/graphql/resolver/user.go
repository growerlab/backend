package resolver

import (
	"context"

	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/service/graphql"
	"github.com/growerlab/backend/app/service/user"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input service.NewUserPayload) (*service.Result, error) {
	var session = graphql.GetSession(ctx)
	var clientIP string
	if session != nil {
		clientIP = session.GetContext().ClientIP()
	}
	ok, err := user.Register(&input, clientIP)
	return &service.Result{Ok: ok}, err
}

func (r *mutationResolver) ActivateUser(ctx context.Context, input service.ActivationCodePayload) (*service.Result, error) {
	ok, err := user.Activate(&input)
	return &service.Result{Ok: ok}, err
}

func (r *mutationResolver) LoginUser(ctx context.Context, input service.LoginUserPayload) (*service.UserToken, error) {
	var session = graphql.GetSession(ctx)
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
