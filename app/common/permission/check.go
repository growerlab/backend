package permission

import (
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/permission/common"
)

func CheckViewRepository(namespaceID int64, repositoryID int64) error {
	c := common.RepositoryContext(repositoryID)
	return checkPermission(namespaceID, c, common.PermissionViewRepository)
}

func checkPermission(namespaceID int64, ctx *ctx.Context, code int) error {
	return permHub.CheckCache(namespaceID, ctx, code, true)
}
