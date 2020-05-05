package events

import (
	"encoding/json"

	"github.com/growerlab/backend/app/common/errors"
)

const DefaultField = "default"

func BuildPushEmailMessage(payload *EmailPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.Trace(err)
	}
	e := NewEmail()
	_, err = e.courier.Add(EmailName, DefaultField, string(body))
	return err
}
