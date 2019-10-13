package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/service/graphql"
	"github.com/growerlab/backend/utils/logger"
	"github.com/growerlab/backend/utils/reader"
)

func GraphQL(ctx *gin.Context) {
	body := ctx.Request.Body
	if body == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer body.Close()

	b, err := reader.LimitReader(body, MaxGraphQLRequestBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, _ := UserID(ctx)
	session := graphql.NewSession(userID)

	result := graphql.Do(session, string(b))
	if result.HasErrors() {
		logger.GraphQLErrors(result.Errors)
	}
	ctx.JSON(200, result)
}
