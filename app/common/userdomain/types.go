package userdomain

import (
	"github.com/growerlab/backend/app/common/context"
)

type Evaluable interface {
	UserDomain() *UserDomain
	Context() *context.Context
	DB() *context.DBContext
}

type UserDomain struct {
	Type  int   `json:"type"`
	Param int64 `json:"param"`
}
