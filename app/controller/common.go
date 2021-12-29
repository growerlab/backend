package controller

import (
	"net/http"

	"github.com/growerlab/backend/app/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/common/errors"
)

const (
	MaxGraphQLRequestBody = int64(1 << 20) // 1MB
)

func LimitGETRequestBody(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		return
	}
	if ctx.Request.ContentLength > MaxGraphQLRequestBody {
		ctx.AbortWithStatus(http.StatusRequestEntityTooLarge)
		return
	}
}

func Render(c *gin.Context, payload interface{}, err error) {
	if err != nil {
		cerr := errors.Cause(err)
		if e, ok := cerr.(*errors.Result); ok {
			c.AbortWithStatusJSON(e.StatusCode, cerr)

			if e2 := errors.Cause(e.Err); e2 != nil {
				logger.Error("render2: %+v\n", e2)
			}
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, cerr)
		}
		logger.Error("render: %+v\n", cerr)

		return
	}
	if payload != nil {
		c.AbortWithStatusJSON(http.StatusOK, payload)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, &errors.Result{
		Code: "ok",
	})
}
