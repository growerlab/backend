package main

import (
	"net/http"

	"github.com/growerlab/backend/app"
	"github.com/growerlab/backend/app/common/notify"
	"github.com/growerlab/backend/app/utils/logger"
)

func main() {
	err := app.Run(":8080")
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info("bye.")
		} else {
			panic(err)
		}
	}
	<-notify.AllOfDone
}
