package env

import "errors"

var (
	ErrNotExists = errors.New("not exists")
)

type Environment struct {
	val map[string]string
}

func NewEnvironment() *Environment {
	return &Environment{
		val: make(map[string]string),
	}
}

func (e *Environment) Set(k, v string) {
	e.val[k] = v
}

func (e *Environment) Get(k string) (v string, err error) {
	if v, ok := e.val[k]; ok {
		return v, nil
	}
	return "", ErrNotExists
}
