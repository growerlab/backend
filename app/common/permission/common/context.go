package common

import (
	"github.com/growerlab/backend/app/common/ctx"
)

func RepositoryContext(repositoryID int64) *ctx.Context {
	return &ctx.Context{
		Type:   ContextRepository,
		Param1: repositoryID,
	}
}
