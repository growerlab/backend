package controller

import (
	"github.com/gin-gonic/gin"
)

const (
	AuthUserToken = "auth-user-token"
)

func GetUserToken(ctx *gin.Context) string {
	token := getValueFromHeaderOrCookie(AuthUserToken, ctx)
	return token
}

func getValueFromHeaderOrCookie(k string, ctx *gin.Context) string {
	v := ctx.GetHeader(k)
	if len(v) == 0 {
		v, _ = ctx.Cookie(k)
	}
	return v
}
