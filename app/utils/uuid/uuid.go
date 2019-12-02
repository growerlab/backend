package uuid

import (
	"strings"

	"github.com/google/uuid"
)

func UUID() string {
	return string(fullUUID()[:8])
}

func UUIDv16() string {
	return string(fullUUID()[:16])
}

func BuildToken() string {
	return uuid.New().String()
}

func fullUUID() string {
	s := strings.ToUpper(uuid.New().String())
	return strings.Replace(s, "-", "", -1)
}
