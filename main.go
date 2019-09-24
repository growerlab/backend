package main

import (
	"github.com/growerlab/codev-back/app"
)

func main() {
	err := app.Run(":8080")
	if err != nil {
		panic(err)
	}
}
