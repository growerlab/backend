package types

import (
	"github.com/graphql-go/graphql"
	"github.com/growerlab/backend/app/common/env"
)

type Session interface {
	Env() *env.Environment
	UserID() int64 // current user
}

type Object interface {
	//
	Name() string
	// Type() graphql.Type
	//
	Description() string
	// Fields
	//
	Types() graphql.Fields

	QueryResolvers(*graphql.Object) graphql.Fields
	MutationResolvers(*graphql.Object) graphql.Fields

	SetSession(Session)
}
