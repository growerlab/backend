package events

type Courier interface {
	Add(consumerName, msgField, msgBody string) (id string, err error)
}

type EventBase struct {
	courier Courier
}

func NewEventBase() EventBase {
	return EventBase{courier: MQ}
}
