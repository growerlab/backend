package db

type Event struct {
	Table  string
	Action string
}

type EventProcessor interface {
	TypeLabel() string
	OnEvent(*Event) error
}

type sqlEvent struct {
	processor map[string]*EventProcessor
}

func NewSqlEvent() *sqlEvent {
	return &sqlEvent{}
}

func (s *sqlEvent) Process() {

}
