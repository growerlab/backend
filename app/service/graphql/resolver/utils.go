package resolver

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/service/graphql"
)

const sessionName = "session"

func GetSession(ctx context.Context) *graphql.Session {
	sess, ok := ctx.Value(sessionName).(*graphql.Session)
	if ok {
		return sess
	}
	return nil
}

func GetContextWithSession(ctx *gin.Context, sess *graphql.Session) context.Context {
	return context.WithValue(ctx.Request.Context(), sessionName, sess)
}
