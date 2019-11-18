package queue

import "github.com/growerlab/backend/app/common/notify"

var queueInstance *Queue

func InitQueue() error {
	queueInstance = NewQueue()
	go queueInstance.Start()

	notify.Subscribe(func() {
		queueInstance.Release()
	})
	return nil
}
