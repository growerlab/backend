package app

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/controller"
)

func Run(addr string) error {
	engine := gin.Default()

	api := engine.Group("/api")
	graphql := api.Group("/graphql")
	{
		graphql.Use(controller.LimitGraphQLRequestBody) // limit request body size
		graphql.POST("/", controller.GraphQL)
		graphql.GET("/playground", controller.GraphQLPlayground())
	}

	return engine.Run(addr)
}
