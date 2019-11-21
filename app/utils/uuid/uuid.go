package uuid

import "github.com/google/uuid"

import "strings"

func UUID() string {
	return string(fullUUID()[:8])
}

func UUIDv16() string {
	return string(fullUUID()[:16])
}

func fullUUID() string {
	s := strings.ToUpper(uuid.New().String())
	return strings.Replace(s, "-", "", -1)
}
