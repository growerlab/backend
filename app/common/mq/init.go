package mq

import (
	"github.com/growerlab/backend/app/model/db"
)

var mq *MessageQueue

func InitMQ() error {
	mq = &MessageQueue{
		memDB:  db.MemDB,
		stream: NewStream(db.MemDB),
	}
	return nil
}
