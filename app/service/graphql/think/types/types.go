package types

import (
	"github.com/graphql-go/graphql"
	"github.com/growerlab/backend/app/common/env"
)

type Session interface {
	Env() env.Environment
	UserID() int // current user
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
	// 查询与创建更新等操作均通过此接口实现
	//
	// Resolvers(Session, *graphql.Object) (query graphql.Fields, mutation graphql.Fields)
	QueryResolver(*graphql.Object) *graphql.Field
	MutationResolvers(*graphql.Object) graphql.Fields

	//
	SetSession(Session)
}
