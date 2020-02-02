package ctx

import "github.com/growerlab/backend/app/common/permission/common"

type Context struct {
	Type   int   `json:"type"`
	Param1 int64 `json:"param1"`
	Param2 int64 `json:"param2"`
}

func RepositoryContext(repositoryID int64) *Context {
	return &Context{
		Type:   common.ContextRepository,
		Param1: repositoryID,
	}
}
