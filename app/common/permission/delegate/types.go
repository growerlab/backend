package delegate

import (
	"github.com/growerlab/backend/app/common/ctx"
)

type EvalArgs struct {
	Ctx *ctx.Context
}

type ContextDelegate interface {
	Type() int
	TypeLabel() string

	// BatchEval 返回权限用户ID
	BatchEval(args *EvalArgs) ([]int64, error)
}

type UserDomainDelegate interface {
	Type() int
	TypeLabel() string
	UserDomainKey() string
	BatchEval(args *EvalArgs) ([]int64, error)
}
