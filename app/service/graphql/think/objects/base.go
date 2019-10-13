package objects

import (
	"github.com/graphql-go/graphql"
	"github.com/growerlab/backend/app/service/graphql/think/types"
)

// 提供一些基础、公共方法等
//
type Base struct {
	session types.Session
}

func (b *Base) BuildSimpleField(
	name string,
	desc string,
	types *graphql.Object,
	resolve graphql.FieldResolveFn,
) *graphql.Field {
	return &graphql.Field{
		Type:        types,
		Name:        name,
		Description: desc,
		Resolve:     resolve,
	}
}

func (b *Base) SetSession(s types.Session) {
	b.session = s
}
