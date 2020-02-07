package context

import (
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/permission/common"
)

var _ common.ContextDelegate = (*Repository)(nil)

type Repository struct {
}

func (r *Repository) Type() int {
	return common.ContextRepository
}

func (r *Repository) TypeLabel() string {
	return "repository"
}

func (r *Repository) Validate(c *ctx.Context) error {
	if c.Param1 <= 0 {
		return errors.Errorf("context param1 is required")
	}
	return nil
}

func (r *Repository) BatchEval(db *ctx.DBContext, args *common.EvalArgs) ([]int64, error) {
	return nil, nil
}
