package app

import (
	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/controller"
)

func Run(addr string) error {
	engine := gin.New()
	engine.Any("/graphql", controller.GraphQL)
	return engine.Run(addr)
}
