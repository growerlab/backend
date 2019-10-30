package resolver

import (
	"context"

	userModel "github.com/growerlab/backend/app/model/user"
	"github.com/growerlab/backend/app/service"
	"github.com/growerlab/backend/app/service/user"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input service.NewUserPayload) (*userModel.User, error) {
	return user.RegisterUser(&input)
}

func (r *queryResolver) Users(ctx context.Context) ([]*userModel.User, error) {

	// fmt.Println("---- ", currentID)

	panic("not implemented")
}
