package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MaxGraphQLRequestBody = int64(1 << 20) // 1MB
)

func LimitGraphQLRequestBody(ctx *gin.Context) {
	if ctx.Request.ContentLength > MaxGraphQLRequestBody {
		ctx.AbortWithStatus(http.StatusRequestEntityTooLarge)
	}
}
