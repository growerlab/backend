package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefine(t *testing.T) {
	e := InvalidParameterError("user", "name", "NotFoundField")
	assert.Equal(t, "InvalidParameter.user.name.NotFoundField", e.Error(), "error!")
}

func TestSqlError(t *testing.T) {
	e := Sql(errors.New("sqlerr-test"))
	assert.Equal(t, "SqlError", e.Error(), "want SqlError", "got "+e.Error())
}
