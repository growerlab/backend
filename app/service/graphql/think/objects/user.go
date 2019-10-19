package objects

import (
	"github.com/graphql-go/graphql"
	"github.com/growerlab/backend/app/service/graphql/think/types"
)

var _ types.Object = (*GQLUser)(nil)

type GQLUser struct {
	Base
}

func NewGQLUser() *GQLUser {
	return &GQLUser{}
}

func (u *GQLUser) Name() string {
	return "user"
}

func (u *GQLUser) Description() string {
	return "graphql user"
}

func (u *GQLUser) Types() graphql.Fields {
	fds := graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
	}
	return fds
}

func (u *GQLUser) QueryResolvers(fields *graphql.Object) (query graphql.Fields) {
	query = graphql.Fields{}
	list := u.BuildSimpleField("list", "list users", fields, u.listUsers)
	query[list.Name] = list

	me := u.BuildSimpleField("me", "about me", fields, u.me)
	query[me.Name] = list

	return
}

func (u *GQLUser) MutationResolvers(fields *graphql.Object) (mutations graphql.Fields) {
	mutations = graphql.Fields{}
	create := u.BuildSimpleField("create", "create user", fields, u.createUser)
	mutations[create.Name] = create
	return
}

func (u *GQLUser) me(p graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (u *GQLUser) listUsers(p graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (u *GQLUser) createUser(p graphql.ResolveParams) (interface{}, error) {
	// TODO call app/service/user/create_user.go
	return nil, nil
}
