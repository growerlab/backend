package errors

import (
	"strings"

	jujuerr "github.com/juju/errors"
)

// 定义错误类型

type Err struct {
	*jujuerr.Err
	code string
}

func (e *Err) Error() string {
	return e.code
}

func New(parts ...string) error {
	code := mustCode(parts)
	err := jujuerr.NewErr(code)
	err.SetLocation(1)
	return &Err{
		code: code,
		Err:  &err,
	}
}

func mustCode(parts []string) string {
	if len(parts) == 0 {
		panic("parts is required")
	}
	return strings.Join(parts, ".")
}

var trace = func(other error) error {
	if other == nil {
		return nil
	}
	e := New().(*Err)
	e.Err = jujuerr.Trace(other).(*jujuerr.Err)
	e.code = other.Error()
	e.SetLocation(1)
	return e
}

var invalidParameterError = func(model, field, reason string) error {
	return New(InvalidParameter, model, field, reason)
}
