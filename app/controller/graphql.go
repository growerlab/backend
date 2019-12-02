package controller

import (
	"context"
	"runtime/debug"
	"strings"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/service/graphql"
	"github.com/growerlab/backend/app/service/graphql/resolver"
	"github.com/growerlab/backend/app/utils/logger"
	"github.com/vektah/gqlparser/gqlerror"
)

func GraphQL(ctx *gin.Context) {
	var userToken = GetUserToken(ctx)
	var session = graphql.NewSession(userToken, ctx)
	var graphqlOpts = make([]handler.Option, 0)

	sessionCtx := resolver.GetContextWithSession(ctx, session)
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

	errorOpt := handler.ErrorPresenter(func(gqlCtx context.Context, err error) *gqlerror.Error {
		logger.Error("graphql err presenter: %+v", err)

		retErr := gql.DefaultErrorPresenter(gqlCtx, err)
		retErr.Message = err.Error()

		// 只返回错误码，而不返回具体的错误信息
		msgParts := strings.Split(retErr.Message, ": ")
		if len(msgParts) > 0 {
			retErr.Message = msgParts[0]
		}

		return retErr
	})
	graphqlOpts = append(graphqlOpts, errorOpt)

	recoverOpt := handler.RecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
		logger.Error("graphql recover err: %v\n%+v", err, string(debug.Stack()))
		return errors.New(errors.GraphQLError())
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
