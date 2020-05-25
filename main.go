package main

import (
	"fmt"
	"net/http"

	"github.com/growerlab/backend/app"
	"github.com/growerlab/backend/app/common/notify"
	"github.com/growerlab/backend/app/utils/conf"
	"github.com/growerlab/backend/app/utils/logger"
)

func main() {
	addr := fmt.Sprintf(":%d", conf.GetConf().Port)
	err := app.Run(addr)
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info("bye.")
		} else {
			panic(err)
		}
	}
	notify.Done()
}
