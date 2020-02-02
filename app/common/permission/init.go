package permission

import (
	"github.com/growerlab/backend/app/common/permission/common"
	"github.com/growerlab/backend/app/common/permission/context"
	"github.com/growerlab/backend/app/common/permission/delegate"
	"github.com/growerlab/backend/app/model/db"
)

var permHub *PermissionHub

func InitPermission() error {
	permHub = NewPermissionHub(db.DB, db.PermissionDB)

	if err := initRules(); err != nil {
		return err
	}
	if err := initUserDomains(); err != nil {
		return err
	}
	if err := initContexts(); err != nil {
		return err
	}
	return nil
}

func initRules() error {
	rules := make([]*Rule, 0)
	rules = append(rules, &Rule{
		Code:                  common.PermissionViewRepository,
		ConstraintUserDomains: []int{common.UserDomainPerson},
		BuiltInUserDomains:    []int{common.UserDomainRepositoryOwner},
	})
	err := permHub.RegisterRules(rules)
	return err
}

func initUserDomains() error {

	return nil
}

func initContexts() error {
	contexts := make([]delegate.ContextDelegate, 0)
	contexts = append(contexts, &context.Repository{})
	permHub.RegisterContexts(contexts)
	return nil
}
