package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	AuthUserId = "auth-user-id"
)

func GetUserID(ctx *gin.Context) (int64, error) {
	id := getValueFromHeaderOrCookie(AuthUserId, ctx)
	userID, err := strconv.ParseInt(id, 10, 64)
	return userID, err
}

func getValueFromHeaderOrCookie(k string, ctx *gin.Context) string {
	v := ctx.GetHeader(k)
	if len(v) == 0 {
		v, _ = ctx.Cookie(k)
	}
	return v
}
