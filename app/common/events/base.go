package events

type Courier interface {
	Add(consumerName, msgField, msgBody string) (id string, err error)
}
