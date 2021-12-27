package events

import (
	"encoding/json"
	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/mq"
)

func async(name, field string, t any) error {
	body, err := json.Marshal(t)
	if err != nil {
		return errors.Trace(err)
	}
	_, err = MQ.Add(name, field, string(body))
	return err
}

func getPayload[T any](pd *mq.Payload, fd string) *T {
	t := new(T)
	if v := pd.Get(fd); v != nil {
		raw := []byte(v.(string))
		if err := json.Unmarshal(raw, t); err != nil {
			return t
		}
	}
	return nil
}
