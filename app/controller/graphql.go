package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/service/graphql"
	"github.com/growerlab/backend/app/utils/logger"
)

func GraphQL(ctx *gin.Context) {
	var req graphql.GQLRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var session *graphql.Session
	userID, err := GetUserID(ctx)
	if err == nil {
		session = graphql.NewSession(userID)
	}

	result := graphql.Do(session, &req)
	if result.HasErrors() {
		logger.GraphQLErrors(result.Errors)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, result)
		return
	}
	ctx.JSON(200, result)
}
