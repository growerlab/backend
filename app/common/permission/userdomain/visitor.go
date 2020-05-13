package userdomain

import (
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/permission/common"
)

var _ common.UserDomainDelegate = (*Visitor)(nil)

type Visitor struct {
}

func (s *Visitor) Type() int {
	return common.UserDomainVisitor
}

func (s *Visitor) TypeLabel() string {
	return "everyone"
}

func (s *Visitor) Validate(*ctx.UserDomain) error {
	return nil
}

func (s *Visitor) BatchEval(db *ctx.DBContext, args *common.EvalArgs) ([]int64, error) {
	return []int64{common.NamespaceVisitor}, nil
}
