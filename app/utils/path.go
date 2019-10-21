package utils

import (
	"os"
	"path/filepath"
)

func BasePath() string {
	dir, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	return dir
}
