package userdomain

import (
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/permission/common"
	"github.com/growerlab/backend/app/model/user"
)

var _ common.UserDomainDelegate = (*Person)(nil)

type Person struct {
}

func (s *Person) Type() int {
	return common.UserDomainPerson
}

func (s *Person) TypeLabel() string {
	return "person"
}

func (s *Person) Validate(ud *ctx.UserDomain) error {
	if ud.Param <= 0 {
		return errors.Errorf("userdomain param is required")
	}
	return nil
}

func (s *Person) BatchEval(db *ctx.DBContext, args *common.EvalArgs) ([]int64, error) {
	u, err := user.GetUser(db.Src, args.UD.Param)
	if err != nil {
		return nil, err
	}
	return []int64{u.NamespaceID}, nil
}
