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

	apiV1 := engine.Group("/api/v1", controller.LimitGETRequestBody)
	repositories := apiV1.Group("/repositories")
	{
		repositories.POST("/:namespace/create", controller.CreateRepository)
		repositories.GET("/:namespace/list", controller.Repositories)
		repositories.GET("/:namespace/:name", controller.Repository)
	}
	users := apiV1.Group("/users")
	{
		users.POST("/register", controller.RegisterUser)
		users.POST("/activate", controller.ActivateUser)
		users.GET("/login", controller.LoginUser)
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
