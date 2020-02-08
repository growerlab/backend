package userdomain

import (
	"github.com/growerlab/backend/app/common/ctx"
	"github.com/growerlab/backend/app/common/permission/common"
	"github.com/growerlab/backend/app/model/repository"
)

var _ common.UserDomainDelegate = (*RepositoryOwner)(nil)

type RepositoryOwner struct {
}

func (s *RepositoryOwner) Type() int {
	return common.UserDomainRepositoryOwner
}

func (s *RepositoryOwner) TypeLabel() string {
	return "repository_owner"
}

func (s *RepositoryOwner) Validate(ud *ctx.UserDomain) error {
	return nil
}

func (s *RepositoryOwner) BatchEval(db *ctx.DBContext, args *common.EvalArgs) ([]int64, error) {
	result := make([]int64, 0)
	repoID := args.Ctx.Param1
	repo, err := repository.GetRepository(db.Src, repoID)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return result, nil
	}

	ownerNamespaceID := repo.Owner().NamespaceID
	return []int64{ownerNamespaceID}, nil
}
