package permission

import (
	"github.com/growerlab/backend/app/common/context"
	"github.com/growerlab/backend/app/common/userdomain"
)

func CheckViewRepository(namespaceID *int64, repositoryID int64) error {
	c := RepositoryContext(repositoryID)
	return checkPermission(namespaceID, c, ViewRepository)
}

func CheckPushRepository(namespaceID int64, repositoryID int64) error {
	c := RepositoryContext(repositoryID)
	return checkPermission(&namespaceID, c, PushRepository)
}

func CheckCloneRepository(namespaceID *int64, repositoryID int64) error {
	c := RepositoryContext(repositoryID)
	return checkPermission(namespaceID, c, CloneRepository)
}

func checkPermission(namespaceID *int64, ctx *context.Context, code int) error {
	if namespaceID == nil || *namespaceID == 0 {
		var visitor = int64(userdomain.NamespaceVisitor)
		namespaceID = &visitor
	}
	return permHub.CheckCache(*namespaceID, ctx, code, true)
}
