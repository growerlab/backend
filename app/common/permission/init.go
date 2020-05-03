package permission

import (
	"github.com/growerlab/backend/app/common/permission/common"
	"github.com/growerlab/backend/app/common/permission/context"
	"github.com/growerlab/backend/app/common/permission/userdomain"
	"github.com/growerlab/backend/app/model/db"
)

var permHub *Hub

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

func initUserDomains() error {
	userDomains := []common.UserDomainDelegate{
		&userdomain.SuperAdmin{},
		&userdomain.Person{},
		&userdomain.RepositoryOwner{},
	}
	return permHub.RegisterUserDomains(userDomains)
}

func initContexts() error {
	contexts := make([]common.ContextDelegate, 0)
	contexts = append(contexts, &context.Repository{})
	return permHub.RegisterContexts(contexts)
}

func initRules() error {
	rules := []*Rule{
		{
			Code:                  common.PermissionViewRepository,
			ConstraintUserDomains: []int{common.UserDomainPerson},
			BuiltInUserDomains:    []int{common.UserDomainRepositoryOwner},
		},
		{
			Code:                  common.PermissionCloneRepository,
			ConstraintUserDomains: []int{common.UserDomainPerson},
			BuiltInUserDomains:    []int{common.UserDomainRepositoryOwner},
		},
		{
			Code:                  common.PermissionPushRepository,
			ConstraintUserDomains: []int{common.UserDomainPerson},
			BuiltInUserDomains:    []int{common.UserDomainRepositoryOwner},
		},
	}
	return permHub.RegisterRules(rules)
}
