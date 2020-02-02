package resolver

import (
	"context"

	"github.com/growerlab/backend/app/model/namespace"
)

func (r *namespaceResolver) ID(ctx context.Context, obj *namespace.Namespace) (string, error) {
	return "", nil
}
