package mq

import (
	"encoding/json"
)

func GetInPayload[T any](pd *Payload, fd string) *T {
	raw := pd.Get(fd)
	if raw == nil {
		return nil
	}

	t := new(T)
	if v := pd.Get(fd); v != nil {
		raw := []byte(v.(string))
		if err := json.Unmarshal(raw, t); err != nil {
			return t
		}
	}
	return nil
}
