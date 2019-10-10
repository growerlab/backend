package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefine(t *testing.T) {
	e := InvalidParameterError("user", "name", "NotFoundField")
	assert.Equal(t, "InvalidParameter.user.name.NotFoundField", e.Error(), "error!")
}
