package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MaxGraphQLRequestBody = int64(1 << 20) // 1MB
)

const (
	AuthUserToken = "auth-user-token"
)

func LimitGraphQLRequestBody(ctx *gin.Context) {
	if ctx.Request.ContentLength > MaxGraphQLRequestBody {
		ctx.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
}

func GetUserToken(ctx *gin.Context) string {
	token := getValueFromHeaderOrCookie(AuthUserToken, ctx)
	return token
}

func getValueFromHeaderOrCookie(k string, ctx *gin.Context) string {
	v := ctx.GetHeader(k)
	if len(v) < 5 {
		v, _ = ctx.Cookie(k)
	}
	return v
}
