package userdomain

import (
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/permission/common"
	"github.com/growerlab/backend/app/model/user"
)

var _ common.UserDomainDelegate = (*SuperAdmin)(nil)

type SuperAdmin struct {
}

func (s *SuperAdmin) Type() int {
	return common.UserDomainSuperAdmin
}

func (s *SuperAdmin) TypeLabel() string {
	return "super_admin"
}

func (s *SuperAdmin) Validate(ud *ctx.UserDomain) error {
	return nil
}

func (s *SuperAdmin) BatchEval(db *ctx.DBContext, args *common.EvalArgs) ([]int64, error) {
	admins, err := user.ListAdminUsers(db.Src)
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, len(admins))
	for i := range admins {
		userIds[i] = admins[i].ID
	}
	return userIds, nil
}
