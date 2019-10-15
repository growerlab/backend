package think

import (
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/growerlab/backend/app/service/graphql/think/objects"
	"github.com/growerlab/backend/app/service/graphql/think/types"
)

func InitGraphQL() error {
	_ = hub.Register(objects.NewGQLUser())
	return nil
}

func BuildSchema(session types.Session) graphql.Schema {

	querySet := graphql.Fields{}
	mutationSet := graphql.Fields{}

	hub.Each(func(obj types.Object) {
		typ := buildType(obj.Name(), obj.Types())

		queryFd := buildObject(obj, session)
		queryFd.Resolve = obj.QueryResolver(typ).Resolve
		querySet[queryFd.Name] = queryFd

		mutations := obj.MutationResolvers(typ)
		for mutName, mut := range mutations {
			action := strings.Join([]string{obj.Name(), mutName}, ".")
			mutationSet[action] = mut
		}
	})

	queryConfig := graphql.ObjectConfig{
		Name:   "Query",
		Fields: querySet,
	}
	mutationConfig := graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: mutationSet,
	}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(queryConfig),
		Mutation: graphql.NewObject(mutationConfig),
	}

	schema, _ := graphql.NewSchema(schemaConfig)
	return schema
}

func buildObject(obj types.Object, sess types.Session) *graphql.Field {
	obj.SetSession(sess)

	typ := buildType(obj.Name(), obj.Types())
	fd := &graphql.Field{}
	fd.Name = obj.Name()
	fd.Description = obj.Description()
	fd.Type = typ
	fd.Args = buildArgs(obj.Types())
	return fd
}

func buildType(typeName string, fds graphql.Fields) *graphql.Object {
	config := graphql.ObjectConfig{
		Name:   typeName,
		Fields: fds,
	}
	return graphql.NewObject(config)
}

func buildArgs(fds graphql.Fields) graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}
	for name, fd := range fds {
		args[name] = &graphql.ArgumentConfig{
			Type: fd.Type,
		}
	}
	return args
}
