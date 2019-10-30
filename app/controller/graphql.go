package controller

import (
	"context"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/service/graphql"
	"github.com/growerlab/backend/app/service/graphql/resolver"
)

func GraphQL(ctx *gin.Context) {
	var session *graphql.Session
	userToken := GetUserToken(ctx)
	session = graphql.NewSession(userToken, ctx)

	sessionCtx := context.WithValue(ctx.Request.Context(), "session", session)
	ctx.Request = ctx.Request.WithContext(sessionCtx)

	fn := handler.GraphQL(resolver.NewExecutableSchema(resolver.Config{Resolvers: &resolver.Resolver{}}))
	fn.ServeHTTP(ctx.Writer, ctx.Request)
}

func GraphQLPlayground(ctx *gin.Context) {
	fn := handler.Playground("GraphQL playground", "/playground")
	fn.ServeHTTP(ctx.Writer, ctx.Request)
}
