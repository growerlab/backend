package events

import (
	"encoding/json"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/common/mq"
)

const (
	EmailName = "send_email"
)

type EmailPayload struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Body   string `json:"body,omitempty"`
	IsHtml bool   `json:"is_html,omitempty"`
}

var _ mq.Consumer = (*Email)(nil)

func NewEmail() *Email {
	return &Email{}
}

type Email struct {}

func (e *Email) Name() string {
	return EmailName
}

func (e *Email) Consume(payload *mq.Payload) error {
	p := new(EmailPayload)
	err := json.Unmarshal([]byte(payload.Values[DefaultField].(string)), p)
	if err != nil {
		return errors.Trace(err)
	}
	return e.Send(p)
}

func (e *Email) Send(payload *EmailPayload) error {
	// TODO 发送邮件的具体逻辑(调用其他的smtp发送库)
	return nil
}
