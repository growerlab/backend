package context

import (
	"github.com/growerlab/backend/app/common/permission/common"
	"github.com/growerlab/backend/app/common/permission/delegate"
)

var _ delegate.ContextDelegate = (*Repository)(nil)

type Repository struct {
}

func (r *Repository) Type() int {
	return common.ContextRepository
}

func (r *Repository) TypeLabel() string {
	return "repository"
}

func (r *Repository) BatchEval(args *delegate.EvalArgs) ([]int64, error) {
	return nil, nil
}
