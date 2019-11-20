package queue

import (
	"encoding/json"
)

const EmailUUID = "send_email"

type EmailPayload struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Body   string `json:"body,omitempty"`
	IsHtml bool   `json:"is_html,omitempty"`
}

func NewEmail() *Email {
	return &Email{}
}

type Email struct {
}

func (e *Email) Name() string {
	return EmailUUID
}

func (e *Email) Eval(payload []byte) (requeue bool, err error) {
	p := new(EmailPayload)
	err = json.Unmarshal(payload, p)
	if err != nil {
		return false, err
	}
	return false, e.Send(p)
}

func (e *Email) Send(payload *EmailPayload) error {
	// TODO 发送邮件的具体逻辑(调用其他的smtp发送库)
	return nil
}
