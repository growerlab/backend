package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/growerlab/backend/app/common/notify"
	"github.com/growerlab/backend/app/controller"
	"github.com/growerlab/backend/app/utils/conf"
)

func Run(addr string) error {
	engine := gin.Default()

	api := engine.Group("/api")
	graphql := api.Group("/graphql", controller.LimitGraphQLRequestBody)
	{
		graphql.POST("", controller.GraphQL)
		graphql.GET("/playground", controller.GraphQLPlayground())
	}

	return runServer(addr, engine)
}

func runServer(addr string, engine *gin.Engine) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	// 平滑关闭
	notify.Subscribe(func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(timeoutCtx)
	})

	// 是否debug
	gin.SetMode(gin.ReleaseMode)
	if conf.GetConf().Debug {
		gin.SetMode(gin.DebugMode)
	}
	return server.ListenAndServe()
}
