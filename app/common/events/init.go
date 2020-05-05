package events

import (
	"github.com/growerlab/backend/app/common/mq"
	"github.com/growerlab/backend/app/common/notify"
	"github.com/growerlab/backend/app/model/db"
)

var MQ *mq.MessageQueue

func InitMQ() error {
	MQ = mq.NewMessageQueue(db.MemDB)
	notify.Subscribe(func() {
		MQ.Release()
	})

	err := MQ.Register(NewEmail())
	if err != nil {
		return err
	}

	return MQ.Run()
}
