package job

import "encoding/json"

type EmailPayload struct {
	From   string
	To     string
	Body   string
	IsHtml bool
}

type PushPayloadFunc func(jobName string, payload []byte) error

func NewEmail(pushFunc PushPayloadFunc) *Email {
	return &Email{
		pushPayloadFunc: pushFunc,
	}
}

type Email struct {
	pushPayloadFunc PushPayloadFunc
}

func (e *Email) Name() string {
	return "send_email"
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
	// TODO 发送邮件的具体逻辑
	return nil
}

func SendEmail(payload *EmailPayload) error {
	return nil
}
