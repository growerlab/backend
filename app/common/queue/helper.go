package queue

import (
	"encoding/json"

	"github.com/growerlab/backend/app/common/errors"
)

func PushSendEmail(payload *EmailPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return errors.WithStack(err)
	}
	return queueInstance.PushPayload(EmailUUID, body)
}
