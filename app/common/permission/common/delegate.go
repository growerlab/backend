package common

import (
	"github.com/growerlab/backend/app/common/ctx"
)

type EvalArgs struct {
	// 上下文
	Ctx *ctx.Context
	// 大部分情况下，用户域依赖上下文
	UD *ctx.UserDomain
}

type ContextDelegate interface {
	Type() int
	TypeLabel() string
	// Validate 用于新增权限时，对context的参数进行验证，以确保其参数是正确或必填的
	Validate(c *ctx.Context) error
	// BatchEval 根据用户域返回相关的namespace IDs
	BatchEval(db *ctx.DBContext, args *EvalArgs) ([]int64, error)
}

type UserDomainDelegate interface {
	Type() int
	TypeLabel() string
	// Validate 用于新增权限时，对userDomain的参数进行验证，以确保其参数是正确或必填的
	Validate(ud *ctx.UserDomain) error
	// BatchEval 根据用户域返回相关的namespace IDs
	BatchEval(db *ctx.DBContext, args *EvalArgs) ([]int64, error)
}
