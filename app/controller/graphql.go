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

	maxFormSize := int64(10 << 20) // 10 MB is a lot of text.
	b, err := reader.LimitReader(body, maxFormSize+1)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result := graphql.Do(string(b))
	if result.HasErrors() {
		logger.GraphQLErrors(result.Errors)
	}
	ctx.JSON(200, result)
}
