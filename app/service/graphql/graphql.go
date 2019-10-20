package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/growerlab/backend/app/service/graphql/think"
	"github.com/growerlab/backend/app/service/graphql/think/types"
)

type GQLRequest struct {
	Query string `json:"query,omitempty"`
}

func Do(session types.Session, req *GQLRequest) *graphql.Result {
	schema := think.BuildSchema(session)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: req.Query,
	})
	return result
}
