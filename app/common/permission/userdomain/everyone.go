package userdomain

import (
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/permission/common"
)

var _ common.UserDomainDelegate = (*Everyone)(nil)

type Everyone struct {
}

func (s *Everyone) Type() int {
	return common.UserDomainEveryone
}

func (s *Everyone) TypeLabel() string {
	return "everyone"
}

func (s *Everyone) Validate(*ctx.UserDomain) error {
	return nil
}

func (s *Everyone) BatchEval(db *ctx.DBContext, args *common.EvalArgs) ([]int64, error) {
	return []int64{common.NamespaceEveryone}, nil
}
