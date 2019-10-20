package main

import (
	"github.com/growerlab/backend/app"
)

func main() {
	err := app.Run(":8080")
	if err != nil {
		panic(err)
	}
}
