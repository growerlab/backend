package main

import (
	"github.com/growerlab/backend/app"
)

func main() {
	app.Init()
	err := app.Run(":8080")
	if err != nil {
		panic(err)
	}
}
