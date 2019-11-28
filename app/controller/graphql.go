package controller

import (
	"context"
	"runtime/debug"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/service/graphql"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/service/graphql/resolver"
	"github.com/growerlab/backend/app/utils/logger"
)

func GraphQL(ctx *gin.Context) {
	var userToken = GetUserToken(ctx)
	var session = graphql.NewSession(userToken, ctx)
	var graphqlOpts = make([]handler.Option, 0)
	
	sessionCtx := context.WithValue(ctx.Request.Context(), "session", session)
	ctx.Request = ctx.Request.WithContext(sessionCtx)
	
	// options
	// reqOpt := handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
	// 	return next(ctx)
	// })
	// graphqlOpts = append(graphqlOpts, reqOpt)

	// resOpt := handler.ResolverMiddleware(func(ctx context.Context, next gql.Resolver) (res interface{}, err error) {
	// 	logger.Info("----ResolverMiddleware")
	// 	return next(ctx)
	// })
	// graphqlOpts = append(graphqlOpts, resOpt)

	recoverOpt := handler.RecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
		logger.Error("graphql recover err: %v\n%+v", err, string(debug.Stack()))
		return errors.New(errors.GraphQLError("ERROR"))
	})
	graphqlOpts = append(graphqlOpts, recoverOpt)

	fn := handler.GraphQL(resolver.NewExecutableSchema(resolver.Config{Resolvers: &resolver.Resolver{}}), graphqlOpts...)
	fn.ServeHTTP(ctx.Writer, ctx.Request)
}

func GraphQLPlayground() gin.HandlerFunc {
	fn := handler.Playground("GraphQL playground", "/api/graphql")
	return func(ctx *gin.Context) {
		fn.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
