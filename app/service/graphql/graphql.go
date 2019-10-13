package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/growerlab/backend/app/service/graphql/think"
	"github.com/growerlab/backend/app/service/graphql/think/types"
)

func Do(session types.Session, query string) *graphql.Result {
	schema := think.BuildSchema(session)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	return result
}
