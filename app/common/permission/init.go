package permission

import (
	"github.com/growerlab/backend/app/common/permission/common"
	"github.com/growerlab/backend/app/common/permission/context"
	"github.com/growerlab/backend/app/common/permission/userdomain"
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
	return permHub.RegisterRules(rules)
}

func initUserDomains() error {
	userDomains := make([]common.UserDomainDelegate, 0)
	userDomains = append(userDomains, &userdomain.SuperAdmin{})
	userDomains = append(userDomains, &userdomain.Person{})
	userDomains = append(userDomains, &userdomain.RepositoryOwner{})
	return permHub.RegisterUserDomains(userDomains)
}

func initContexts() error {
	contexts := make([]common.ContextDelegate, 0)
	contexts = append(contexts, &context.Repository{})
	return permHub.RegisterContexts(contexts)
}
