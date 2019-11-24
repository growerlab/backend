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

func trace(other error) error {
	return wrapError(other, other.Error())
}

func invalidParameterError(model, field, reason string) error {
	return New(InvalidParameter, model, field, reason)
}

func sqlError(sqlErr error) error {
	return wrapError(sqlErr, SqlError)
}

func notFoundError(model, field string) error {
	return New(NotFound, model, field)
}

func wrapError(err error, code string) error {
	if err == nil {
		return nil
	}
	e := New(code).(*Err)
	e.Err = jujuerr.Trace(err).(*jujuerr.Err)
	e.SetLocation(1)
	return e
}
