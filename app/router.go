package app

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/controller"
)

func Run(addr string) error {
	engine := gin.New()
	api := engine.Group("/api")
	{
		api.Use(controller.LimitGraphQLRequestBody) // limit request body size
		api.POST("/graphql", controller.GraphQL)
	}
	return engine.Run(addr)
}
