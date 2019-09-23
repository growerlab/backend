package controller

import (
	"github.com/gin-gonic/gin"
)

func GraphQL(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"hello": "world"})
}
