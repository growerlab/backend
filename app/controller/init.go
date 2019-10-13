package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	userIdKey = "user-id"
)

func UserID(ctx *gin.Context) (int, error) {
	id := getValueFromHeaderOrCookie(userIdKey, ctx)
	userID, err := strconv.Atoi(id)
	return userID, err
}

func getValueFromHeaderOrCookie(k string, ctx *gin.Context) string {
	v := ctx.GetHeader(k)
	if len(v) == 0 {
		v, _ = ctx.Cookie(k)
	}
	return v
}
