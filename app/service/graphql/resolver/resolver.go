package resolver

import (
	"context"

	"github.com/growerlab/backend/app/model/namespace"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Namespace() NamespaceResolver {
	return &namespaceResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

type namespaceResolver struct{ *Resolver }

func (r *namespaceResolver) ID(ctx context.Context, obj *namespace.Namespace) (string, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }
